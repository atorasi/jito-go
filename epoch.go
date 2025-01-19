package jito

import (
	"context"

	"github.com/gagliardetto/solana-go/rpc"
)

func (j *JitoManager) fetchEpochInfo() error {
	schedule, err := j.RpcClient.GetEpochInfo(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return err
	}

	j.SlotIndex = schedule.SlotIndex
	if j.Epoch != schedule.Epoch {
		if err = j.fetchLeaderSchedule(); err != nil {
			return err
		}

		j.Epoch = schedule.Epoch
	}

	return nil
}
