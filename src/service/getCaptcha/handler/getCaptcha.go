package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"getCaptcha/model"
	pb "getCaptcha/proto"
	"github.com/afocus/captcha"
	"go-micro.dev/v4/logger"
	"image/color"
)

type GetCaptcha struct{}

func (e *GetCaptcha) Call(ctx context.Context, req *pb.CallRequest, rsp *pb.CallResponse) error {
	logger.Infof("Received GetCaptcha.Call request: %v", req)

	ctc := captcha.New()
	// 设置字体
	err := ctc.SetFont("./conf/comic.ttf")
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 设置验证码大小
	ctc.SetSize(128, 64)
	// 设置干扰强度，越高验证码越难以辨认
	ctc.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	ctc.SetFrontColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	ctc.SetBkgColor(color.RGBA{R: 255, A: 255}, color.RGBA{B: 255, A: 255}, color.RGBA{G: 153, A: 255})
	img, str := ctc.Create(6, captcha.NUM)
	// 将图片验证码和 uuid 存到 redis 中
	err = model.SaveImgCode(str, req.GetUuid())
	if err != nil {
		return err
	}

	imgBuf, err := json.Marshal(img)
	if err != nil {
		return err
	}
	rsp.Img = imgBuf
	return nil
}
