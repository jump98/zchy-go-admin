package collections

import "time"

// 形变数据（心跳）
const Table_Deformation_Point = "deformation_point"

type DeformationPointModel struct {
	SvrTime     time.Time //服务器时间
	RadarId     int64     //雷达ID
	TimeStamp   time.Time //时间戳(设备时间)
	Index       int       //下标
	Deformation int       //形变值(毫米) 已乘100
	Distance    int       //距离值(毫米) 已乘100
}
