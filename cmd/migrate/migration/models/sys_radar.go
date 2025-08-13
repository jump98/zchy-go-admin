package models

import (
	"go-admin/common/models"
)

type SysRadar struct {
	//models.Model
	RadarId int64 `json:"radarId" gorm:"primaryKey;autoIncrement;column:radar_id;comment:主键编码"`

	RadarName   string   `json:"radarName" gorm:"type:varchar(64);comment:雷达名称"`
	RadarKey    string   `json:"radarKey" gorm:"type:varchar(100);comment:雷达编号"`
	SpecialKey  string   `json:"specialKey" gorm:"type:varchar(100);comment:雷达特殊编号"`
	DeptId      int64    `json:"deptId" gorm:"type:bigint;comment:部门"`
	Lng         string   `json:"lng" gorm:"type:bigint;comment:经度"`
	Lat         string   `json:"lat" gorm:"type:bigint;comment:纬度"`
	Alt         string   `json:"alt" gorm:"type:bigint;comment:高度"`
	Remark      string   `json:"remark" gorm:"type:varchar(255);comment:备注"`
	Status      string   `json:"status" gorm:"type:varchar(4);comment:状态"`
	Vender      string   `json:"vender" gorm:"size:100;comment:设备厂家名"`
	Secret      string   `json:"secret" gorm:"size:100;comment:密钥"`
	FromProject int      `json:"fromProject" gorm:"column:from_project;size:4;"` //是否是自动创建，当来自项目时为1
	Dept        *SysDept `json:"dept"`
	models.ModelTime
	models.ControlBy
}

func (SysRadar) TableName() string {
	return "sys_radar"
}
