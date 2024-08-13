package constant

const (
	UserAccessTokenExpires      = 3600 * 24 * 30  // 用户访问令牌过期时间
	UserRefreshTokenExpires     = 3600 * 24 * 365 // 用户刷新令牌过期时间
	UserAccessTokenRedisExpires = 3600
)
