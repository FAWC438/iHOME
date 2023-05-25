package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"iHome/src/utils"
	"log"
	"net/http"
)

// DeleteSession 删除（退出）登陆信息
func DeleteSession(ctx *gin.Context) {
	s := sessions.Default(ctx)
	s.Delete("userName")
	err := s.Save()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SERVERERR, "errmsg": utils.RecodeText(utils.RECODE_SERVERERR)})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_OK, "errmsg": utils.RecodeText(utils.RECODE_OK)})
}
