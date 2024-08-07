package types

type TransferType int

const (
	TransferTypeIn      TransferType = 1 // 收款
	TransferTypeOut     TransferType = 2 // 付款
	TransferTypeFee     TransferType = 3 // 手续费
	TransferTypeArrange TransferType = 4 // 零钱整理
)
