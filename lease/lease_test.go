package main

import (
	"sync"
	"testing"
	"time"
)

func TestLeaseAcquire(t *testing.T) {
	lease := NewLease("", 1*time.Second)

	if !lease.Acquire("holder1", 2*time.Second) {
		t.Error("Expected Acquire to return true")
	}

	if lease.Acquire("holder2", 2*time.Second) {
		t.Error("Expected Acquire to return false")
	}
}

func TestLeaseRelease(t *testing.T) {
	lease := NewLease("holder", 1*time.Second)

	lease.Release()

	if lease.Holder() != "" {
		t.Error("Expected Holder to return empty string after Release")
	}

	if !lease.Expiry().IsZero() {
		t.Error("Expected Expiry to return zero time after Release")
	}
}

func TestLeaseRenew(t *testing.T) {
	lease := NewLease("holder", 1*time.Second)

	lease.Renew()

	if !lease.Expiry().After(time.Now()) {
		t.Error("Expected Expiry to be in the future after Renew")
	}
}

func TestLeaseStartAndStopRenewalLoop(t *testing.T) {
	lease := NewLease("", 4*time.Second)

	lease.Acquire("holder", 8*time.Second)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		lease.StartRenewalLoop()
	}()

	// Wait for 3 minutes or until the test times out
	select {
	case <-time.After(4 * time.Second):
		t.Error("RenewalLoop did not stop within the expected time")
	case <-time.After(30 * time.Second):
		t.Error("Test timed out")
	}

	lease.StopRenewalLoop()

	wg.Wait()

	if !lease.Expiry().After(time.Now()) {
		t.Error("Expected Expiry to be in the future after RenewalLoop")
	}

	t.Logf("Holder: %s", lease.Holder())
	t.Logf("Expiry: %s", lease.Expiry())
}
