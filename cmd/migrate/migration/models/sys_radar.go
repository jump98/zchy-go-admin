package models

import (
	m "go-admin/app/radar/models"
)

type SysRadar struct {
	m.Radar
}

func (mm *SysRadar) TableName() string {
	return m.Radar{}.TableName()
}
