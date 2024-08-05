package handler

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"token-payment/internal/chain"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
)

// GetChainRpcClient 获取链的rpc client
//
//	@Description:  获取链的rpc client
//	@param ctx
//	@param ch 链
//	@return client 链的rpc client
//	@return err
func GetChainRpcClient(ctx context.Context, ch *sqlmodel.Chain) (client chain.BaseChain, err error) {
	var (
		chainRPCs []sqlmodel.ChainRPC
		chainRpcQ = sqlmodel.ChainRPCColumns
		rpcUrls   []string
	)
	// 获取链的rpc
	err = dao.FetchAllChainRPC(ctx, &chainRPCs, dao.And(
		chainRpcQ.ChainSymbol.Eq(ch.ChainSymbol),
		chainRpcQ.Disable.Eq(0),
	), 0, 0)
	if err != nil {
		zap.S().Errorw("fetch all chain rpc error", "chain", ch.ChainSymbol, "error", err)
		return
	}
	for _, rpc := range chainRPCs {
		rpcUrls = append(rpcUrls, rpc.RPCURL)
	}
	client, err = chain.NewChain(ctx, chain.Config{
		Name:        ch.Name,
		ChainType:   ch.ChainType,
		ChainID:     ch.ChainID,
		ChainSymbol: ch.ChainSymbol,
		Currency:    ch.Currency,
		RpcURLs:     rpcUrls,
		GasPrice:    ch.GasPrice,
	})
	if err != nil {
		zap.S().Errorw("new chain client error", "chain", ch.ChainSymbol, "error", err)
		return
	}
	return
}

// ReadNextBlock 读取后续区块
//
//	@Description: 读取后续区块
//	@param ctx
//	@param ch
func ReadNextBlock(ctx context.Context, ch *sqlmodel.Chain) {
	var (
		chainQ    = sqlmodel.ChainColumns
		blockQ    = sqlmodel.ChainBlockColumns
		unChecked int64
		err       error
	)
	// 读取未检查的区块数量
	unChecked, err = dao.CountChainBlock(ctx, dao.And(blockQ.ChainSymbol.Eq(ch.ChainSymbol), blockQ.Checked.Eq(0)))
	if err != nil {
		zap.S().Errorw("count chain block error", "chain", ch.ChainSymbol, "error", err)
		return
	}
	if unChecked > int64(ch.Concurrent) {
		// 未检查的区块数量大于并发数 暂停读取
		return
	}
	if err != nil && !errors.Is(err, dao.ErrNotFound) {
		zap.S().Errorw("fetch chain block error", "chain", ch.ChainSymbol, "error", err)
		return
	}
	// TODO: 获取链的rpc client
	chainClient, err := GetChainRpcClient(ctx, ch)
	if err != nil {
		zap.S().Errorw("new chain client error", "chain", ch.ChainSymbol, "error", err)
		return
	}
	// TODO: 读取下一批区块
	var (
		chainBlocks         = make([]sqlmodel.ChainBlock, 0)
		lastBlockNum        int64
		latestChainBlockNum int64
	)
	latestChainBlockNum, err = chainClient.GetLatestBlockNumber(ctx)
	if err != nil {
		zap.S().Errorw("get latest block number error", "chain", ch.ChainSymbol, "error", err)
		return
	}
	if ch.LatestBlock == 0 {
		lastBlockNum = latestChainBlockNum - int64(ch.Concurrent)
	} else {
		lastBlockNum = ch.LatestBlock
	}
	if latestChainBlockNum == lastBlockNum {
		// 没有新的区
		return
	}
	// 并发读取区块
	for i := 0; i < int(latestChainBlockNum-lastBlockNum); i++ {
		lastBlockNum++
		chainBlocks = append(chainBlocks, sqlmodel.ChainBlock{
			ChainSymbol: ch.ChainSymbol,
			BlockNumber: lastBlockNum,
		})
	}
	if len(chainBlocks) == 0 { // 没有新的区块
		return
	}
	// 存储区块
	err = dao.GetDB(ctx).Transaction(func(tx *gorm.DB) (txErr error) {
		c := dao.CtxWithTransaction(ctx, tx)
		_, txErr = dao.AddsChainBlock(c, &chainBlocks)
		if txErr != nil {
			return
		}
		// 更新链的最新区块
		_, txErr = dao.UpdatesChain(c,
			dao.And(chainQ.ChainSymbol.Eq(ch.ChainSymbol)),
			dao.M{chainQ.LatestBlock.Name: lastBlockNum})
		return
	})
	return
}

// CheckRebase
//
//	@Description: 检查rebase
//	@param ctx
//	@param ch
func CheckRebase(ctx context.Context, ch *sqlmodel.Chain) {
	var (
		chainBlocks []sqlmodel.ChainBlock
		blockQ      = sqlmodel.ChainBlockColumns
		chainQ      = sqlmodel.ChainColumns
		err         error
	)
	// 获取需要检查的区块
	err = dao.FetchAllChainBlock(ctx, &chainBlocks,
		dao.And(
			blockQ.Checked.Eq(1),
			blockQ.ChainSymbol.Eq(ch.ChainSymbol),
			blockQ.BlockNumber.Gte(ch.RebaseBlock)),
		0, int(ch.Concurrent)*2, blockQ.BlockNumber.Asc())
	if err != nil || len(chainBlocks) == 0 {
		return
	}
	for i := 1; i < len(chainBlocks); i++ {
		ch.RebaseBlock = chainBlocks[i-1].BlockNumber
		block := chainBlocks[i]
		lastBlock := chainBlocks[i-1]
		if block.Checked == 0 || lastBlock.Checked == 0 {
			// 存在未检查的区块
			return
		}
		if block.ParentHash != lastBlock.BlockHash {
			ch.HasBranch = 1
			break
		}
	}
	_, err = dao.UpdatesChain(ctx, dao.And(chainQ.ChainSymbol.Eq(ch.ChainSymbol)),
		dao.M{chainQ.HasBranch.Name: ch.HasBranch, chainQ.RebaseBlock.Name: ch.RebaseBlock})
	if err != nil {
		zap.S().Errorw("update chain error", "chain", ch.ChainSymbol, "error", err)
		return
	}
	_ = UpdateTransactionsConfirm(ctx, ch)
}

