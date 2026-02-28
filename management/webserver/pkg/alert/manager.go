package alert

import (
	"time"

	"chaitin.cn/dev/go/log"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
	"gorm.io/gorm"
)

// AlertManager 告警管理器
type AlertManager struct {
	lastCheckTime  int64
	checkInterval  int
}

// NewAlertManager 创建告警管理器
func NewAlertManager() *AlertManager {
	return &AlertManager{
		lastCheckTime: time.Now().Unix(),
		checkInterval: config.GlobalConfig.Alert.CheckInterval,
	}
}

// Start 启动告警管理器
func (am *AlertManager) Start() {
	// 使用配置文件中的检查间隔，如果未配置则使用默认值10秒
	interval := am.checkInterval
	if interval <= 0 {
		interval = 10
	}
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			am.checkAndSendAlerts()
		}
	}
}

// checkAndSendAlerts 检查新的检测日志并发送告警
func (am *AlertManager) checkAndSendAlerts() {
	now := time.Now().Unix()
	logger := log.GetLogger("alert")
	
	// 检查全局告警开关
	db := database.GetDB()
	var alertEnabled model.Options
	db.Where("key = ?", constants.AlertEnabled).First(&alertEnabled)
	if alertEnabled.Value != "true" {
		// 全局告警开关已关闭
		return
	}
	
	// 检查配置文件中的全局开关
	if !config.GlobalConfig.Alert.Enabled {
		// 配置文件中已关闭告警
		return
	}

	// 获取所有启用的告警配置
	var alertConfigs []model.AlertConfig
	db.Where("enabled = ?", true).Find(&alertConfigs)

	if len(alertConfigs) == 0 {
		return
	}

	// 获取默认告警级别
	defaultLevel := config.GlobalConfig.Alert.DefaultLevel
	if defaultLevel <= 0 {
		defaultLevel = 2 // 默认高风险
	}

	// 获取上次检查时间以来的所有符合条件的检测日志
	var detectLogBasics []model.DetectLogBasic
	var detectLogDetails []model.DetectLogDetail

	db.Where("timestamp > ? AND risk_level >= ?", am.lastCheckTime, defaultLevel).Find(&detectLogBasics)
	
	// 检查是否超过最大告警数量
	if len(detectLogBasics) > config.GlobalConfig.Alert.MaxAlertCount {
		logger.Warnf("检测到 %d 条告警，超过最大限制 %d，将只发送部分告警", len(detectLogBasics), config.GlobalConfig.Alert.MaxAlertCount)
		detectLogBasics = detectLogBasics[:config.GlobalConfig.Alert.MaxAlertCount]
	}
	
	if len(detectLogBasics) == 0 {
		am.lastCheckTime = now
		return
	}

	// 获取对应的详细日志
	var eventIDs []string
	for _, basic := range detectLogBasics {
		eventIDs = append(eventIDs, basic.EventId)
	}
	db.Where("event_id IN ?", eventIDs).Find(&detectLogDetails)

	// 构建事件ID到详细日志的映射
	detailMap := make(map[string]model.DetectLogDetail)
	for _, detail := range detectLogDetails {
		detailMap[detail.EventId] = detail
	}

	// 为每个告警配置创建告警服务
	for _, config := range alertConfigs {
		alertService, err := NewAlertService(&config)
		if err != nil {
			logger.Warnf("创建告警服务失败: %v", err)
			continue
		}

		// 为每个检测日志发送告警
		for _, basic := range detectLogBasics {
			detail, exists := detailMap[basic.EventId]
			if !exists {
				logger.Warnf("事件ID %s 对应的详细日志不存在", basic.EventId)
				continue
			}

			// 获取攻击类型名称
			attackType, ok := constants.AttackType[basic.AttackType]
			if !ok {
				attackType = constants.AttackType[62] // unknown
			}

			// 构建告警消息
			alertMsg := &AlertMessage{
				EventID:    basic.EventId,
				SiteUUID:   basic.SiteUUID,
				SrcIP:      basic.SrcIp,
				Host:       basic.Host,
				UrlPath:    basic.UrlPath,
				AttackType: attackType,
				RiskLevel:  basic.RiskLevel,
				Action:     basic.Action,
				RuleID:     basic.RuleId,
				Timestamp:  basic.Timestamp,
				Payload:    detail.Payload,
			}

			// 发送告警
			if err := alertService.SendAlert(alertMsg); err != nil {
				logger.Warnf("发送告警失败: %v", err)
				continue
			}

			logger.Infof("发送告警成功，事件ID: %s, 攻击类型: %s", basic.EventId, attackType)
		}
	}

	// 更新上次检查时间
	am.lastCheckTime = now
}
