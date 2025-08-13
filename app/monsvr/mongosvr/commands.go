package mongosvr

import (
	"context"
	"fmt"
	"time"

	mongodto "go-admin/app/monsvr/mongosvr/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//形变数据（心跳）

const mongoCollectionCommand = "commands"

const (
	CMD_RD_REBOOT       = 100
	CMD_RD_GETSTATEINFO = 101
	CMD_RD_GETDEVINFO   = 102
	CMD_RD_GETRAWDATA   = 300
	CMD_RD_ADDPOINT     = 400
	CMD_RD_DELETEPOINT  = 401
)

type CommandData struct {
	SvrTime     time.Time
	RadarId     int64
	TimeStamp   int64
	CommandCode int
	Message     string
	Parameters  map[string]interface{}
	Send        int //默认0表示 没有发出去，1表示已经发送出去了
}

func InitCommand() error {
	// 确保连接有效
	if err := ensureConnection(mongoUri); err != nil {
		return fmt.Errorf("failed to ensure connection: %v", err)
	}

	// 插入数据
	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionCommand)
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

func InsertCommandData(data *CommandData) error {
	// 先根据RadarId和CommandCode查询数据，如果存在则不插入
	filter := bson.M{
		"radarid":     data.RadarId,
		"commandcode": data.CommandCode,
	}

	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionCommand)
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("查询数据失败: %v", err)
	}

	// 如果数据已存在，则不插入
	if count > 0 {
		return nil
	}

	data.SvrTime = time.Now()
	data.Send = 0
	e := insertDocumentData(mongoUri, mongoRadarDBName,
		mongoCollectionCommand, data)

	return e
}

func QueryRadarComandData(radarId int64) ([]mongodto.CommandDataDto, error) {
	// 查询某个时间段的数据
	//startTime := time.Now().Add(-24 * time.Hour) // 24小时前
	//endTime := time.Now()                        // 当前时间

	filter := bson.M{
		"radarid": radarId,
		"send":    0,
		// "svrtime": bson.M{
		// 	"$gte": startTime, // 大于等于开始时间
		// 	//"$lte": endTime,   // 小于等于结束时间
		// },
	}

	// 按时间字段排序
	opts := options.Find().SetSort(bson.M{"svrtime": 1}) // 1表示升序，-1表示降序

	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionCommand)
	// 执行查询
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// 解析查询结果
	var results []CommandData
	if err = cursor.All(context.TODO(), &results); err != nil {
		//log.Fatal(err)
		return nil, err
	}

	reses := make([]mongodto.CommandDataDto, 0)
	for _, v := range results {
		reses = append(reses, mongodto.CommandDataDto{
			CommandCode: v.CommandCode,
			Message:     v.Message,
			Parameters:  v.Parameters,
		})
	}
	// 批量更新send字段为1
	// update := bson.M{"$set": bson.M{"send": 1}}
	// _, err = collection.UpdateMany(context.TODO(), filter, update)
	// if err != nil {
	// 	return nil, err
	// }
	//删除旧命令
	_, err = collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return reses, nil
}
