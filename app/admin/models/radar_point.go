package models

import (
	"go-admin/common/models"

	"gorm.io/gorm"
)

type RadarPoint struct {
	models.Model
	//Id int64 `json:"Id" gorm:"primaryKey;autoIncrement;column:id;comment:主键编码"`

	PointName  string `json:"pointName" gorm:"type:varchar(64);comment:监测点名称"`
	PointKey   string `json:"pointKey" gorm:"type:varchar(100);comment:监测点编号"`
	PointType  string `json:"pointType" gorm:"type:varchar(10);comment:监测点类型"`
	RadarId    int64  `json:"radarId" gorm:"type:bigint;comment:雷达ID"`
	Lng        string `json:"lng" gorm:"type:bigint;comment:经度"`
	Lat        string `json:"lat" gorm:"type:bigint;comment:纬度"`
	Alt        string `json:"alt" gorm:"type:bigint;comment:高度"`
	Distance   string `json:"distance" gorm:"type:bigint;comment:距离"`
	PointIndex int    `json:"pointIndex" gorm:"type:int;comment:下标"`
	Remark     string `json:"remark" gorm:"type:varchar(255);comment:备注"`
	AStatus    string `json:"aStatus" gorm:"type:bigint;comment:激活状态"`
	XStatus    string `json:"xStatus" gorm:"type:bigint;comment:消警状态"`
	MTypeId    string `json:"mTypeId" gorm:"type:bigint;comment:门限类型"`
	models.ModelTime
	models.ControlBy
}

func (RadarPoint) TableName() string {
	return "radar_point"
}

func (e *RadarPoint) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *RadarPoint) GetId() interface{} {
	return e.Id
	//return e.Id
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
