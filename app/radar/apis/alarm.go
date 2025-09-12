package apis

import (
	"go-admin/app/radar/models"
	"go-admin/app/radar/service"
	"go-admin/app/radar/service/dto"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/utils"
)

type Alarm struct {
	api.Api
}

// 获取预警规则
func (e Alarm) GetAlarmRules(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			e.Logger.Error("获取预警规则出错:", err)
		}
	}()

	req := dto.GetAlarmRulesReq{}
	s := service.Alarm{}
	if err = e.MakeContext(c).MakeOrm().Bind(&req).MakeService(&s.Service).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 之后可以在这里验证一下权限
	// deptId := user.GetDeptId(c)

	var alarmRuleList []models.AlarmRule
	var alarmRuleLevelList []models.AlarmRuleLevel
	deptId := req.DeptId
	if alarmRuleList, alarmRuleLevelList, err = s.GetAlarmRules(deptId); err != nil {
		e.Error(500, err, "获取警报规则失败")
		return
	}

	resp := dto.GetAlarmRulesResp{
		AlarmRuleList:      alarmRuleList,
		AlarmRuleLevelList: alarmRuleLevelList,
	}

	e.OK(resp, "success")
}

// 增加预警规则
func (e Alarm) AddAlarmRule(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			e.Logger.Error("增加预警规则出错:", err)
		}
	}()

	req := dto.AddAlarmRuleReq{}
	s := service.Alarm{}
	if err = e.MakeContext(c).MakeOrm().Bind(&req).MakeService(&s.Service).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	levelItems := req.AlarmRuleLevelItem
	if len(levelItems) != 4 {
		// e.Error(400, errors.New("参数错误"), "参数错误")
		utils.ParameterError("参数错误")
		return
	}
	var levels = make([]models.AlarmLevel, 4)
	for i, item := range levelItems {
		levels[i] = item.AlarmLevel
	}
	if !checkAlarmLevel(levels) {
		utils.ParameterError("预警等级参数错误")
		return
	}

	// 之后可以在这里验证一下权限
	// deptId := user.GetDeptId(c)

	// var alarmRuleList []models.AlarmRule
	// var alarmRuleLevelList []models.AlarmRuleLevel
	// deptId := req.DeptId
	if err = s.AddAlarmRule(req); err != nil {
		e.Error(500, err, "获取警报规则失败")
		return
	}

	resp := dto.AddAlarmRuleResp{
		// AlarmRuleList:      alarmRuleList,
		// AlarmRuleLevelList: alarmRuleLevelList,
	}

	e.OK(resp, "success")
}

// 修改预警规则
func (e Alarm) UpdateAlarmRule(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			e.Logger.Error("修改预警规则出错:", err)
		}
	}()

	req := dto.UpdateAlarmRuleReq{}
	s := service.Alarm{}
	if err = e.MakeContext(c).MakeOrm().Bind(&req).MakeService(&s.Service).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	levelItems := req.AlarmRuleLevelItem
	if len(levelItems) != 4 {
		utils.ParameterError("参数错误")
		return
	}
	var levels = make([]models.AlarmLevel, 4)
	for i, item := range levelItems {
		levels[i] = item.AlarmLevel
	}

	if !checkAlarmLevel(levels) {
		utils.ParameterError("预警等级参数错误")
		return
	}

	// 之后可以在这里验证一下权限
	if err = s.UpdateAlarmRule(req); err != nil {
		e.Error(500, err, "获取警报规则失败")
		return
	}

	resp := dto.AddAlarmRuleResp{
		// AlarmRuleList:      alarmRuleList,
		// AlarmRuleLevelList: alarmRuleLevelList,
	}

	e.OK(resp, "success")
}

// 删除预警规则
func (e Alarm) DeleteAlarmRule(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			e.Logger.Error("删除预警规则出错:", err)
		}
	}()

	req := dto.DeleteAlarmRuleReq{}
	s := service.Alarm{}
	if err = e.MakeContext(c).MakeOrm().Bind(&req).MakeService(&s.Service).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return

	}

	alarmRuleId := req.AlarmRuleId

	// 之后可以在这里验证一下权限
	if err = s.DeleteAlarmRule(alarmRuleId); err != nil {
		e.Error(500, err, "获取警报规则失败")
		return
	}

	resp := dto.DeleteAlarmRuleResp{}
	e.OK(resp, "success")
}

func checkAlarmLevel(levels []models.AlarmLevel) bool {
	var rightLevels = []models.AlarmLevel{1, 2, 3, 4}
	return !slices.Equal(levels, rightLevels)
}
