package models

import (
	"time"
)

type AlarmRuleCondition struct {
	Field string  `json:"field"` //字段
	Op    string  `json:"op"`
	Value float64 `json:"value"`
}

type AlarmRuleOption struct {
	AlarmType AlarmType          `json:"alarmType"` //预警类型
	Cond      AlarmRuleCondition `json:"cond"`      //条件
}

// AlarmRuleLevel 预警规则等级表
type AlarmRuleLevel struct {
	Id          int64                 `json:"id"          gorm:"primaryKey;autoIncrement"`
	DeptId      int64                 `json:"deptId"      gorm:"comment:机构Id"` //机构ID
	AlarmRuleId int64                 `json:"alarmRuleId" gorm:"uniqueIndex:idx_rule_id_level_key;"`
	AlarmLevel  AlarmLevel            `json:"alarmLevel"  gorm:"uniqueIndex:idx_rule_id_level_key; type:tinyint;not null;comment:预警等级 1=红 2=橙 3=黄 4=蓝"` //预警等级 1=红 2=橙 3=黄 4=蓝
	Option      JSON[AlarmRuleOption] `json:"option"      gorm:"serializer:json;comment:条件"`                                                            //预警条件
	OptionMode  AlarmRuleOptionMode   `json:"optionMode"  gorm:"type:tinyint;    default:1;comment:条件组合方式"`                                             //预警满足条件:all、or
	Suggestion  string                `json:"suggestion"  gorm:"size:255;  comment:处理建议"`
	Horn        bool                  `json:"horn"        gorm:"type:tinyint(1) ;default:0;comment:喇叭是否开启"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (AlarmRuleLevel) TableName() string {
	return "alarm_rule_level"
}
