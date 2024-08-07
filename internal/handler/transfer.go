package handler

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"math"
	"math/big"
	"time"
	"token-payment/internal/chain"
	"token-payment/internal/config"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/utils"
)

// SendTransferTransactions
//
//	@Description: 发送交易
//	@param ctx
//	@param ch
func SendTransferTransactions(ctx context.Context, ch *sqlmodel.Chain) {
	var (
		sendTxQ = sqlmodel.ChainSendTxColumns
		sendTxs = make([]sqlmodel.ChainSendTx, 0)
	)
	err := dao.FetchAllChainSendTx(ctx, &sendTxs, dao.And(
		sendTxQ.ChainSymbol.Eq(ch.ChainSymbol),
		sendTxQ.TransferSuccess.Eq(0),                   // 未发送成功
		sendTxQ.TransferNextTime.Lte(time.Now().Unix()), // 下次发送时间小于当前时间
	), 0, int(ch.Concurrent), sendTxQ.TransferFailedTimes.Asc(), sendTxQ.ID.Asc())
	if err != nil || len(sendTxs) == 0 {
		return
	}
	for _, sendTx := range sendTxs {
		_ = SendTransferTransaction(ctx, ch, &sendTx)
	}
}

// SendTransferTransaction
//
//	@Description: 发送转账交易
//	@param ctx
//	@param ch
//	@param tx
//	@return err
func SendTransferTransaction(ctx context.Context, ch *sqlmodel.Chain, tx *sqlmodel.ChainSendTx) (err error) {
	var (
		tokenQ   = sqlmodel.ChainTokenColumns
		token    sqlmodel.ChainToken
		addressQ = sqlmodel.ChainAddressColumns
		address  sqlmodel.ChainAddress
	)
	err = dao.FetchChainToken(ctx, &token,
		dao.And(
			tokenQ.ChainSymbol.Eq(ch.ChainSymbol),
			tokenQ.ContractAddress.Eq(tx.ContractAddress),
			tokenQ.Symbol.Eq(tx.Symbol)))
	if err != nil {
		return
	}
	err = dao.FetchChainAddress(ctx, &address, dao.And(
		addressQ.ChainSymbol.Eq(ch.ChainSymbol),
		addressQ.Address.Eq(tx.FromAddress),
	))
	if err != nil {
		return
	}
	pk, err := utils.AesDecrypt(address.EncKey, config.C.Secret)
	if err != nil {
		return
	}
	transferOrder := chain.TransferOrder{
		TxHash:          tx.TxHash,
		From:            tx.FromAddress,
		FromPrivateKey:  pk,
		To:              tx.ToAddress,
		ContractAddress: tx.ContractAddress,
		Value:           big.NewInt(int64(tx.Value * math.Pow10(int(token.Decimals)))),
		TokenID:         big.NewInt(tx.TokenID),
		Gas:             uint64(ch.Gas),
		GasPrice:        big.NewInt(tx.GasPrice),
		Nonce:           uint64(tx.Nonce),
	}
	client, err := GetChainRpcClient(ctx, ch)
	if err != nil {
		return
	}
	txHash, _err := client.Transfer(ctx, &transferOrder)
	if _err != nil {
		zap.S().Warnf("transfer failed, err: %v", _err)
		tx.TransferSuccess = 0
		tx.TransferFailedTimes++
		tx.TransferNextTime = time.Now().Unix() + 60*5 // 5分钟后重试
	} else {
		tx.TransferSuccess = 1
	}
	tx.TxHash = txHash                // 更新txHash
	tx.TransferAt = time.Now().Unix() // 记录发送时间 方便查询pending时间过长的交易 防止交易卡住 后面可以定时检索并处理
	_, err = dao.UpdateChainSendTx(ctx, tx)
	return
}

func GetTransferNonce(ctx context.Context, ch *sqlmodel.Chain, address string) (nonce uint64, err error) {
	var (
		txQ    = sqlmodel.ChainSendTxColumns
		lastTx sqlmodel.ChainSendTx
		client chain.BaseChain
	)
	err = dao.FetchChainSendTx(ctx, &lastTx, dao.And(
		txQ.ChainSymbol.Eq(ch.ChainSymbol),
		txQ.FromAddress.Eq(address),
	), txQ.Nonce.Desc())
	if err != nil {
		if errors.Is(err, dao.ErrNotFound) {
			client, err = GetChainRpcClient(ctx, ch)
			if err != nil {
				return 0, err
			}
			nonce, err = client.GetNonce(ctx, address)
			if err != nil {
				return
			}
		}
		return
	}
	return uint64(lastTx.Nonce + 1), nil
}
