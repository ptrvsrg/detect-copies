package detectcopies

import (
	"log"
	"net"
	"time"
)

type TableManager struct {
	dateMap          map[string]time.Time
	copyMap          map[string]*net.UDPAddr
	garbageCollector bool
}

func (tableManager TableManager) addCopy(id string, address *net.UDPAddr) {
	_, ok := tableManager.dateMap[id]

	tableManager.dateMap[id] = time.Now()
	tableManager.dateMap[id].Add(10 * time.Second)
	tableManager.copyMap[id] = address

	if !ok {
		log.Printf("Copy %s added", address.String())
	} else {
		log.Printf("Copy %s updated", address.String())
	}
}

func (tableManager TableManager) removeExpiredCopy() {
	for id, expiredTime := range tableManager.dateMap {
		if expiredTime.Before(time.Now()) {
			address := tableManager.copyMap[id]

			delete(tableManager.dateMap, id)
			delete(tableManager.copyMap, id)

			log.Printf("Copy %s removed", address.String())
		}
	}
}

func (tableManager TableManager) startGarbageCollector() {
	if tableManager.garbageCollector {
		return
	}

	tableManager.garbageCollector = true
	go func() {
		for tableManager.garbageCollector {
			tableManager.removeExpiredCopy()
			time.Sleep(5 * time.Second)
		}
	}()
}

func (tableManager TableManager) stopGarbageCollector() {
	tableManager.garbageCollector = false
}
