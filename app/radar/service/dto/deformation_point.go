package dto

import (
	"time"
)

// TimeUnit 用于表示时间单位的类型
type TimeUnit string

// 定义合法的时间类型常量
const (
	TimeUnitSeconds TimeUnit = "seconds"
	TimeUnitMinutes TimeUnit = "minutes"
	TimeUnitHours   TimeUnit = "hours"
	TimeUnitDays    TimeUnit = "days"
)

// 变形点数据查询参数
type GetDeformationDataReq struct {
	RadarId   int64    `json:"radarId"   validate:"required"` // 设备ID
	Index     int64    `json:"index"     validate:"required"` // 索引
	StartTime string   `json:"startTime" validate:"required"` // 开始时间 (格式: 2006-01-02 15:04:05)
	EndTime   string   `json:"endTime"   validate:"required"` // 结束时间 (格式: 2006-01-02 15:04:05)
	Hours     int64    `json:"hours"`                         // 查询最近几小时（单位：小时）
	TimeUnit  TimeUnit `json:"timeUnit"`                      // 时间单位（seconds,minutes,hours,days）
}

// 变形点数据查询结果
type GetDeformationDataResp struct {
	LastTime time.Time             `json:"lastTime"` //最后一条数据的时间
	List     []DeformationDataItem `json:"list"`     //形变数据
}

// 形变数据列表
type DeformationDataItem struct {
	Time           time.Time `json:"time"`           //时间（可能是一个区间范围值）
	DeformationMax int64     `json:"deformationMax"` //最大形变值(毫米) 已乘100 (最大形变值)
	DeformationMin int64     `json:"deformationMin"` //最小形变值(毫米) 已乘100 (最小形变值)
	DeformationAvg int64     `json:"deformationAvg"` //最小形变值(毫米) 已乘100 (平均形变值)
	Distance       int64     `json:"distance"`       //距离值(毫米) 已乘100
}

// 变形点速度参数
// type GetDeformationVelocityReq struct {
// 	RadarId   int64    `json:"radarId" validate:"required"`   // 设备ID
// 	Index     int      `json:"index" validate:"required"`     // 索引
// 	StartTime string   `json:"startTime" validate:"required"` // 开始时间 (格式: 2006-01-02 15:04:05)
// 	EndTime   string   `json:"endTime" validate:"required"`   // 结束时间 (格式: 2006-01-02 15:04:05)
// 	Hours     int64    `json:"hours"`                         // 查询最近几小时（单位：小时）
// 	TimeUnit  TimeUnit `json:"timeUnit"`                      // 时间单位（seconds,minutes,hours,days）
// }

// 变形速度返回数据
type GetDeformationVelocityResp struct {
	LastTime time.Time                 `json:"last_time"` //最后一条数据的时间
	List     []DeformationVelocityItem `json:"list"`      //形变速度数据
	Unit     string                    `json:"unit"`      // "mm/s" "mm/m" "mm/h" "mm/d"
}

type DeformationVelocityItem struct {
	Time time.Time `json:"time"` //时间
	Avg  float64   `json:"avg"`  //平均速度
	Max  float64   `json:"max"`  //最大速度
	Min  float64   `json:"min"`  //最小速度
}

// 变形点速度查询参数
// type GetDeformationAccelerationReq struct {
// 	RadarId   int64    `json:"radarId" validate:"required"`   // 设备ID
// 	Index     int      `json:"index" validate:"required"`     // 索引
// 	StartTime string   `json:"startTime" validate:"required"` // 开始时间 (格式: 2006-01-02 15:04:05)
// 	EndTime   string   `json:"endTime" validate:"required"`   // 结束时间 (格式: 2006-01-02 15:04:05)
// 	Hours     int64    `json:"hours"`                         // 查询最近几小时（单位：小时）
// 	TimeUnit  TimeUnit `json:"timeUnit"`                      // 时间单位（seconds,minutes,hours,days）
// }

// 变形加速度返回数据
type GetDeformationAccelerationResp struct {
	LastTime time.Time                     `json:"last_time"` //最后一条数据的时间
	List     []DeformationAccelerationItem `json:"list"`      //形变速度数据
	Unit     string                        `json:"unit"`      // "mm/s" "mm/m" "mm/h" "mm/d"
}

// 形变加速度数据
type DeformationAccelerationItem struct {
	Time time.Time `json:"time"` //时间
	Avg  float64   `json:"avg"`  //平均速度
	Max  float64   `json:"max"`  //最大速度
	Min  float64   `json:"min"`  //最小速度
}

// GetDeformCurveListReq 形变曲线
type GetDeformCurveListReq struct {
	Kind      int64    `json:"kind"`                          // 数据类型：0=形变、1=速度、2=加速度
	RadarId   int64    `json:"radarId"   validate:"required"` // 设备ID
	Index     []int64  `json:"index"     validate:"required"` // 监测点
	StartTime string   `json:"startTime" validate:"required"` // 开始时间 (格式: 2006-01-02 15:04:05)
	EndTime   string   `json:"endTime"   validate:"required"` // 结束时间 (格式: 2006-01-02 15:04:05)
	Hours     int64    `json:"hours"`                         // 查询最近几小时（单位：小时）
	TimeUnit  TimeUnit `json:"timeUnit"`                      // 时间单位（seconds,minutes,hours,days）
}
type GetDeformCurveListResp struct {
	LastTime time.Time                   `json:"lastTime"` //最后一条数据的时间
	List     map[int64][]DeformCurveItem `json:"list"`     //形变数据
}
type DeformCurveItem struct {
	T time.Time `json:"t"` //时间
	V int64     `json:"v"` //值  （已乘100）
}
