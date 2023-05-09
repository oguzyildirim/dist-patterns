package main

import (
	"fmt"
	"time"
)

type HeartBeatScheduler struct {
	action            func()
	heartBeatInterval time.Duration
	scheduledTask     *time.Timer
}

func NewHeartBeatScheduler(action func(), heartBeatInterval time.Duration) *HeartBeatScheduler {
	return &HeartBeatScheduler{
		action:            action,
		heartBeatInterval: heartBeatInterval,
	}
}

func (scheduler *HeartBeatScheduler) Start() {
	scheduler.scheduledTask = time.AfterFunc(scheduler.heartBeatInterval, func() {
		scheduler.action()
		scheduler.Start()
	})
}

func (scheduler *HeartBeatScheduler) Stop() {
	scheduler.scheduledTask.Stop()
}

type SendingServer struct {
	serverId string
}

func (sendingServer *SendingServer) SendHeartbeat() {
	fmt.Printf("Sending Heartbeat from server %v\n", sendingServer.serverId)
}

type FailureDetector interface {
	HeartBeatCheck()
	HeartBeatReceived(serverId string)
}

type AbstractFailureDetector struct {
	heartbeatScheduler *HeartBeatScheduler
}

func NewAbstractFailureDetector(heartbeatInterval time.Duration) *AbstractFailureDetector {
	return &AbstractFailureDetector{
		heartbeatScheduler: NewHeartBeatScheduler(func() {
			fmt.Println("AbstractFailureDetector HeartBeatCheck")
		}, heartbeatInterval),
	}
}

func (failureDetector *AbstractFailureDetector) Start() {
	failureDetector.heartbeatScheduler.Start()
}

func (failureDetector *AbstractFailureDetector) Stop() {
	failureDetector.heartbeatScheduler.Stop()
}

type ReceivingServer struct {
	failureDetector FailureDetector
}

func (receivingServer *ReceivingServer) HandleRequest(request string) {
	fmt.Printf("ReceivingServer HandleRequest: %v\n", request)
}

type TimeoutBasedFailureDetector struct {
	AbstractFailureDetector
	heartbeatReceivedTimes map[string]time.Time
	timeoutDuration        time.Duration
}

func NewTimeoutBasedFailureDetector(heartbeatInterval, timeoutDuration time.Duration) *TimeoutBasedFailureDetector {
	return &TimeoutBasedFailureDetector{
		AbstractFailureDetector: *NewAbstractFailureDetector(heartbeatInterval),
		heartbeatReceivedTimes:  make(map[string]time.Time),
		timeoutDuration:         timeoutDuration,
	}
}

func (failureDetector *TimeoutBasedFailureDetector) HeartBeatCheck() {
	fmt.Println("TimeoutBasedFailureDetector HeartBeatCheck")
	now := time.Now()
	for serverId, lastHeartbeatReceivedTime := range failureDetector.heartbeatReceivedTimes {
		timeSinceLastHeartbeat := now.Sub(lastHeartbeatReceivedTime)
		if timeSinceLastHeartbeat >= failureDetector.timeoutDuration {
			fmt.Printf("TimeoutBasedFailureDetector Marking server %v down\n", serverId)
		}
	}
}

func (failureDetector *TimeoutBasedFailureDetector) HeartBeatReceived(serverId string) {
	fmt.Printf("TimeoutBasedFailureDetector HeartBeatReceived from server %v\n", serverId)
	failureDetector.heartbeatReceivedTimes[serverId] = time.Now()
}
