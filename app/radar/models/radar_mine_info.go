package models

import (
	"time"
)

// RadarSideInfo 边坡基本信息表
type RadarSideInfo struct {
	Id        int64     `json:"id"          gorm:"primaryKey;autoIncrement;comment:主键编码"`
	DeptId    int64     `json:"deptId"      gorm:"uniqueIndex:idx_rule_key; comment:机构Id"` //机构ID
	SideName  string    `json:"SideName"    gorm:"comment:隐患点名称"`
	SideType  string    `json:"sideType"    gorm:"comment:隐患点类型"`
	Address   string    `json:"address"     gorm:"comment:地址"`
	CreatedAt time.Time `json:"createdAt"   gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updatedAt"   gorm:"comment:最后更新时间"`
}

func (*RadarSideInfo) TableName() string {
	return "radar_side_info"
}
