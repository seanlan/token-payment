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
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/zap"
	"math/big"
	"strings"
	"token-payment/pkg/evmclient"
)

const (
	// EVMErc20ABI PolygonErc20ABI erc20 abi TokenERC20
	EVMErc20ABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenOwner\",\"type\":\"address\"},{\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"tokenOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"spender\",\"type\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	// EVMErc721ABI PolygonErc721ABI erc721 abi TokenERC721
	EVMErc721ABI = "[{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"
	// EVMErc1155ABI PolygonErc1155ABI erc1155
	EVMErc1155ABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"URI\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	// EVMErcZeroAddress PolygonErcZeroAddress erc20 zero address
	EVMErcZeroAddress = "0x0000000000000000000000000000000000000000"
)

const (
	EVMTransactionStatusFail = iota
	EVMTransactionStatusSuccess
)

const (
	EVMErc20TransferEvent         = "Transfer"
	EVMErc721TransferEvent        = "Transfer"
	EVMErc1155TransferSingleEvent = "TransferSingle"
	EVMErc1155TransferBatchEvent  = "TransferBatch"
)

func init() {
	addChainFactory("evm", func(ctx context.Context, c Config) (BaseChain, error) {
		if c.ChainType != "evm" {
			return nil, errors.New("chain type error")
		}
		if len(c.RpcURLs) == 0 {
			return nil, errors.New("rpc urls error")
		}
		equalizerName := "chain_" + c.ChainSymbol + "_equalizer"
		return &EvmChain{
			Name:        c.Name,
			ChainType:   c.ChainType,
			ChainID:     c.ChainID,
			ChainSymbol: c.ChainSymbol,
			Currency:    c.Currency,
			RpcURLs:     c.RpcURLs,
			GasPrice:    c.GasPrice,
			equalizer:   NewEqualizer(ctx, equalizerName, c.RpcURLs),
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
	equalizer   *Equalizer
}

// selectRpc
//
//	@Description: 获取一个rpc url
//	@receiver e
//	@param ctx
//	@return string
//	@return error
func (e *EvmChain) selectRpc(ctx context.Context) (string, error) {
	return e.equalizer.Get(ctx)
}

// GetLatestBlockNumber
//
//	@Description: 获取最新区块号
//	@receiver e
//	@param ctx
//	@return int64
//	@return error
func (e *EvmChain) GetLatestBlockNumber(ctx context.Context) (int64, error) {
	rpcUrl, err := e.selectRpc(ctx)
	if err != nil {
		return 0, err
	}
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		e.equalizer.Skip(ctx, rpcUrl)
		return 0, err
	}
	var latestID uint64
	latestID, err = client.BlockNumber(ctx)
	if err != nil {
		zap.S().Errorw("get block number error", "rpc", rpcUrl, "error", err)
		e.equalizer.Skip(ctx, rpcUrl) // 跳过
		return 0, err
	}
	return int64(latestID), nil
}

// _receiptToTransaction
//
//	@Description: 将交易回执转换为交易
//	@receiver e
//	@param ctx
//	@param txReceipt 交易回执
//	@param tx 交易
//	@return *Transaction 交易
//	@return error 错误
func (e *EvmChain) _receiptToTransaction(ctx context.Context, txReceipt *types.Receipt, tx *evmclient.Transaction) (*Transaction, error) {
	var (
		erc20Abi, erc721Abi, erc1155Abi abi.ABI
		err                             error
	)
	erc20Abi, err = abi.JSON(strings.NewReader(EVMErc20ABI))
	if err != nil {
		zap.S().Errorf("erc20 abi error: %v", err)
		return nil, err
	}
	erc721Abi, err = abi.JSON(strings.NewReader(EVMErc721ABI))
	if err != nil {
		zap.S().Errorf("erc721 abi error: %v", err)
		return nil, err
	}
	erc1155Abi, err = abi.JSON(strings.NewReader(EVMErc1155ABI))
	if err != nil {
		zap.S().Errorf("erc1155 abi error: %v", err)
		return nil, err
	}
	if txReceipt.Status == EVMTransactionStatusFail { // 交易失败
		return nil, nil
	}
	var (
		transaction = Transaction{
			BlockNumber: txReceipt.BlockNumber.Int64(),
			BlockHash:   txReceipt.BlockHash.String(),
			Hash:        tx.Hash.String(),
			Bills:       nil,
			Time:        tx.Time,
		}
		transferBills = make([]*TransferBill, 0)
	)
	// 交易成功
	var fromAddress common.Address
	fromAddress = tx.From
	// 1. 普通转账
	if tx.To.String() != EVMErcZeroAddress && tx.Value != nil && tx.Value.Int64() != 0 {
		transferBills = append(transferBills, &TransferBill{
			From:            strings.ToLower(fromAddress.String()),
			To:              strings.ToLower(tx.To.String()),
			ContractAddress: "",
			Index:           -1,
			Value:           tx.Value,
		})
	}
	// 2. 合约转账
	for index, log := range txReceipt.Logs {
		if log.Removed {
			continue
		}
		var (
			toAddress      common.Address
			value, tokenID *big.Int
		)
		if len(log.Topics) == 0 {
			continue
		}
		if log.Topics[0].String() == erc20Abi.Events[EVMErc20TransferEvent].ID.String() && len(log.Data) > 0 {
			if len(log.Topics) < 3 {
				continue
			}
			fromAddress = common.HexToAddress(log.Topics[1].String())
			toAddress = common.HexToAddress(log.Topics[2].String())
			value = new(big.Int).SetBytes(log.Data)
			transferBills = append(transferBills, &TransferBill{
				From:            strings.ToLower(fromAddress.String()),
				To:              strings.ToLower(toAddress.String()),
				ContractAddress: strings.ToLower(log.Address.String()),
				Index:           index,
				Value:           value,
			})
		} else if log.Topics[0].String() == erc721Abi.Events[EVMErc721TransferEvent].ID.String() {
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
		} else if log.Topics[0].String() == erc1155Abi.Events[EVMErc1155TransferSingleEvent].ID.String() {
			var event struct {
				Id    *big.Int
				Value *big.Int
			}
			err = erc1155Abi.UnpackIntoInterface(&event, EVMErc1155TransferSingleEvent, log.Data)
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
		} else if log.Topics[0].String() == erc1155Abi.Events[EVMErc1155TransferBatchEvent].ID.String() {
			var event struct {
				Ids    []*big.Int
				Values []*big.Int
			}
			err = erc1155Abi.UnpackIntoInterface(&event, EVMErc1155TransferBatchEvent, log.Data)
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

// GetTransaction
//
//	@Description: 获取交易
//	@receiver e
//	@param ctx
//	@param hash 交易hash
//	@return *Transaction 交易
//	@return error
func (e *EvmChain) GetTransaction(ctx context.Context, hash string) (*Transaction, error) {
	rpcUrl, err := e.selectRpc(ctx)
	if err != nil {
		return nil, err
	}
	cli := evmclient.NewEvmClient(rpcUrl)
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		e.equalizer.Skip(ctx, rpcUrl)
		return nil, err
	}
	tx, err := cli.TransactionByHash(ctx, common.HexToHash(hash))
	if err != nil {
		zap.S().Errorw("get transaction error", "hash", hash, "rpc", rpcUrl, "error", err)
		e.equalizer.Skip(ctx, rpcUrl) // 跳过
		return nil, err
	}
	txReceipt, err := client.TransactionReceipt(ctx, tx.Hash)
	if err != nil {
		zap.S().Errorw("get transaction receipt error", "hash", tx.Hash.String(), "rpc", rpcUrl, "error", err)
		e.equalizer.Skip(ctx, rpcUrl) // 跳过
		return nil, err
	}
	return e._receiptToTransaction(ctx, txReceipt, tx)
}

// GetBlock
//
//	@Description: 获取区块
//	@receiver e
//	@param ctx
//	@param number 区块号
//	@return block 区块
//	@return err 错误
func (e *EvmChain) GetBlock(ctx context.Context, number int64) (block *Block, err error) {
	rpcUrl, err := e.selectRpc(ctx)
	if err != nil {
		return nil, err
	}
	cli := evmclient.NewEvmClient(rpcUrl)
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		e.equalizer.Skip(ctx, rpcUrl)
		return nil, err
	}
	b, err := cli.BlockByNumber(ctx, number)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			err = ErrorNotFound
		} else {
			zap.S().Errorw("get block error", "number", number, "rpc", rpcUrl, "error", err)
			e.equalizer.Skip(ctx, rpcUrl)
		}
		return nil, err
	}
	var (
		receipts   []*types.Receipt
		ReceiptMap = make(map[string]*types.Receipt)
		trans      []*Transaction
	)
	receipts, err = client.BlockReceipts(ctx, rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(number)))
	if err != nil {
		zap.S().Errorw("get block receipts error", "number", number, "rpc", rpcUrl, "error", err)
		e.equalizer.Skip(ctx, rpcUrl)
		return nil, err
	}
	for _, receipt := range receipts {
		ReceiptMap[receipt.TxHash.String()] = receipt
	}
	for _, tx := range b.Transactions {
		receipt := ReceiptMap[tx.Hash.String()]
		if receipt == nil {
			hash := tx.Hash.String()
			zap.S().Errorw("receipt not found", "hash", hash)
			return nil, errors.New("receipt not found")
		}
		_tx, _err := e._receiptToTransaction(ctx, receipt, tx)
		if _err != nil {
			return nil, _err
		}
		trans = append(trans, _tx)
	}
	block = &Block{
		Number:       b.Number,
		Hash:         b.Hash.String(),
		ParentHash:   b.ParentHash.String(),
		ReceiveAt:    b.Timestamp,
		Transactions: trans,
	}
	return block, nil
}

