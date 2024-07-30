package tokenpay

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

type Client struct {
	AppName string
	AppKey  string
	Gateway string
}

func MapToUrlencoded(m map[string]interface{}, secretKey string) string {
	var keys []string
	var _source []string
	for k := range m {
		keys = append(keys, k)
	}
	//字符串排序
	sort.Strings(keys)
	for _, k := range keys {
		_source = append(_source, fmt.Sprintf("%s=%v", k, m[k]))
	}
	//map URL加入密钥拼接
	_source = append(_source, fmt.Sprintf("%s=%s", "key", secretKey))
	sourceStr := strings.Join(_source, "&")
	//MD5加密
	h := md5.New()
	h.Write([]byte(sourceStr))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func NewClient(appName, appKey, gate string) *Client {
	return &Client{
		AppName: appName,
		AppKey:  appKey,
		Gateway: appKey,
	}
}

// GetSign API参数签名
func (c *Client) GetSign(jsonObject map[string]interface{}) string {
	return MapToUrlencoded(jsonObject, c.AppKey)
}
