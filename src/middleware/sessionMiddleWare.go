package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"iHome/src/utils"
	"net/http"
)

// SessionAuthFilter 中间件，用于在 session 中判断用户是否登陆
func SessionAuthFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := sessions.Default(ctx)
		if s.Get("userName") == nil {
			// session 中没有用户数据
			ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SESSIONERR, "errmsg": utils.RecodeText(utils.RECODE_SESSIONERR)})
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}
