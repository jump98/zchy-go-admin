package mongosvr

import (
	"context"
	"fmt"
	"go-admin/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//mongoDB文档：https://www.mongodb.com/zh-cn/docs/drivers/go/current/crud/insert/

var (
	client *mongo.Client
	MDB    *mongo.Database //mongoDB 客户端
)

var mongoUri string         //mongoDB URL
var mongoRadarDBName string //雷达数据库DB name

func Init() {
	config := config.ExtConfig
	mongoUri = config.MongoDB.Source
	mongoRadarDBName = config.MongoDB.RadarDBName
	// 连接到MongoDB
	if err := initMongoDB(mongoUri); err != nil {
		fmt.Println("连接mangoDB出错：", err)
		panic(err)
	}

	// dis := DistanceData{}
	// dis.DevType = 0x00000100
	// dis.DevID = 0x00000001
	// dis.Cmd = 0x05
	// dis.ExtCmd = 0x00
	// dis.DataLen = 2000 * 4
	// InsertDistanceData(&dis)
	// queryDistanceDataTest()

	// def := DeformationData{}
	// def.DevType = 0x00000100
	// def.DevID = 0x00000001
	// def.Cmd = 0x05
	// def.ExtCmd = 0x00
	// def.DataLen = 2000 * 4
	// def.Option = 1
	// def.TimeStamp = 123
	// def.DetectCount = 1
	// def.TestNum = 2
	// def.FrameNo = 1
	// def.FrameTimes = 2

	// InsertDeformationData(&def)
	// queryDeformationDataTest()

}

// 初始化MongoDB连接
func initMongoDB(uri string) error {

	// loggerOptions := options.
	// 	Logger().
	// 	SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	// 设置客户端连接配置
	clientOpts := options.Client().
		ApplyURI(uri).
		// SetLoggerOptions(loggerOptions).     // 设置日志级别
		SetMaxPoolSize(100).                 // 最大连接数
		SetMinPoolSize(10).                  // 最小保持连接数
		SetMaxConnIdleTime(30 * time.Second) // 连接最大空闲时间

	var err error
	// 连接 MongoDB
	if client, err = mongo.Connect(context.Background(), clientOpts); err != nil {
		fmt.Println("创建mongoDB连接出错:", err)
		return err
	}
	MDB = client.Database(mongoRadarDBName)

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	// InitDistance()
	InitDistanceV2()
	InitCommand()
	InitRadarStatus()
	InitRadarDevInfo()
	InitAlarm()

	fmt.Println("Connected to MongoDB!")
	return nil
}

// 重连MongoDB
func reconnectMongoDB(uri string) error {
	if client != nil {
		_ = client.Disconnect(context.TODO()) // 关闭旧连接
	}
	return initMongoDB(uri)
}

// 检查连接状态并重连
func ensureConnection(uri string) error {
	if client == nil {
		return initMongoDB(uri)
	}

	// 检查连接是否有效
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		fmt.Println("Connection lost, reconnecting...")
		return reconnectMongoDB(uri)
	}

	return nil
}
