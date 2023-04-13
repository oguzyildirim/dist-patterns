package main

import (
	"log"
)

type Lease struct {
	Name         string
	TTL          int64
	ExpiresAt    int64
	AttachedKeys []string
}

func NewLease(name string, ttl int64, now int64) *Lease {
	return &Lease{
		Name:         name,
		TTL:          ttl,
		ExpiresAt:    now + ttl,
		AttachedKeys: make([]string, 0),
	}
}

func (l *Lease) Refresh(now int64) {
	l.ExpiresAt = now + l.TTL
	log.Printf("Refreshing lease %s. Expiration time is %d", l.Name, l.ExpiresAt)
}

func (l *Lease) AttachKey(key string) {
	l.AttachedKeys = append(l.AttachedKeys, key)
}

func (l *Lease) GetAttachedKeys() []string {
	return l.AttachedKeys
}
