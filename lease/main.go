package main

import (
	"fmt"
	"time"
)

func main() {
	// Create a new key-value store
	store := NewKVStore()

	// Register a new lease with a TTL of 10 seconds
	lease, err := store.RegisterLease("my-lease", 10*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Attach a key to the lease
	lease.AttachKey("my-key")

	// Refresh the lease every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err := store.RefreshLease("my-lease")
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Lease refreshed %q at %v /n", lease.name, time.Now())
	}
}
