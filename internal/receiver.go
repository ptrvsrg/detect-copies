package detectcopies

import (
	"github.com/google/uuid"
	"log/slog"
	"net"
	"os"
	"sync"
)

const bufferSize = 1024

type receiver struct {
	ID            uuid.UUID
	MulticastAddr *net.UDPAddr
	TableManager  *tableManager
}

func (receiver receiver) start(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		listener, err := net.ListenUDP("udp", receiver.MulticastAddr)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		slog.Info("receiver listening to " + receiver.MulticastAddr.String())

		buffer := make([]byte, bufferSize)
		for {
			count, addr, err := listener.ReadFrom(buffer)
			if err != nil {
				slog.Error(err.Error())
			}

			udpAddr, err := net.ResolveUDPAddr(addr.Network(), addr.String())
			if err != nil {
				slog.Error(err.Error())
			}

			copyId := uuid.UUID{}
			err = copyId.UnmarshalBinary(buffer[:count])
			if err != nil {
				slog.Error(err.Error())
			}
			if receiver.ID == copyId {
				continue
			}

			receiver.TableManager.addCopy(copyId, udpAddr)
		}
	}()
}
