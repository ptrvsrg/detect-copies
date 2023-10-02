package detectcopies

import (
	"net"
	"strconv"
	"sync"

	"github.com/google/uuid"

	"detect-copies/internal/log"
)

// Start is the entry point for the detectcopies package. It initializes and starts
// the sender and receiver components for a multicast network communication system.
func Start(addr string, port int) {
	// Generate a new UUID for this instance.
	id, err := uuid.NewUUID()
	if err != nil {
		log.Log.Errorf("UUID creation error: %v", err)
		return
	}

	// Resolve the multicast address using the provided IP address and port.
	multicastAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		log.Log.Errorf("Resolving error: %v", err)
		return
	}

	// Check if the resolved address is a multicast group address.
	if !multicastAddr.IP.IsMulticast() {
		log.Log.Errorf("UDP address error: %v is not a multicast group address", multicastAddr)
		return
	}

	// Create a new table manager to handle data tables.
	tableManager := newTableManager()

	// Create a sender instance with the generated UUID and multicast address.
	sender := NewSender(id, multicastAddr)

	// Create a receiver instance with the generated UUID, multicast address, and table manager.
	receiver := NewReceiver(id, multicastAddr, tableManager)

	// Create a wait group to ensure that both sender and receiver goroutines complete before exiting.
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	// Start the receiver and sender goroutines concurrently.
	go func() { receiver.Start(&waitGroup) }()
	go func() { sender.Start(&waitGroup) }()

	// Wait for both sender and receiver goroutines to complete before exiting.
	waitGroup.Wait()
}
