package models

import "time"

// AlarmPoint 预警规则
type AlarmPoint struct {
	Id int64 `json:"id"             gorm:"primaryKey;autoIncrement;comment:主键编码"`

	AlarmCheckType AlarmCheckType `json:"alarmCheckType" gorm:"type:tinyint;  comment:监测类型"`                          //监测类型
	AlarmName      string         `json:"alarmName"      gorm:"size:64;   comment:判据名称"`                              //判据名称
	RadarId        int64          `json:"radarId"        gorm:"type:tinyint;comment:雷达Id"`                            //雷达Id
	Mode           int64          `json:"mode"           gorm:"type:tinyint;comment:全局/局部"`                           //监测点ID (0=对机构全局生效)
	DeptId         int64          `json:"deptId"         gorm:"uniqueIndex:idx_rule_key; comment:机构Id"`               //机构ID
	RadarPointId   int64          `json:"radarPointId"   gorm:"uniqueIndex:idx_rule_key; type:tinyint;comment:监测点ID"` //监测点ID (0=对机构全局生效)
	AlarmType      AlarmType      `json:"alarmType"      gorm:"uniqueIndex:idx_rule_key; type:tinyint;comment:预警类型"`  //预警类型
	RedOption      string         `json:"redOption"      gorm:"size:255;comment:红色预警条件"`                              //预警条件
	OrangeOption   string         `json:"orangeOption"   gorm:"size:255;comment:红色预警条件"`                              //预警条件
	YellowOption   string         `json:"yellowOption"   gorm:"size:255;comment:红色预警条件"`                              //预警条件
	BlueOption     string         `json:"blueOption"     gorm:"size:255;comment:红色预警条件"`                              //预警条件
	Interval       uint64         `json:"interval"       gorm:"comment:预警间隔时间（分）"`                                    //预警间隔时间（分）
	Duration       uint64         `json:"duration"       gorm:"comment:连续预警次数"`                                       //连续预警次数
	CreatedAt      time.Time      `json:"createdAt"      gorm:"comment:创建时间"`
	UpdatedAt      time.Time      `json:"updatedAt"      gorm:"comment:最后更新时间"`
}

func (AlarmPoint) TableName() string {
	return "alarm_point"
}
