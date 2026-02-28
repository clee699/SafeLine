package alert

import (
	"encoding/json"
	"fmt"

	"chaitin.cn/patronus/safeline-2/management/webserver/model"
)

// AlertService 告警服务接口
type AlertService interface {
	SendAlert(alert *AlertMessage)error
}

// AlertMessage 告警消息结构
type AlertMessage struct {
	EventID    string `json:"event_id"`
	SiteUUID   string `json:"site_uuid"`
	SrcIP      string `json:"src_ip"`
	Host       string `json:"host"`
	UrlPath    string `json:"url_path"`
	AttackType string `json:"attack_type"`
	RiskLevel  int    `json:"risk_level"`
	Action     int    `json:"action"`
	RuleID     string `json:"rule_id"`
	Timestamp  int64  `json:"timestamp"`
	Payload    string `json:"payload"`
}

// NewAlertService 根据告警类型创建对应的告警服务
func NewAlertService(config *model.AlertConfig) (AlertService, error) {
	var cfg map[string]interface{}
	if err := json.Unmarshal([]byte(config.Config), &cfg); err != nil {
		return nil, fmt.Errorf("解析告警配置失败: %v", err)
	}

	switch config.Type {
	case model.AlertTypeEmail:
		return NewEmailAlertService(cfg)
	case model.AlertTypeWeChat:
		return NewWeChatAlertService(cfg)
	default:
		return nil, fmt.Errorf("不支持的告警类型: %d", config.Type)
	}
}
