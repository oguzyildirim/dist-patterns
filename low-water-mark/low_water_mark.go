package main

import (
	"sync"
	"time"
)

// LogCleaner represents the log cleaner that continuously checks which portion of the log can be safely discarded
type LogCleaner struct {
	wal           *WriteAheadLog
	config        *Config
	cleaningMutex sync.Mutex
}

type Config struct {
	CleanTaskIntervalMs int
	LogMaxDurationMs    int
}

// NewLogCleaner creates a new log cleaner instance
func NewLogCleaner(wal *WriteAheadLog, config *Config) *LogCleaner {
	return &LogCleaner{
		wal:           wal,
		config:        config,
		cleaningMutex: sync.Mutex{},
	}
}

// Start starts the log cleaner
func (lc *LogCleaner) Start() {
	go lc.scheduleLogCleaning()
}

// scheduleLogCleaning schedules the log cleaning task
func (lc *LogCleaner) scheduleLogCleaning() {
	for {
		time.Sleep(time.Duration(lc.config.CleanTaskIntervalMs) * time.Millisecond)
		lc.cleanLogs()
	}
}

// cleanLogs cleans the logs that can be safely discarded
func (lc *LogCleaner) cleanLogs() {
	lc.cleaningMutex.Lock()
	defer lc.cleaningMutex.Unlock()

	// Get the low water mark using snapshot-based approach
	snapshot := lc.wal.TakeSnapshot()
	lowWaterMark := snapshot.SnapshotIndex

	// Alternatively, get the low water mark using time-based approach
	// lowWaterMark := lc.getLowWaterMark()

	// Get the segments that can be safely discarded
	segments := lc.getSegmentsBefore(lowWaterMark)

	// Delete the segments that can be safely discarded
	for _, segment := range segments {
		lc.wal.DeleteSegment(segment)
	}
}

// getSegmentsBefore returns the segments that can be safely discarded based on the snapshot-based approach
func (lc *LogCleaner) getSegmentsBefore(snapshotIndex int64) []*WALSegment {
	var markedForDeletion []*WALSegment
	for _, segment := range lc.wal.sortedSavedSegments {
		if segment.LastLogEntryIndex < snapshotIndex {
			markedForDeletion = append(markedForDeletion, segment)
		}
	}
	return markedForDeletion
}

// getLowWaterMark returns the low water mark based on the time-based approach
func (lc *LogCleaner) getLowWaterMark() int64 {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	logMaxDurationMs := lc.config.LogMaxDurationMs
	lastLogEntryTimestamp := lc.wal.GetLastLogEntryTimestamp()
	timeElapsed := now - lastLogEntryTimestamp
	return lastLogEntryTimestamp - logMaxDurationMs
}

// WALSegment represents a segment of the write-ahead log
type WALSegment struct {
	// fields omitted for brevity
}

// Snapshot represents a snapshot taken by the storage engine
type Snapshot struct {
	Data          []byte
	SnapshotIndex int64
}

// WriteAheadLog represents the write-ahead log
type WriteAheadLog struct {
	// fields omitted for brevity
}
