package token

import (
	"context"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"math/big"
)

func GetJettonBalance(ctx context.Context, api ton.APIClientWrapped, tokenAddress *address.Address, userAddress *address.Address) (*big.Int, error) {
	var b, err = api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch lastblock: %w", err)
	}
	c := cell.BeginCell().MustStoreAddr(userAddress).ToSlice()
	jettonAddrResult, err := api.RunGetMethod(ctx, b, tokenAddress, "get_wallet_address", c)
	if err != nil {
		return nil, fmt.Errorf("failed to call get_wallet_address: %w", err)
	}

	jettonWalletAddr := jettonAddrResult.MustSlice(0).MustLoadAddr()
	jettonWalletDataResult, err := api.RunGetMethod(ctx, b, jettonWalletAddr, "get_wallet_data")
	if err != nil {
		return nil, fmt.Errorf("failed to call get_wallet_data: %w", err)
	}
	ban := jettonWalletDataResult.MustInt(0)
	return ban, nil
}

func GetTonBalance(ctx context.Context, api ton.APIClientWrapped, userAddress *address.Address) (*big.Int, error) {
	var b, err = api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch lastblock: %w", err)
	}
	account, err := api.GetAccount(ctx, b, userAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to call getAccount: %w", err)
	}

	balance := account.State.Balance.Nano()

	return balance, nil
}
