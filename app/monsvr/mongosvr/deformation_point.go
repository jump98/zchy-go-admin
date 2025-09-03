package mongosvr

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 形变数据（心跳）
const Table_Deformation_Point = "deformation_point"

type DeformationPointData struct {
	SvrTime     time.Time //服务器时间
	RadarId     int64     //雷达ID
	TimeStamp   time.Time //时间戳(设备时间)
	Index       int       //下标
	Deformation int       //形变值(毫米) 已乘100
	Distance    int       //距离值(毫米) 已乘100
}

// func InsertDeformationPointData(data *DeformationPointData) error {
// 	data.SvrTime = time.Now()
// 	return insertDocumentData(mongoUri, mongoRadarDBName, Table_Deformation_Point, data)
// }

func InsertArrayDeformationPointData(data []interface{}) error {
	return insertArrayDocumentData(mongoUri, mongoRadarDBName, Table_Deformation_Point, data)
}

// 根据时间范围查询距离像形变数据列表
func QueryDeformationPointData(devid int64, index int, startTimeStr, endTimeStr string, hours int64) ([]DeformationPointData, time.Time, error) {
	var err error
	var startTime, endTime time.Time
	var lastTime time.Time // 最后一条数据的时间

	// 使用本地时区解析时间
	loc, _ := time.LoadLocation("Local") // 或 "Asia/Shanghai"
	startTime, err = time.ParseInLocation("2006-01-02 15:04:05", startTimeStr, loc)
	if err != nil {
		return nil, lastTime, fmt.Errorf("开始时间格式错误: %v", err)
	}

	endTime, err = time.ParseInLocation("2006-01-02 15:04:05", endTimeStr, loc)
	if err != nil {
		return nil, lastTime, fmt.Errorf("结束时间格式错误: %v", err)
	}

	// 打印本地时间
	fmt.Println("查询开始时间:", startTime.Format("2006-01-02 15:04:05"))
	fmt.Println("查询结束时间:", endTime.Format("2006-01-02 15:04:05"))

	// 构建查询条件，MongoDB 内部存储 UTC 时间
	filter := bson.M{
		"radarid": devid,
		"index":   index,
		"svrtime": bson.M{
			"$gt":  startTime.UTC(), //大于
			"$lte": endTime.UTC(),   //小于等于
		},
	}

	// 按时间字段排序
	opts := options.Find().SetSort(bson.M{"svrtime": 1}) // 1升序，-1降序

	// 执行查询
	cursor, err := MDB.Collection(Table_Deformation_Point).Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, lastTime, err
	}
	defer cursor.Close(context.TODO())

	// 解析查询结果
	var data []DeformationPointData
	if err = cursor.All(context.TODO(), &data); err != nil {
		return nil, lastTime, err
	}

	// 获取最后一条数据的时间（本地时间）
	if len(data) > 0 {
		lastTime = data[len(data)-1].SvrTime.Local()
	}

	maxPoints := getMaxPointsForRange(hours)
	// 采样
	sampledData := sampleDeformationData(data, maxPoints)

	fmt.Printf("原始数据条数: %d\n", len(data))
	fmt.Printf("采样后数据条数: %d\n", len(sampledData))

	return sampledData, lastTime, nil
}

// 采样函数：保留极值 + 均匀抽样
// 严格控制采样结果不超过 maxPoints
func sampleDeformationData(data []DeformationPointData, maxPoints int) []DeformationPointData {
	if len(data) <= maxPoints {
		return data
	}

	step := len(data) / maxPoints
	if step < 1 {
		step = 1
	}

	sampled := make([]DeformationPointData, 0, maxPoints)
	for i := 0; i < len(data); i += step {
		end := i + step
		if end > len(data) {
			end = len(data)
		}

		// 在区间内找极值
		maxVal := data[i]
		minVal := data[i]
		for j := i; j < end; j++ {
			if data[j].Deformation > maxVal.Deformation {
				maxVal = data[j]
			}
			if data[j].Deformation < minVal.Deformation {
				minVal = data[j]
			}
		}

		// 选波动幅度较大的点（离均值更远）
		mid := (data[i].Deformation + data[end-1].Deformation) / 2
		if abs(maxVal.Deformation-mid) > abs(minVal.Deformation-mid) {
			sampled = append(sampled, maxVal)
		} else {
			sampled = append(sampled, minVal)
		}

		// 控制上限
		if len(sampled) >= maxPoints {
			break
		}
	}

	return sampled
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// 根据时间范围计算最大采样点数
func getMaxPointsForRange(hours int64) int {
	// 用 map 定义时间范围对应的最大采样点数
	var maxPointsMap = map[int64]int{
		1:       100,
		3:       200,
		6:       250,
		12:      300,
		24:      400,
		3 * 24:  500,
		7 * 24:  500,
		30 * 24: 600,
		90 * 24: 600,
	}
	if val, ok := maxPointsMap[hours]; ok {
		return val
	}
	return 500 // 默认值

}
