package models

import (
	m "go-admin/app/admin/models"
)

type RadarPoint struct {
	m.RadarPoint
}

func (*RadarPoint) TableName() string {
	return "radar_point"
}
