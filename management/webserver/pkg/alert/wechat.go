package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// WeChatConfig 微信告警配置
type WeChatConfig struct {
	WebhookURL string `json:"webhook_url"`
	Secret     string `json:"secret"`
}

// WeChatAlertService 微信告警服务
type WeChatAlertService struct {
	config *WeChatConfig
	client *http.Client
}

// NewWeChatAlertService 创建微信告警服务
func NewWeChatAlertService(cfg map[string]interface{}) (AlertService, error) {
	wechatCfg := &WeChatConfig{}

	if webhookURL, ok := cfg["webhook_url"].(string); ok {
		wechatCfg.WebhookURL = webhookURL
	} else {
		return nil, fmt.Errorf("缺少必填配置: webhook_url")
	}

	if secret, ok := cfg["secret"].(string); ok {
		wechatCfg.Secret = secret
	}

	return &WeChatAlertService{
		config: wechatCfg,
		client: &http.Client{},
	}, nil
}

// SendAlert 发送微信告警
func (s *WeChatAlertService) SendAlert(alert *AlertMessage) error {
	// 构建微信消息
	msg := s.buildWeChatMessage(alert)

	// 发送HTTP请求
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化微信消息失败: %v", err)
	}

	resp, err := s.client.Post(s.config.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送微信告警失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("发送微信告警失败，HTTP状态码: %d", resp.StatusCode)
	}

	return nil
}

// wechatMarkdownMessage 企业微信Markdown消息结构
type wechatMarkdownMessage struct {
	MsgType  string                  `json:"msgtype"`
	Markdown wechatMarkdownContent   `json:"markdown"`
}

// wechatMarkdownContent 企业微信Markdown内容
type wechatMarkdownContent struct {
	Content string `json:"content"`
}

// buildWeChatMessage 构建微信消息
func (s *WeChatAlertService) buildWeChatMessage(alert *AlertMessage) *wechatMarkdownMessage {
	content := fmt.Sprintf(`**SafeLine WAF 检测到攻击告警**

> **事件ID**: %s
> **站点UUID**: %s
> **攻击IP**: %s
> **目标主机**: %s
> **请求路径**: %s
> **攻击类型**: %s
> **风险级别**: %d
> **处理动作**: %d
> **触发规则**: %s
> **攻击时间**: %d
> **攻击载荷**: %s

此消息由 SafeLine WAF 自动发送`,
		alert.EventID,
		alert.SiteUUID,
		alert.SrcIP,
		alert.Host,
		alert.UrlPath,
		alert.AttackType,
		alert.RiskLevel,
		alert.Action,
		alert.RuleID,
		alert.Timestamp,
		alert.Payload,
	)

	return &wechatMarkdownMessage{
		MsgType: "markdown",
		Markdown: wechatMarkdownContent{
			Content: content,
		},
	}
}
