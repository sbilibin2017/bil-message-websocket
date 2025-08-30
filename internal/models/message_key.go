package models

import (
	"time"

	"github.com/google/uuid"
)

// MessageKeyDB представляет запись из таблицы message_keys
type MessageKeyDB struct {
	MessageKeyUUID uuid.UUID `json:"message_key_uuid" db:"message_key_uuid"` // Уникальный идентификатор записи (ключа)
	MessageUUID    uuid.UUID `json:"message_uuid" db:"message_uuid"`         // UUID сообщения, к которому относится ключ
	DeviceUUID     uuid.UUID `json:"device_uuid" db:"device_uuid"`           // UUID устройства, для которого зашифрован ключ
	EncryptedKey   string    `json:"encrypted_key" db:"encrypted_key"`       // Зашифрованный симметричный ключ сообщения
	CreatedAt      time.Time `json:"created_at" db:"created_at"`             // Время создания записи
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`             // Время последнего обновления записи
}
