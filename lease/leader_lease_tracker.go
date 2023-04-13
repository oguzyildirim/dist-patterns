package main

import (
	"errors"
	"log"
	"time"
)

type LeaderLeaseTracker struct {
	leases      map[string]*Lease
	executor    *time.Ticker
	quitChannel chan struct{}
}

func NewLeaderLeaseTracker(leases map[string]*Lease) *LeaderLeaseTracker {
	return &LeaderLeaseTracker{
		leases:      leases,
		executor:    time.NewTicker(5 * time.Millisecond),
		quitChannel: make(chan struct{}),
	}
}

func (t *LeaderLeaseTracker) AddLease(name string, ttl int64) error {
	if _, ok := t.leases[name]; ok {
		return errors.New("duplicate lease")
	}
	now := time.Now().UnixNano()
	lease := &Lease{
		Name:         name,
		TTL:          ttl,
		ExpiresAt:    now + ttl,
		AttachedKeys: make([]string, 0),
	}
	t.leases[name] = lease
	return nil
}

func (t *LeaderLeaseTracker) Start() {
	go func() {
		for {
			select {
			case <-t.executor.C:
				t.checkAndExpireLeases()
			case <-t.quitChannel:
				t.executor.Stop()
				return
			}
		}
	}()
}

func (t *LeaderLeaseTracker) Stop() {
	close(t.quitChannel)
}

func (t *LeaderLeaseTracker) checkAndExpireLeases() {
	for _, lease := range t.expiredLeases() {
		t.expireLease(lease)
		t.submitExpireLeaseRequest(lease)
	}
}

func (t *LeaderLeaseTracker) expiredLeases() []string {
	now := time.Now().UnixNano()
	expiredLeases := make([]string, 0)
	for leaseID, lease := range t.leases {
		if lease.ExpiresAt < now {
			expiredLeases = append(expiredLeases, leaseID)
		}
	}
	return expiredLeases
}

func (t *LeaderLeaseTracker) RefreshLease(name string) {
	if lease, ok := t.leases[name]; ok {
		now := time.Now().UnixNano()
		lease.ExpiresAt = now + lease.TTL
		log.Printf("Refreshing lease %v Expiration time is %v\n", name, lease.ExpiresAt)
	}
}

func (t *LeaderLeaseTracker) GetLeaseDetails(name string) *Lease {
	if lease, ok := t.leases[name]; ok {
		return lease
	}
	return nil
}

func (t *LeaderLeaseTracker) expireLease(name string) {
	if lease, ok := t.leases[name]; ok {
		log.Printf("Expiring lease %v\n", name)
		delete(t.leases, name)
		t.removeAttachedKeys(lease)
	}
}

func (t *LeaderLeaseTracker) removeAttachedKeys(lease *Lease) {
	if lease == nil {
		return
	}
	for _, attachedKey := range lease.AttachedKeys {
		log.Printf("Removing %v with lease %v\n", attachedKey, lease)
		// TODO: implement kvStore.Remove(attachedKey)
	}
}

func (t *LeaderLeaseTracker) submitExpireLeaseRequest(name string) {
	// TODO: submit request to notify followers about expired lease
}
