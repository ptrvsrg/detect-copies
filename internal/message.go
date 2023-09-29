package detectcopies

import (
	"time"

	"github.com/google/uuid"
)

const name = "ptrvsrg"

type Message struct {
	Name      string    `json:"name"`
	ID        uuid.UUID `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

func NewMessage(name string, id uuid.UUID) *Message {
	return &Message{
		Name:      name,
		ID:        id,
		Timestamp: time.Now(),
	}
}
