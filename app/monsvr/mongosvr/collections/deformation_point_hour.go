package collections

import "time"

// TableDeformationPointHour 形变数据-时
const TableDeformationPointHour = "deformation_point_hour"

type DeformationPointHourModel struct {
	Time        time.Time //时间(精确到分)
	RadarId     int64     //雷达ID
	PointIndex  int64     //监测点下标
	Deformation int64     //形变值(毫米) 已乘100
	Distance    int64     //距离值(毫米) 已乘100
}
