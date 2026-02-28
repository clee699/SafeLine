package config

import (
	"chaitin.cn/dev/go/settings"
)

// AlertConfig 告警配置结构体
type AlertConfig struct {
	Enabled        bool   `yaml:"enabled"`        // 告警全局开关
	CheckInterval  int    `yaml:"check_interval"` // 告警检查间隔（秒）
	DefaultLevel   int    `yaml:"default_level"`  // 默认告警级别
	MaxAlertCount  int    `yaml:"max_alert_count"` // 最大告警数量
}

// DefaultAlertConfig 创建默认告警配置
func DefaultAlertConfig() AlertConfig {
	return AlertConfig{
		Enabled:        true,
		CheckInterval:  10,
		DefaultLevel:   2,
		MaxAlertCount:  1000,
	}
}

// Load 加载告警配置
func (ac *AlertConfig) Load(setting *settings.Setting) error {
	if err := setting.Unmarshal("alert", ac); err != nil {
		return err
	}

	return nil
}
