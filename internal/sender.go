package detectcopies

import (
	"net"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"

	"detect-copies/internal/log"
)

const senderTimeout = 3 * time.Second

type Sender struct {
	id            uuid.UUID
	multicastAddr *net.UDPAddr
}

func NewSender(id uuid.UUID, multicastAddr *net.UDPAddr) Sender {
	return Sender{
		id:            id,
		multicastAddr: multicastAddr,
	}
}

func (sender Sender) Start(waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	// Connect to multicast group
	conn, err := net.DialUDP("udp", nil, sender.multicastAddr)
	if err != nil {
		log.Log.Errorf("Multicast sender creation error: %v", err)
		os.Exit(1)
	}
	log.Log.Infof("Sender running on %v", conn.LocalAddr())

	// Create message
	message, err := sender.id.MarshalBinary()
	if err != nil {
		log.Log.Errorf("Unmarshalling error: %v", err)
		return
	}

	// Start sending
	for {
		// Write message
		_, err := conn.Write(message)
		if err != nil {
			log.Log.Errorf("Writing error: %v", err)
			return
		}

		// Time out
		time.Sleep(senderTimeout)
	}
}
