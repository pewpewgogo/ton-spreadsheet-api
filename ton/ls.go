package ton

import (
	"context"
	"os"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

func NewLsConnection(ctx context.Context) (ton.APIClientWrapped, error) {
	client := liteclient.NewConnectionPool()
	globalConfig, err := liteclient.GetConfigFromUrl(ctx, "https://ton.org/global.config.json")
	if err != nil {
		return nil, err
	}

	usePublicLs := os.Getenv("PRIVATE_LS_HOST") == ""
	if !usePublicLs {
		if err = client.AddConnection(ctx, os.Getenv("PRIVATE_LS_HOST"), os.Getenv("PRIVATE_LS_KEY")); err != nil {
			return nil, err
		}
	} else {
		err = client.AddConnectionsFromConfig(ctx, globalConfig)
		if err != nil {
			return nil, err
		}
	}

	api := ton.NewAPIClient(client, ton.ProofCheckPolicyUnsafe).WithRetry()
	api.SetTrustedBlockFromConfig(globalConfig)

	_, err = api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, err
	}

	return api, nil
}
