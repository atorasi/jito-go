package jito

import (
	"context"
	"fmt"

	utils "github.com/weeaa/jito-go/pkg"

	"github.com/gagliardetto/solana-go"
)

func (j *JitoManager) GenerateTipInstruction() (solana.Instruction, error) {
	tipAmount := j.generateTipAmount()
	j.status(fmt.Sprintf("Generating tip instruction for %.5f SOL", float64(tipAmount)/1e9))
	return j.JitoClient.GenerateTipRandomAccountInstruction(tipAmount, j.PrivateKey.PublicKey())
}

func (j *JitoManager) generateTipAmount() uint64 {
	return j.JitoTip
}

func (j *JitoManager) manageTipStream() {
	go func() {
		for {
			if err := j.subscribeTipStream(); err != nil {
				j.statusr("Error reading tip stream: " + err.Error())
			}
		}
	}()
}

func (j *JitoManager) subscribeTipStream() error {
	infoChan, errChan, err := utils.SubscribeTipStream(context.TODO())
	if err != nil {
		return err
	}

	for {
		select {
		case info := <-infoChan:
			j.status(fmt.Sprintf("Received tip stream (75th percentile=%.3fSOL, 95th percentile=%.3fSOL, 99th percentile=%.3fSOL)", info.LandedTips75ThPercentile, info.LandedTips95ThPercentile, info.LandedTips99ThPercentile))
			j.TipInfo = info
		case err = <-errChan:
			return err
		}
	}
}
