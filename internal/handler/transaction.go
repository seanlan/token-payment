package handler

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"math"
	"strings"
	"sync"
	"token-payment/internal/chain"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/types"
)

// CheckTransactions
//
//	@Description: 检查交易
//	@param ctx
//	@param txs
//	@return err
func CheckTransactions(ctx context.Context, ch *sqlmodel.Chain, txs []*chain.Transaction) (err error) {
	for _, tx := range txs {
		if tx == nil {
			continue
		}
		err = CheckRechargeTransaction(ctx, ch, tx) // 检查充值交易
		if err != nil {
			return
		}
		err = CheckSendTxTransaction(ctx, ch, tx) // 检查转账交易
		if err != nil {
			return
		}
	}
	return
}

// CheckRechargeTransaction
//
//	@Description: 检查充值交易
//	@param ctx
//	@param ch
//	@param tx
//	@return err
func CheckRechargeTransaction(ctx context.Context, ch *sqlmodel.Chain, tx *chain.Transaction) (err error) {
	var (
		addressQ                = sqlmodel.ChainAddressColumns
		tokenQ                  = sqlmodel.ChainTokenColumns
		appChainQ               = sqlmodel.ApplicationChainColumns
		bills                   = make([]sqlmodel.ChainTx, 0)
		tokenContracts          = make([]string, 0)
		tokens                  = make([]sqlmodel.ChainToken, 0)
		tokenMap                = make(map[string]sqlmodel.ChainToken)
		fromAddressList         = make([]string, 0)
		toAddressList           = make([]string, 0)
		toAddresses             = make([]sqlmodel.ChainAddress, 0)
		toAddressMap            = make(map[string]sqlmodel.ChainAddress)
		chainFeeWallets         = make([]sqlmodel.ApplicationChain, 0)
		chainFeeWalletAddresses = make(map[string]sqlmodel.ApplicationChain)
	)
	for _, bill := range tx.Bills {
		tokenContracts = append(tokenContracts, bill.ContractAddress)
		toAddressList = append(toAddressList, strings.ToLower(bill.To))
		fromAddressList = append(fromAddressList, strings.ToLower(bill.From))
	}
	err = dao.FetchAllChainToken(ctx, &tokens, dao.And(
		tokenQ.ChainSymbol.Eq(ch.ChainSymbol),
		tokenQ.ContractAddress.Eq(tokenContracts),
	), 0, 0)
	if err != nil {
		return
	}
	for _, token := range tokens {
		tokenMap[token.ContractAddress] = token
	}
	err = dao.FetchAllChainAddress(ctx, &toAddresses, dao.And(
		addressQ.ChainSymbol.Eq(ch.ChainSymbol),
		addressQ.Watch.Eq(1), // 是否是监控地址
		addressQ.Address.In(toAddressList),
	), 0, 0)
	if err != nil || len(toAddresses) == 0 {
		return
	}
	for _, addr := range toAddresses {
		toAddressMap[addr.Address] = addr
	}
	err = dao.FetchAllApplicationChain(ctx, &chainFeeWallets, dao.And(
		appChainQ.ChainSymbol.Eq(ch.ChainSymbol),
		appChainQ.FeeWallet.In(fromAddressList),
	), 0, 0)
	if err != nil {
		return
	}
	for _, chainFeeWallet := range chainFeeWallets {
		chainFeeWalletAddresses[chainFeeWallet.FeeWallet] = chainFeeWallet
	}
	for _, bill := range tx.Bills {
		var (
			address    sqlmodel.ChainAddress
			token      sqlmodel.ChainToken
			fromAddr   = strings.ToLower(bill.From)
			toAddr     = strings.ToLower(bill.To)
			nftTokenID int64
			ok         bool
		)
		if address, ok = toAddressMap[toAddr]; !ok {
			// 不是监控地址
			continue
		}
		if token, ok = tokenMap[bill.ContractAddress]; !ok {
			// 不是监控token
			continue
		}
		if _, ok = chainFeeWalletAddresses[fromAddr]; ok {
			// 是手续费钱包转账的交易
			continue
		}
		if bill.TokenID != nil {
			nftTokenID = bill.TokenID.Int64()
		}
		chainTx := sqlmodel.ChainTx{
			ApplicationID:   address.ApplicationID,
			ChainSymbol:     ch.ChainSymbol,
			BlockNumber:     tx.BlockNumber,
			BlockHash:       tx.BlockHash,
			TxHash:          tx.Hash,
			FromAddress:     strings.ToLower(bill.From),
			ToAddress:       strings.ToLower(bill.To),
			ContractAddress: bill.ContractAddress,
			Symbol:          token.Symbol,
			Value:           float64(bill.Value.Int64()) / math.Pow10(int(token.Decimals)),
			TokenID:         nftTokenID,
			TxIndex:         int64(bill.Index),
			BatchIndex:      int64(bill.BatchIndex),
			TransferType:    int32(types.TransferTypeIn),
			CreateAt:        tx.Time.Unix(),
		}
		bills = append(bills, chainTx)
	}
	if len(bills) > 0 {
		_, err = dao.AddsChainTx(ctx, &bills)
	}
	return
}

