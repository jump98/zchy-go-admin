package mongosvr

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 告警数据
const mongoCollectionAlarm = "alarms"

type AlarmData struct {
	SvrTime     time.Time `bson:"svrtime"`
	RadarId     int64     `bson:"radarid"`
	TimeStamp   int64     `bson:"timestamp"`
	Voltage     int       `bson:"voltage"`
	Temperature int       `bson:"temperature"`
	Battery     int       `bson:"battery"`
	SolarPanel  int       `bson:"solar_panel"`
	RadarData   int       `bson:"radar_data"`
}

func InitAlarm() error {
	// 确保连接有效
	if err := ensureConnection(mongoUri); err != nil {
		return fmt.Errorf("failed to ensure connection: %v", err)
	}

	// 插入数据
	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionAlarm)
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

func InsertAlarmData(data *AlarmData) error {
	data.SvrTime = time.Now()
	return insertDocumentData(mongoUri, mongoRadarDBName, mongoCollectionAlarm, data)
}

func QueryAlarmsDataTimeBefore(radarID int64, startTime time.Time, num int) ([]AlarmData, error) {
	// 查询某个时间段的数据
	filter := bson.M{
		"radarid": radarID,
		"svrtime": bson.M{
			"$lte": startTime, // 小于等于结束时间
		},
	}

	// 按时间字段排序，按时间降序排列
	opts := options.Find().SetSort(bson.M{"svrtime": -1})

	// 限制返回记录数
	if num > 0 {
		opts.SetLimit(int64(num))
	}

	collection := client.Database(mongoRadarDBName).Collection(mongoCollectionAlarm)
	// 执行查询
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// 解析查询结果
	var results []AlarmData
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

// QueryLastAlarmForRadarIDs 查询每个雷达ID的最后一条告警记录
func QueryLastAlarmForRadarIDs(radarIDs []int64) ([]AlarmData, error) {
	var results []AlarmData

	for _, radarID := range radarIDs {
		filter := bson.M{"radarid": radarID}
		opts := options.Find().SetSort(bson.M{"svrtime": -1}).SetLimit(1)

		collection := client.Database(mongoRadarDBName).Collection(mongoCollectionAlarm)
		cursor, err := collection.Find(context.TODO(), filter, opts)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.TODO())

		var alarms []AlarmData
		if err = cursor.All(context.TODO(), &alarms); err != nil {
			return nil, err
		}

		if len(alarms) > 0 {
			results = append(results, alarms[0])
		}
	}

	return results, nil
}
