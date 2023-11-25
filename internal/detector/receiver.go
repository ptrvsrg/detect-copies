package detector

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/google/uuid"

	"detect-copies/internal/log"
)

const BUFFER_SIZE = 1024

type Receiver struct {
	id            uuid.UUID     // Unique identifier for the receiver
	multicastAddr *net.UDPAddr  // Multicast UDP address to listen to
	tableManager  *tableManager // Reference to the tableManager to manage received copies
}

// NewReceiver creates a new Receiver instance with the given parameters.
func NewReceiver(id uuid.UUID, multicastAddr *net.UDPAddr, tableManager *tableManager) Receiver {
	return Receiver{
		id:            id,
		multicastAddr: multicastAddr,
		tableManager:  tableManager,
	}
}

// Start begins listening for incoming UDP messages on the multicast address.
func (receiver Receiver) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	// Listen to multicast group
	listener, err := net.ListenMulticastUDP("udp", nil, receiver.multicastAddr)
	if err != nil {
		log.Log.Errorf("Multicast listener creation error: %v", err)
		return
	}
	log.Log.Info("Receiver listening to " + receiver.multicastAddr.String())
	defer listener.Close()

	// Start receiving
	jsonBytes := make([]byte, BUFFER_SIZE)
	message := Message{}
	for {
		// Receive message
		count, addr, err := listener.ReadFromUDP(jsonBytes)
		if err != nil {
			log.Log.Errorf("Reading error: %v", err)
			return
		}
		log.Log.Debugf("Received message: %v", string(jsonBytes[:count]))

		// Convert JSON to Message
		err = json.Unmarshal(jsonBytes[:count], &message)
		if err != nil {
			log.Log.Errorf("Unmarshalling error: %v", err)
			continue
		}

		// Validate message
		if !receiver.validateMessage(message) {
			continue
		}

		// Add record to copy table
		receiver.tableManager.addCopy(message.ID, addr)
	}
}

// validateMessage checks if the received message is valid based on certain criteria.
func (receiver Receiver) validateMessage(message Message) bool {
	if message.Version != VERSION {
		return false
	}
	if message.ID.String() == receiver.id.String() {
		return false
	}
	return true
}
