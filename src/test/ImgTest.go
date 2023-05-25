package main

import (
	"fmt"
	"github.com/afocus/captcha"
	"image/color"
	"image/png"
	"net/http"
)

/*
	参考 https://github.com/afocus/captcha/blob/master/examples/main.go
*/

func main() {
	ctc := captcha.New()
	// 设置字体
	err := ctc.SetFont("/comic.ttf")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 设置验证码大小
	ctc.SetSize(128, 64)
	// 设置干扰强度，越高验证码越难以辨认
	ctc.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	ctc.SetFrontColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	ctc.SetBkgColor(color.RGBA{R: 255, A: 255}, color.RGBA{B: 255, A: 255}, color.RGBA{G: 153, A: 255})

	// 根据上述配置生成
	http.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
		img, str := ctc.Create(6, captcha.ALL)
		err := png.Encode(w, img)
		if err != nil {
			return
		}
		fmt.Println(str)
	})

	// 根据 http 请求的负载生成
	http.HandleFunc("/c", func(w http.ResponseWriter, r *http.Request) {
		str := r.URL.RawQuery
		img := ctc.CreateCustom(str)
		err := png.Encode(w, img)
		if err != nil {
			return
		}
		fmt.Println(str)
	})

	err = http.ListenAndServe(":8085", nil)
	if err != nil {
		return
	}
}