// CheckSendTxTransaction
//
//	@Description: 检查转账交易
//	@param ctx
//	@param ch
//	@param tx
//	@return err
func CheckSendTxTransaction(ctx context.Context, ch *sqlmodel.Chain, tx *chain.Transaction) (err error) {
	// TODO: 检查提现交易
	var (
		sendTxQ = sqlmodel.ChainSendTxColumns
		sendTx  sqlmodel.ChainSendTx
		bills   = make([]sqlmodel.ChainTx, 0)
	)
	err = dao.FetchChainSendTx(ctx, &sendTx, dao.And(
		sendTxQ.ChainSymbol.Eq(ch.ChainSymbol),
		sendTxQ.TxHash.Eq(tx.Hash),
	))
	if err != nil {
		if errors.Is(err, dao.ErrNotFound) {
			return nil
		}
		return
	}
	if sendTx.ID == 0 {
		return
	}
	for _, bill := range tx.Bills {
		if bill.ContractAddress == sendTx.ContractAddress && // 合约地址相同
			strings.ToLower(bill.To) == strings.ToLower(sendTx.ToAddress) && // 提现地址相同
			strings.ToLower(bill.From) == strings.ToLower(sendTx.FromAddress) { // 发送地址是热钱包
			var nftTokenID int64
			if bill.TokenID != nil {
				nftTokenID = bill.TokenID.Int64()
			}
			chainTx := sqlmodel.ChainTx{
				ApplicationID:   sendTx.ApplicationID,
				ChainSymbol:     ch.ChainSymbol,
				BlockNumber:     tx.BlockNumber,
				BlockHash:       tx.BlockHash,
				TxHash:          tx.Hash,
				FromAddress:     strings.ToLower(bill.From),
				ToAddress:       strings.ToLower(bill.To),
				ContractAddress: bill.ContractAddress,
				Symbol:          sendTx.Symbol,
				Value:           sendTx.Value,
				TokenID:         nftTokenID,
				TxIndex:         int64(bill.Index),
				BatchIndex:      int64(bill.BatchIndex),
				TransferType:    sendTx.TransferType,
				CreateAt:        tx.Time.Unix(),
			}
			bills = append(bills, chainTx)
		}
	}
	// 更新提现订单状态
	sendTx.TransferSuccess = 1
	sendTx.Received = 1
	sendTx.ReceiveAt = tx.Time.Unix()
	_, err = dao.UpdateChainSendTx(ctx, &sendTx)
	if err != nil {
		return
	}
	if len(bills) > 0 {
		_, err = dao.AddsChainTx(ctx, &bills)
	}
	return
}

