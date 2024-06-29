package chain

import (
	"context"
	"errors"
	"math/big"
	"time"
)

var (
	chainFactory = make(map[string]newChainFunc)
)

type Config struct {
	Name        string
	ChainType   string
	ChainID     int64
	ChainSymbol string
	Currency    string
	RpcURLs     []string
	GasPrice    int64
}

var (
	ErrorNotFound = errors.New("not found")
	ErrorChain    = errors.New("chain error")
)

type newChainFunc func(Config) (BaseChain, error)

func addChainFactory(chainType string, f newChainFunc) {
	chainFactory[chainType] = f
}

func NewChain(c Config) (BaseChain, error) {
	if chainFactory[c.ChainType] == nil {
		return nil, ErrorChain
	}
	return chainFactory[c.ChainType](c)
}

type Block struct {
	Number     int64     // 区块高度
	Hash       string    // 区块hash
	ParentHash string    // 父区块hash
	ReceiveAt  time.Time // 区块时间
}

type Transaction struct {
	BlockNumber int64           // 区块高度
	BlockHash   string          // 区块hash
	Hash        string          // 交易hash
	Bills       []*TransferBill // 交易账单
	Time        time.Time       // 交易时间
}

type TransferBill struct {
	From            string   // 转出地址
	To              string   // 转入地址
	ContractAddress string   // 合约地址
	Value           *big.Int // 转账金额
	TokenID         *big.Int // token id (erc721)
	Index           int      // 交易序号
	BatchIndex      int      // 交易批次序号
}

type BaseChain interface {
	GetLatestBlockNumber(ctx context.Context) (int64, error)                        // 获取最新区块
	GetBlock(ctx context.Context, number int64) (*Block, error)                     // 获取区块
	GetBlockTransactions(ctx context.Context, number int64) ([]*Transaction, error) // 获取区块交易
	GetTransaction(ctx context.Context, hash string) (*Transaction, error)          // 获取交易
	GenerateAddress(ctx context.Context) (string, string, error)                    // 生成地址
}
