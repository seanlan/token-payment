package evmclient

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"time"
)

func heXtoInt64(hex string) int64 {
	n := new(big.Int)
	n.SetString(hex, 0)
	return n.Int64()
}

func hexToBigInt(hex string) *big.Int {
	n := new(big.Int)
	n.SetString(hex, 0)
	return n
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	// It's negative.
	if number.IsInt64() {
		return rpc.BlockNumber(number.Int64()).String()
	}
	// It's negative and large, which is invalid.
	return fmt.Sprintf("<invalid %d>", number)
}

type rpcBlock struct {
	BaseFeePerGas    string            `json:"baseFeePerGas"`
	Difficulty       string            `json:"difficulty"`
	ExtraData        string            `json:"extraData"`
	GasLimit         string            `json:"gasLimit"`
	GasUsed          string            `json:"gasUsed"`
	Hash             string            `json:"hash"`
	LogsBloom        string            `json:"logsBloom"`
	Miner            string            `json:"miner"`
	MixHash          string            `json:"mixHash"`
	Nonce            string            `json:"nonce"`
	Number           string            `json:"number"`
	ParentHash       string            `json:"parentHash"`
	ReceiptsRoot     string            `json:"receiptsRoot"`
	Sha3Uncles       string            `json:"sha3Uncles"`
	Size             string            `json:"size"`
	StateRoot        string            `json:"stateRoot"`
	Timestamp        string            `json:"timestamp"`
	TotalDifficulty  string            `json:"totalDifficulty"`
	Transactions     []*rpcTransaction `json:"transactions"`
	TransactionsRoot string            `json:"transactionsRoot"`
	Uncles           []interface{}     `json:"uncles"`
}

type rpcTransaction struct {
	BlockHash            string        `json:"blockHash"`
	BlockNumber          string        `json:"blockNumber"`
	From                 string        `json:"from"`
	Gas                  string        `json:"gas"`
	GasPrice             string        `json:"gasPrice"`
	MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerGas         string        `json:"maxFeePerGas,omitempty"`
	Hash                 string        `json:"hash"`
	Input                string        `json:"input"`
	Nonce                string        `json:"nonce"`
	To                   string        `json:"to"`
	TransactionIndex     string        `json:"transactionIndex"`
	Value                string        `json:"value"`
	Type                 string        `json:"type"`
	AccessList           []interface{} `json:"accessList,omitempty"`
	ChainId              string        `json:"chainId"`
	V                    string        `json:"v"`
	YParity              string        `json:"yParity,omitempty"`
	R                    string        `json:"r"`
	S                    string        `json:"s"`
}

type Block struct {
	BaseFeePerGas    int64
	Difficulty       int64
	ExtraData        []byte
	GasLimit         int64
	GasUsed          int64
	Hash             common.Hash
	Miner            common.Hash
	MixHash          string
	Nonce            int64
	Number           int64
	ParentHash       common.Hash
	ReceiptsRoot     common.Hash
	Sha3Uncles       common.Hash
	Size             int64
	StateRoot        common.Hash
	Timestamp        time.Time
	TotalDifficulty  int64
	Transactions     []*Transaction
	TransactionsRoot common.Hash
}

type Transaction struct {
	BlockHash        common.Hash
	BlockNumber      int64
	From             common.Address
	Gas              int64
	GasPrice         int64
	Hash             common.Hash
	Nonce            int64
	To               common.Address
	TransactionIndex int64
	Value            *big.Int
	Type             int64
	ChainId          int64
	Time             time.Time
}

func (t *rpcTransaction) ToTransaction() *Transaction {
	return &Transaction{
		BlockHash:        common.HexToHash(t.BlockHash),
		BlockNumber:      heXtoInt64(t.BlockNumber),
		From:             common.HexToAddress(t.From),
		Gas:              heXtoInt64(t.Gas),
		GasPrice:         heXtoInt64(t.GasPrice),
		Hash:             common.HexToHash(t.Hash),
		Nonce:            heXtoInt64(t.Nonce),
		To:               common.HexToAddress(t.To),
		TransactionIndex: heXtoInt64(t.TransactionIndex),
		Value:            hexToBigInt(t.Value),
		Type:             heXtoInt64(t.Type),
		ChainId:          heXtoInt64(t.ChainId),
	}
}

func (b *rpcBlock) ToBlock() *Block {
	block := &Block{
		BaseFeePerGas:    heXtoInt64(b.BaseFeePerGas),
		Difficulty:       heXtoInt64(b.Difficulty),
		ExtraData:        []byte(b.ExtraData),
		GasLimit:         heXtoInt64(b.GasLimit),
		GasUsed:          heXtoInt64(b.GasUsed),
		Hash:             common.HexToHash(b.Hash),
		Miner:            common.HexToHash(b.Miner),
		MixHash:          b.MixHash,
		Nonce:            heXtoInt64(b.Nonce),
		Number:           heXtoInt64(b.Number),
		ParentHash:       common.HexToHash(b.ParentHash),
		ReceiptsRoot:     common.HexToHash(b.ReceiptsRoot),
		Sha3Uncles:       common.HexToHash(b.Sha3Uncles),
		Size:             heXtoInt64(b.Size),
		StateRoot:        common.HexToHash(b.StateRoot),
		Timestamp:        time.Unix(heXtoInt64(b.Timestamp), 0),
		TotalDifficulty:  heXtoInt64(b.TotalDifficulty),
		TransactionsRoot: common.HexToHash(b.TransactionsRoot),
	}
	for _, t := range b.Transactions {
		tx := t.ToTransaction()
		tx.Time = block.Timestamp
		block.Transactions = append(block.Transactions, tx)
	}
	return block
}
