package models

import (
	"go-admin/common/models"
	"strconv"

	"gorm.io/gorm"
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

func (e *SysRadar) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysRadar) GetId() interface{} {
	//return e.Model.Id
	return e.RadarId
}

func ConvertStringFloat(str string, bMulti bool) string {
	s, err := strconv.ParseFloat(str, 64)
	if err == nil {
		if bMulti {
			return strconv.FormatFloat(s*10000000, 'f', -1, 64)
		} else {
			return strconv.FormatFloat(s/10000000.0, 'f', -1, 64)
		}
	} else {
		return str
	}
}

func (e *SysRadar) ConvertLatLngAlt() error {
	e.Lng = ConvertStringFloat(e.Lng, true)
	e.Lat = ConvertStringFloat(e.Lat, true)
	e.Alt = ConvertStringFloat(e.Alt, true)
	return nil
}

func (e *SysRadar) BeforeCreate(_ *gorm.DB) error {
	return e.ConvertLatLngAlt()
}

func (e *SysRadar) BeforeUpdate(_ *gorm.DB) error {
	return e.ConvertLatLngAlt()
}

func (e *SysRadar) AfterFind(_ *gorm.DB) error {
	e.Lng = ConvertStringFloat(e.Lng, false)
	e.Lat = ConvertStringFloat(e.Lat, false)
	e.Alt = ConvertStringFloat(e.Alt, false)
	return nil
}
