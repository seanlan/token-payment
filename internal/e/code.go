package e

import "token-payment/pkg/xlerror"

var (
	ErrRequest = xlerror.ErrRequest
	ErrServer  = xlerror.ErrServer
	ErrForWait = xlerror.ErrForWait
	ErrorToken = xlerror.ErrToken

	ErrorUsernameOrPassword = xlerror.New(10001, "用户名或密码错误")

	ErrorPermission = xlerror.New(20001, "权限不足")
)
