package collections

import "time"

// TableDeformationPointMinute 形变数据-分
const TableDeformationPointMinute = "deformation_point_minute"

type DeformationPointMinuteModel struct {
	Time        time.Time //时间(精确到分)
	RadarId     int64     //雷达ID
	PointIndex  int64     //监测点下标
	Deformation int64     //形变值(厘米) 已乘100
}
