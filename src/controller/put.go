package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"iHome/src/model"
	"iHome/src/utils"
	"log"
	"net/http"
)

// PutUserName 更新用户名
func PutUserName(ctx *gin.Context) {
	// 从 session 获得当前用户的用户名
	s := sessions.Default(ctx)
	userNameInSession := s.Get("userName")

	// 中间件保证 session 合法
	//if userNameInSession == nil {
	//	// session 中没有用户数据
	//	ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SESSIONERR, "errmsg": utils.RecodeText(utils.RECODE_SESSIONERR)})
	//	return
	//}

	// 从 payload 拿新用户名的数据
	var userNameData struct {
		Name string `json:"name,omitempty"`
	}
	err := ctx.Bind(&userNameData)
	if err != nil {
		log.Println("PayLoad 读取错误：", err)
		ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_PARAMERR, "errmsg": utils.RecodeText(utils.RECODE_PARAMERR)})
		return
	}

	// 更新数据库
	err = model.UpdateUserName(userNameData.Name, userNameInSession.(string))
	if err != nil {
		// 数据库更新失败
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_DBERR, "errmsg": utils.RecodeText(utils.RECODE_DBERR)})
		return
	}

	// 更新 session
	s.Set("userName", userNameData.Name)
	err = s.Save() // 不要忘了 save
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SERVERERR, "errmsg": utils.RecodeText(utils.RECODE_SERVERERR)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
		"data":   userNameData,
	})
}
