package chain

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"math/big"
	"math/rand"
	"strings"
)

const (
	// Erc20ABI erc20 abi TokenERC20
	Erc20ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"tokenOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	// Erc721ABI erc721 abi TokenERC721
	Erc721ABI = "[{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"
	// Erc1155ABI erc1155
	Erc1155ABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"URI\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	// ErcZeroAddress erc20 zero address
	ErcZeroAddress = "0x0000000000000000000000000000000000000000"
)

const (
	TransactionStatusFail = iota
	TransactionStatusSuccess
)

const (
	Erc20TransferEvent         = "Transfer"
	Erc721TransferEvent        = "Transfer"
	Erc1155TransferSingleEvent = "TransferSingle"
	Erc1155TransferBatchEvent  = "TransferBatch"
)

func init() {
	addChainFactory("evm", func(c Config) (BaseChain, error) {
		if c.ChainType != "evm" {
			return nil, errors.New("chain type error")
		}
		if len(c.RpcURLs) == 0 {
			return nil, errors.New("rpc urls error")
		}
		return &EvmChain{
			Name:        c.Name,
			ChainType:   c.ChainType,
			ChainID:     c.ChainID,
			ChainSymbol: c.ChainSymbol,
			Currency:    c.Currency,
			RpcURLs:     c.RpcURLs,
			GasPrice:    c.GasPrice,
		}, nil
	})

}

type EvmChain struct {
	Name        string
	ChainType   string
	ChainID     int64
	ChainSymbol string
	Currency    string
	RpcURLs     []string
	GasPrice    int64
	currentRpc  string
}

// selectRpc 随机选择一个rpc url
func (e *EvmChain) selectRpc() string {
	e.currentRpc = e.RpcURLs[rand.Intn(len(e.RpcURLs))]
	return e.currentRpc
}

func (e *EvmChain) _getTransaction(ctx context.Context, tx *types.Transaction) (*Transaction, error) {
	client, err := ethclient.Dial(e.selectRpc())
	if err != nil {
		return nil, err
	}
	txReceipt, err := client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return nil, err
	}
	var (
		erc20Abi, erc721Abi, erc1155Abi abi.ABI
	)
	erc20Abi, err = abi.JSON(strings.NewReader(Erc20ABI))
	if err != nil {
		zap.S().Errorf("erc20 abi error: %v", err)
		return nil, err
	}
	erc721Abi, err = abi.JSON(strings.NewReader(Erc721ABI))
	if err != nil {
		zap.S().Errorf("erc721 abi error: %v", err)
		return nil, err
	}
	erc1155Abi, err = abi.JSON(strings.NewReader(Erc1155ABI))
	if err != nil {
		zap.S().Errorf("erc1155 abi error: %v", err)
		return nil, err
	}
	if txReceipt.Status == TransactionStatusFail { // 交易失败
		return nil, nil
	}
	tx.Time()
	var (
		transaction = Transaction{
			BlockNumber: txReceipt.BlockNumber.Int64(),
			BlockHash:   txReceipt.BlockHash.String(),
			Hash:        tx.Hash().String(),
			Bills:       nil,
			Time:        tx.Time(),
		}
		transferBills = make([]*TransferBill, 0)
	)
	// 交易成功
	var fromAddress common.Address
	fromAddress, err = client.TransactionSender(ctx, tx, txReceipt.BlockHash, txReceipt.TransactionIndex)
	if err != nil {
		zap.S().Errorf("get transaction sender error: %v", err)
		return nil, err
	}
	// 1. 普通转账
	if tx.To() != nil && tx.Value() != nil && tx.Value().String() != "0" {
		transferBills = append(transferBills, &TransferBill{
			From:            fromAddress.String(),
			To:              tx.To().String(),
			ContractAddress: ErcZeroAddress,
			Index:           -1,
			Value:           tx.Value(),
		})
	}
	for index, log := range txReceipt.Logs {
		if log.Removed {
			continue
		}
		var (
			toAddress      common.Address
			value, tokenID *big.Int
		)
		zap.S().Infof("log topics 0: %v", log.Topics[0].String())
		if log.Topics[0].String() == erc20Abi.Events[Erc20TransferEvent].ID.String() && len(log.Data) > 0 {
			fromAddress = common.HexToAddress(log.Topics[1].String())
			toAddress = common.HexToAddress(log.Topics[2].String())
			value = new(big.Int).SetBytes(log.Data)
			transferBills = append(transferBills, &TransferBill{
				From:            fromAddress.String(),
				To:              toAddress.String(),
				ContractAddress: log.Address.String(),
				Index:           index,
				Value:           value,
			})
		} else if log.Topics[0].String() == erc721Abi.Events[Erc721TransferEvent].ID.String() {
			fromAddress = common.HexToAddress(log.Topics[1].String())
			toAddress = common.HexToAddress(log.Topics[2].String())
			tokenID = new(big.Int).SetBytes(log.Topics[3].Bytes())
			transferBills = append(transferBills, &TransferBill{
				From:            fromAddress.String(),
				To:              toAddress.String(),
				ContractAddress: log.Address.String(),
				Index:           index,
				TokenID:         tokenID,
				Value:           big.NewInt(1),
			})
		} else if log.Topics[0].String() == erc1155Abi.Events[Erc1155TransferSingleEvent].ID.String() {
			var event struct {
				Id    *big.Int
				Value *big.Int
			}
			err = erc1155Abi.UnpackIntoInterface(&event, Erc1155TransferSingleEvent, log.Data)
			if err != nil {
				zap.S().Errorf("erc1155 abi unpack error: %v", err)
				return nil, err
			}
			fromAddress = common.HexToAddress(log.Topics[2].String())
			toAddress = common.HexToAddress(log.Topics[3].String())
			transferBills = append(transferBills, &TransferBill{
				From:            fromAddress.String(),
				To:              toAddress.String(),
				ContractAddress: log.Address.String(),
				Index:           index,
				TokenID:         event.Id,
				Value:           event.Value,
			})
		} else if log.Topics[0].String() == erc1155Abi.Events[Erc1155TransferBatchEvent].ID.String() {
			var event struct {
				Ids    []*big.Int
				Values []*big.Int
			}
			err = erc1155Abi.UnpackIntoInterface(&event, Erc1155TransferBatchEvent, log.Data)
			if err != nil {
				zap.S().Errorf("erc1155 abi unpack error: %v", err)
				return nil, err
			}
			fromAddress = common.HexToAddress(log.Topics[2].String())
			toAddress = common.HexToAddress(log.Topics[3].String())
			tokenID = new(big.Int).SetBytes(log.Data)
			for i := 0; i < len(event.Ids); i++ {
				transferBills = append(transferBills, &TransferBill{
					From:            fromAddress.String(),
					To:              toAddress.String(),
					ContractAddress: log.Address.String(),
					Index:           index,
					BatchIndex:      i,
					TokenID:         event.Ids[i],
					Value:           event.Values[i],
				})
			}
		}
	}
	transaction.Bills = transferBills
	return &transaction, nil
}

