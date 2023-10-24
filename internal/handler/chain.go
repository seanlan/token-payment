package handler

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"token-payment/internal/chain"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
)

func ReadNextBlock(ctx context.Context, ch *sqlmodel.Chain) {
	var (
		lastBlock sqlmodel.ChainBlock
		blockQ    = sqlmodel.ChainBlockColumns
		chainRPCs []sqlmodel.ChainRPC
		chainRpcQ = sqlmodel.ChainRPCColumns
		rpcUrls   []string
	)
	// 获取最后一个区块的数据
	err := dao.FetchChainBlock(ctx, &lastBlock,
		dao.And(
			blockQ.ChainSymbol.Eq(ch.ChainSymbol),
			blockQ.BlockNumber.Eq(ch.LatestBlock),
			blockQ.Removed.Eq(0),
		),
		blockQ.ID.Desc())
	if err != nil && !errors.Is(err, dao.ErrNotFound) {
		return
	}
	// TODO: 读取下一个区块
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
	chainClient, err := chain.NewChain(chain.Config{
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
	newBlock, err := chainClient.GetBlock(ctx, ch.LatestBlock+1)
	if err != nil {
		if errors.Is(err, chain.ErrorNotFound) {
			return
		}
		zap.S().Errorw("get block error", "chain", ch.ChainSymbol, "error", err)
		return
	}
	if lastBlock.ID != 0 && newBlock.ParentHash != lastBlock.BlockHash { // 有分叉
		ch.RebaseBlock = lastBlock.BlockNumber
		_, err = dao.UpdateChain(ctx, ch)
		if err != nil {
			zap.S().Errorw("save chain error", "chain", ch.ChainSymbol, "error", err)
			return
		}
	} else {
		err = dao.GetDB(ctx).Transaction(func(tx *gorm.DB) (txErr error) {
			c := dao.CtxWithTransaction(ctx, tx)
			// 保存区块
			_, txErr = dao.AddChainBlock(c, &sqlmodel.ChainBlock{
				BlockHash:   newBlock.Hash,
				BlockNumber: newBlock.Number,
				ChainSymbol: ch.ChainSymbol,
				ParentHash:  newBlock.ParentHash,
			})
			if txErr != nil {
				zap.S().Errorw("add chain block error", "chain", ch.ChainSymbol, "error", txErr)
				return
			}
			// 更新链的最新区块
			ch.LatestBlock = newBlock.Number
			_, txErr = dao.UpdateChain(c, ch)
			if txErr != nil {
				zap.S().Errorw("save chain error", "chain", ch.ChainSymbol, "error", txErr)
				return
			}
			return
		})
	}
}
