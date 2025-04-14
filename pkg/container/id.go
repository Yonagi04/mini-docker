package container

import (
	"fmt"
	"sync"
	"time"
)

type SnowflakeIDGenerator struct {
	mutex      sync.Mutex
	timestamp  int64
	workerID   int64
	datacenter int64
	sequence   int64
}

const (
	workerIDBits     = 5
	datacenterIDBits = 5
	sequanceBits     = 12

	maxWorkerId     = -1 ^ (-1 << workerIDBits)
	maxDatacenterID = -1 ^ (-1 << datacenterIDBits)
	maxSequance     = -1 ^ (-1 << sequanceBits)

	timeshift       = workerIDBits + datacenterIDBits + sequanceBits
	workershift     = datacenterIDBits + sequanceBits
	datacentershift = sequanceBits

	twepoch = int64(1704067200000)
)

type IDGenerator interface {
	NextID() string
}

var _ IDGenerator = (*SnowflakeIDGenerator)(nil)

func NewSnowflakeIDGenerator(workerID, datacenterID int64) (*SnowflakeIDGenerator, error) {
	if workerID > maxWorkerId || workerID < 0 {
		return nil, fmt.Errorf("worker ID %d 超过范围 [0, %d]", workerID, maxWorkerId)
	}
	if datacenterID > maxDatacenterID || datacenterID < 0 {
		return nil, fmt.Errorf("datacenter ID %d 超过范围 [0, %d]", datacenterID, maxDatacenterID)
	}

	return &SnowflakeIDGenerator{
		timestamp:  0,
		workerID:   workerID,
		datacenter: datacenterID,
		sequence:   0,
	}, nil
}

func (sg *SnowflakeIDGenerator) NextID() string {
	sg.mutex.Lock()
	defer sg.mutex.Unlock()

	now := time.Now().UnixMilli()

	if now < sg.timestamp {
		for now <= sg.timestamp {
			now = time.Now().UnixMilli()
		}
	}

	if sg.timestamp == now {
		sg.sequence = (sg.sequence + 1) & maxSequance
		if sg.sequence == 0 {
			for now <= sg.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		sg.sequence = 0
	}

	sg.timestamp = now

	id := ((now - twepoch) << timeshift) |
		(sg.datacenter << datacentershift) |
		(sg.workerID << workershift) |
		sg.sequence

	return fmt.Sprintf("%012x", id)
}