func (e *EvmChain) GetTransaction(ctx context.Context, hash string) (*Transaction, error) {
	client, err := ethclient.Dial(e.selectRpc())
	if err != nil {
		return nil, err
	}
	tx, _, err := client.TransactionByHash(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, err
	}
	return e._getTransaction(ctx, tx)
}

func (e *EvmChain) GetBlock(ctx context.Context, number int64) (block *Block, err error) {
	client, err := ethclient.Dial(e.selectRpc())
	if err != nil {
		return nil, err
	}
	b, err := client.BlockByNumber(ctx, big.NewInt(number))
	if err != nil {
		if err != ethereum.NotFound {
			zap.S().Errorw("get block error", "number", number, "rpc", e.currentRpc, "error", err)
			err = ErrorNotFound
		}
		return nil, err
	}
	block = &Block{
		Number:     b.Number().Int64(),
		Hash:       b.Hash().String(),
		ParentHash: b.ParentHash().String(),
	}
	return block, nil
}

func (e *EvmChain) GetBlockTransactions(ctx context.Context, number int64) ([]*Transaction, error) {
	var bills = make([]*Transaction, 0)
	client, err := ethclient.Dial(e.selectRpc())
	if err != nil {
		return bills, err
	}
	b, err := client.BlockByNumber(ctx, big.NewInt(number))
	if err != nil {
		return bills, err
	}
	for _, tx := range b.Transactions() {
		_tx, _err := e._getTransaction(ctx, tx)
		if _err != nil {
			return nil, _err
		}
		bills = append(bills, _tx)
	}
	return bills, nil
}

func (e *EvmChain) GenerateAddress(ctx context.Context) (string, string, error) {
	// 创建一个新的随机私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyStr := hexutil.Encode(privateKeyBytes)
	if err != nil {
		return "", "", err
	}
	// 获取地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("can't change public key")
	}
	// 地址全部储存为小写方便处理
	address := strings.ToLower(crypto.PubkeyToAddress(*publicKeyECDSA).String())
	return address, privateKeyStr, nil
}
