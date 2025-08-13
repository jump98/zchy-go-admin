package models

import (
	"go-admin/common/models"
)

type SysDeviceRadar struct {
	RadarId    int    `gorm:"primaryKey;autoIncrement;comment:RadarID"  json:"radarId"`
	RadarName  string `json:"radarName" gorm:"type:varchar(64);comment:雷达名称"`
	RadarKey   string `json:"radarKey" gorm:"type:varchar(100);comment:雷达编号"`
	SpecialKey string `json:"specialKey" gorm:"type:varchar(100);comment:雷达特殊编号"`
	DeptId     int    `json:"deptId" gorm:"type:bigint;comment:部门"`
	Lng        int    `json:"lng" gorm:"type:bigint;comment:经度"`
	Lat        int    `json:"lat" gorm:"type:bigint;comment:纬度"`
	Alt        int    `json:"alt" gorm:"type:bigint;comment:高度"`
	Remark     string `json:"remark" gorm:"type:varchar(255);comment:备注"`
	Status     string `json:"status" gorm:"type:varchar(4);comment:状态"`
	models.ControlBy
	models.ModelTime
}

func (*SysDeviceRadar) TableName() string {
	return "sys_radar"
}

func (e *SysDeviceRadar) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysDeviceRadar) GetId() interface{} {
	return e.RadarId
}
