package jito

import (
	"context"

	"github.com/gagliardetto/solana-go/rpc"
)

func (j *JitoManager) IsJitoLeader() bool {
	j.Lock.Lock()
	defer j.Lock.Unlock()

	validator, ok := j.SlotLeader[j.SlotIndex]
	if !ok {
		return false
	}

	j.status("Checking if validator is a Jito leader: " + validator)
	isLeader := j.JitoValidators[j.VoteAccounts[validator]]

	return isLeader
}

func (j *JitoManager) fetchLeaderSchedule() error {
	j.status("Fetching leader schedule")

	scheduleResult, err := j.RpcClient.GetLeaderSchedule(context.Background())
	if err != nil {
		return err
	}

	j.buildLeaderSchedule(&scheduleResult)

	return nil
}

func (j *JitoManager) buildLeaderSchedule(scheduleResult *rpc.GetLeaderScheduleResult) {
	j.Lock.Lock()
	defer j.Lock.Unlock()

	j.SlotLeader = make(map[uint64]string)
	for validator, slots := range *scheduleResult {
		for _, slot := range slots {
			j.SlotLeader[slot] = validator.String()
		}
	}
}
