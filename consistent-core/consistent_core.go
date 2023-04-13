package main

type WatchEventType int

const (
	Created WatchEventType = iota
	Updated
	Deleted
)

type WatchEvent struct {
	Key   string
	Value string
	Type  WatchEventType
}

type ConsistentCore interface {
	put(key string, value string) error
	get(keyPrefix string) ([]string, error)
	registerLease(name string, ttl int64) error
	refreshLease(name string) error
	watch(name string, watchCallback func(event WatchEvent))
}

type RaftCore struct {
	// Fields for Raft consensus algorithm
	// ...
}

func (rc *RaftCore) put(key string, value string) error {
	// Implement Raft consensus algorithm for key-value storage
	// ...
	return nil
}

func (rc *RaftCore) get(keyPrefix string) ([]string, error) {
	// Implement Raft consensus algorithm for key-value storage
	// ...
	return []string{}, nil
}

func (rc *RaftCore) registerLease(name string, ttl int64) error {
	// Implement Raft consensus algorithm for lease registration
	// ...
	return nil
}

func (rc *RaftCore) refreshLease(name string) error {
	// Implement Raft consensus algorithm for lease refresh
	// ...
	return nil
}

func (rc *RaftCore) watch(name string, watchCallback func(event WatchEvent)) {
	// Implement Raft consensus algorithm for key-value storage
	// ...
}

type ConsistentCoreImpl struct {
	raftCore *RaftCore
	// Other fields for metadata storage, hierarchical storage, etc.
	// ...
}

func (cc *ConsistentCoreImpl) put(key string, value string) error {
	return cc.raftCore.put(key, value)
}

func (cc *ConsistentCoreImpl) get(keyPrefix string) ([]string, error) {
	return cc.raftCore.get(keyPrefix)
}

func (cc *ConsistentCoreImpl) registerLease(name string, ttl int64) error {
	return cc.raftCore.registerLease(name, ttl)
}

func (cc *ConsistentCoreImpl) refreshLease(name string) error {
	return cc.raftCore.refreshLease(name)
}

func (cc *ConsistentCoreImpl) watch(name string, watchCallback func(event WatchEvent)) {
	cc.raftCore.watch(name, watchCallback)
}