// UpdateTransactionsConfirm
//
//	@Description: 更新交易确认数
//	@param ctx
//	@param ch
//	@param currentBlockNumber
//	@return err
func UpdateTransactionsConfirm(ctx context.Context, ch *sqlmodel.Chain) (err error) {
	var (
		transQ = sqlmodel.ChainTxColumns
		txList = make([]sqlmodel.ChainTx, 0)
		wg     sync.WaitGroup
	)
	err = dao.FetchAllChainTx(ctx, &txList, dao.And(
		transQ.ChainSymbol.Eq(ch.ChainSymbol),
		transQ.Removed.Eq(0),
		transQ.BlockNumber.Lte(ch.RebaseBlock),
		transQ.Confirm.Lt(ch.Confirm)), 0, 0)
	if err != nil {
		return
	}
	for _, tx := range txList {
		wg.Add(1)
		go func(tx sqlmodel.ChainTx) {
			defer wg.Done()
			_ = UpdateTransactionConfirm(ctx, ch, &tx)
		}(tx)
	}
	return
}

// UpdateTransactionConfirm
//
//	@Description: 更新交易确认数
//	@param ctx
//	@param ch
//	@param tx
//	@return err
func UpdateTransactionConfirm(ctx context.Context, ch *sqlmodel.Chain, tx *sqlmodel.ChainTx) (err error) {
	var (
		blockQ = sqlmodel.ChainBlockColumns
		count  int64
	)
	count, err = dao.CountChainBlock(ctx, blockQ.BlockHash.Eq(tx.BlockHash))
	if err != nil || count == 0 { // 未找到区块
		return
	}
	if count == 0 {
		tx.Removed = 1
		tx.Confirmed = 0
	} else {
		confirm := ch.RebaseBlock - tx.BlockNumber
		tx.Confirm = int32(math.Min(float64(ch.Confirm), float64(confirm)))
		tx.NotifyNextTime = 0
		tx.NotifySuccess = 0
		tx.NotifyFailedTimes = 0
		if tx.Confirm >= ch.Confirm {
			tx.Confirmed = 1
		}
	}
	err = dao.GetDB(ctx).Transaction(func(db *gorm.DB) (txErr error) {
		c := dao.CtxWithTransaction(ctx, db)
		_, txErr = dao.UpdateChainTx(c, tx)
		if tx.TransferType != int32(types.TransferTypeIn) && tx.Confirmed == 1 {
			// 不是充值交易
			_ = ConfirmTransferOrder(c, ch, tx)
		}
		return
	})
	return
}

func ConfirmTransferOrder(ctx context.Context, ch *sqlmodel.Chain, tx *sqlmodel.ChainTx) (err error) {
	var (
		withdrawQ = sqlmodel.ApplicationWithdrawOrderColumns
		arrangeQ  = sqlmodel.ApplicationArrangeTxColumns
		feeQ      = sqlmodel.ApplicationArrangeFeeTxColumns
		sendTxQ   = sqlmodel.ChainSendTxColumns
		sendTx    sqlmodel.ChainSendTx
	)
	err = dao.FetchChainSendTx(ctx, &sendTx,
		dao.And(
			sendTxQ.ChainSymbol.Eq(tx.ChainSymbol),
			sendTxQ.TxHash.Eq(tx.TxHash),
			sendTxQ.TransferType.Eq(tx.TransferType)))
	if err != nil {
		return
	}
	switch types.TransferType(tx.TransferType) {
	case types.TransferTypeOut:
		// 提现交易
		_, err = dao.UpdatesApplicationWithdrawOrder(ctx, withdrawQ.SendTxID.Eq(sendTx.ID), dao.M{withdrawQ.Confirmed.FieldName: 1})
	case types.TransferTypeFee:
		// 手续费交易
		_, err = dao.UpdatesApplicationArrangeFeeTx(ctx, feeQ.SendTxID.Eq(sendTx.ID), dao.M{feeQ.Confirmed.FieldName: 1})
	case types.TransferTypeArrange:
		// 整理交易
		_, err = dao.UpdatesApplicationArrangeTx(ctx, arrangeQ.SendTxID.Eq(sendTx.ID), dao.M{arrangeQ.Confirmed.FieldName: 1})
	}
	return
}
