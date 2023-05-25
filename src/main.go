package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	ctrl "iHome/src/controller"
	"iHome/src/middleware"
	"iHome/src/model"
	"log"
	"net/http"
)

// initSession 初始化 session 容器
func initSession() (sessions.Store, error) {
	// 最后的加密密钥是必须的
	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("ihome"))
	if err != nil {
		return nil, err
	}
	// 设置容器中 session 的过期时间
	store.Options(sessions.Options{
		// 12小时过期
		MaxAge: 12 * 60 * 60,
	})
	return store, nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 初始化 MySQL
	_, err := model.InitMySQL()
	if err != nil {
		log.Println(err)
		return
	}

	// 初始化 session 容器
	store, err := initSession()
	if err != nil {
		log.Println(err)
		return
	}

	r := gin.Default()

	// Recovery 中间件保证发生 panic 时返回 500 错误
	r.Use(gin.Recovery())
	// 使用 session 容器（中间件）
	r.Use(sessions.Sessions("ihomeSession", store))

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
	r.Static("/home", "view")
	// 路由分组
	apiFunc := r.Group("/api/v1.0/")
	{
		apiFunc.GET("/imagecode/:uuid", ctrl.GetImageCd)
		apiFunc.GET("/smscode/:mobile", ctrl.GetSmsCd)
		apiFunc.GET("/areas", ctrl.GetAreas)
		apiFunc.POST("/users", ctrl.PostRegister)
		apiFunc.POST("/sessions", ctrl.PostLogin)
		apiFunc.Use(middleware.SessionAuthFilter()) // 以下的路由都需要通过 session 用户登陆验证
		{
			apiFunc.GET("/session", ctrl.GetSession)
			apiFunc.DELETE("/session", ctrl.DeleteSession)
			userFunc := apiFunc.Group("/user")
			{
				userFunc.GET("/", ctrl.GetUserInfo)
				userFunc.POST("/avatar", ctrl.PostAvatar)
				userFunc.PUT("/name", ctrl.PutUserName)
			}
		}
	}

	err = r.Run(":8080")
	if err != nil {
		log.Println(err)
		return
	}
}
