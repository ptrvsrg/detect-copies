package detectcopies

import (
	"net"
	"sync"

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

	// Create buffer
	buffer := make([]byte, bufferSize)

	// Start receiving
	for {
		// Read message
		count, addr, err := listener.ReadFrom(buffer)
		if err != nil {
			log.Log.Errorf("Reading error: %v", err)
			return
		}

		// Get source address
		udpAddr, err := net.ResolveUDPAddr(addr.Network(), addr.String())
		if err != nil {
			log.Log.Errorf("Resolving error: %v", err)
			continue
		}

		// Check UUID
		copyId := uuid.UUID{}
		err = copyId.UnmarshalBinary(buffer[:count])
		if err != nil {
			log.Log.Errorf("Unmarshalling error: %v", err)
			continue
		}
		if receiver.id == copyId {
			continue
		}

		// Add record to copy table
		receiver.tableManager.addCopy(copyId, udpAddr)
	}
}
