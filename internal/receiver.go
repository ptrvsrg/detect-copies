package detectcopies

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"

	"detect-copies/internal/log"
)

const bufferSize = 1024

type Receiver struct {
	id            uuid.UUID
	multicastAddr *net.UDPAddr
	tableManager  *tableManager
}

func NewReceiver(id uuid.UUID, multicastAddr *net.UDPAddr, tableManager *tableManager) Receiver {
	return Receiver{
		id:            id,
		multicastAddr: multicastAddr,
		tableManager:  tableManager,
	}
}

func (receiver Receiver) Start(wg *sync.WaitGroup) {
	defer wg.Done()

	// Listen to multicast group
	listener, err := net.ListenUDP("udp", receiver.multicastAddr)
	if err != nil {
		log.Log.Errorf("Multicast listener creation error: %v", err)
		return
	}
	log.Log.Info("Receiver listening to " + receiver.multicastAddr.String())
	defer listener.Close()

	// Start receiving
	jsonBytes := make([]byte, bufferSize)
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

func (receiver Receiver) validateMessage(message Message) bool {
	if message.Name != name {
		return false
	}
	if message.Timestamp.After(time.Now()) {
		return false
	}
	if message.ID.String() == receiver.id.String() {
		return false
	}
	return true
}
