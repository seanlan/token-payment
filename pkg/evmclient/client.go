package evmclient

import (
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

type EvmClient struct {
	Url string
}

func NewEvmClient(url string) *EvmClient {
	return &EvmClient{Url: url}
}

func (e *EvmClient) client(ctx context.Context) (*rpc.Client, error) {
	return rpc.DialContext(ctx, e.Url)
}

func (e *EvmClient) BlockByNumber(ctx context.Context, number int64) (*Block, error) {
	cli, err := e.client(ctx)
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	var raw json.RawMessage
	err = cli.CallContext(ctx, &raw, "eth_getBlockByNumber", toBlockNumArg(big.NewInt(number)), true)
	if err != nil {
		return nil, err
	}
	var b rpcBlock
	err = json.Unmarshal(raw, &b)
	if err != nil {
		return nil, err
	}
	return b.ToBlock(), nil
}

func (e *EvmClient) TransactionByHash(ctx context.Context, hash common.Hash) (*Transaction, error) {
	var rt *rpcTransaction
	cli, err := e.client(ctx)
	if err != nil {
		return nil, err
	}
	err = cli.CallContext(ctx, &rt, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, err
	} else if rt == nil {
		return nil, ethereum.NotFound
	}
	return rt.ToTransaction(), nil
}
