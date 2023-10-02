package detectcopies

import (
	"net"
	"time"

	"github.com/google/uuid"

	"detect-copies/internal/log"
)

const garbageCollectorTimeout = 5 * time.Second

// tableManager is a struct that manages a map of UUIDs to timestamps and another map of UUIDs to UDP addresses.
type tableManager struct {
	dateMap map[uuid.UUID]time.Time    // Keeps track of the timestamp when each UUID was added.
	copyMap map[uuid.UUID]*net.UDPAddr // Stores the associated UDP address for each UUID.
}

// newTableManager creates a new tableManager instance with initialized maps and starts a garbage collector routine.
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

// addCopy adds or updates a copy in the tableManager with the given UUID and UDP address.
func (tableManager *tableManager) addCopy(id uuid.UUID, address *net.UDPAddr) {
	_, ok := tableManager.dateMap[id]

	// Update the timestamp and associated address for the UUID.
	tableManager.dateMap[id] = time.Now()
	tableManager.copyMap[id] = address

	// Log a message indicating whether the copy was added or updated.
	if !ok {
		log.Log.Infof("Copy %v added", address)
	} else {
		log.Log.Infof("Copy %v updated", address)
	}
}

// removeExpiredCopy removes copies from the tableManager that have expired based on a 10-second threshold.
func (tableManager *tableManager) removeExpiredCopy() {
	for id, expiredTime := range tableManager.dateMap {
		if time.Since(expiredTime).Seconds() > 10 {
			address := tableManager.copyMap[id]

			// Delete the expired copy from both maps.
			delete(tableManager.dateMap, id)
			delete(tableManager.copyMap, id)

			// Log a message indicating the removal of the expired copy.
			log.Log.Infof("Copy %v removed", address)
		}
	}
}
