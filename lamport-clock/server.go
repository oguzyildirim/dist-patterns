package main

type Server struct {
	mvccStore *MVCCStore
	clock     *LamportClock
}

func NewServer(mvccStore *MVCCStore) *Server {
	return &Server{
		mvccStore: mvccStore,
		clock:     NewLamportClock(1),
	}
}

func (s *Server) Write(key, value string, requestTimestamp int) int {
	writeAtTimestamp := s.clock.Tick(requestTimestamp)
	s.mvccStore.Put(&VersionedKey{Key: key, Timestamp: writeAtTimestamp}, value)
	return writeAtTimestamp
}
