package mongosvr

import (
	"context"
	"fmt"
	"go-admin/app/monsvr/mongosvr/collections"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type deformationPointMinuteService struct {
}

var DeformationPointMinuteService = deformationPointMinuteService{}

// CreateIndex 创建索引
func (deformationPointMinuteService) CreateIndex() error {
	//if err := ensureConnection(mongoUri); err != nil {
	//	return fmt.Errorf("failed to ensure connection: %v", err)
	//}

	fmt.Println("准备创建索引")
	collection := client.Database(mongoRadarDBName).Collection(collections.TableDeformationPointMinute)
	// 初始化时创建索引
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "radarid", Value: 1},
			{Key: "pointindex", Value: 1},
			{Key: "time", Value: -1},
		},
		Options: options.Index().SetName("radarid_pointindex_time_idx").SetUnique(true), // 按需决定是否唯一
	}

	var err error
	if _, err = collection.Indexes().CreateOne(context.TODO(), indexModel); err != nil {
		return fmt.Errorf("创建索引失败: %v", err)
	}
	return err
}

// SaveDeformMinuteByTime 根据时间颗粒度存储形变数据
func (s deformationPointMinuteService) SaveDeformMinuteByTime(ctx context.Context, startTime, endTime time.Time) error {
	var err error
	var data []collections.DeformationPointModel
	if data, err = DeformationPointService.QueryDeformationPointByTime(ctx, startTime, endTime); err != nil {
		return err
	}
	if len(data) == 0 {
		fmt.Printf("[%s] 这一分钟没有数据\n", startTime.Format("2006-01-02 15:04"))
		return nil
	}

	// map 按监测点分组
	lastMap := make(map[string]collections.DeformationPointModel)
	for _, d := range data {
		key := fmt.Sprintf("%d_%d", d.RadarId, d.Index)
		if old, ok := lastMap[key]; !ok || d.SvrTime.After(old.SvrTime) {
			lastMap[key] = d
		}
	}

	// 存入分表
	for _, last := range lastMap {
		dayDoc := collections.DeformationPointMinuteModel{
			Time:        startTime,
			RadarId:     last.RadarId,
			PointIndex:  last.Index,
			Deformation: last.Deformation,
		}
		_, err = MDB.Collection(collections.TableDeformationPointMinute).InsertOne(ctx, dayDoc)
		if err != nil {
			return fmt.Errorf("保存形变分钟表数据失败: %v", err)
		}
	}

	return nil
}

// FindDeformMinuteByTime 根据时间范围查询距离像形变数据列表
func (s deformationPointMinuteService) FindDeformMinuteByTime(ctx context.Context, radarId, pointIndex int64, startTime, endTime time.Time) ([]collections.DeformationPointMinuteModel, error) {
	var err error
	// 构建查询条件，MongoDB 内部存储 UTC 时间
	var filter bson.M
	if radarId != 0 && pointIndex != 0 {
		filter = bson.M{
			"radarid":    radarId,
			"pointindex": pointIndex,
			"time": bson.M{
				"$gt":  startTime.UTC(), //大于
				"$lte": endTime.UTC(),   //小于等于
			},
		}
	} else {
		filter = bson.M{
			"time": bson.M{
				"$gt":  startTime.UTC(), //大于
				"$lte": endTime.UTC(),   //小于等于
			},
		}
	}
	// 按时间字段排序
	opts := options.Find().SetSort(bson.M{"time": 1}) // 1升序，-1降序
	var cursor *mongo.Cursor
	if cursor, err = MDB.Collection(collections.TableDeformationPointMinute).Find(ctx, filter, opts); err != nil {
		return nil, err
	}
	var data []collections.DeformationPointMinuteModel
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s deformationPointMinuteService) FindDeformMinuteByTimeV2(ctx context.Context, radarId int64, pointIndex []int64, startTime, endTime time.Time) ([]collections.DeformationPointMinuteModel, error) {
	var err error
	// 构建查询条件，MongoDB 内部存储 UTC 时间
	var filter bson.M
	filter = bson.M{
		"radarid": radarId,
		"pointindex": bson.M{
			"$in": pointIndex, // 多个值
		},
		"time": bson.M{
			"$gt":  startTime.UTC(), //大于
			"$lte": endTime.UTC(),   //小于等于
		},
	}

	// 按时间字段排序
	opts := options.Find().SetSort(bson.M{"time": 1}) // 1升序，-1降序
	var cursor *mongo.Cursor
	if cursor, err = MDB.Collection(collections.TableDeformationPointMinute).Find(ctx, filter, opts); err != nil {
		return nil, err
	}
	var data []collections.DeformationPointMinuteModel
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}
