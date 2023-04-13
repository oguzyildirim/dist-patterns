package main

import (
	"log"
	"time"
)

func main() {
	// Create a new lease tracker
	leases := make(map[string]*Lease)
	tracker := NewLeaderLeaseTracker(leases)

	// Add a lease to the tracker
	if err := tracker.AddLease("example_lease", 10*time.Second.Nanoseconds()); err != nil {
		log.Fatalf("Error adding lease: %v", err)
	}

	// Start the lease tracker
	tracker.Start()

	// Wait for some time
	time.Sleep(20 * time.Second)

	// Stop the lease tracker
	tracker.Stop()
}
