package models

import (
	"time"
)

// 预警等级枚举
type AlarmLevel int

const (
	AlarmLevelRed    AlarmLevel = 1 // 红色
	AlarmLevelOrange AlarmLevel = 2 // 橙色
	AlarmLevelYellow AlarmLevel = 3 // 黄色
	AlarmLevelBlue   AlarmLevel = 4 // 蓝色
)

// 预警满足条件
type AlarmRuleOptionMode int

const (
	//预警满足条件-全部满足
	AlarmRuleOptionMode_All AlarmRuleOptionMode = 1
	//预警满足条件-满足其一
	AlarmRuleOptionMode_Or AlarmRuleOptionMode = 2
)

type Condition struct {
	Field string  `json:"field"` //字段
	Op    string  `json:"op"`
	Value float64 `json:"value"`
}

// 预警规则等级表
type RadarAlarmRuleLevel struct {
	Id         int64               `json:"id"         gorm:"primaryKey;autoIncrement"`
	RuleId     int64               `json:"ruleId"     gorm:"uniqueIndex:idx_ruleid_level_key;"`
	AlarmLevel AlarmLevel          `json:"alarmLevel" gorm:"uniqueIndex:idx_ruleid_level_key; type:tinyint;not null;comment:预警等级 1=红 2=橙 3=黄 4=蓝"` //预警等级 1=红 2=橙 3=黄 4=蓝
	Option     JSON[Condition]     `json:"option"     gorm:"serializer:json;comment:条件"`                                                           // GORM v2                                          //预警条件
	OptionMode AlarmRuleOptionMode `json:"optionMode" gorm:"type:tinyint;    default:1;comment:条件组合方式"`                                            //预警满足条件:all、or
	Suggestion string              `json:"suggestion" gorm:"size:255;  comment:处理建议"`
	Horn       bool                `json:"horn"       gorm:"type:tinyint(1) ;default:0;comment:喇叭是否开启"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (RadarAlarmRuleLevel) TableName() string {
	return "alarm_rule_level"
}
