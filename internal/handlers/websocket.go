package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"

	"github.com/sbilibin2017/bil-message-websocket/internal/models"
)

// RoomAccessValidator описывает сервис для проверки доступа к комнате
type RoomAccessValidator interface {
	ValidateAccess(ctx context.Context, userUUID, roomUUID uuid.UUID) bool
}

// TokenParser интерфейс для работы с токенами
type TokenParser interface {
	GetFromHeader(header http.Header) (string, error)
	Parse(ctx context.Context, tokenString string) (userUUID uuid.UUID, deviceUUID uuid.UUID, err error)
}

// NewRoomWebSocketHandler создаёт HTTP-обработчик для WebSocket соединений
// @Summary WebSocket соединение для комнаты
// @Description Устанавливает WebSocket соединение для отправки и получения зашифрованных сообщений в комнате
// @Tags Rooms
// @Param Authorization header string true "JWT токен в формате Bearer <token>"
// @Param room-uuid path string true "UUID комнаты"
// @Success 101 {string} string "Успешное подключение по WebSocket"
// @Failure 400 {string} string "Некорректные данные запроса или неверный токен"
// @Failure 403 {string} string "Нет доступа к комнате"
// @Failure 500 {string} string "Ошибка на сервере при установке соединения"
// @Router /websocket/{room-uuid} [get]
func NewRoomWebSocketHandler(
	svc RoomAccessValidator,
	tp TokenParser,
	nc *nats.Conn,
) http.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем и парсим токен
		tokenString, err := tp.GetFromHeader(r.Header)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		userUUID, _, err := tp.Parse(r.Context(), tokenString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Получаем UUID комнаты
		roomUUID, err := uuid.Parse(chi.URLParam(r, "room-uuid"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Проверяем доступ
		if !svc.ValidateAccess(r.Context(), userUUID, roomUUID) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Апгрейд до WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		// Чтение сообщений от клиента и публикация
		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				break // клиент отключился
			}

			var msgReq models.MessageRequest
			if err := json.Unmarshal(msgBytes, &msgReq); err != nil {
				continue
			}

			// Создаём объект для хранения сообщения
			message := models.MessageDB{
				MessageUUID:      uuid.New(),
				RoomUUID:         roomUUID,
				SenderUUID:       &userUUID,
				SenderDeviceUUID: msgReq.DeviceUUID,
				EncryptedText:    msgReq.EncryptedText,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}

			// Создаём объект для хранения ключа сообщения
			messageKey := models.MessageKeyDB{
				MessageKeyUUID: uuid.New(),
				MessageUUID:    message.MessageUUID,
				DeviceUUID:     msgReq.DeviceUUID,
				EncryptedKey:   msgReq.EncryptedKey,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			// Публикуем сообщение
			if payload, err := json.Marshal(message); err == nil {
				nc.Publish("messages", payload)
			}

			// Публикуем ключ
			if payload, err := json.Marshal(messageKey); err == nil {
				nc.Publish("message_keys", payload)
			}
		}
	}
}
