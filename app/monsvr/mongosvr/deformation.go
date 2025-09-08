package mongosvr

import (
	"go-admin/app/monsvr/mongosvr/collections"
	"time"
)

// 形变数据（心跳）
const mongoCollectionDeformation = "deformation"

// 形变数据
type DeformationData struct {
	SvrTime   time.Time
	RadarKey  string
	RadarId   int64
	TimeStamp time.Time
	Interval  int
	DefData   []DeformationDefData
}

// 形变数据 详细
type DeformationDefData struct {
	Index       int //下标
	Deformation int //形变值(毫米) *100
	Distance    int //距离值(毫米) *100
}

// 插入形变数据
func InsertDeformationData(data *DeformationData) error {
	data.SvrTime = time.Now()
	e := insertDocumentData(mongoUri, mongoRadarDBName, mongoCollectionDeformation, data)
	ds := []any{}
	for _, v := range data.DefData {
		ds = append(ds, collections.DeformationPointModel{
			SvrTime:     time.Now(),
			RadarId:     data.RadarId,
			TimeStamp:   data.TimeStamp,
			Index:       v.Index,
			Deformation: v.Deformation,
			Distance:    v.Distance,
		})
	}
	if e == nil {
		e = DeformationPointService.InsertArrayDeformationPointData(ds)
	}

	//插入监测点的形变数据
	return e
}
