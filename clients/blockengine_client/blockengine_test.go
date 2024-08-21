package blockengine_client

import (
	"context"
	"github.com/gagliardetto/solana-go"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/weeaa/jito-go"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	_, filename, _, _ := runtime.Caller(0)
	godotenv.Load(filepath.Join(filepath.Dir(filename), "..", "..", "..", "jito-go", ".env"))
	os.Exit(m.Run())
}

func Test_BlockEngineRelayerClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	privKey, ok := os.LookupEnv("PRIVATE_KEY")
	if !assert.True(t, ok, "getting PRIVATE_KEY from .env") {
		t.FailNow()
	}

	if !assert.NotEqualf(t, "", privKey, "PRIVATE_KEY shouldn't be equal to [%s]", privKey) {
		t.FailNow()
	}

	privateKey, err := solana.PrivateKeyFromBase58(privKey)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	validator, err := NewValidator(
		ctx,
		jito_go.Amsterdam.BlockEngineURL,
		privateKey,
		nil,
	)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	defer validator.GrpcConn.Close()

	t.Run("Validator_SubscribePackets", func(t *testing.T) {
		sub, err := validator.SubscribePackets()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		resp, err := sub.Recv()
		assert.NoError(t, err)

		assert.NotNil(t, resp.Msg)
	})

	t.Run("Validator_SubscribeBundles", func(t *testing.T) {
		sub, err := validator.SubscribeBundles()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		resp, err := sub.Recv()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		for _, bundle := range resp.Bundles {
			assert.NotNil(t, bundle.Bundle)
		}
	})

	t.Run("Validator_GetBlockBuilderFeeInfo", func(t *testing.T) {
		resp, err := validator.GetBlockBuilderFeeInfo()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		resp.String()
	})
}

func Test_BlockEngineValidatorClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	privKey, ok := os.LookupEnv("PRIVATE_KEY")
	if !assert.True(t, ok, "getting PRIVATE_KEY from .env") {
		t.FailNow()
	}

	privateKey, err := solana.PrivateKeyFromBase58(privKey)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	relayer, err := NewRelayer(
		ctx,
		jito_go.NewYork.RelayerURL,
		privateKey,
		nil,
	)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	defer relayer.GrpcConn.Close()

	t.Run("SubscribeAccountsOfInterest", func(t *testing.T) {
		sub, err := relayer.SubscribeAccountsOfInterest()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		resp, err := sub.Recv()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		for _, account := range resp.Accounts {
			assert.NotEqual(t, "", account)
		}
	})

	t.Run("StartExpiringPacketStream", func(t *testing.T) {
		sub, err := relayer.StartExpiringPacketStream()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		resp, err := sub.Recv()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		assert.NotNil(t, resp.Heartbeat)
	})

	t.Run("SubscribeProgramsOfInterest", func(t *testing.T) {
		sub, err := relayer.SubscribeProgramsOfInterest()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		resp, err := sub.Recv()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		for _, program := range resp.Programs {
			assert.NotEqual(t, "", program)
		}
	})
}
