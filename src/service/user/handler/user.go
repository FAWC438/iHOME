package handler

import (
	"context"
	"go-micro.dev/v4/logger"
	"user/model"
	"user/utils"

	pb "user/proto"
)

type User struct{}

func (e *User) SendSms(ctx context.Context, req *pb.SendSmsRequest, rsp *pb.SendSmsResponse) error {
	logger.Infof("Received User.Call request: %v", req)

	// 检查图片验证码并判断短信发送是否成功
	checkOk := model.CheckImgCode(req.GetUuid(), req.GetImgCode())
	var errorCode, errMsg string
	if checkOk {
		// 图片验证码通过
		smsCode, err := utils.SendSMSAliyun(req.GetPhone()) // 发送验证码
		if err != nil || smsCode == "" {
			// 短信验证码发送失败
			errorCode = utils.RECODE_SMSERR
			errMsg = utils.RecodeText(utils.RECODE_SMSERR)
			if err != nil {
				logger.Info(err)
				// rsp.Error, rsp.Errmsg = errorCode, errMsg
				// return err
			}
		} else {
			// 短信验证码发送成功
			errorCode = utils.RECODE_OK
			errMsg = utils.RecodeText(utils.RECODE_OK)
			// redis 中存储手机号-手机验证码键值对
			err := model.SaveSmsCode(req.GetPhone(), smsCode)
			if err != nil {
				logger.Info(err)
				errorCode = utils.RECODE_DATAERR
				errMsg = utils.RecodeText(utils.RECODE_DATAERR)
				// rsp.Error, rsp.Errmsg = errorCode, errMsg
				// return err
			}
		}
	} else {
		// 图片验证码不通过
		// TODO: 图片验证码不通过时应当在 redis 立即销毁过期图片验证码而不是等待 keepalive 过期
		errorCode = utils.RECODE_DATAERR
		errMsg = utils.RecodeText(utils.RECODE_DATAERR)
	}
	rsp.Errno, rsp.Errmsg = errorCode, errMsg

	return nil
}

func (e *User) Register(ctx context.Context, req *pb.RegisterRequest, rsp *pb.RegisterResponse) error {
	// 先校验短信验证码是否正确，如果正确再将数据写入 MySQL
	err, ok := model.CheckSmsCode(req.Phone, req.SmsCode)
	if err != nil {
		// redis 查询错误
		logger.Info("redis 查询错误")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
	} else if !ok {
		// 无效短信验证码
		logger.Info("无效短信验证码")
		rsp.Errno = utils.RECODE_SMSERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_SMSERR)
	} else {
		exist, err := model.RegisterUserInMySQL(req.Phone, req.Password)
		if exist {
			// 用户已存在
			logger.Info("用户已存在")
			rsp.Errno = utils.RECODE_USERONERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_USERONERR)
		} else if err != nil {
			// 用户数据存入 MySQL 失败
			logger.Info("用户数据存入 MySQL 失败")
			rsp.Errno = utils.RECODE_DATAERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
		} else {
			// 成功注册
			rsp.Errno = utils.RECODE_OK
			rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
		}
	}

	return err
}
