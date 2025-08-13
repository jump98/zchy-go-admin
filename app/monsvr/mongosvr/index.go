package mongosvr

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
)

const mongoUri = "mongodb://localhost:27017"
const mongoRadarDBName = "radardata"

func Init() {
	// 连接到MongoDB
	err := initMongoDB(mongoUri)
	for err != nil {
		fmt.Println("mongodb reconnect after 10 seconds...")
		time.Sleep(time.Second * 10)
		err = initMongoDB(mongoUri)
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
	clientOptions := options.Client().ApplyURI(uri)
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	InitDistance()
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
