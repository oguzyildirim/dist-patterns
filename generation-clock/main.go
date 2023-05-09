package main

import (
	"fmt"
	"time"
)

func main() {
	// Create a new replicated log
	log := NewReplicatedLog()

	// Start the leader election process
	go log.StartLeaderElection()

	// Wait for some time to simulate leader/follower disconnection
	time.Sleep(time.Second * 5)

	// Append some data to the log
	log.AppendToLocalLog([]byte("Hello, world!"))

	// Print the last log entry
	lastEntry := log.GetLastLogEntry()
	fmt.Printf("Last log entry: %s\n", string(lastEntry.Data))
}
