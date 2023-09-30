package detectcopies

import (
	"github.com/google/uuid"
)

const version = "1.0.0"

type Message struct {
	Version string    `json:"version"`
	ID      uuid.UUID `json:"id"`
}

func NewMessage(id uuid.UUID) *Message {
	return &Message{
		Version: version,
		ID:      id,
	}
}
