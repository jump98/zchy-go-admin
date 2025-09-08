package models

import (
	m "go-admin/app/admin/models"
)

type RadarPoint struct {
	m.RadarPoint
}

func (mm RadarPoint) TableName() string {
	return m.RadarPoint{}.TableName()
}
