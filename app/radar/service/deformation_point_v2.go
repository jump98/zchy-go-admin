package service

import (
	"context"
	"errors"
	"go-admin/app/monsvr/mongosvr"
	"go-admin/app/monsvr/mongosvr/collections"
	"go-admin/app/radar/service/dto"
	"sort"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/service"
)

type DeformationPointV2 struct {
	service.Service
}

// GetDeformCurveList 获取形变数据列表
func (s DeformationPointV2) GetDeformCurveList(ctx context.Context, req dto.GetDeformCurveListReq) (*dto.GetDeformCurveListResp, error) {
	var err error
	//hours := req.Hours
	//timeUnit := req.TimeUnit
	radarId := req.RadarId
	index := req.Index

	var start, end time.Time
	loc, _ := time.LoadLocation("Local") // 或 "Asia/Shanghai"
	if start, err = time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, loc); err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	if end, err = time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, loc); err != nil {
		return nil, errors.New("结束时间格式错误")
	}

	// 获取数据
	var data []collections.DeformationPointMinuteModel
	if data, err = mongosvr.DeformationPointMinuteService.FindDeformMinuteByTimeV2(ctx, radarId, index, start, end); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("暂无形变数据")
	}
	indexMap := map[int64][]dto.DeformCurveItem{}
	for _, v := range data {
		item := dto.DeformCurveItem{
			T: v.Time,
			V: v.Deformation,
		}
		indexMap[v.PointIndex] = append(indexMap[v.PointIndex], item)
	}

	switch req.Kind {
	case 1:
		//速度
		if indexMap, err = s.calcDeformVelocity(indexMap); err != nil {
			return nil, err
		}
	case 2:
		//加速度
		if indexMap, err = s.calcDeformAcceleration(indexMap); err != nil {
			return nil, err
		}
	default:
		//形变
	}

	resp := &dto.GetDeformCurveListResp{
		LastTime: data[len(data)-1].Time.Local(),
		List:     indexMap,
	}
	return resp, nil
}

// 计算形变速度
func (s DeformationPointV2) calcDeformVelocity(indexMap map[int64][]dto.DeformCurveItem) (map[int64][]dto.DeformCurveItem, error) {
	stats := make(map[int64][]dto.DeformCurveItem)

	for pointID, items := range indexMap {
		if len(items) < 2 {
			stats[pointID] = []dto.DeformCurveItem{}
			continue
		}

		// 按时间排序，避免乱序
		sort.Slice(items, func(i, j int) bool {
			return items[i].T.Before(items[j].T)
		})

		var velocityList []dto.DeformCurveItem
		for i := 1; i < len(items); i++ {
			dt := items[i].T.Sub(items[i-1].T).Seconds()
			if dt <= 0 {
				continue
			}

			dv := float64(items[i].V-items[i-1].V) / 100.0 // 转成 mm
			velocityMM := dv / dt                          // mm/s

			velocityList = append(velocityList, dto.DeformCurveItem{
				T: items[i].T,
				V: int64(velocityMM), // 保存成整数（mm/h * 100 可以自己决定）
			})
		}

		stats[pointID] = velocityList
	}

	return stats, nil
}

// 计算形变加速度
func (s DeformationPointV2) calcDeformAcceleration(indexMap map[int64][]dto.DeformCurveItem) (map[int64][]dto.DeformCurveItem, error) {
	stats := make(map[int64][]dto.DeformCurveItem)

	// 先计算速度
	velocityMap, err := s.calcDeformVelocity(indexMap)
	if err != nil {
		return nil, err
	}

	for pointID, velocityItems := range velocityMap {
		if len(velocityItems) < 2 { // 至少需要两个速度点才能计算加速度
			stats[pointID] = []dto.DeformCurveItem{}
			continue
		}

		// 按时间排序，避免乱序
		sort.Slice(velocityItems, func(i, j int) bool {
			return velocityItems[i].T.Before(velocityItems[j].T)
		})

		var accelList []dto.DeformCurveItem
		for i := 1; i < len(velocityItems); i++ {
			dt := velocityItems[i].T.Sub(velocityItems[i-1].T).Seconds()
			if dt <= 0 {
				continue
			}

			dv := float64(velocityItems[i].V-velocityItems[i-1].V) / 100.0 // 速度差 mm/h
			accel := dv / dt                                               // mm/h/s

			accelList = append(accelList, dto.DeformCurveItem{
				T: velocityItems[i].T,
				V: int64(accel * 100),
			})
		}

		stats[pointID] = accelList
	}

	return stats, nil
}
