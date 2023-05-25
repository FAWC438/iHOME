package utils

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// createClient 使用AK&SK初始化账号Client
func createClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

// SendSMSAliyun
//
//	@Description: 通过阿里云发送短信验证码
//	@param phone 目标手机号
//	@return smsCode 发送到目标手机号的手机验证码，一个四位整数
//	@return _err
func SendSMSAliyun(phone string) (smsCode string, _err error) {
	// accessKeyId and accessKeySecret
	client, _err := createClient(tea.String("****"), tea.String("****"))
	if _err != nil {
		return "", _err
	}

	// 1. 正常业务中，生成的随机验证码
	//randRes, err := rand.Int(rand.Reader, big.NewInt(9999))
	//if err != nil {
	//	return err
	//}
	//strCode := fmt.Sprintf("%04d", randRes)

	// 2. 阿里云个人用户限制，强制使用验证码
	strCode := "1234"

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:     tea.String("阿里云短信测试"),
		TemplateCode: tea.String("SMS_154950909"),
		// PhoneNumbers:  tea.String(phone),	// 阿里云个人用户限制，只能发给特定电话 18605958778
		PhoneNumbers:  tea.String("****"),
		TemplateParam: tea.String(`{"code":"` + strCode + `"}`),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		resp, _err := client.SendSmsWithOptions(sendSmsRequest, runtime)
		if _err != nil {
			return _err
		}
		fmt.Println("SMS response: ", resp)

		return nil
	}()

	if tryErr != nil {
		var e = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			e = _t
		} else {
			e.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 e
		_, _err = util.AssertAsString(e.Message)
		if _err != nil {
			return "", _err
		}
	}
	return strCode, _err
}
