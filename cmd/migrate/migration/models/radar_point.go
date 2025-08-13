package models

import (
	"go-admin/common/models"
)

type RadarPoint struct {
	ID         int    `json:"id" gorm:"primaryKey;type:bigint;autoIncrement;comment:PointId"`
	PointName  string `json:"pointName" gorm:"type:varchar(64);comment:监测点名称"`
	PointKey   string `json:"radarKey" gorm:"type:varchar(100);comment:监测点编号"`
	PointType  string `json:"pointType" gorm:"type:varchar(10);comment:监测点类型"`
	RadarId    int    `json:"radarId" gorm:"type:bigint;comment:雷达ID"`
	Lng        int    `json:"lng" gorm:"type:bigint;comment:经度"`
	Lat        int    `json:"lat" gorm:"type:bigint;comment:纬度"`
	Alt        int    `json:"alt" gorm:"type:bigint;comment:高度"`
	Distance   int    `json:"distance" gorm:"type:bigint;comment:距离"`
	PointIndex int    `json:"pointIndex" gorm:"type:int;comment:下标"`
	Remark     string `json:"remark" gorm:"type:varchar(255);comment:备注"`
	AStatus    string `json:"aStatus" gorm:"type:varchar(4);comment:激活状态"`
	XStatus    string `json:"xStatus" gorm:"type:varchar(4);comment:消警状态"`
	MTypeID    string `json:"mTypeId" gorm:"type:varchar(4);comment:门限类型"`
	models.ControlBy
	models.ModelTime
}

func (*RadarPoint) TableName() string {
	return "radar_point"
}