func (e *EvmChain) GetBalance(ctx context.Context, address, contract string) (*big.Int, error) {
	rpcUrl, err := e.selectRpc(ctx)
	if err != nil {
		return nil, err
	}
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		e.equalizer.Skip(ctx, rpcUrl)
		return nil, err
	}
	if len(contract) > 0 {
		erc20Abi, _err := abi.JSON(strings.NewReader(EVMErc20ABI))
		if _err != nil {
			return nil, _err
		}
		targetAddress := common.HexToAddress(address)
		data, _err := erc20Abi.Pack("balanceOf", targetAddress)
		if _err != nil {
			return nil, _err
		}
		// 构造调用参数
		contractAddress := common.HexToAddress(contract)
		callMsg := ethereum.CallMsg{
			To:   &contractAddress,
			Data: data,
		}
		// 进行调用
		result, _err := client.CallContract(context.Background(), callMsg, nil)
		if err != nil {
			return nil, _err
		}
		// 解析返回的结果
		var balance *big.Int
		err = erc20Abi.UnpackIntoInterface(&balance, "balanceOf", result)
		if err != nil {
			return nil, err
		}
		return balance, nil
	} else {
		return client.PendingBalanceAt(ctx, common.HexToAddress(address))
	}
}

// GenerateAddress
//
//	@Description: 生成一个新的地址
//	@receiver e
//	@param ctx
//	@return address 地址
//	@return privateKey 私钥
//	@return err 错误
func (e *EvmChain) GenerateAddress(ctx context.Context) (address string, privateKey string, err error) {
	// 创建一个新的随机私钥
	_privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	privateKeyBytes := crypto.FromECDSA(_privateKey)
	privateKey = hexutil.Encode(privateKeyBytes)
	if err != nil {
		return "", "", err
	}
	// 获取地址
	publicKey := _privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("can't change public key")
	}
	// 地址全部储存为小写方便处理
	address = strings.ToLower(crypto.PubkeyToAddress(*publicKeyECDSA).String())
	return
}

