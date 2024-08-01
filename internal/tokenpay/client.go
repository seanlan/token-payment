package tokenpay

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"strings"
)

type Client struct {
	AppKey    string
	AppSecret string
	Gateway   string
}

func NewClient(appKey, appSecret, gate string) *Client {
	return &Client{
		AppKey:    appKey,
		AppSecret: appSecret,
		Gateway:   gate,
	}
}

// GetSign API参数签名
func (c *Client) GetSign(data string) string {
	sourceStr := data + c.AppSecret
	h := md5.New()
	h.Write([]byte(sourceStr))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func (c *Client) NotifyTransaction(tx NotifyTx, notifyUrl string) (success bool, err error) {
	request := gorequest.New()
	txBytes, _ := json.Marshal(tx)
	jsonObject := make(map[string]interface{})
	txData := string(txBytes)
	jsonObject["data"] = txData
	jsonObject["sign"] = c.GetSign(txData)
	_, body, errs := request.Post(notifyUrl).
		Type("json").
		Send(jsonObject).End()
	if len(errs) > 0 {
		err = errs[0]
	} else {
		err = nil
	}
	if strings.Contains(strings.ToLower(body), "success") { // 通知成功
		success = true
	} else {
		success = false
	}
	return
}

func (c *Client) request(method string, data interface{}, resp interface{}) (err error) {
	request := gorequest.New()
	reqBytes, _ := json.Marshal(data)
	jsonObject := make(map[string]interface{})
	reqData := string(reqBytes)
	jsonObject["app_key"] = c.AppKey
	jsonObject["data"] = reqData
	jsonObject["sign"] = c.GetSign(reqData)
	_, body, errs := request.Post(c.Gateway + method).
		Type("json").
		Send(jsonObject).End()
	if len(errs) > 0 {
		err = errs[0]
	} else {
		err = nil
	}
	err = json.Unmarshal([]byte(body), resp)
	return
}

func (c *Client) Address(data CreatePaymentAddressReqData) (resp ApiResponse[AddressRespData], err error) {
	err = c.request("/api/v1/payment/address", data, &resp)
	return
}

func (c *Client) Withdraw(data WithdrawReqData) (resp ApiResponse[AddressRespData], err error) {
	err = c.request("/api/v1/payment/withdraw", data, &resp)
	return
}
