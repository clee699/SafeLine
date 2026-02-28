package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chaitin.cn/dev/go/log"
	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
)

var (
	alertLogger = log.GetLogger("alert_api")
)

// AlertConfigRequest 告警配置请求结构
type AlertConfigRequest struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" binding:"required"`
	Type        int    `json:"type" binding:"required"`
	Enabled     bool   `json:"enabled"`
	Config      string `json:"config" binding:"required"`
	AlertLevels string `json:"alert_levels"`
	AttackTypes string `json:"attack_types"`
}

// GetAlertConfig 获取告警配置列表
func GetAlertConfig(c *gin.Context) {
	var alertConfigs []model.AlertConfig
	db := database.GetDB()
	db.Find(&alertConfigs)

	response.Success(c, alertConfigs)
}

// PostAlertConfig 创建告警配置
func PostAlertConfig(c *gin.Context) {
	var req AlertConfigRequest
	if err := c.BindJSON(&req); err != nil {
		alertLogger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusBadRequest)
		return
	}

	alertConfig := model.AlertConfig{
		Name:        req.Name,
		Type:        model.AlertType(req.Type),
		Enabled:     req.Enabled,
		Config:      req.Config,
		AlertLevels: req.AlertLevels,
		AttackTypes: req.AttackTypes,
	}

	db := database.GetDB()
	if err := db.Create(&alertConfig).Error; err != nil {
		alertLogger.Error(err)
		response.Error(c, response.ErrorInternal, http.StatusInternalServerError)
		return
	}

	response.Success(c, alertConfig)
}

// PutAlertConfig 更新告警配置
func PutAlertConfig(c *gin.Context) {
	var req AlertConfigRequest
	if err := c.BindJSON(&req); err != nil {
		alertLogger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusBadRequest)
		return
	}

	if req.ID == 0 {
		response.Error(c, response.ErrorParamNotOK, http.StatusBadRequest)
		return
	}

	var alertConfig model.AlertConfig
	db := database.GetDB()
	if err := db.First(&alertConfig, req.ID).Error; err != nil {
		alertLogger.Error(err)
		response.Error(c, response.ErrorDataNotExist, http.StatusNotFound)
		return
	}

	// 更新告警配置
	alertConfig.Name = req.Name
	alertConfig.Type = model.AlertType(req.Type)
	alertConfig.Enabled = req.Enabled
	alertConfig.Config = req.Config
	alertConfig.AlertLevels = req.AlertLevels
	alertConfig.AttackTypes = req.AttackTypes

	if err := db.Save(&alertConfig).Error; err != nil {
		alertLogger.Error(err)
		response.Error(c, response.ErrorInternal, http.StatusInternalServerError)
		return
	}

	response.Success(c, alertConfig)
}

// DeleteAlertConfig 删除告警配置
func DeleteAlertConfig(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.Error(c, response.ErrorParamNotOK, http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	if err := db.Delete(&model.AlertConfig{}, id).Error; err != nil {
		alertLogger.Error(err)
		response.Error(c, response.ErrorInternal, http.StatusInternalServerError)
		return
	}

	response.Success(c, nil)
}
