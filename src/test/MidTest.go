package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func midTest(ctx *gin.Context) {
	log.Println("MT-1")
	// return     // 中止执行中间件，将直接执行下一个中间件
	ctx.Next() // 直接执行下一个中间件，并将本中间件在 Next 以下的代码跳过并入栈，最后再以出栈顺序执行
	log.Println("MT-2")
}

func anotherMidTest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("AMT-1")
		// ctx.Next()
		ctx.Abort() // 只执行当前中间件，操作完成后以出栈顺序直接执行所有栈内中间件代码（包括调用 Abort 的这个中间件），不再继续调用其它中间件
		// 如果在此处使用 Abort ，将不会输出 Router
		log.Println("AMT-2")
	}
}

func main() {
	r := gin.Default()
	r.Use(midTest)
	r.Use(anotherMidTest())
	r.GET("/test", func(ctx *gin.Context) {
		// 路由 controller 实际上也实现了 gin.HandlerFunc ，也可以看作是一个中间件
		log.Println("Router")
		_, err := ctx.Writer.WriteString("hello world")
		if err != nil {
			log.Println(err)
			return
		}
	})
	err := r.Run(":9999")
	if err != nil {
		log.Println(err)
		return
	}
}
