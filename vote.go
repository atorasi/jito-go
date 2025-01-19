package jito

import (
	"context"

	"github.com/gagliardetto/solana-go/rpc"
)

func (j *JitoManager) fetchVoteAccounts() error {
	j.status("Fetching vote accounts")

	voteAccounts, err := j.RpcClient.GetVoteAccounts(context.Background(), nil)
	if err != nil {
		return err
	}

	j.buildVoteAccounts(voteAccounts.Current)

	return nil
}

func (j *JitoManager) buildVoteAccounts(voteAccounts []rpc.VoteAccountsResult) {
	j.Lock.Lock()
	defer j.Lock.Unlock()

	for _, account := range voteAccounts {
		j.VoteAccounts[account.NodePubkey.String()] = account.VotePubkey.String()
	}
}
