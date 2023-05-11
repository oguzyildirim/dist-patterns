package main

import "math"

type LamportClock struct {
	LatestTime int
}

func NewLamportClock(timestamp int) *LamportClock {
	return &LamportClock{LatestTime: timestamp}
}

func (c *LamportClock) Tick(requestTime int) int {
	c.LatestTime = int(math.Max(float64(c.LatestTime), float64(requestTime)))
	c.LatestTime++
	return c.LatestTime
}
