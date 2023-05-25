package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	r := gin.Default()

	// 初始化 session 容器，最后的加密密钥是必须的
	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("ihome"))
	if err != nil {
		log.Println(err)
		return
	}
	// 设置容器中 session 的过期时间
	store.Options(sessions.Options{
		MaxAge: 5 * 60,
	})
	// 使用该容器
	r.Use(sessions.Sessions("testSession", store))

	r.GET("/ping", func(c *gin.Context) {
		// 设置 session
		s := sessions.Default(c)
		s.Set("mySession", "test info")
		err := s.Save()
		if err != nil {
			log.Println(err)
			return
		}

		sessionText := s.Get("mySession").(string)
		if sessionText != "" {
			log.Printf(sessionText)
		}

		// 设置 cookie
		c.SetCookie("test", "this-is-a-test", 5*60, "", "", true, true)
		cookie, _ := c.Cookie("test")
		if cookie != "" {
			log.Println(cookie)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	err = r.Run(":9999")
	if err != nil {
		log.Println(err)
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
