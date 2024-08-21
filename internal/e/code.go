package e

import "token-payment/pkg/xlerror"

var (
	ErrRequest = xlerror.ErrRequest
	ErrServer  = xlerror.ErrServer
	ErrForWait = xlerror.ErrForWait
	ErrorToken = xlerror.ErrToken

	ErrorUsernameOrPassword   = xlerror.New(10001, "用户名或密码错误")
	ErrorGroupIDNotExist      = xlerror.New(10003, "群组不存在")
	ErrorPermissionIDNotExist = xlerror.New(10004, "权限不存在")
	ErrorAccountExist         = xlerror.New(10005, "账号已存在")
	ErrorSuperUserNotEdit     = xlerror.New(10006, "超级管理员不可被编辑")
	ErrorAccountFrozen        = xlerror.New(10007, "账号已被冻结")
	ErrorStaffHandover        = xlerror.New(10008, "交接失败")

	ErrorPermission = xlerror.New(20001, "权限不足")

	ErrorApplicationNotFound = xlerror.New(30001, "应用不存在")
	ErrorSign                = xlerror.New(30002, "签名错误")
	ErrorDataParam           = xlerror.New(30003, "参数错误")
	ErrorChainNotSupport     = xlerror.New(30004, "不支持的链类型")
	ErrorTokenNotSupport     = xlerror.New(30005, "不支持的币种")
)
