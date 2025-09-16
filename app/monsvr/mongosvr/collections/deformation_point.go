package collections

import "time"

// TableDeformationPoint 形变数据（心跳）
const TableDeformationPoint = "deformation_point"

type DeformationPointModel struct {
	SvrTime     time.Time //服务器时间
	RadarId     int64     //雷达ID
	TimeStamp   time.Time //时间戳(设备时间)
	Index       int64     //监测点下标
	Deformation int64     //形变值(毫米) 已乘100
	Distance    int64     //距离值(毫米) 已乘100
}
