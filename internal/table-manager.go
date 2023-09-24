package detectcopies

import (
	"github.com/google/uuid"
	"log/slog"
	"net"
	"time"
)

const garbageCollectorTimeout = 5 * time.Second

type tableManager struct {
	dateMap map[uuid.UUID]time.Time
	copyMap map[uuid.UUID]*net.UDPAddr
}

func newTableManager() *tableManager {
	tableManager := &tableManager{
		dateMap: make(map[uuid.UUID]time.Time),
		copyMap: make(map[uuid.UUID]*net.UDPAddr),
	}

	// Run garbage collector
	go func() {
		for {
			tableManager.removeExpiredCopy()
			time.Sleep(garbageCollectorTimeout)
		}
	}()

	return tableManager
}

func (tableManager *tableManager) addCopy(id uuid.UUID, address *net.UDPAddr) {
	_, ok := tableManager.dateMap[id]

	tableManager.dateMap[id] = time.Now()
	tableManager.copyMap[id] = address

	if !ok {
		slog.Info("Copy " + address.String() + " added")
	} else {
		slog.Info("Copy " + address.String() + " updated")
	}
}

func (tableManager *tableManager) removeExpiredCopy() {
	for id, expiredTime := range tableManager.dateMap {
		if time.Since(expiredTime).Seconds() > 10 {
			address := tableManager.copyMap[id]

			delete(tableManager.dateMap, id)
			delete(tableManager.copyMap, id)

			slog.Info("Copy " + address.String() + " removed")
		}
	}
}
