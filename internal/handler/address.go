package handler

import (
	"context"
	"time"
	"token-payment/internal/config"
	"token-payment/internal/dao"
	"token-payment/internal/dao/sqlmodel"
	"token-payment/internal/utils"
)

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
		_ = GenerateAddress(ctx, ch, int(ch.AddressPool))
	}
}

func GenerateAddress(ctx context.Context, ch *sqlmodel.Chain, count int) (err error) {
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
