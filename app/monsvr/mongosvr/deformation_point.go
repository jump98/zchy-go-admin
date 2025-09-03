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
	Deformation float32   //形变值(毫米)
	Distance    float32   //距离值(毫米)
}

func InsertDeformationPointData(data *DeformationPointData) error {
	data.SvrTime = time.Now()
	return insertDocumentData(mongoUri, mongoRadarDBName, Table_Deformation_Point, data)
}

func InsertArrayDeformationPointData(data []interface{}) error {
	return insertArrayDocumentData(mongoUri, mongoRadarDBName, Table_Deformation_Point, data)
}

// 采样函数，当记录数量大于100时，采样为100条记录，尽量保留波峰和波谷
func SampleDeformationDataFunc(data []DeformationPointData) []DeformationPointData {
	return sampleDeformationData(data)
}

// 采样函数，当记录数量大于100时，采样为100条记录，尽量保留波峰和波谷
func sampleDeformationData(data []DeformationPointData) []DeformationPointData {
	const MAX = 1000
	// 如果记录数量不超过100，直接返回所有数据
	if len(data) <= MAX {
		return data
	}

	// 如果记录数量超过100，进行采样
	result := make([]DeformationPointData, 0, MAX)
	used := make([]bool, len(data)) // 标记哪些数据点已被使用

	// 首先找到波峰和波谷
	peaksAndValleys := findPeaksAndValleys(data)

	// 将波峰和波谷添加到结果中
	for _, idx := range peaksAndValleys {
		if len(result) < MAX {
			result = append(result, data[idx])
			used[idx] = true
		}
	}

	// 如果还没有达到100个点，用等间距的点填充
	needed := MAX - len(result)
	if needed > 0 {
		step := float64(len(data)-1) / float64(needed)
		for i := 0; i < needed; i++ {
			index := int(float64(i) * step)
			if index >= len(data) {
				index = len(data) - 1
			}

			// 如果这个点已经被使用，寻找最近的未使用点
			if used[index] {
				// 向前搜索
				found := false
				for j := index; j >= 0; j-- {
					if !used[j] {
						result = append(result, data[j])
						used[j] = true
						found = true
						break
					}
				}

				// 向后搜索
				if !found {
					for j := index + 1; j < len(data); j++ {
						if !used[j] {
							result = append(result, data[j])
							used[j] = true
							break
						}
					}
				}
			} else {
				result = append(result, data[index])
				used[index] = true
			}
		}
	}

	// 如果仍然超过100个点，截取前100个
	if len(result) > MAX {
		result = result[:MAX]
	}

	return result
}

// 查找波峰和波谷的位置
func findPeaksAndValleys(data []DeformationPointData) []int {
	indices := make([]int, 0)

	for i := 1; i < len(data)-1; i++ {
		// 检查是否为波峰
		if data[i].Deformation > data[i-1].Deformation && data[i].Deformation > data[i+1].Deformation {
			indices = append(indices, i)
		}
		// 检查是否为波谷
		if data[i].Deformation < data[i-1].Deformation && data[i].Deformation < data[i+1].Deformation {
			indices = append(indices, i)
		}
	}

	return indices
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
	var results []DeformationPointData
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, lastTime, err
	}

	fmt.Println("查询的记录总数:", len(results))

	// 获取最后一条数据的时间（本地时间）
	if len(results) > 0 {
		lastTime = results[len(results)-1].SvrTime.Local()
	}

	// if hours <= 6 {
	// 	return results, lastTime, nil
	// }
	// 对查询结果进行采样处理
	sampleData := sampleDeformationData(results)
	return sampleData, lastTime, nil
}

// func queryDeformationPointDataTest() {
// 	// 查询某个时间段的数据
// 	startTime := time.Now().Add(-24 * time.Hour) // 24小时前
// 	endTime := time.Now()                        // 当前时间

// 	filter := bson.M{
// 		"devid": 1,
// 		"svrtime": bson.M{
// 			"$gte": startTime, // 大于等于开始时间
// 			"$lte": endTime,   // 小于等于结束时间
// 		},
// 	}

// 	// 按时间字段排序
// 	opts := options.Find().SetSort(bson.M{"svrtime": 1}) // 1表示升序，-1表示降序

// 	collection := client.Database(mongoRadarDBName).Collection(Table_Deformation_Point)
// 	// 执行查询
// 	cursor, err := collection.Find(context.TODO(), filter, opts)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cursor.Close(context.TODO())

// 	// 解析查询结果
// 	var results []DeformationPointData
// 	if err = cursor.All(context.TODO(), &results); err != nil {
// 		log.Fatal(err)
// 	}

// 	// 对查询结果进行采样处理
// 	sampledData := sampleDeformationData(results)

// 	// 打印采样后的结果
// 	fmt.Println("Sampled deformation data:")
// 	for i, data := range sampledData {
// 		fmt.Printf("Index: %d, Deformation: %.2f\n", i, data.Deformation)
// 	}

// 	// 打印查询结果
// 	fmt.Println("Documents within the specified time range:")
// 	for _, result := range results {
// 		// 转换为北京时间
// 		fmt.Printf("ID: %d, Content: %d, CreatedAt: %s\n", result.RadarId, result.Index, result.SvrTime.In(time.FixedZone("CST", 8*3600)))
// 	}
// }
