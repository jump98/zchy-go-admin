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

type deformationPointService struct {
}

var DeformationPointService = deformationPointService{}

// InsertArrayDeformationPointData 插入多条形变数据
func (s deformationPointService) InsertArrayDeformationPointData(data []interface{}) error {
	return insertArrayDocumentData(mongoUri, mongoRadarDBName, collections.TableDeformationPoint, data)
}

// QueryDeformationPointData 根据时间范围查询距离像形变数据列表
func (s deformationPointService) QueryDeformationPointData(ctx context.Context, radarId, pointIndex int64, startTimeStr, endTimeStr string) ([]collections.DeformationPointModel, error) {
	var err error
	var startTime, endTime time.Time

	// 使用本地时区解析时间
	loc, _ := time.LoadLocation("Local") // 或 "Asia/Shanghai"
	startTime, err = time.ParseInLocation("2006-01-02 15:04:05", startTimeStr, loc)
	if err != nil {
		return nil, fmt.Errorf("开始时间格式错误: %v", err)
	}
	endTime, err = time.ParseInLocation("2006-01-02 15:04:05", endTimeStr, loc)
	if err != nil {
		return nil, fmt.Errorf("结束时间格式错误: %v", err)
	}

	// 打印本地时间
	fmt.Println("查询开始时间:", startTime.Format("2006-01-02 15:04:05"))
	fmt.Println("查询结束时间:", endTime.Format("2006-01-02 15:04:05"))

	// 构建查询条件，MongoDB 内部存储 UTC 时间
	filter := bson.M{
		"radarid": radarId,
		"index":   pointIndex,
		"svrtime": bson.M{
			"$gt":  startTime.UTC(), //大于
			"$lte": endTime.UTC(),   //小于等于
		},
	}

	// 按时间字段排序
	opts := options.Find().SetSort(bson.M{"svrtime": 1}) // 1升序，-1降序
	// 执行查询
	var cursor *mongo.Cursor
	if cursor, err = MDB.Collection(collections.TableDeformationPoint).Find(ctx, filter, opts); err != nil {
		return nil, err
	}
	//defer cursor.Close(ctx)
	// 解析查询结果
	var data []collections.DeformationPointModel
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// QueryDeformationPointByTime 根据时间范围查询距离像形变数据列表
func (s deformationPointService) QueryDeformationPointByTime(ctx context.Context, startTime, endTime time.Time) ([]collections.DeformationPointModel, error) {
	var err error
	fmt.Println("查询开始时间:", startTime.Format("2006-01-02 15:04:05"))
	fmt.Println("查询结束时间:", endTime.Format("2006-01-02 15:04:05"))

	// 构建查询条件，MongoDB 内部存储 UTC 时间
	filter := bson.M{
		"svrtime": bson.M{
			"$gt":  startTime.UTC(), //大于
			"$lte": endTime.UTC(),   //小于等于
		},
	}
	// 按时间字段排序
	opts := options.Find().SetSort(bson.M{"svrtime": 1}) // 1升序，-1降序
	// 执行查询
	var cursor *mongo.Cursor
	if cursor, err = MDB.Collection(collections.TableDeformationPoint).Find(ctx, filter, opts); err != nil {
		return nil, err
	}
	//defer cursor.Close(ctx)
	// 解析查询结果
	var data []collections.DeformationPointModel
	if err = cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}
