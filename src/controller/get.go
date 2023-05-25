package controller

import (
	"context"
	"encoding/json"
	"fmt"
	getCaptchaProto "getCaptcha/proto"
	"github.com/afocus/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"iHome/src/model"
	"iHome/src/utils"
	"image/png"
	"log"
	"net/http"
	"regexp"
	userProto "user/proto"
)

// GetSession 获取 session 信息
func GetSession(ctx *gin.Context) {
	s := sessions.Default(ctx)
	userName := s.Get("userName")
	// 中间件保证 session 合法
	//if userName == nil {
	//	// session 中不存在这个用户名
	//	ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SESSIONERR, "errmsg": utils.RecodeText(utils.RECODE_SESSIONERR)})
	//}

	ctx.JSON(http.StatusOK, gin.H{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
		"data":   gin.H{"name": userName},
	})
}

// GetImageCd 获取图片信息
func GetImageCd(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	// 得到 micro 客户端实例
	microClient := getCaptchaProto.NewGetCaptchaService("getCaptcha", utils.GetMicroClientFromConsul())
	// 调用 rpc 服务获得响应
	response, err := microClient.Call(context.Background(), &getCaptchaProto.CallRequest{Uuid: uuid})
	if err != nil {
		log.Println(err)
		return
	}
	// 将数据通过 json 反序列化
	var img captcha.Image
	err = json.Unmarshal(response.GetImg(), &img)
	if err != nil {
		log.Println(err)
		return
	}

	// 编码为 png 并写入 gin 上下文
	err = png.Encode(ctx.Writer, img)
	if err != nil {
		log.Println(err)
		return
	}
}

// GetSmsCd 获取短信验证码
func GetSmsCd(ctx *gin.Context) {
	// 获取手机号码
	mobile := ctx.Param("mobile")
	// 得到当前的图片验证码的值
	imgCode := ctx.Query("text")
	// 得到图片验证码的 uuid
	uuid := ctx.Query("id")

	fmt.Println("------------- SmsCd out: ", mobile, imgCode, uuid)

	// 判断手机号输入格式
	regRule := regexp.MustCompile("^1[345789]\\d{9}$")
	phoneNumCheck := regRule.MatchString(mobile)
	if !phoneNumCheck {
		// 手机号输入格式不正确
		ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_MOBILEERR, "errmsg": utils.RecodeText(utils.RECODE_MOBILEERR)})
		return
	}

	// 得到 micro 客户端实例
	microClient := userProto.NewUserService("user", utils.GetMicroClientFromConsul())
	// 调用 rpc 服务获得响应
	response, err := microClient.SendSms(context.Background(), &userProto.SendSmsRequest{Phone: mobile, ImgCode: imgCode, Uuid: uuid})
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SERVERERR, "errmsg": utils.RecodeText(utils.RECODE_SERVERERR)})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetAreas 获取地域信息
func GetAreas(ctx *gin.Context) {
	// 如果 redis 中有数据直接返回
	var areas []model.Area
	redisConn := model.RedisPool.Get()
	areasByteSlices, err := redis.Bytes(redisConn.Do("get", "areasData"))
	if err == nil {
		log.Println("redis 缓存命中")
		err := json.Unmarshal(areasByteSlices, &areas)
		if err != nil {
			log.Println("areasByteSlices decode error: ", err)
			ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SERVERERR, "errmsg": utils.RecodeText(utils.RECODE_SERVERERR)})
			return
		}
	} else if err == redis.ErrNil {
		log.Println("redis 缓存未命中")
		// 从 MySQL 中获取数据
		model.MySQLPool.Find(&areas)
		areasBuffer, _ := json.Marshal(areas) // json 序列化

		// 再将数据写入 redis
		// 直接存入对象数组，在后面获取的时候 redigo 没有 api 来获取它
		// _, err := redisConn.Do("set", "areasData", areas)
		// 需要存入对象数组的 json 序列化数据
		_, err = redisConn.Do("set", "areasData", areasBuffer)
		if err != nil {
			log.Println("redis set error: ", err)
			ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_DATAERR, "errmsg": utils.RecodeText(utils.RECODE_DATAERR)})
			return
		}
	} else {
		log.Println("redis get areas json byte error: ", err)
		ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_DATAERR, "errmsg": utils.RecodeText(utils.RECODE_DATAERR)})
		return
	}

	jsonData := gin.H{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
		"data":   areas,
	}
	ctx.JSON(http.StatusOK, jsonData)
}

// GetUserInfo 获取用户信息
func GetUserInfo(ctx *gin.Context) {
	s := sessions.Default(ctx)
	userName := s.Get("userName")
	// 中间件保证 session 合法
	//if userName == nil {
	//	// session 中没有用户数据
	//	ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_SESSIONERR, "errmsg": utils.RecodeText(utils.RECODE_SESSIONERR)})
	//	return
	//}

	// 通过用户名查 MySQL
	userInfo, err := model.GetUserInfo(userName.(string))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"errno": utils.RECODE_DBERR, "errmsg": utils.RecodeText(utils.RECODE_DBERR)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"errno":  utils.RECODE_OK,
		"errmsg": utils.RecodeText(utils.RECODE_OK),
		// 此处不直接传 User 对象的原因在于，难以将 User 对象的密码字段用 `json:"-"` 忽略，因为 User 的使用范围很广，这可能会造成难以预估的错误
		"data": gin.H{
			"user_id":    userInfo.ID,
			"name":       userInfo.Name,
			"mobile":     userInfo.Mobile,
			"real_name":  userInfo.Real_name,
			"id_card":    userInfo.Id_card,
			"avatar_url": userInfo.Avatar_url,
		},
	})
}
