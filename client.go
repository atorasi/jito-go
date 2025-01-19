package jito

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/weeaa/jito-go/clients/searcher_client"
	"github.com/weeaa/jito-go/pkg"

	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type JitoManager struct {
	EnableJito bool
	Client     *http.Client
	RpcClient  *rpc.Client
	Proxy      string
	PrivateKey solana.PrivateKey

	SlotIndex uint64
	Epoch     uint64

	// jitoValidators is a map of validator IDs that are running Jito.
	JitoValidators map[string]bool

	// slotLeader maps slot to validator ID.
	SlotLeader map[uint64]string

	// voteAccounts maps nodeAccount to voteAccount
	VoteAccounts map[string]string

	Lock *sync.Mutex

	// tipInfo maps the latest tip information from Jito.
	TipInfo    *pkg.TipStreamInfo
	JitoTip    uint64
	JitoClient *searcher_client.Client
}

func NewJitoManager(rpcClient *rpc.Client, privateKey solana.PrivateKey, bribe float64, jitoProxy string) (*JitoManager, error) {
	jitoClient, err := searcher_client.NewNoAuth(
		context.Background(),
		NewYork.BlockEngineURL,
		rpcClient,
		rpcClient,
		// privateKey,
		nil,
	)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}
	if jitoProxy != "" {
		proxyParts := strings.Split(jitoProxy, "@")
		auth := strings.Split(proxyParts[1], ":")
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(&url.URL{
					Scheme: "http",
					Host:   proxyParts[0],
					User:   url.UserPassword(auth[0], auth[1]),
				}),
			},
		}
	}

	return &JitoManager{
		EnableJito: true,
		JitoTip:    uint64(bribe * float64(solana.LAMPORTS_PER_SOL)),
		Client:     httpClient,
		RpcClient:  rpcClient,
		JitoClient: jitoClient,

		JitoValidators: make(map[string]bool),
		SlotLeader:     make(map[uint64]string),
		VoteAccounts:   make(map[string]string),

		Lock: &sync.Mutex{},

		Proxy:      jitoProxy,
		PrivateKey: privateKey,
	}, nil
}

func (j *JitoManager) Start() error {
	if j.JitoClient == nil {
		return nil
	}

	j.manageTipStream()

	if err := j.fetchJitoValidators(); err != nil {
		return err
	}

	if err := j.fetchLeaderSchedule(); err != nil {
		return err
	}

	if err := j.fetchVoteAccounts(); err != nil {
		return err
	}

	if err := j.fetchEpochInfo(); err != nil {
		return err
	}

	go func() {
		for {
			if err := j.fetchEpochInfo(); err != nil {
				fmt.Println("Failed to fetch epoch info: %s", err)
			}

			time.Sleep(10 * time.Millisecond)
		}
	}()

	go func() {
		for {
			if err := j.fetchLeaderSchedule(); err != nil {
				fmt.Println("Failed to fetch epoch info: %s", err)
			}

			time.Sleep(10 * time.Minute)
		}
	}()

	go func() {
		for {
			if err := j.fetchJitoValidators(); err != nil {
				fmt.Println("Failed to fetch epoch info: %s", err)
			}

			time.Sleep(10 * time.Minute)
		}
	}()

	go func() {
		for {
			if err := j.fetchVoteAccounts(); err != nil {
				fmt.Println("Failed to fetch epoch info: %s", err)
			}

			time.Sleep(10 * time.Minute)
		}
	}()

	return nil
}

func (j *JitoManager) status(msg string) {
	// logger.Log.Info("Jito Manager %s", msg)
}

func (j *JitoManager) statusr(msg string) {
	// logger.Log.Info("Jito Manager (R) %s", msg)
}
