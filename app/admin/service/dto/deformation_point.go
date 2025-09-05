package dto

import (
	"go-admin/app/monsvr/mongosvr"
	"time"
)

// 变形点数据查询参数
type DeformationPointQueryReq struct {
	Devid     int64  `json:"devid" validate:"required"`     // 设备ID
	Index     int    `json:"index" validate:"required"`     // 索引
	StartTime string `json:"startTime" validate:"required"` // 开始时间 (格式: 2006-01-02 15:04:05)
	EndTime   string `json:"endTime" validate:"required"`   // 结束时间 (格式: 2006-01-02 15:04:05)
	Hours     int64  `json:"hours"`                         // 查询最近几小时（单位：小时）
	TimeType  string `json:"timeType"`                      // 时间类型（seconds,minutes,hours,days）
	// LastTime  string `json:"lastTime"`                      // 追加查询时，需要最后的一个
}

// 变形点数据查询结果
type DeformationPointQueryResp struct {
	LastTime time.Time                       `json:"last_time"` //最后一条数据的时间
	List     []mongosvr.DeformationPointData `json:"list"`      //形变数据
}
