package main

import (
	"errors"
	"fmt"
	"time"
)

type ReplicatedKVStore struct {
	leases             map[string]*Lease
	leaderLeaseTracker *LeaderLeaseTracker
}

func NewReplicatedKVStore(leases map[string]*Lease) *ReplicatedKVStore {
	return &ReplicatedKVStore{
		leases: leases,
	}
}

func (r *ReplicatedKVStore) onBecomingLeader() {
	r.leaderLeaseTracker = NewLeaderLeaseTracker(r.leases)
	r.leaderLeaseTracker.Start()
}

func (r *ReplicatedKVStore) onCandidateOrFollower() {
	if r.leaderLeaseTracker != nil {
		r.leaderLeaseTracker.Stop()
	}
	// r.leaseTracker = FollowerLeaseTracker(r.leases)
}

func (r *ReplicatedKVStore) RegisterLease(name string, ttl int64) {
	if r.leaseExists(name) {
		fmt.Println(errors.New("lease exist"))
	}
	// Register lease, The request is complete only when the High-Water Mark reaches the log index of the request entry in the replicated log.
	return
}

func (r *ReplicatedKVStore) leaseExists(name string) bool {
	_, exists := r.leases[name]
	return exists
}

func (r *ReplicatedKVStore) applySetValueCommand(walEntryId int64, key string) {
	lease, exists := r.leases[key]
	if !exists {
		//The lease to attach is not available with the Consistent Core
		fmt.Println(errors.New("no lease exist"))
	}
	lease.AttachKey(key)
	r.leases[key] = &Lease{
		Name:         "name",
		TTL:          int64(5 * time.Second),
		ExpiresAt:    0,
		AttachedKeys: []string{key},
	}
}
