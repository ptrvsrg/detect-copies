package detectcopies

import (
	"encoding/json"
	"net"
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
		return
	}
	log.Log.Infof("Sender running on %v", conn.LocalAddr())
	defer conn.Close()

	// Create JSON message
	msg := NewMessage(name, sender.id)
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Log.Errorf("Marshalling error: %v", err)
		return
	}
	log.Log.Debugf("Message to send: %v", string(jsonMsg))

	for {
		// Write message
		_, err := conn.Write(jsonMsg)
		if err != nil {
			log.Log.Errorf("Writing error: %v", err)
			return
		}

		// Time out
		time.Sleep(senderTimeout)
	}
}
