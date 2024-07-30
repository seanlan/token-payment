package tokenpay

type NotifyTx struct {
	ApplicationID   int64   `json:"application_id"`
	ChainSymbol     string  `json:"chain_symbol"`
	TxHash          string  `json:"tx_hash"`
	FromAddress     string  `json:"from_address"`
	ToAddress       string  `json:"to_address"`
	ContractAddress string  `json:"contract_address"`
	Symbol          string  `json:"symbol"`
	Decimals        int32   `json:"decimals"`
	Value           float64 `json:"value"`
	TxIndex         int64   `json:"tx_index"`
	BatchIndex      int64   `json:"batch_index"`
	Confirm         int32   `json:"confirm"`
	MaxConfirm      int32   `json:"max_confirm"`
	TransferType    int32   `json:"transfer_type"`
	SerialNo        string  `json:"serial_no"`
	CreateAt        int64   `json:"create_at"`
}
