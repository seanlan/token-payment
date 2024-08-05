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

type newChainFunc func(context.Context, Config) (BaseChain, error)

func addChainFactory(chainType string, f newChainFunc) {
	chainFactory[chainType] = f
}

func NewChain(ctx context.Context, c Config) (BaseChain, error) {
	if chainFactory[c.ChainType] == nil {
		return nil, ErrorChain
	}
	return chainFactory[c.ChainType](ctx, c)
}

type Block struct {
	Number       int64          // 区块高度
	Hash         string         // 区块hash
	ParentHash   string         // 父区块hash
	ReceiveAt    time.Time      // 区块时间
	Transactions []*Transaction // 交易
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

type TransferOrder struct {
	TxHash          string   // 交易hash
	From            string   // 转出地址
	FromPrivateKey  string   // 转出地址私钥
	To              string   // 转入地址
	ContractAddress string   // 合约地址
	Value           *big.Int // 转账金额
	TokenID         *big.Int // token id (erc721)
	Gas             uint64   // gas
	GasPrice        *big.Int // gas price
	Nonce           uint64   // nonce
}

type BaseChain interface {
	GetLatestBlockNumber(ctx context.Context) (int64, error)                    // 获取最新区块
	GetBlock(ctx context.Context, number int64) (*Block, error)                 // 获取区块
	GetTransaction(ctx context.Context, hash string) (*Transaction, error)      // 获取交易
	GetBalance(ctx context.Context, address, contract string) (*big.Int, error) // 获取余额
	GenerateAddress(ctx context.Context) (string, string, error)                // 生成地址
	GetNonce(ctx context.Context, address string) (uint64, error)               // 获取nonce
	GenerateTransaction(ctx context.Context, order *TransferOrder) error        // 生成交易订单
	Transfer(ctx context.Context, order *TransferOrder) (string, error)         // 转账
}
