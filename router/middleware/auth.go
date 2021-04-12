package middleware

import (
	"apiserver/handler"
	"apiserver/pkg/errno"
	"apiserver/pkg/token"

	"github.com/gin-gonic/gin"
)

//AuthMiddleware jwt 验证中间件
//通过 HTTP Header Authorization: Bearer $token 来携带 token。携带 token 后不需要再次查询数据库核对密码，即可完成授权
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the json web token.
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
