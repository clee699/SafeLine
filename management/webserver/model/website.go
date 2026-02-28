package model

import "gorm.io/datatypes"

type Website struct {
	Base
	Comment     string         `gorm:"comment"           json:"comment"`
	ServerNames datatypes.JSON `gorm:"server_names"      json:"server_names"`
	Ports       datatypes.JSON `gorm:"ports"             json:"ports"`
	Upstreams   datatypes.JSON `gorm:"upstreams"         json:"upstreams"`

	CertFilename string `gorm:"cert_filename"    json:"cert_filename"`
	KeyFilename  string `gorm:"key_filename"     json:"key_filename"`

	IsEnabled bool `gorm:"is_enabled;default=true"       json:"is_enabled"`
}

type UpstreamConfig struct {
	URL        string `json:"url"`
	Weight     int    `json:"weight,omitempty"`     // 权重，默认1
	MaxFails   int    `json:"max_fails,omitempty"`   // 最大失败次数，默认1
	FailTimeout int   `json:"fail_timeout,omitempty"` // 失败超时时间，默认10秒
}

type LoadBalancerConfig struct {
	Enabled   bool             `json:"enabled"`
	Algorithm string           `json:"algorithm"` // 负载均衡算法：round_robin, least_conn, ip_hash
	Upstreams []UpstreamConfig `json:"upstreams"`
}

type ForbiddenPage struct {
	Base
	Content     string `gorm:"content"     json:"content"`
	IsEnabled   bool   `gorm:"is_enabled;default=false" json:"is_enabled"`
	WebsiteID   uint   `gorm:"website_id"  json:"website_id"`
}

type WebsiteAdvancedConfig struct {
	Base
	WebsiteID           uint           `gorm:"website_id;uniqueIndex" json:"website_id"`
	RateLimit           datatypes.JSON `gorm:"rate_limit"            json:"rate_limit"`
	LoadBalancer        datatypes.JSON `gorm:"load_balancer"         json:"load_balancer"`
	AdvancedSettings    datatypes.JSON `gorm:"advanced_settings"     json:"advanced_settings"`
	IsRateLimitEnabled  bool           `gorm:"is_rate_limit_enabled;default=false" json:"is_rate_limit_enabled"`
	IsLoadBalancerEnabled bool          `gorm:"is_load_balancer_enabled;default=false" json:"is_load_balancer_enabled"`
}

type RateLimitConfig struct {
	Enabled          bool   `json:"enabled"`
	Period           int    `json:"period"`           // 时间窗口（秒）
	Limit            int    `json:"limit"`            // 限制数量
	BlockDuration    int    `json:"block_duration"`   // 封禁时长（秒）
	Strategy         string `json:"strategy"`         // 限流策略：ip、session、user_agent
	WhitelistIPs     []string `json:"whitelist_ips"`  // 白名单IP
}
