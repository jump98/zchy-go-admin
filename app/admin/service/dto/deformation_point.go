package dto

// DeformationPointQueryReq 变形点数据查询参数
type DeformationPointQueryReq struct {
	Devid     int64  `json:"devid" validate:"required"`     // 设备ID
	Index     int    `json:"index" validate:"required"`     // 索引
	StartTime string `json:"startTime" validate:"required"` // 开始时间 (格式: 2006-01-02 15:04:05)
	EndTime   string `json:"endTime" validate:"required"`   // 结束时间 (格式: 2006-01-02 15:04:05)
}