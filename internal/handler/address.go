package handler

import (
	"context"
	"time"
	"token-payment/internal/config"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/utils"
)

// CheckAddressPool
//
//	@Description: 检查地址池
//	@param ctx
//	@param ch
func CheckAddressPool(ctx context.Context, ch *sqlmodel.Chain) {
	var (
		addressQ = sqlmodel.ChainAddressColumns
	)
	// 获取地址池
	addressCount, _ := dao.CountChainAddress(ctx, dao.And(
		addressQ.ChainSymbol.Eq(ch.ChainSymbol),
		addressQ.Used.Eq(0),
	))
	if addressCount < int64(ch.AddressPool) {
		// 生成地址
		_ = GenerateAddressBatch(ctx, ch, int(ch.AddressPool))
	}
}

// GenerateAddress
//
//	@Description: 生成地址
//	@param ctx
//	@param ch
//	@return address
//	@return err
func GenerateAddress(ctx context.Context, ch *sqlmodel.Chain) (address sqlmodel.ChainAddress, err error) {
	client, err := GetChainRpcClient(ctx, ch)
	if err != nil {
		return
	}
	var (
		privateKey, encKey string
	)
	// 生成地址
	address.Address, privateKey, err = client.GenerateAddress(ctx)
	if err != nil {
		return
	}
	encKey, err = utils.AesEncrypt(privateKey, config.C.Secret)
	address.ChainSymbol = ch.ChainSymbol
	address.EncKey = encKey
	address.Hook = ""
	address.CreateAt = time.Now().Unix()
	_, err = dao.AddChainAddress(ctx, &address)
	return
}

// GenerateAddressBatch
//
//	@Description: 批量生成地址
//	@param ctx
//	@param ch
//	@param count
//	@return err
func GenerateAddressBatch(ctx context.Context, ch *sqlmodel.Chain, count int) (err error) {
	client, err := GetChainRpcClient(ctx, ch)
	if err != nil {
		return
	}
	var addressList = make([]sqlmodel.ChainAddress, 0)
	for i := 0; i < count; i++ {
		var (
			address, privateKey, encKey string
			_err                        error
		)
		// 生成地址
		address, privateKey, _err = client.GenerateAddress(ctx)
		if _err != nil {
			return
		}
		encKey, _err = utils.AesEncrypt(privateKey, config.C.Secret)
		addressList = append(addressList, sqlmodel.ChainAddress{
			ChainSymbol: ch.ChainSymbol,
			Address:     address,
			EncKey:      encKey,
			Hook:        "",
			CreateAt:    time.Now().Unix(),
		})
	}
	_, err = dao.AddsChainAddress(ctx, &addressList)
	return
}
