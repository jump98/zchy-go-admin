package models

import (
	"database/sql"
	"go-admin/common/models"

	"gorm.io/gorm"
)

// RadarPointMType 监测点门限类型
type RadarPointMType int

const (
	RadarPointMTypeGlobal RadarPointMType = iota //预警-检测点
	RadarPointMTypeAlone                         //预警-检测点
)

// RadarPoint 雷达检测点表
type RadarPoint struct {
	Id            int64           `json:"id"          gorm:"primaryKey;autoIncrement;comment:主键编码"`
	PointName     string          `json:"pointName"   gorm:"type:varchar(64);  comment:监测点名称"`
	PointKey      string          `json:"pointKey"    gorm:"type:varchar(100); comment:监测点编号"`
	PointType     string          `json:"pointType"   gorm:"type:varchar(10);  comment:监测点类型"`
	RadarId       int64           `json:"radarId"     gorm:"uniqueIndex:idx_radarid_pointindex_key; comment:雷达ID"`
	PointIndex    int64           `json:"pointIndex"  gorm:"uniqueIndex:idx_radarid_pointindex_key; comment:监测点Index"`
	PoseDepth     int64           `json:"PoseDepth"   gorm:"DEFAULT:0;comment:位置滤波器缓存深度"`
	PhaseDepth    int64           `json:"PhaseDepth"  gorm:"DEFAULT:0;comment:相位滤波器缓存深度"`
	Lng           string          `json:"lng"         gorm:"type:varchar(20);  comment:经度"` //经度
	Lat           string          `json:"lat"         gorm:"type:varchar(20);  comment:纬度"` //纬度
	Alt           string          `json:"alt"         gorm:"type:varchar(20);  comment:高度"` //高度
	Distance      string          `json:"distance"    gorm:"type:varchar(20);  comment:距离"` //距离
	Remark        string          `json:"remark"      gorm:"type:varchar(255); comment:备注"`
	AStatus       string          `json:"aStatus"     gorm:"type:varchar(20);  comment:激活状态"`
	AlarmLevel    AlarmLevel      `json:"alarmLevel"  gorm:"type:tinyint; DEFAULT:0;  comment:预警状态"` //告警等级
	MTypeId       RadarPointMType `json:"mTypeId"     gorm:"type:tinyint; DEFAULT:0;  comment:门限类型"` //0=全局门限 1=独立门限
	LastAlarmTime sql.NullTime    `json:"last_time"   gorm:"comment:最近一次检测预警的时间"`                    //最近一次检测预警的时间
	models.ModelTime
	models.ControlBy
}

func (*RadarPoint) TableName() string {
	return "radar_point"
}

func (e *RadarPoint) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *RadarPoint) GetId() interface{} {
	return e.Id
}

func (e *RadarPoint) ConvertLatLngAlt() error {
	e.Lng = ConvertStringFloat(e.Lng, true)
	e.Lat = ConvertStringFloat(e.Lat, true)
	e.Alt = ConvertStringFloat(e.Alt, true)
	return nil
}

func (e *RadarPoint) BeforeCreate(_ *gorm.DB) error {
	return e.ConvertLatLngAlt()
}

func (e *RadarPoint) BeforeUpdate(_ *gorm.DB) error {
	return e.ConvertLatLngAlt()
}

func (e *RadarPoint) AfterFind(_ *gorm.DB) error {
	e.Lng = ConvertStringFloat(e.Lng, false)
	e.Lat = ConvertStringFloat(e.Lat, false)
	e.Alt = ConvertStringFloat(e.Alt, false)
	return nil
}
