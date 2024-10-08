//generated by lazy
//author: seanlan

package model

type CreatePaymentAddressReqData struct {
	Chain     string `json:"chain,required"`
	NotifyUrl string `json:"notify_url,required"`
}

type CreatePaymentAddressReq struct {
	BaseReq
	AppKey string `json:"app_key" form:"app_key" binding:"required"`
	Data   string `json:"data" form:"data" binding:"required"`
	Sign   string `json:"sign" form:"sign" binding:"required"`
}

type CreatePaymentAddressResp struct {
	Address string `json:"address"`
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

type WithdrawReq struct {
	BaseReq
	AppKey string `json:"app_key" form:"app_key" binding:"required"`
	Data   string `json:"data" form:"data" binding:"required"`
	Sign   string `json:"sign" form:"sign" binding:"required"`
}

type WithdrawResp struct {
	Exist bool `json:"exist"`
}
