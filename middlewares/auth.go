package middlewares

import (
	"bluebell/controller"
	"bluebell/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带token有三种方式 1. 放在请求头 2. 放在请求体 3.放在URL
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// Authorization: Bearer xxxxxxx.xxx.xxx  / X-TOKEN: xxx.xxx.xx
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}

		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		c.Set(controller.CtxUserIDKey, mc.UserID)
		c.Next()
	}
}
