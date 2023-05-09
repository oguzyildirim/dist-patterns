package main

type ReplicationState struct {
	generation int64
	// other replication state data...
}

type ReplicatedLog struct {
	replicationState *ReplicationState
	// other log data...
}

func (rlog *ReplicatedLog) startLeaderElection() {
	rlog.replicationState.generation++
	rlog.registerSelfVote()
	rlog.requestVoteFromFollowers()
}

func (rlog *ReplicatedLog) appendToLocalLog(data []byte) int64 {
	generation := rlog.replicationState.generation
	return rlog.appendToLocalLogWithGeneration(data, generation)
}

func (rlog *ReplicatedLog) appendToLocalLogWithGeneration(data []byte, generation int64) int64 {
	logEntryId := rlog.getLastLogIndex() + 1
	logEntry := &WALEntry{Id: logEntryId, Data: data, Type: DATA, Generation: generation}
	rlog.writeEntryToWAL(logEntry)
	return logEntryId
}

func (rlog *ReplicatedLog) becomeFollower(leaderId int, generation int64) {
	rlog.resetReplicationState()
	rlog.replicationState.generation = generation
	rlog.replicationState.leaderId = leaderId
	rlog.transitionTo(FOLLOWER)
}

func (rlog *ReplicatedLog) handleReplicationRequest(request *ReplicationRequest) *ReplicationResponse {
	currentGeneration := rlog.replicationState.generation
	if request.generation < currentGeneration {
		return &ReplicationResponse{Status: FAILED, ServerId: rlog.serverId, Generation: currentGeneration, LastLogIndex: rlog.getLastLogIndex()}
	}
	// replicate entries...
	return &ReplicationResponse{Status: SUCCESS, ServerId: rlog.serverId, Generation: currentGeneration, LastLogIndex: rlog.getLastLogIndex()}
}

func (rlog *ReplicatedLog) handleReplicationResponse(response *ReplicationResponse) {
	if !response.IsSuccessful() {
		if response.generation > rlog.replicationState.generation {
			rlog.becomeFollower(LEADER_NOT_KNOWN, response.generation)
			return
		}
		// handle other cases...
	}
	// handle successful response...
}
