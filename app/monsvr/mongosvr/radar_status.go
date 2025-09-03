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

const mongoCollectionStatus = "radar_status"

// RadarStatus 设备状态请求
type RadarStatus struct {
	SvrTime time.Time
	RadarId int64
	RadarStatusRequest
}

type RadarStatusRequest struct {
	RadarKey    string `json:"radarkey"`
	Timestamp   int64  `json:"timestamp"`    //时间戳
	CommandCode int    `json:"command_code"` //命令码
	DiskTotal   uint64 `json:"disk_total"`   //磁盘总容量
	DiskFree    uint64 `json:"disk_free"`    //磁盘剩余容量
	RamTotal    uint64 `json:"ram_total"`    //内存总容量
	RamFree     uint64 `json:"ram_free"`     //内存剩余容量
	SimState    int    `json:"sim_state"`    //SIM卡状态 :0=正常 1异常
	SimRSSI     int    `json:"sim_RSSI"`     //SIM接收信号强度 单位（dBm）
	Battery     int    `json:"battery"`      //电池状态 0 充电中 1 放电中
	Voltage     string `json:"voltage"`      //电压：   {"12V":11.686076, "5V3":5.285750, "2V1":2.123250}V
	Current     string `json:"current"`      //供电电流：    {\"12V\":0.888750}"  {电流名：电流值(单位安)}
	Temperature string `json:"temperature"`  //设备温度：    {\"local\":44.625000,\"PCB\":44.375000,\"ZYNQ\":49.500000}" //设备温度，可有多个值 键 温度名 值 温度值(单位摄氏度)
}

// 12V电流：:0.888750A

// 12V电压: 11.686076V
// 5.3V电压: 11.686076V
// 2.1V电压: 11.686076V

// 设备温度
// 主板温度：xxx
// 设备外壳温度：xxx
// 处理器温度：xxx

func InitRadarStatus() error {
	// 确保连接有效
	if err := ensureConnection(mongoUri); err != nil {
		return fmt.Errorf("failed to ensure connection: %v", err)
	}

	// 插入数据
	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionStatus)
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

func InsertRadarStatus(data *RadarStatus) error {
	data.SvrTime = time.Now()
	return insertDocumentData(mongoUri, mongoRadarDBName, mongoCollectionStatus, data)
}

func GetLatestRadarStatus(radarid int64) (*RadarStatus, error) {
	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionStatus)

	// 构建查询条件
	filter := bson.D{
		{Key: "radarid", Value: radarid}, // 指定设备ID
	}
	// 设置查询选项：按时间降序排序，取第一条
	opts := options.FindOne().
		SetSort(bson.D{{Key: "svrtime", Value: -1}}) // -1 表示降序

	var result RadarStatus
	err := collection.FindOne(context.TODO(), filter, opts).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("没有找到任何数据")
		}
		return nil, fmt.Errorf("查询失败: %v", err)
	}

	return &result, nil
}

func queryRadarStatusDataTest() {
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

	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionStatus)
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
