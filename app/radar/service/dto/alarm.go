package dto

import (
	"go-admin/app/radar/models"
)

// GetAlarmRulesReq  获取预警规则列表
type GetAlarmRulesReq struct {
	DeptId int64 `json:"deptId"` //机构ID
}
type GetAlarmRulesResp struct {
	AlarmRuleList      []models.AlarmRule      `json:"alarmRule"`      //机构ID
	AlarmRuleLevelList []models.AlarmRuleLevel `json:"alarmRuleLevel"` //机构ID
}

// AddAlarmRuleReq 增加预警规则列表
type AddAlarmRuleReq struct {
	DeptId             int64                 `json:"deptId"`         //机构ID
	AlarmCheckType     models.AlarmCheckType `json:"alarmCheckType"` //监测类型
	AlarmName          string                `json:"alarmName"`      //预警规则名称
	Remark             string                `json:"remark"`         //预警规则介绍
	AlarmRuleLevelItem []AlarmRuleLevelItem
}
type AddAlarmRuleResp struct {
	AlarmRuleId       int64   `json:"alarmRuleId"`       //预警规则ID
	AlarmRuleLevelIds []int64 `json:"alarmRuleLevelIds"` //预警规则级别ID
}

// AlarmRuleLevelItem  预警规则等级表
type AlarmRuleLevelItem struct {
	AlarmLevel models.AlarmLevel          `json:"alarmLevel"` //预警等级
	Option     []models.AlarmRuleOption   `json:"option"`     //预警条件
	OptionMode models.AlarmRuleOptionMode `json:"optionMode"` //预警满足条件:all、or
	Suggestion string                     `json:"suggestion"`
	Horn       bool                       `json:"horn"`
}

// UpdateAlarmRuleReq 修改预警规则列表
type UpdateAlarmRuleReq struct {
	AlarmRuleId        int64                 `json:"alarmRuleId"`    //预警规矩ID
	DeptId             int64                 `json:"deptId"`         //机构ID
	AlarmCheckType     models.AlarmCheckType `json:"alarmCheckType"` //监测类型
	AlarmName          string                `json:"alarmName"`      //预警规则名称
	Remark             string                `json:"remark"`         //预警规则介绍
	AlarmRuleLevelItem []AlarmRuleLevelItem
}
type UpdateAlarmRuleResp struct {
}

// DeleteAlarmRuleReq  增加预警规则列表
type DeleteAlarmRuleReq struct {
	AlarmRuleId int64 `json:"alarmRuleId"` //预警规矩ID
}
type DeleteAlarmRuleResp struct {
}
