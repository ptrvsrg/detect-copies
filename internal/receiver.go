package detectcopies

import (
	"log"
	"net"
	"sync"
)

type Receiver struct {
	id           string
	address      *net.UDPAddr
	tableManager TableManager
}

func (receiver Receiver) start(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		listener, err := net.ListenUDP("udp", receiver.address)
		if err != nil {
			log.Fatal(err)
		}

		buffer := make([]byte, 1024)
		for {
			_, addr, err := listener.ReadFrom(buffer)
			if err != nil {
				log.Fatal(err)
			}

			udpAddr, err := net.ResolveUDPAddr(addr.Network(), addr.String())
			if err != nil {
				log.Fatal(err)
			}

			copyId := string(buffer)
			if !isValidId(receiver.id, copyId) {
				continue
			}

			receiver.tableManager.addCopy(copyId, udpAddr)
		}
	}()
}

func isValidId(id, copyId string) bool {
	return id == copyId
}
