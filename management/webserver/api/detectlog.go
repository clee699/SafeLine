package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"chaitin.cn/dev/go/errors"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg"

	"github.com/gin-gonic/gin"

	"chaitin.cn/dev/go/log"

	"chaitin.cn/patronus/safeline-2/management/webserver/api/response"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
	"chaitin.cn/patronus/safeline-2/management/webserver/utils"
)

type (
	GetDetectLogDetailRequest struct {
		EventId string `json:"event_id"   form:"event_id"`
	}

	PostFalsePositivesRequest struct {
		EventId string `json:"event_id"`
	}

	telemetryFalsePositives struct {
		Telemetry struct {
			Id string `json:"id"`
		} `json:"telemetry"`

		Safeline struct {
			Id        string          `json:"id"`
			Type      string          `json:"type"`
			DetectLog model.DetectLog `json:"detect_log"`
		} `json:"safeline"`
	}
)

func getDetectLog(eventId string) (*model.DetectLog, error) {
	db := database.GetDB()
	var detectLogBasic model.DetectLogBasic
	res := db.Where(&model.DetectLogBasic{EventId: eventId}).First(&detectLogBasic)
	if res.RowsAffected == 0 {
		return nil, errors.New("Data queried does not exist")
	}

	var detectLogDetail model.DetectLogDetail
	db.Where(&model.DetectLogDetail{EventId: eventId}).First(&detectLogDetail)

	detectLog, err := model.TransformDetectLog(&detectLogBasic, &detectLogDetail)
	if err != nil {
		return nil, err
	}

	return detectLog, nil
}

func GetDetectLogList(c *gin.Context) {
	var params pageRequest
	if err := c.BindQuery(&params); err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}
	db := database.GetDB()

	tx := db.Where("")
	// 按 ip 搜索条件
	if ip := c.Query("ip"); ip != "" {
		tx = tx.Where("src_ip = ?", ip)
	}
	// 按 url 搜索条件
	if url := c.Query("url"); url != "" {
		tx = tx.Where("url_path like ?", "%"+url+"%")
	}
	// 按 type 搜索条件
	if at := c.Query("attack_type"); at != "" {
		ns := make([]int, 0)
		for _, s := range strings.Split(at, ",") {
			n, err := strconv.Atoi(s)
			if err == nil {
				ns = append(ns, n)
			}
		}
		tx = tx.Where("attack_type in (?)", ns)
	}
	// 按验证类型搜索条件
	if verifyType := c.Query("verify_type"); verifyType != "" {
		verifyTypes := strings.Split(verifyType, ",")
		attackTypes := make([]int, 0)
		for _, vt := range verifyTypes {
			switch vt {
			case "captcha":
				attackTypes = append(attackTypes, 65) // 人机验证
			case "auth":
				attackTypes = append(attackTypes, 66) // 身份认证
			case "waiting_room":
				attackTypes = append(attackTypes, 67) // 等候室
			}
		}
		if len(attackTypes) > 0 {
			tx = tx.Where("attack_type in (?)", attackTypes)
		}
	}

	var total int64
	tx.Model(&model.DetectLogBasic{}).Count(&total)

	var basicList []model.DetectLogBasic
	tx.Limit(params.PageSize).Offset(params.PageSize * (params.Page - 1)).Order("id desc").Find(&basicList)
	var dLogList []*model.DetectLog
	for _, basic := range basicList {
		dLog, err := model.TransformDetectLog(&basic, nil)
		if err != nil {
			logger.Warn(err)
			continue
		}
		dLogList = append(dLogList, dLog)
	}
	response.Success(c, gin.H{"data": dLogList, "total": total})
}

// GetVerificationLogList 获取人机、身份认证、等候室日志
func GetVerificationLogList(c *gin.Context) {
	var params pageRequest
	if err := c.BindQuery(&params); err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}
	db := database.GetDB()

	// 只查询验证类型的日志
	tx := db.Where("attack_type in (65, 66, 67)")
	// 按 ip 搜索条件
	if ip := c.Query("ip"); ip != "" {
		tx = tx.Where("src_ip = ?", ip)
	}
	// 按 url 搜索条件
	if url := c.Query("url"); url != "" {
		tx = tx.Where("url_path like ?", "%"+url+"%")
	}
	// 按具体验证类型搜索
	if verifyType := c.Query("verify_type"); verifyType != "" {
		var attackType int
		switch verifyType {
		case "captcha":
			attackType = 65 // 人机验证
		case "auth":
			attackType = 66 // 身份认证
		case "waiting_room":
			attackType = 67 // 等候室
		}
		tx = tx.Where("attack_type = ?", attackType)
	}

	var total int64
	tx.Model(&model.DetectLogBasic{}).Count(&total)

	var basicList []model.DetectLogBasic
	tx.Limit(params.PageSize).Offset(params.PageSize * (params.Page - 1)).Order("id desc").Find(&basicList)
	var dLogList []*model.DetectLog
	for _, basic := range basicList {
		dLog, err := model.TransformDetectLog(&basic, nil)
		if err != nil {
			logger.Warn(err)
			continue
		}
		dLogList = append(dLogList, dLog)
	}
	response.Success(c, gin.H{"data": dLogList, "total": total})
}