// RebaseBlock
//
//	@Description: 重新设置区块
//	@param ctx
//	@param ch
func RebaseBlock(ctx context.Context, ch *sqlmodel.Chain) {
	var (
		rebaseChainBlock sqlmodel.ChainBlock
		nextChainBlock   sqlmodel.ChainBlock
		chainQ           = sqlmodel.ChainColumns
		blockQ           = sqlmodel.ChainBlockColumns
	)
	// 获取rebase的区块
	err := dao.FetchChainBlock(ctx, &rebaseChainBlock,
		dao.And(
			blockQ.ChainSymbol.Eq(ch.ChainSymbol),
			blockQ.BlockNumber.Eq(ch.RebaseBlock),
		))
	if err != nil && !errors.Is(err, dao.ErrNotFound) {
		return
	}
	// 获取下一个区块
	err = dao.FetchChainBlock(ctx, &nextChainBlock,
		dao.And(
			blockQ.ChainSymbol.Eq(ch.ChainSymbol),
			blockQ.BlockNumber.Eq(ch.RebaseBlock+1),
		))
	if err != nil && !errors.Is(err, dao.ErrNotFound) { // 发生查询错误
		return
	}
	if rebaseChainBlock.ID == 0 || rebaseChainBlock.BlockHash == nextChainBlock.ParentHash { // 没有发生rebase
		_, err = dao.UpdatesChain(ctx,
			dao.And(chainQ.ChainSymbol.Eq(ch.ChainSymbol)),
			dao.M{
				chainQ.HasBranch.Name: 0,
			})
		if err != nil {
			zap.S().Errorw("save chain error", "chain", ch.ChainSymbol, "error", err)
			return
		}
	} else {
		err = dao.GetDB(ctx).Transaction(func(tx *gorm.DB) (txErr error) {
			c := dao.CtxWithTransaction(ctx, tx)
			// 保存区块
			txErr = CheckChainBlock(c, ch, &rebaseChainBlock)
			if txErr != nil {
				return
			}
			// 更新链的最新区块
			_, err = dao.UpdatesChain(c,
				dao.And(chainQ.ChainSymbol.Eq(ch.ChainSymbol)),
				dao.M{
					chainQ.RebaseBlock.Name: rebaseChainBlock.BlockNumber - 1,
				})
			if txErr != nil {
				zap.S().Errorw("update chain error", "chain", ch.ChainSymbol, "error", txErr)
				return
			}
			return
		})
	}
}

// CheckBlocks
//
//	@Description: 检查区块
//	@param ctx
//	@param ch
func CheckBlocks(ctx context.Context, ch *sqlmodel.Chain) {
	var (
		blocks []sqlmodel.ChainBlock
		blockQ = sqlmodel.ChainBlockColumns
		err    error
		wg     sync.WaitGroup
	)
	// 获取需要检查的区块
	err = dao.FetchAllChainBlock(ctx, &blocks,
		dao.And(
			blockQ.ChainSymbol.Eq(ch.ChainSymbol),
			blockQ.Checked.Eq(0),
		),
		0, int(ch.Concurrent), blockQ.ID.Asc())
	if err != nil {
		zap.S().Errorw("fetch all chain block error", "chain", ch.ChainSymbol, "error", err)
		return
	}
	// 检查区块
	for _, block := range blocks {
		wg.Add(1)
		go func(b sqlmodel.ChainBlock) {
			defer wg.Done()
			_ = CheckChainBlock(ctx, ch, &b)
		}(block)
	}
	wg.Wait()
}

// CheckChainBlock
//
//	@Description: 检查链的区块
//	@param ctx
//	@param ch
//	@param block
//	@return err
func CheckChainBlock(ctx context.Context, ch *sqlmodel.Chain, block *sqlmodel.ChainBlock) (err error) {
	// 获取链的rpc client
	chainClient, err := GetChainRpcClient(ctx, ch)
	if err != nil {
		zap.S().Errorw("new chain client error", "chain", ch.ChainSymbol, "error", err)
		return
	}
	chainBlock, err := chainClient.GetBlock(ctx, block.BlockNumber)
	if err != nil {
		return
	}
	err = dao.GetDB(ctx).Transaction(func(tx *gorm.DB) (txErr error) {
		c := dao.CtxWithTransaction(ctx, tx)
		txErr = CheckTransactions(c, ch, chainBlock.Transactions)
		if txErr != nil {
			return
		}
		block.BlockHash = chainBlock.Hash
		block.ParentHash = chainBlock.ParentHash
		block.Checked = 1
		_, txErr = dao.UpdateChainBlock(c, block)
		return
	})
	return
}
