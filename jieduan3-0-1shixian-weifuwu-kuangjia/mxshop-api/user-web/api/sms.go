package api

import (
	"context"
	"fmt"
	"math/rand"
	"mxshop-api/user-web/forms"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"

	// redis 链接库
	"github.com/go-redis/redis/v8"

	"mxshop-api/user-web/global"
)

// 定义函数：生成短信验证码，参数是长度，返回字符串
func GenerateSmsCode(witdh int) string {
	//生成width长度的短信验证码

	// 1. 定义一个固定长度为 10 的 byte 数组，存放数字 0-9
	// byte类型的数组初始化，直接用数值赋值
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	// 2. 获取数组的长度（结果是 10），赋值给变量 r
	r := len(numeric)

	// 3. 随机数种子：用当前纳秒级时间戳做种子，保证每次生成的随机数不一样
	// 如果不写这句，每次生成的验证码都会是相同的
	// 给随机数生成器设置种子
	rand.Seed(time.Now().UnixNano())
	// 4. 声明一个 strings.Builder，用来高效拼接字符串
	var sb strings.Builder
	// 5. 循环：生成 width 位的验证码
	for i := 0; i < witdh; i++ {
		// 随机取 0~9 之间一个数字
		// 然后写入 sb 里（拼接字符串）
		// rand.Intn( 生成随机整数
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	// 6. 把拼接好的字符串返回出去
	return sb.String()
}

// 发送短信
func SendSms(ctx *gin.Context) {
	// 短信验证吗校验
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.ServerConfig.AliSmsInfo.ApiKey, global.ServerConfig.AliSmsInfo.ApiSecrect)
	if err != nil {
		panic(err)
	}
	smsCode := GenerateSmsCode(6)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = sendSmsForm.Mobile            //手机号
	request.QueryParams["SignName"] = "慕学在线"                            //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "SMS_181850725"               //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。
	// 调用 发送短信验证码
	response, err := client.ProcessCommonRequest(request)
	fmt.Print(client.DoAction(request, response))
	if err != nil {
		fmt.Print(err.Error())
	}
	// 短信验证码发送完后，用户注册提交时，还需要把验证吗提交进来，进行验证，所以需要把验证码保存起来，等用户注册提交时进行对比验证。
	// 保存的时候：一定要把验证吗和手机号进行绑定起来，一般手机号为key，验证码为val
	//将验证码保存起来 - redis，redis服务我们用docker启动redis服务
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	// 用代码连接 Redis，把 “手机号” 作为 key，“验证码” 作为 value 存进去，并且设置过期时间。
	// 短信肯定有过期时间的
	rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
