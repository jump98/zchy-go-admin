package mongosvr

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoCollectionDevInfo = "radardevinfo"

// RadarStatus 设备状态请求
type RadarDevInfo struct {
	SvrTime time.Time
	RadarId int64
	RadarDevInfoRequest
}

type RadarDevInfoRequest struct {
	RadarKey     string `json:"radarkey"`
	Timestamp    int64  `json:"timestamp"`
	CommandCode  int    `json:"command_code"`
	FpgaVersion  string `json:"fpga_version"`
	PsVersion    string `json:"ps_version"`
	ImageVersion string `json:"image_version"`
	SerialNumber string `json:"serial_number"`
	SimIMEI      string `json:"sim_IMEI"`
}

func InitRadarDevInfo() error {
	// 确保连接有效
	if err := ensureConnection(mongoUri); err != nil {
		return fmt.Errorf("failed to ensure connection: %v", err)
	}

	// 插入数据
	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionDevInfo)
	// 初始化时创建索引
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "radarid", Value: 1},  // 设备ID升序
			{Key: "svrtime", Value: -1}, // 时间降序
		},
		Options: options.Index().SetName("radarid_svrtime_idx"),
	}
	if _, err := collection.Indexes().CreateOne(context.TODO(), indexModel); err != nil {
		return (fmt.Errorf("创建索引失败: %v", err))
	}
	return nil
}

func InsertRadarDevInfo(data *RadarDevInfo) error {
	data.SvrTime = time.Now()
	return insertDocumentData(mongoUri, mongoRadarDBName, mongoCollectionDevInfo, data)
}

func GetLatestRadarDevInfo(radarid int64) (*RadarDevInfo, error) {
	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionDevInfo)

	// 构建查询条件
	filter := bson.D{
		{Key: "radarid", Value: radarid}, // 指定设备ID
	}
	// 设置查询选项：按时间降序排序，取第一条
	opts := options.FindOne().
		SetSort(bson.D{{Key: "svrtime", Value: -1}}) // -1 表示降序

	var result RadarDevInfo
	err := collection.FindOne(context.TODO(), filter, opts).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("没有找到任何数据")
		}
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	return &result, nil
}

func queryRadarDevInfoDataTest() {
	// 查询某个时间段的数据
	startTime := time.Now().Add(-24 * time.Hour) // 24小时前
	endTime := time.Now()                        // 当前时间

	filter := bson.M{
		"radarid": 1,
		"svrtime": bson.M{
			"$gte": startTime, // 大于等于开始时间
			"$lte": endTime,   // 小于等于结束时间
		},
	}

	// 按时间字段排序
	opts := options.Find().SetSort(bson.M{"svrtime": 1}) // 1表示升序，-1表示降序

	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionDevInfo)
	// 执行查询
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	// 解析查询结果
	var results []DistanceDataV2
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	// 打印查询结果
	fmt.Println("Documents within the specified time range:")
	for _, result := range results {
		// 转换为北京时间
		fmt.Printf("ID: %d, Content: %d, CreatedAt: %s\n", result.RadarID, result.CommandCode, result.SvrTime.In(time.FixedZone("CST", 8*3600)))
	}
}