func (e *EvmChain) GetNonce(ctx context.Context, address string) (uint64, error) {
	rpcUrl, err := e.selectRpc(ctx)
	if err != nil {
		return 0, err
	}
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		e.equalizer.Skip(ctx, rpcUrl)
		return 0, err
	}
	return client.PendingNonceAt(ctx, common.HexToAddress(address))
}

func (e *EvmChain) GenerateTransaction(ctx context.Context, order *TransferOrder) error {
	urlRpc, err := e.selectRpc(ctx)
	if err != nil {
		return err
	}
	// 创建一个rpc client
	client, err := ethclient.Dial(urlRpc)
	if err != nil {
		e.equalizer.Skip(ctx, urlRpc)
		return err
	}
	// 创建交易
	var (
		tx        *types.Transaction
		txData    []byte
		toAddress common.Address
		value     *big.Int
	)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		zap.S().Warnf("ItemSettle SuggestGasPrice err : %v", err)
		return err
	}
	if len(order.ContractAddress) > 0 { // ERC20 token 交易
		toAddress = common.HexToAddress(order.ContractAddress)
		erc20Abi, _err := abi.JSON(strings.NewReader(EVMErc20ABI))
		if _err != nil {
			return _err
		}
		txData, _err = erc20Abi.Pack("transfer", common.HexToAddress(order.To), order.Value)
		if _err != nil {
			return _err
		}
		value = big.NewInt(0)
	} else { // 普通转账
		toAddress = common.HexToAddress(order.To)
		value = order.Value
	}

	tx = types.NewTx(&types.LegacyTx{
		Nonce:    order.Nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      order.Gas,
		GasPrice: gasPrice,
		Data:     txData,
	})
	_privateKey, err := crypto.HexToECDSA(strings.Replace(order.FromPrivateKey, "0x", "", 1))
	if err != nil {
		return err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(e.ChainID)), _privateKey)
	order.TxHash = signedTx.Hash().String()
	order.GasPrice = gasPrice
	return nil
}

