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
	client, err = chain.NewChain(chain.Config{
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
		lastChainBlock sqlmodel.ChainBlock
		blockQ         = sqlmodel.ChainBlockColumns
		err            error
	)
	// 获取最后一个区块的数据
	err = dao.FetchChainBlock(ctx, &lastChainBlock,
		dao.And(
			blockQ.ChainSymbol.Eq(ch.ChainSymbol),
			blockQ.BlockNumber.Eq(ch.LatestBlock),
			blockQ.Removed.Eq(0),
		),
		blockQ.ID.Desc())
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
		wg           sync.WaitGroup
		chainBlocks  = make([]sqlmodel.ChainBlock, 0)
		lastBlockNum int64
		newBlocks    = make([]*chain.Block, 0)
		newBlocksMap = make(map[int64]*chain.Block)
		rebaseBlock  int64
		rpcErr       error
	)
	if ch.LatestBlock == 0 {
		lastBlockNum, err = chainClient.GetLatestBlockNumber(ctx)
	} else {
		lastBlockNum = ch.LatestBlock + 1
	}
	// 并发读取区块
	for i := 0; i < int(ch.Concurrent); i++ {
		wg.Add(1)
		go func(blockNumber int64) {
			defer wg.Done()
			newBlock, _err := chainClient.GetBlock(ctx, blockNumber)
			if _err != nil {
				if errors.Is(_err, chain.ErrorNotFound) {
					return
				}
				rpcErr = _err
				zap.S().Errorw("get block error", "chain", ch.ChainSymbol, "error", _err)
				return
			}
			if newBlock != nil {
				newBlocks = append(newBlocks, newBlock)
			}
		}(lastBlockNum + int64(i))
	}
	wg.Wait()
	// 判断是否有新区块
	if len(newBlocks) == 0 {
		zap.S().Infow("no new block", "chain", ch.ChainSymbol)
		return
	}
	if rpcErr != nil {
		return
	}
	for _, newBlock := range newBlocks {
		newBlocksMap[newBlock.Number] = newBlock
	}
	// 保存区块
	for i := 0; i < len(newBlocks); i++ {
		zap.S().Infow("read block", "chain", ch.ChainSymbol, "block", lastBlockNum+int64(i))
		newBlock := newBlocksMap[lastBlockNum+int64(i)]
		if newBlock == nil {
			zap.S().Errorw("new block is nil", "chain", ch.ChainSymbol, "block", lastBlockNum+int64(i))
			break
		}
		chainBlocks = append(chainBlocks, sqlmodel.ChainBlock{
			BlockHash:   newBlock.Hash,
			BlockNumber: newBlock.Number,
			ChainSymbol: ch.ChainSymbol,
			ParentHash:  newBlock.ParentHash,
		})
		if i == 0 { // 存在 uncle block
			if lastChainBlock.ID != 0 && newBlock.ParentHash != lastChainBlock.BlockHash {
				rebaseBlock = lastChainBlock.BlockNumber
				break
			}
		} else {
			lastBlock := newBlocksMap[lastBlockNum+int64(i-1)]
			if newBlock.ParentHash != lastBlock.Hash { // 存在 uncle block
				rebaseBlock = lastBlock.Number
				break
			}
		}
	}
	if len(chainBlocks) == 0 {
		return
	}
	err = dao.GetDB(ctx).Transaction(func(tx *gorm.DB) (txErr error) {
		c := dao.CtxWithTransaction(ctx, tx)
		// 保存区块
		_, txErr = dao.AddsChainBlock(c, &chainBlocks)
		if txErr != nil {
			// 打印错误的sql语句
			zap.S().Errorw("ReadNextBlock add chain block error", "block", chainBlocks, "error", txErr)
			return
		}
		// 更新链的最新区块
		ch.LatestBlock = chainBlocks[len(chainBlocks)-1].BlockNumber
		ch.RebaseBlock = rebaseBlock
		_, txErr = dao.UpdateChain(c, ch)
		if txErr != nil {
			zap.S().Errorw("save chain error", "chain", ch.ChainSymbol, "error", txErr)
			return
		}
		return
	})
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
		chainClient      chain.BaseChain
		blockQ           = sqlmodel.ChainBlockColumns
	)
	// 获取rebase的区块
	err := dao.FetchChainBlock(ctx, &rebaseChainBlock,
		dao.And(
			blockQ.ChainSymbol.Eq(ch.ChainSymbol),
			blockQ.BlockNumber.Eq(ch.RebaseBlock),
			blockQ.Removed.Eq(0),
		),
		blockQ.ID.Desc())
	if err != nil {
		return
	}
	// 获取上一个区块
	err = dao.FetchChainBlock(ctx, &nextChainBlock,
		dao.And(
			blockQ.ChainSymbol.Eq(ch.ChainSymbol),
			blockQ.BlockNumber.Eq(ch.RebaseBlock+1),
			blockQ.Removed.Eq(0),
		))
	if err != nil && !errors.Is(err, dao.ErrNotFound) { // 发生查询错误
		return
	}
	if nextChainBlock.BlockHash == "" || nextChainBlock.ParentHash == rebaseChainBlock.BlockHash { // 没有发生rebase
		ch.RebaseBlock = 0 // 解决了
		_, err = dao.UpdateChain(ctx, ch)
		if err != nil {
			zap.S().Errorw("save chain error", "chain", ch.ChainSymbol, "error", err)
			return
		}
	} else {
		// 获取链的rpc
		chainClient, err = GetChainRpcClient(ctx, ch)
		if err != nil {
			zap.S().Errorw("new chain client error", "chain", ch.ChainSymbol, "error", err)
			return
		}
		var newBlock *chain.Block
		newBlock, err = chainClient.GetBlock(ctx, ch.RebaseBlock)
		if err != nil {
			if errors.Is(err, chain.ErrorNotFound) {
				return
			}
			zap.S().Errorw("get block error", "chain", ch.ChainSymbol, "error", err)
			return
		}
		err = dao.GetDB(ctx).Transaction(func(tx *gorm.DB) (txErr error) {
			c := dao.CtxWithTransaction(ctx, tx)
			// 保存区块
			_, txErr = dao.AddChainBlock(c, &sqlmodel.ChainBlock{
				BlockHash:   newBlock.Hash,
				BlockNumber: newBlock.Number,
				ChainSymbol: ch.ChainSymbol,
				ParentHash:  newBlock.ParentHash,
				Removed:     0,
			})
			if txErr != nil {
				zap.S().Errorw("add chain block error", "chain", ch.ChainSymbol, "error", txErr)
				return
			}
			// 移除旧的区块
			rebaseChainBlock.Removed = 1
			_, txErr = dao.UpdateChainBlock(c, &rebaseChainBlock)
			if txErr != nil {
				zap.S().Errorw("remove chain block error", "chain", ch.ChainSymbol, "error", txErr)
				return
			}
			// 更新链的最新区块
			ch.RebaseBlock = rebaseChainBlock.BlockNumber - 1 // 重新设置
			_, txErr = dao.UpdateChain(c, ch)
			if txErr != nil {
				zap.S().Errorw("update chain error", "chain", ch.ChainSymbol, "error", txErr)
				return
			}
			return
		})
	}
}

