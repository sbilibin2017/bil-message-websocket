package models

import (
	"time"

	"github.com/google/uuid"
)

// MessageRequest представляет данные сообщения, которые отправляет клиент через WebSocket
type MessageRequest struct {
	DeviceUUID    uuid.UUID `json:"device_uuid"`    // уникальный идентификатор устройства
	EncryptedText string    `json:"encrypted_text"` // Зашифрованный текст сообщения
	EncryptedKey  string    `json:"encrypted_key"`  // Зашифрованный ключ сообщения
}

// MessageDB представляет запись сообщения из таблицы messages
type MessageDB struct {
	MessageUUID      uuid.UUID  `json:"message_uuid" db:"message_uuid"`             // Уникальный идентификатор сообщения
	RoomUUID         uuid.UUID  `json:"room_uuid" db:"room_uuid"`                   // UUID комнаты, к которой относится сообщение
	SenderUUID       *uuid.UUID `json:"sender_uuid" db:"sender_uuid"`               // UUID отправителя (может быть NULL)
	SenderDeviceUUID uuid.UUID  `json:"sender_device_uuid" db:"sender_device_uuid"` // UUID устройства отправителя
	EncryptedText    string     `json:"encrypted_text" db:"encrypted_text"`         // Зашифрованный текст сообщения
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`                 // Время создания записи
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`                 // Время последнего обновления записи
}
