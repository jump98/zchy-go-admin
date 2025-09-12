package models

import "time"

// AlarmRule 预警规则
type AlarmRule struct {
	Id             int64          `json:"id"             gorm:"primaryKey;autoIncrement;comment:主键编码"`
	DeptId         int64          `json:"deptId"         gorm:"comment:机构Id"`                //机构ID
	AlarmCheckType AlarmCheckType `json:"alarmCheckType" gorm:"type:tinyint;  comment:监测类型"` //监测类型
	AlarmName      string         `json:"alarmName"      gorm:"size:64;   comment:判据名称"`     //判据名称
	Remark         string         `json:"remark"         gorm:"size:255;  comment:判据简介"`     //判据简介
	CreatedAt      time.Time      `json:"createdAt"      gorm:"comment:创建时间"`
	UpdatedAt      time.Time      `json:"updatedAt"      gorm:"comment:最后更新时间"`
}

func (AlarmRule) TableName() string {
	return "alarm_rule"
}
