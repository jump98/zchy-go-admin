package dto

import (
	"time"
)

// 变形点数据查询参数
type GetDeformationDataReq struct {
	RadarId   int64  `json:"radarId"   validate:"required"` // 设备ID
	Index     int    `json:"index"     validate:"required"` // 索引
	StartTime string `json:"startTime" validate:"required"` // 开始时间 (格式: 2006-01-02 15:04:05)
	EndTime   string `json:"endTime"   validate:"required"` // 结束时间 (格式: 2006-01-02 15:04:05)
	Hours     int64  `json:"hours"`                         // 查询最近几小时（单位：小时）
	TimeType  string `json:"timeType"`                      // 时间类型（seconds,minutes,hours,days）
}

// 变形点数据查询结果
type GetDeformationDataResp struct {
	LastTime time.Time             `json:"lastTime"` //最后一条数据的时间
	List     []DeformationDataItem `json:"list"`     //形变数据
}

// 形变数据列表
type DeformationDataItem struct {
	Time           time.Time `json:"time"`           //时间（可能是一个区间范围值）
	DeformationMax int       `json:"deformationMax"` //最大形变值(毫米) 已乘100 (最大形变值)
	DeformationMin int       `json:"deformationMin"` //最小形变值(毫米) 已乘100 (最小形变值)
	DeformationAvg int       `json:"deformationAvg"` //最小形变值(毫米) 已乘100 (平均形变值)
	Distance       int       `json:"distance"`       //距离值(毫米) 已乘100
}

// 变形点速度参数
type GetDeformationVelocityReq struct {
	RadarId   int64  `json:"radarId" validate:"required"`   // 设备ID
	Index     int    `json:"index" validate:"required"`     // 索引
	StartTime string `json:"startTime" validate:"required"` // 开始时间 (格式: 2006-01-02 15:04:05)
	EndTime   string `json:"endTime" validate:"required"`   // 结束时间 (格式: 2006-01-02 15:04:05)
	Hours     int64  `json:"hours"`                         // 查询最近几小时（单位：小时）
	TimeType  string `json:"timeType"`                      // 时间类型（seconds,minutes,hours,days）
	// LastTime  string `json:"lastTime"`                      // 追加查询时，需要最后的一个
}

// 变形点速度返回数据
type GetDeformationVelocityResp struct {
	LastTime time.Time                 `json:"last_time"` //最后一条数据的时间
	List     []DeformationVelocityItem `json:"list"`      //形变速度数据
}

type DeformationVelocityItem struct {
	Time time.Time `json:"time"` //时间
	Avg  float64   `json:"avg"`  //平均速度
	Max  float64   `json:"max"`  //最大速度
	Min  float64   `json:"min"`  //最小速度
}

// 变形点速度查询参数
type GetDeformationAccelerationReq struct {
	RadarId   int64  `json:"radarId" validate:"required"`   // 设备ID
	Index     int    `json:"index" validate:"required"`     // 索引
	StartTime string `json:"startTime" validate:"required"` // 开始时间 (格式: 2006-01-02 15:04:05)
	EndTime   string `json:"endTime" validate:"required"`   // 结束时间 (格式: 2006-01-02 15:04:05)
	Hours     int64  `json:"hours"`                         // 查询最近几小时（单位：小时）
	TimeType  string `json:"timeType"`                      // 时间类型（seconds,minutes,hours,days）
}
type GetDeformationAccelerationResp struct {
	RadarId   int64  `json:"radarId" validate:"required"`   // 设备ID
	Index     int    `json:"index" validate:"required"`     // 索引
	StartTime string `json:"startTime" validate:"required"` // 开始时间 (格式: 2006-01-02 15:04:05)
	EndTime   string `json:"endTime" validate:"required"`   // 结束时间 (格式: 2006-01-02 15:04:05)
	Hours     int64  `json:"hours"`                         // 查询最近几小时（单位：小时）
	TimeType  string `json:"timeType"`                      // 时间类型（seconds,minutes,hours,days）
}
