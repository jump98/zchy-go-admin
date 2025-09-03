package mongosvr

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		ds = append(ds, DeformationPointData{
			SvrTime:     time.Now(),
			RadarId:     data.RadarId,
			TimeStamp:   data.TimeStamp,
			Index:       v.Index,
			Deformation: v.Deformation,
			Distance:    v.Distance,
		})
	}
	if e == nil {
		e = InsertArrayDeformationPointData(ds)
	}

	//插入监测点的形变数据
	return e
}

func queryDeformationDataTest() {
	// 查询某个时间段的数据
	startTime := time.Now().Add(-24 * time.Hour) // 24小时前
	endTime := time.Now()                        // 当前时间

	filter := bson.M{
		"devid": 1,
		"svrtime": bson.M{
			"$gte": startTime, // 大于等于开始时间
			"$lte": endTime,   // 小于等于结束时间
		},
	}

	// 按时间字段排序
	opts := options.Find().SetSort(bson.M{"svrtime": 1}) // 1表示升序，-1表示降序

	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionDeformation)
	// 执行查询
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	// 解析查询结果
	var results []DeformationData
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	// 打印查询结果
	fmt.Println("Documents within the specified time range:")
	for _, result := range results {
		// 转换为北京时间
		fmt.Printf("ID: %s, Interval: %d, CreatedAt: %s\n", result.RadarKey, result.Interval, result.SvrTime.In(time.FixedZone("CST", 8*3600)))
	}
}
