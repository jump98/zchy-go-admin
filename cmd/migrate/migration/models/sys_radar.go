package models

import (
	m "go-admin/app/admin/models"
)

type SysRadar struct {
	m.SysRadar
}

func (mm *SysRadar) TableName() string {
	return m.SysRadar{}.TableName()
}
