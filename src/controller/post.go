package controller

import (
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"iHome/src/model"
	"iHome/src/utils"
	"log"
	"net/http"
	userProto "user/proto"
)

// PostRegister 用户注册
func PostRegister(ctx *gin.Context) {
	// 以下只能处理 form 数据，无法处理 payload
	//mobile := ctx.PostForm("mobile")
	//pwd := ctx.PostForm("password")
	//smsCode := ctx.PostForm("sms_code")
	//fmt.Println(mobile, pwd, smsCode)

	var regData struct {
		Mobile   string `json:"mobile,omitempty"`
		PassWord string `json:"password,omitempty"`
		SmsCode  string `json:"sms_code,omitempty"`
	}
	err := ctx.Bind(&regData)
	if err != nil {
		log.Println("PayLoad 读取错误：", err)
		ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SERVERERR, "errmsg": utils.RecodeText(utils.RECODE_SERVERERR)})
		return
	}
	// fmt.Println(regData)
	microClient := userProto.NewUserService("user", utils.GetMicroClientFromConsul())
	response, err := microClient.Register(context.Background(), &userProto.RegisterRequest{
		Phone:    regData.Mobile,
		SmsCode:  regData.SmsCode,
		Password: regData.PassWord,
	})
	if err != nil {
		log.Println("远程注册服务错误：", err)
		ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SERVERERR, "errmsg": utils.RecodeText(utils.RECODE_SERVERERR)})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// PostLogin 用户登陆
func PostLogin(ctx *gin.Context) {
	var loginData struct {
		Mobile   string `json:"mobile,omitempty"`
		PassWord string `json:"password,omitempty"`
	}
	err := ctx.Bind(&loginData)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SERVERERR, "errmsg": utils.RecodeText(utils.RECODE_SERVERERR)})
		return
	}

	userName, err := model.LoginJudgement(loginData.Mobile, loginData.PassWord)
	if err != nil {
		if err == model.UserNotExistErr {
			ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_USERERR, "errmsg": utils.RecodeText(utils.RECODE_USERERR)})
		} else if err == model.UserPasswordErr {
			ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_PWDERR, "errmsg": utils.RecodeText(utils.RECODE_PWDERR)})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_LOGINERR, "errmsg": utils.RecodeText(utils.RECODE_LOGINERR)})
		}
	} else {
		// 登陆成功，将用户的电话号码存入 session 中
		s := sessions.Default(ctx)
		s.Set("userName", userName)
		err := s.Save() // 不要忘了 save
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SERVERERR, "errmsg": utils.RecodeText(utils.RECODE_SERVERERR)})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_OK, "errmsg": utils.RecodeText(utils.RECODE_OK)})
	}
}

// PostAvatar 上传头像
func PostAvatar(ctx *gin.Context) {
	// gin 框架内建方法，获取请求所发送的图片文件，在实际 Web 项目中基本使用云存储而非该方案
	//// 获取用户上传的图片文件
	//avatarFile, err := ctx.FormFile("avatar")
	//if err != nil {
	//	log.Println(err)
	//	ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SERVERERR, "errmsg": utils.RecodeText(utils.RECODE_SERVERERR)})
	//	return
	//}
	//
	//err = ctx.SaveUploadedFile(avatarFile, "test/"+avatarFile.Filename) // 目录默认路径是项目 main.go 文件所在路径
	//if err != nil {
	//	ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SERVERERR, "errmsg": utils.RecodeText(utils.RECODE_SERVERERR)})
	//	return
	//}
}
