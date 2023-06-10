package middlewares

type Middlewares struct {
	GetAuthenticationInfo *GetAuthenticationInfo
	RequireAuthentication *RequireAuthentication
	RateLimiter           RateLimiter
}
