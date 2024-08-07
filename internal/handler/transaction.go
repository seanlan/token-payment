package handler

import (
	"context"
	"errors"
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
		addressList             = make([]string, 0)
		addresses               = make([]sqlmodel.ChainAddress, 0)
		addressMap              = make(map[string]sqlmodel.ChainAddress)
		chainFeeWallets         = make([]sqlmodel.ApplicationChain, 0)
		chainFeeWalletAddresses = make(map[string]sqlmodel.ApplicationChain)
	)
	for _, bill := range tx.Bills {
		tokenContracts = append(tokenContracts, bill.ContractAddress)
		addressList = append(addressList, strings.ToLower(bill.To))
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
	err = dao.FetchAllChainAddress(ctx, &addresses, dao.And(
		addressQ.ChainSymbol.Eq(ch.ChainSymbol),
		addressQ.Watch.Eq(1), // 是否是监控地址
		addressQ.Address.In(addressList),
	), 0, 0)
	if err != nil {
		return
	}
	for _, addr := range addresses {
		addressMap[addr.Address] = addr
	}
	err = dao.FetchAllApplicationChain(ctx, &chainFeeWallets, dao.And(
		appChainQ.ChainSymbol.Eq(ch.ChainSymbol),
		appChainQ.FeeWallet.In(addressList),
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
			toAddr     = strings.ToLower(bill.To)
			nftTokenID int64
			ok         bool
		)
		if address, ok = addressMap[toAddr]; !ok {
			// 不是监控地址
			continue
		}
		if token, ok = tokenMap[bill.ContractAddress]; !ok {
			// 不是监控token
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
		if _, ok = chainFeeWalletAddresses[toAddr]; ok {
			// 是手续费钱包转账的交易
			chainTx.TransferType = int32(types.TransferTypeFee)
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
		//appChainQ = sqlmodel.ApplicationChainColumns
		//appChain  sqlmodel.ApplicationChain
		bills = make([]sqlmodel.ChainTx, 0)
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
				TransferType:    int32(types.TransferTypeOut),
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
	} else {
		confirm := ch.RebaseBlock - tx.BlockNumber
		tx.Confirm = int32(math.Min(float64(ch.Confirm), float64(confirm)))
		tx.NotifyNextTime = 0
		tx.NotifySuccess = 0
		tx.NotifyFailedTimes = 0
	}
	_, err = dao.UpdateChainTx(ctx, tx)
	return
}