func (e *EvmChain) Transfer(ctx context.Context, order *TransferOrder) (string, error) {
	urlRpc, err := e.selectRpc(ctx)
	if err != nil {
		return "", err
	}
	// 创建一个新的私钥
	_privateKey, err := crypto.HexToECDSA(strings.Replace(order.FromPrivateKey, "0x", "", 1))
	if err != nil {
		return "", err
	}
	// 创建一个rpc client
	client, err := ethclient.Dial(urlRpc)
	if err != nil {
		e.equalizer.Skip(ctx, urlRpc)
		return "", err
	}
	// 创建交易
	var (
		tx        *types.Transaction
		txData    []byte
		toAddress common.Address
		value     *big.Int
	)
	if len(order.ContractAddress) > 0 { // ERC20 token 交易
		toAddress = common.HexToAddress(order.ContractAddress)
		erc20Abi, _err := abi.JSON(strings.NewReader(EVMErc20ABI))
		if _err != nil {
			return "", _err
		}
		txData, _err = erc20Abi.Pack("transfer", common.HexToAddress(order.To), order.Value)
		if _err != nil {
			return "", _err
		}
		value = big.NewInt(0)
	} else { // 普通转账
		toAddress = common.HexToAddress(order.To)
		value = order.Value
	}
	tx = types.NewTx(&types.LegacyTx{
		Nonce:    order.Nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      order.Gas,
		GasPrice: order.GasPrice,
		Data:     txData,
	})
	zap.S().Infof("tx hash: %s", tx.Hash().String())
	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(e.ChainID)), _privateKey)
	if err != nil {
		return "", err
	}
	zap.S().Infof("signed tx hash: %s", signedTx.Hash().String())
	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil && strings.Contains(err.Error(), "known transaction") {
		e.equalizer.Skip(ctx, urlRpc)
		return "", err
	}
	zap.S().Infof("send tx hash: %s", signedTx.Hash().String())
	return signedTx.Hash().String(), err
}
