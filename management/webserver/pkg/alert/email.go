package alert

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailConfig 邮件告警配置
type EmailConfig struct {
	SMTPServer string   `json:"smtp_server"`
	SMTPPort   int      `json:"smtp_port"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	From       string   `json:"from"`
	To         []string `json:"to"`
	Subject    string   `json:"subject"`
}

// EmailAlertService 邮件告警服务
type EmailAlertService struct {
	config *EmailConfig
}

// NewEmailAlertService 创建邮件告警服务
func NewEmailAlertService(cfg map[string]interface{}) (AlertService, error) {
	emailCfg := &EmailConfig{}

	if server, ok := cfg["smtp_server"].(string); ok {
		emailCfg.SMTPServer = server
	} else {
		return nil, fmt.Errorf("缺少必填配置: smtp_server")
	}

	if port, ok := cfg["smtp_port"].(float64); ok {
		emailCfg.SMTPPort = int(port)
	} else {
		emailCfg.SMTPPort = 25 // 默认端口
	}

	if username, ok := cfg["username"].(string); ok {
		emailCfg.Username = username
	} else {
		return nil, fmt.Errorf("缺少必填配置: username")
	}

	if password, ok := cfg["password"].(string); ok {
		emailCfg.Password = password
	} else {
		return nil, fmt.Errorf("缺少必填配置: password")
	}

	if from, ok := cfg["from"].(string); ok {
		emailCfg.From = from
	} else {
		emailCfg.From = emailCfg.Username // 默认使用用户名作为发件人
	}

	if to, ok := cfg["to"].([]interface{}); ok {
		for _, t := range to {
			if email, ok := t.(string); ok {
				emailCfg.To = append(emailCfg.To, email)
			}
		}
	} else {
		return nil, fmt.Errorf("缺少必填配置: to")
	}

	if subject, ok := cfg["subject"].(string); ok {
		emailCfg.Subject = subject
	} else {
		emailCfg.Subject = "SafeLine WAF 告警通知" // 默认主题
	}

	return &EmailAlertService{config: emailCfg}, nil
}

// SendAlert 发送邮件告警
func (s *EmailAlertService) SendAlert(alert *AlertMessage) error {
	// 构建邮件内容
	body := s.buildEmailBody(alert)

	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = s.config.From
	headers["To"] = strings.Join(s.config.To, ",")
	headers["Subject"] = s.config.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=utf-8"

	// 构建完整邮件
	message := s.buildMessage(headers, body)

	// 发送邮件
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.SMTPServer)
	addr := fmt.Sprintf("%s:%d", s.config.SMTPServer, s.config.SMTPPort)

	return smtp.SendMail(addr, auth, s.config.From, s.config.To, []byte(message))
}

// buildEmailBody 构建邮件正文
func (s *EmailAlertService) buildEmailBody(alert *AlertMessage) string {
	return fmt.Sprintf(`SafeLine WAF 检测到攻击告警

事件ID: %s
站点UUID: %s
攻击IP: %s
目标主机: %s
请求路径: %s
攻击类型: %s
风险级别: %d
处理动作: %d
触发规则: %s
攻击时间: %d
攻击载荷: %s

此邮件由 SafeLine WAF 自动发送`,
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
}

// buildMessage 构建完整邮件消息
func (s *EmailAlertService) buildMessage(headers map[string]string, body string) string {
	var message strings.Builder

	// 构建邮件头
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	// 添加空行分隔邮件头和正文
	message.WriteString("\r\n")
	// 添加邮件正文
	message.WriteString(body)

	return message.String()
}
