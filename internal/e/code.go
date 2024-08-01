package e

import "token-payment/pkg/xlerror"

var (
	ErrRequest = xlerror.ErrRequest
	ErrServer  = xlerror.ErrServer
	ErrForWait = xlerror.ErrForWait
	ErrorToken = xlerror.ErrToken

	ErrorUsernameOrPassword = xlerror.New(10001, "用户名或密码错误")

	ErrorPermission = xlerror.New(20001, "权限不足")

	ErrorApplicationNotFound = xlerror.New(30001, "应用不存在")
	ErrorSign                = xlerror.New(30002, "签名错误")
	ErrorDataParam           = xlerror.New(30003, "参数错误")
	ErrorChainNotSupport     = xlerror.New(30004, "不支持的链类型")
	ErrorTokenNotSupport     = xlerror.New(30005, "不支持的币种")
)
