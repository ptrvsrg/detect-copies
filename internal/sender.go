package detectcopies

import (
	"log"
	"net"
	"sync"
	"time"
)

type Sender struct {
	id      string
	address *net.UDPAddr
}

func (sender Sender) start(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		conn, err := net.DialUDP("udp", nil, sender.address)
		if err != nil {
			log.Fatal(err)
		}

		for {
			_, err := conn.Write([]byte(sender.id))
			if err != nil {
				log.Fatal(err)
			}

			time.Sleep(3 * time.Second)
		}
	}()
}
