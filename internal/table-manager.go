package detectcopies

import (
	"net"
	"time"

	"github.com/google/uuid"

	"detect-copies/internal/log"
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

	// Start garbage collector
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
		log.Log.Infof("Copy %v added", address)
	} else {
		log.Log.Infof("Copy %v updated", address)
	}
}

func (tableManager *tableManager) removeExpiredCopy() {
	for id, expiredTime := range tableManager.dateMap {
		if time.Since(expiredTime).Seconds() > 10 {
			address := tableManager.copyMap[id]

			delete(tableManager.dateMap, id)
			delete(tableManager.copyMap, id)

			log.Log.Infof("Copy %v removed", address)
		}
	}
}