func GetDetectLogDetail(c *gin.Context) {
	var params GetDetectLogDetailRequest
	if err := c.BindQuery(&params); err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	detectLog, err := getDetectLog(params.EventId)
	if err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorDataNotExist, http.StatusNotFound)
		return
	}

	response.Success(c, detectLog)
}

func PostFalsePositives(c *gin.Context) {
	var params PostFalsePositivesRequest
	if err := c.BindJSON(&params); err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}

	detectLog, err := getDetectLog(params.EventId)
	if err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorDataNotExist, http.StatusNotFound)
		return
	}

	db := database.GetDB()
	var option model.Options
	db.Where(&model.Options{Key: constants.MachineID}).First(&option)

	var jsonData telemetryFalsePositives
	jsonData.Telemetry.Id = constants.TelemetryId
	jsonData.Safeline.Id = option.Value
	jsonData.Safeline.Type = constants.FalsePositives
	jsonData.Safeline.DetectLog = *detectLog

	data, err := json.Marshal(jsonData)
	if err != nil {
		log.Warn(err)
		response.Success(c, nil)
		return
	}
	reader := bytes.NewReader(data)

	client := utils.GetHTTPClient()
	addr := config.GlobalConfig.Telemetry.Addr

	rsp, err := pkg.DoPostTelemetry(client, addr, reader)
	if err != nil {
		log.Warn(err)
		response.Success(c, nil)
		return
	}

	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusCreated {
		log.Warn(fmt.Sprintf("Transfer telemetry failed, status code = %d", rsp.StatusCode), err)
		response.Success(c, nil)
		return
	}

	response.Success(c, nil)
}

// GetDetectLogDownload 下载攻击日志
func GetDetectLogDownload(c *gin.Context) {
	var params pageRequest
	if err := c.BindQuery(&params); err != nil {
		logger.Error(err)
		response.Error(c, response.ErrorParamNotOK, http.StatusInternalServerError)
		return
	}
	db := database.GetDB()

	tx := db.Where("")
	// 按 ip 搜索条件
	if ip := c.Query("ip"); ip != "" {
		tx = tx.Where("src_ip = ?", ip)
	}
	// 按 url 搜索条件
	if url := c.Query("url"); url != "" {
		tx = tx.Where("url_path like ?", "%"+url+"%")
	}
	// 按 type 搜索条件
	if at := c.Query("attack_type"); at != "" {
		ns := make([]int, 0)
		for _, s := range strings.Split(at, ",") {
			n, err := strconv.Atoi(s)
			if err == nil {
				ns = append(ns, n)
			}
		}
		tx = tx.Where("attack_type in (?)", ns)
	}
	// 按验证类型搜索条件
	if verifyType := c.Query("verify_type"); verifyType != "" {
		verifyTypes := strings.Split(verifyType, ",")
		attackTypes := make([]int, 0)
		for _, vt := range verifyTypes {
			switch vt {
			case "captcha":
				attackTypes = append(attackTypes, 65) // 人机验证
			case "auth":
				attackTypes = append(attackTypes, 66) // 身份认证
			case "waiting_room":
				attackTypes = append(attackTypes, 67) // 等候室
			}
		}
		if len(attackTypes) > 0 {
			tx = tx.Where("attack_type in (?)", attackTypes)
		}
	}

	// 获取所有符合条件的数据，不进行分页
	var basicList []model.DetectLogBasic
	tx.Order("id desc").Find(&basicList)

	// 转换为CSV格式
	var csvBuffer bytes.Buffer
	// 写入CSV表头
	csvBuffer.WriteString("事件ID,时间,IP地址,URL,攻击类型,处理动作,规则模块,规则描述\n")

	for _, basic := range basicList {
		dLog, err := model.TransformDetectLog(&basic, nil)
		if err != nil {
			logger.Warn(err)
			continue
		}

		// 获取攻击类型名称
		attackTypeName := "未知"
		if name, ok := constants.AttackType[dLog.AttackType]; ok {
			attackTypeName = name
		}

		// 获取处理动作名称
		actionName := "未知"
		switch dLog.Action {
		case 0:
			actionName = "放行"
		case 1:
			actionName = "拦截"
		case 2:
			actionName = "告警"
		}

		// 获取规则模块名称
		ruleModuleName := "未知"
		if name, ok := constants.RuleModule[dLog.RuleModule]; ok {
			ruleModuleName = name
		}

		// 写入CSV行
		csvBuffer.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n",
			dLog.EventId,
			dLog.Timestamp.Format("2006-01-02 15:04:05"),
			dLog.SrcIP,
			dLog.UrlPath,
			attackTypeName,
			actionName,
			ruleModuleName,
			dLog.RuleReason,
		))
	}

	// 设置HTTP头，使浏览器能够下载文件
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=attack_logs_%s.csv", time.Now().Format("20060102150405")))
	c.Header("Content-Length", fmt.Sprintf("%d", csvBuffer.Len()))

	// 写入响应体
	c.String(http.StatusOK, csvBuffer.String())
}
