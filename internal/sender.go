package detectcopies

import (
	"github.com/google/uuid"
	"log/slog"
	"net"
	"os"
	"sync"
	"time"
)

const senderTimeout = 3 * time.Second

type sender struct {
	ID            uuid.UUID
	MulticastAddr *net.UDPAddr
}

func (sender sender) start(waitGroup *sync.WaitGroup) {
	go func() {
		defer waitGroup.Done()

		conn, err := net.DialUDP("udp", nil, sender.MulticastAddr)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		slog.Info("sender running on " + conn.LocalAddr().String())

		message, err := sender.ID.MarshalBinary()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		for {
			_, err := conn.Write(message)
			if err != nil {
				slog.Error(err.Error())
			}

			time.Sleep(senderTimeout)
		}
	}()
}
