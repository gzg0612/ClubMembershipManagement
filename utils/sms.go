package utils

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// SMSConfig 短信配置
type SMSConfig struct {
	SecretID   string
	SecretKey  string
	AppID      string
	SignName   string
	TemplateID string
}

// SMSClient 短信客户端
type SMSClient struct {
	client *sms.Client
	config *SMSConfig
}

// NewSMSClient 创建短信客户端
func NewSMSClient(config *SMSConfig) (*SMSClient, error) {
	credential := common.NewCredential(
		config.SecretID,
		config.SecretKey,
	)

	prof := profile.NewClientProfile()
	client, err := sms.NewClient(credential, "ap-guangzhou", prof)
	if err != nil {
		return nil, fmt.Errorf("创建短信客户端失败: %v", err)
	}

	return &SMSClient{
		client: client,
		config: config,
	}, nil
}

// Send 发送短信
func (c *SMSClient) Send(phoneNumber string, templateParam []string) error {
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(c.config.AppID)
	request.SignName = common.StringPtr(c.config.SignName)
	request.TemplateId = common.StringPtr(c.config.TemplateID)
	request.TemplateParamSet = common.StringPtrs(templateParam)
	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + phoneNumber})

	_, err := c.client.SendSms(request)
	if err != nil {
		return fmt.Errorf("发送短信失败: %v", err)
	}

	return nil
}

// SendVerificationCode 发送验证码
func (c *SMSClient) SendVerificationCode(phoneNumber string) error {
	// 生成6位随机验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 发送验证码
	err := c.Send(phoneNumber, []string{code})
	if err != nil {
		return err
	}

	// 将验证码保存到Redis,设置5分钟过期
	ctx := context.Background()
	key := fmt.Sprintf("sms:code:%s", phoneNumber)
	err = GetRedis().Set(ctx, key, code, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("保存验证码失败: %v", err)
	}

	return nil
}

// VerifyCode 验证验证码
func (c *SMSClient) VerifyCode(phoneNumber, code string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("sms:code:%s", phoneNumber)

	// 从Redis获取验证码
	savedCode, err := GetRedis().Get(ctx, key).Result()
	if err != nil {
		return false
	}

	// 验证码匹配
	if savedCode == code {
		// 验证成功后删除验证码
		GetRedis().Del(ctx, key)
		return true
	}

	return false
}
