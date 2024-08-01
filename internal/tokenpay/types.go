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
	TokenID         int64   `json:"token_id"`
	Value           float64 `json:"value"`
	TxIndex         int64   `json:"tx_index"`
	BatchIndex      int64   `json:"batch_index"`
	Confirm         int32   `json:"confirm"`
	MaxConfirm      int32   `json:"max_confirm"`
	TransferType    int32   `json:"transfer_type"`
	SerialNo        string  `json:"serial_no"`
	CreateAt        int64   `json:"create_at"`
}

type CreatePaymentAddressReqData struct {
	Chain     string `json:"chain,required"`
	NotifyUrl string `json:"notify_url,required"`
}

type WithdrawReqData struct {
	Chain           string `json:"chain,required"`
	SerialNo        string `json:"serial_no,required"`
	Symbol          string `json:"symbol,required"`
	ContractAddress string `json:"contract_address,omitempty"`
	TokenID         int64  `json:"token_id,omitempty"`
	Value           string `json:"value,required"`
	ToAddress       string `json:"to_address,required"`
	NotifyUrl       string `json:"notify_url,required"`
}

type ApiResponse[T any] struct {
	Error   int    `json:"error,required"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,required"`
	Link    string `json:"link,omitempty"`
}

type AddressRespData struct {
	Address string `json:"address"`
}

type WithdrawRespData struct {
	TxHash string `json:"tx_hash"`
}
