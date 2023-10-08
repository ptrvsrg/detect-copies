package detector

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"

	"detect-copies/internal/log"
)

// Constant defining the sender timeout duration
const senderTimeout = 3 * time.Second

// Sender represents a sender instance that sends UDP multicast messages.
type Sender struct {
	id            uuid.UUID    // Unique identifier for the sender
	multicastAddr *net.UDPAddr // Multicast address to send messages to
}

// NewSender creates a new Sender instance with the provided UUID and multicast address.
func NewSender(id uuid.UUID, multicastAddr *net.UDPAddr) Sender {
	return Sender{
		id:            id,
		multicastAddr: multicastAddr,
	}
}

// Start initiates the sender's operation.
func (sender Sender) Start(waitGroup *sync.WaitGroup) {
	defer waitGroup.Done() // Decrement the wait group when this function exits

	// Connect to the multicast group using UDP
	conn, err := net.DialUDP("udp", nil, sender.multicastAddr)
	if err != nil {
		log.Log.Errorf("Multicast sender creation error: %v", err)
		return
	}
	log.Log.Infof("Sender running on %v", conn.LocalAddr())
	defer conn.Close() // Close the connection when this function exits

	// Create a JSON message
	msg := NewMessage(sender.id)
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Log.Errorf("Marshalling error: %v", err)
		return
	}
	log.Log.Debugf("Message to send: %v", string(jsonMsg))

	for {
		// Write the JSON message to the multicast group
		_, err := conn.Write(jsonMsg)
		if err != nil {
			log.Log.Errorf("Writing error: %v", err)
			return
		}

		// Sleep for a defined timeout before sending the next message
		time.Sleep(senderTimeout)
	}
}
