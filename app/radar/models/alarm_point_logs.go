package models

import "time"

// AlarmPointLogs 预警记录
type AlarmPointLogs struct {
	Id           int64      `json:"id"             gorm:"primaryKey;autoIncrement;comment:主键编码"`
	AlarmType    AlarmType  `json:"alarmType"      gorm:"type:tinyint;comment:预警类型"`  //预警类型
	RadarId      int64      `json:"radarId"        gorm:"type:tinyint;comment:雷达Id"`  //雷达Id
	RadarPointId int64      `json:"radarPointId"   gorm:"type:tinyint;comment:监测点ID"` //监测点ID
	AlarmLevel   AlarmLevel `json:"alarmLevel"     gorm:"type:tinyint;comment:报警等级"`  //报警等级
	DeptId       int64      `json:"deptId"         gorm:"comment:机构Id"`               //机构ID
	LimitValue   string     `json:"limitValue"     gorm:"comment:门限值"`                //门限值
	AlarmValue   string     `json:"alarmValue"     gorm:"comment:记录值"`                //记录值
	Interval     int64      `json:"interval"       gorm:"comment:预警间隔时间"`             //预警间隔时间（分）
	Duration     uint64     `json:"duration"       gorm:"comment:报警次数"`               //连续预警次数
	OperatorId   int64      `json:"operatorId"     gorm:"comment:操作人ID"`
	//Processed     bool      `json:"processed" gorm:"comment:是否处理完成"`
	ProcessRemark string    `json:"processRemark" gorm:"type:text;comment:处理备注"`
	CreatedAt     time.Time `json:"createdAt"      gorm:"comment:创建时间"`
	UpdatedAt     time.Time `json:"updatedAt"      gorm:"comment:最后更新时间"`
}

func (AlarmPointLogs) TableName() string {
	return "alarm_point_logs"
}
