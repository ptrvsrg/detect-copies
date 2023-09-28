package detectcopies

import (
	"fmt"
	"net"
	"sync"

	"github.com/google/uuid"

	"detect-copies/internal/log"
)

func Start(addr string, port int) {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Log.Errorf("UUID creation error: %v", err)
		return
	}

	multicastAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Log.Errorf("Resolving error: %v", err)
		return
	}
	if !multicastAddr.IP.IsMulticast() {
		log.Log.Errorf("UDP address error: %v is not multicast group address", multicastAddr)
		return
	}

	tableManager := newTableManager()
	sender := NewSender(id, multicastAddr)
	receiver := NewReceiver(id, multicastAddr, tableManager)

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func() { receiver.Start(&waitGroup) }()
	go func() { sender.Start(&waitGroup) }()

	waitGroup.Wait()
}
