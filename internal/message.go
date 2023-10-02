package detectcopies

import (
	"github.com/google/uuid"
)

// VERSION represents the current version of the package.
const VERSION = "1.0.0"

// Message represents a data structure for a message with version and ID.
type Message struct {
	Version string    `json:"version"` // Version field to store the message version.
	ID      uuid.UUID `json:"id"`      // ID field to store the universally unique identifier (UUID) of the message.
}

// NewMessage creates a new Message with the provided UUID as ID and the package's VERSION.
func NewMessage(id uuid.UUID) *Message {
	return &Message{
		Version: VERSION,
		ID:      id,
	}
}
