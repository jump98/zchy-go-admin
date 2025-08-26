package mongosvr

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//形变数据（心跳）

const mongoCollectionDeformationPoint = "deformationpoint"

type DeformationPointData struct {
	SvrTime     time.Time
	RadarID     int64
	TimeStamp   time.Time //时间戳
	Index       int       //下标
	Deformation float32   //形变值(毫米)
	Distance    float32   //距离值(毫米)
}

func InsertDeformationPointData(data *DeformationPointData) error {
	data.SvrTime = time.Now()
	return insertDocumentData(mongoUri, mongoRadarDBName,
		mongoCollectionDeformationPoint, data)
}

func InsertArrayDeformationPointData(data []interface{}) error {
	return insertArrayDocumentData(mongoUri,
		mongoRadarDBName, mongoCollectionDeformationPoint,
		data)
}

// SampleDeformationDataFunc 采样函数，当记录数量大于100时，采样为100条记录，尽量保留波峰和波谷
func SampleDeformationDataFunc(data []DeformationPointData) []DeformationPointData {
	return sampleDeformationData(data)
}

// sampleDeformationData 采样函数，当记录数量大于100时，采样为100条记录，尽量保留波峰和波谷
func sampleDeformationData(data []DeformationPointData) []DeformationPointData {
	if len(data) <= 100 {
		// 如果记录数量不超过100，直接返回所有数据
		return data
	}

	// 如果记录数量超过100，进行采样
	result := make([]DeformationPointData, 0, 100)
	used := make([]bool, len(data)) // 标记哪些数据点已被使用

	// 首先找到波峰和波谷
	peaksAndValleys := findPeaksAndValleys(data)

	// 将波峰和波谷添加到结果中
	for _, idx := range peaksAndValleys {
		if len(result) < 100 {
			result = append(result, data[idx])
			used[idx] = true
		}
	}

	// 如果还没有达到100个点，用等间距的点填充
	needed := 100 - len(result)
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
	if len(result) > 100 {
		result = result[:100]
	}

	return result
}

// findPeaksAndValleys 查找波峰和波谷的位置
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

func QueryDeformationPointData(devid int64, index int, startTimeStr, endTimeStr string) ([]DeformationPointData, error) {
	// 解析时间字符串
	startTime, err := time.Parse("2006-01-02 15:04:05", startTimeStr)
	if err != nil {
		return nil, fmt.Errorf("开始时间格式错误: %v", err)
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		return nil, fmt.Errorf("结束时间格式错误: %v", err)
	}

	// 查询某个时间段的数据
	filter := bson.M{
		"radarid": devid,
		"index":   index,
		"svrtime": bson.M{
			"$gte": startTime, // 大于等于开始时间
			"$lte": endTime,   // 小于等于结束时间
		},
	}

	// 按时间字段排序
	opts := options.Find().SetSort(bson.M{"svrtime": -1}) // 1表示升序，-1表示降序

	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionDeformationPoint)
	// 执行查询
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// 解析查询结果
	var results []DeformationPointData
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	// 对查询结果进行采样处理
	return sampleDeformationData(results), nil
}

func queryDeformationPointDataTest() {
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

	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionDeformationPoint)
	// 执行查询
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	// 解析查询结果
	var results []DeformationPointData
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	// 对查询结果进行采样处理
	sampledData := sampleDeformationData(results)

	// 打印采样后的结果
	fmt.Println("Sampled deformation data:")
	for i, data := range sampledData {
		fmt.Printf("Index: %d, Deformation: %.2f\n", i, data.Deformation)
	}

	// 打印查询结果
	fmt.Println("Documents within the specified time range:")
	for _, result := range results {
		// 转换为北京时间
		fmt.Printf("ID: %d, Content: %d, CreatedAt: %s\n", result.RadarID, result.Index, result.SvrTime.In(time.FixedZone("CST", 8*3600)))
	}
}
