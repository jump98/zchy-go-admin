package models

import "time"

// 预警规则类型
type AlarmRuleType int

const (
	//雷达-状态信息
	AlarmRuleType_RadarStatus AlarmRuleType = 1
	//雷达点位-地表位移
	AlarmRuleType_PointMove AlarmRuleType = 2
	//雷达点位-速度
	AlarmRuleType_PointVelocity AlarmRuleType = 3
	//雷达点位-加速度
	AlarmRuleType_PointAcceleration AlarmRuleType = 4
)

// 预警规则
type AlarmRule struct {
	Id        int64         `json:"id"             gorm:"primaryKey;autoIncrement;comment:主键编码"`
	DeptId    int64         `json:"deptId"         gorm:"comment:机构Id"`                //机构ID
	AlarmType AlarmRuleType `json:"alarmType"      gorm:"type:tinyint;  comment:预警类型"` //预警类型
	AlarmName string        `json:"alarmName"      gorm:"size:64;   comment:规则名称"`     //预警规则名称
	Remark    string        `json:"remark"         gorm:"size:255;  comment:规则介绍"`     //预警规则介绍
	CreatedAt time.Time     `json:"createdAt"      gorm:"comment:创建时间"`
	UpdatedAt time.Time     `json:"updatedAt"      gorm:"comment:最后更新时间"`
}

func (AlarmRule) TableName() string {
	return "alarm_rule"
}
