package main

type Client struct {
	clock *LamportClock
}

func NewClient() *Client {
	return &Client{
		clock: NewLamportClock(1),
	}
}

func (c *Client) Write(server1, server2 *Server) {
	server1WrittenAt := server1.Write("name", "Alice", c.clock.LatestTime())
	c.clock.UpdateTo(server1WrittenAt)

	server2WrittenAt := server2.Write("title", "Microservices", c.clock.LatestTime())
	c.clock.UpdateTo(server2WrittenAt)

	if server2WrittenAt <= server1WrittenAt {
		panic("server2WrittenAt <= server1WrittenAt")
	}
}
