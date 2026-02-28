package constants

const (
	SuperUser      = "admin"
	ProductName    = "长亭雷池 WAF 社区版"
	ProductVersion = ""
	ConfigFilePath = "config.yml"
	CertsPath      = "certs"
)

const (
	NotUpgrade int = iota
	RecommendedUpgrade
	MustUpgrade
)

// 告警相关常量
const (
	// 告警全局开关
	AlertEnabled = "alert_enabled"
	// 告警检查间隔（秒）
	AlertCheckInterval = "alert_check_interval"
	// 告警类型
	AlertTypeEmail = "email"
	AlertTypeWeChat = "wechat"
)