// CheckBlock
//
//	@Description: 检查区块
//	@param ctx
//	@param ch
func CheckBlock(ctx context.Context, ch *sqlmodel.Chain) {
	var (
		chainBlocks []sqlmodel.ChainBlock
		blockQ      = sqlmodel.ChainBlockColumns
		err         error
		wg          sync.WaitGroup
	)
	// 获取需要检查的区块
	err = dao.FetchAllChainBlock(ctx, &chainBlocks,
		dao.And(
			blockQ.ChainSymbol.Eq(ch.ChainSymbol),
			blockQ.Checked.Eq(0),
			blockQ.Removed.Eq(0),
		),
		0, int(ch.Concurrent), blockQ.ID.Asc())
	if err != nil {
		zap.S().Errorw("fetch all chain block error", "chain", ch.ChainSymbol, "error", err)
		return
	}
	// 获取链的rpc client
	chainClient, err := GetChainRpcClient(ctx, ch)
	if err != nil {
		zap.S().Errorw("new chain client error", "chain", ch.ChainSymbol, "error", err)
		return
	}
	// 检查区块
	for _, block := range chainBlocks {
		wg.Add(1)
		go func(block sqlmodel.ChainBlock) {
			zap.S().Infow("check block", "chain", ch.ChainSymbol, "block", block.BlockNumber)
			// TODO: 检查区块
			transactions, _err := chainClient.GetBlockTransactions(ctx, block.BlockNumber)
			if _err != nil {
				return
			}
			_ = dao.GetDB(ctx).Transaction(func(tx *gorm.DB) (txErr error) {
				c := dao.CtxWithTransaction(ctx, tx)
				block.Checked = 1
				txErr = CheckTransactions(c, ch, transactions)
				if txErr != nil {
					return
				}
				_, txErr = dao.UpdateChainBlock(c, &block)
				return
			})
		}(block)
	}
}
