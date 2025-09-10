package models

import (
	"go-admin/common/models"
)

// 预警通知类型
type AlarmNotifyType int

const (
	//预警通知类型-邮件
	AlarmNotifyType_EMail AlarmNotifyType = 1
	//预警通知类型-短信
	AlarmNotifyType_SMS AlarmNotifyType = 2
)

// 预警配置-
type AlarmConfig struct {
	Id           int64           `json:"id"               gorm:"primaryKey;autoIncrement;comment:主键编码"`
	DeptId       int64           `json:"deptId"           gorm:"uniqueIndex:uniq_deptid_tyle_ruleid;  comment:机构"`                  //机构ID
	AlarmType    AlarmRuleType   `json:"alarmType"        gorm:"uniqueIndex:uniq_deptid_tyle_ruleid;  type:tinyint;  comment:预警类型"` //预警类型
	AlarmRuleId  int64           `json:"alarmRuleId"      gorm:"uniqueIndex:uniq_deptid_tyle_ruleid;  comment:预警规则ID"`              //预警规则ID
	RadarId      int64           `json:"radarId"          gorm:"comment:雷达ID"`                                                      //雷达ID
	RadarPointId int64           `json:"radarPointId"     gorm:"comment:监控点ID"`                                                     //雷达监控的ID
	NotifyType   AlarmNotifyType `json:"notifyType"       gorm:"type:tinyint;comment:预警通知类型"`                                       //预警通知类型
	Interval     uint64          `json:"interval"         gorm:"comment:预警间隔时间（分）"`                                                 //预警间隔时间（分）
	Duration     uint64          `json:"duration"         gorm:"comment:连续预警次数"`                                                    //连续预警次数
	models.ModelTime
}

func (AlarmConfig) TableName() string {
	return "alarm_config"
}
