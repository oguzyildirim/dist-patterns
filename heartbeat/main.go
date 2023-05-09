package main

import (
	"time"
)

func main() {
	sendingServer := SendingServer{serverId: "A"}
	receivingServer := ReceivingServer{}
	failureDetector := NewTimeoutBasedFailureDetector(time.Millisecond*100, time.Millisecond*1000)
	receivingServer.failureDetector = failureDetector

	failureDetector.Start()

	sendingServer.SendHeartbeat()
	receivingServer.HandleRequest("Request 1")
	time.Sleep(time.Millisecond * 500)

	sendingServer.SendHeartbeat()
	receivingServer.HandleRequest("Request 2")
	time.Sleep(time.Millisecond * 500)

	sendingServer.SendHeartbeat()
	receivingServer.HandleRequest("Request 3")
	time.Sleep(time.Millisecond * 500)

	failureDetector.Stop()
}
