package tokenpay

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"strings"
)

type Client struct {
	AppName string
	AppKey  string
	Gateway string
}

func NewClient(appName, appKey, gate string) *Client {
	return &Client{
		AppName: appName,
		AppKey:  appKey,
		Gateway: appKey,
	}
}

// GetSign API参数签名
func (c *Client) GetSign(data string) string {
	sourceStr := data + c.AppKey
	h := md5.New()
	h.Write([]byte(sourceStr))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func (c *Client) NotifyTransaction(tx NotifyTx, notifyUrl string) (success bool, err error) {
	request := gorequest.New()
	txBytes, _ := json.Marshal(tx)
	jsonObject := make(map[string]interface{})
	txData := string(txBytes)
	jsonObject["tx_data"] = txData
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
