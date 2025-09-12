package models

import "time"

// AlarmContactGroup 预警联系人组
type AlarmContactGroup struct {
	Id          int64  `gorm:"primaryKey;autoIncrement"`
	DeptId      int64  `json:"deptId"  gorm:"comment:机构"` //机构ID
	Name        string `gorm:"uniqueIndex:uniq_name; size:50;not null;comment:联系人组名称"`
	Description string `gorm:"size:255;comment:描述"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (AlarmContactGroup) TableName() string {
	return "alarm_contact_group"
}
