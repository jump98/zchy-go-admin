package models

import (
	adminModel "go-admin/app/admin/models"
	"go-admin/common/models"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Radar 雷达基本信息表
type Radar struct {
	RadarId    int64  `json:"radarId"     gorm:"column:radar_id;     primaryKey;  autoIncrement;comment:主键编码"`
	RadarName  string `json:"radarName"   gorm:"column:radar_name;   type:varchar(64);  comment:雷达名称"`
	RadarKey   string `json:"radarKey"    gorm:"column:radar_key;    uniqueIndex:idx_radar_key;   type:varchar(100); comment:雷达编号"`
	SpecialKey string `json:"specialKey"  gorm:"column:special_key;  type:varchar(100); comment:雷达特殊编号"`
	DeptId     int64  `json:"deptId"      gorm:"column:dept_id;      type:bigint;       comment:部门"`
	Lng        string `json:"lng"         gorm:"column:lng;          type:bigint;       comment:经度"`
	Lat        string `json:"lat"         gorm:"column:lat;          type:bigint;       comment:纬度"`
	Alt        string `json:"alt"         gorm:"column:alt;          type:bigint;       comment:高度"`
	Remark     string `json:"remark"      gorm:"column:remark;       type:varchar(255); comment:备注"`
	Status     string `json:"status"      gorm:"column:status;       type:varchar(4);   comment:状态"`
	Vender     string `json:"vender"      gorm:"column:vender;       size:100;          comment:设备厂家名"`
	Secret     string `json:"secret"      gorm:"column:secret;       size:100;          comment:密钥"`
	// FromProject int64               `json:"fromProject" gorm:"column:from_project; size:4;"` //是否是自动创建，当来自项目时为1
	Dept      *adminModel.SysDept `json:"dept"`
	CreatedAt time.Time           `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt time.Time           `json:"updatedAt" gorm:"comment:最后更新时间"`
	models.ControlBy
}

func (*Radar) TableName() string {
	return "radar"
}

func (e *Radar) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *Radar) GetId() interface{} {
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

func (e *Radar) ConvertLatLngAlt() error {
	e.Lng = ConvertStringFloat(e.Lng, true)
	e.Lat = ConvertStringFloat(e.Lat, true)
	e.Alt = ConvertStringFloat(e.Alt, true)
	return nil
}

func (e *Radar) BeforeCreate(_ *gorm.DB) error {
	return e.ConvertLatLngAlt()
}

func (e *Radar) BeforeUpdate(_ *gorm.DB) error {
	return e.ConvertLatLngAlt()
}

func (e *Radar) AfterFind(_ *gorm.DB) error {
	e.Lng = ConvertStringFloat(e.Lng, false)
	e.Lat = ConvertStringFloat(e.Lat, false)
	e.Alt = ConvertStringFloat(e.Alt, false)
	return nil
}
