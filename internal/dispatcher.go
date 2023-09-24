package detectcopies

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net"
	"os"
	"sync"
)

const title = "\n" +
	"    ____       __            __                      _         \n" +
	"   / __ \\___  / /____  _____/ /_   _________  ____  (_)__  _____\n" +
	"  / / / / _ \\/ __/ _ \\/ ___/ __/  / ___/ __ \\/ __ \\/ / _ \\/ ___/\n" +
	" / /_/ /  __/ /_/  __/ /__/ /_   / /__/ /_/ / /_/ / /  __(__  )\n" +
	"/_____/\\___/\\__/\\___/\\___/\\__/   \\___/\\____/ .___/_/\\___/____/ \n" +
	"                                          /_/                  \n" +
	"Detect-Copies: 1.0.0"

func Run(addr string, port int) {
	fmt.Println(title)

	id, err := uuid.NewUUID()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	if !udpAddr.IP.IsMulticast() {
		slog.Error(udpAddr.String() + " is not multicast group address")
		os.Exit(1)
	}

	tableManager := newTableManager()
	sender := sender{
		ID:            id,
		MulticastAddr: udpAddr,
	}
	receiver := receiver{
		ID:            id,
		MulticastAddr: udpAddr,
		TableManager:  tableManager,
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	receiver.start(&waitGroup)
	sender.start(&waitGroup)

	waitGroup.Wait()
}
