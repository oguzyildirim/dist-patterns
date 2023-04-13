package main

import (
	"fmt"
	"sync"
)

func main() {
	// Initialize ConsistentCoreImpl with RaftCore
	rc := &RaftCore{}
	cc := &ConsistentCoreImpl{raftCore: rc}

	// Test ConsistentCore functionality
	err := cc.put("key1", "value1")
	if err != nil {
		fmt.Println("Error putting key-value pair:", err)
	}
	values, err := cc.get("key")
	if err != nil {
		fmt.Println("Error getting key:", err)
	} else {
		fmt.Println("Values with prefix 'key':", values)
	}
	err = cc.registerLease("lease1", 10)
	if err != nil {
		fmt.Println("Error registering lease:", err)
	}
	err = cc.refreshLease("lease1")
	if err != nil {
		fmt.Println("Error refreshing lease:", err)
	}
	cc.watch("key1", func(event WatchEvent) {
		fmt.Println("Watch event:", event)
	})

	// Wait for WatchEvents to be processed
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
