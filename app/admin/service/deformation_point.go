package service

import (
	"context"
	"errors"
	"fmt"
	"go-admin/app/admin/service/dto"
	"go-admin/app/monsvr/mongosvr"
	"go-admin/app/monsvr/mongosvr/collections"
	"math"
	"sort"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk/service"
)

type DeformationPoint struct {
	service.Service
}

// 获取形变数据列表
func (s DeformationPoint) GetDeformationPoinList(ctx context.Context, req dto.GetDeformationDataReq) (*dto.GetDeformationDataResp, error) {
	var err error
	hours := req.Hours
	timeType := req.TimeType
	// 获取数据
	var data []collections.DeformationPointModel
	if data, err = mongosvr.DeformationPointService.QueryDeformationPointData(ctx, req); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("暂无形变数据")
	}
	// 根据时间范围计算最大采样点数
	maxPoints := s.getMaxPointsForRange(hours)
	// 采样
	sampledData := s.sampleDeformationData(data, maxPoints)
	// 根据颗粒度聚合
	aggregatedData := s.aggregateByTimeType(sampledData, timeType)

	fmt.Printf("原始数据条数: %d\n", len(data))
	fmt.Printf("采样后数据条数: %d\n", len(sampledData))
	fmt.Printf("aggregatedData: %d\n", len(aggregatedData))

	resp := &dto.GetDeformationDataResp{
		LastTime: data[len(data)-1].SvrTime.Local(),
		List:     aggregatedData,
	}
	return resp, nil
}

// 获取形变速度
func (s DeformationPoint) GetDeformationVelocity(ctx context.Context, req dto.GetDeformationDataReq) (*dto.GetDeformationVelocityResp, error) {
	var err error
	resp := &dto.GetDeformationVelocityResp{}
	var deformData *dto.GetDeformationDataResp
	if deformData, err = s.GetDeformationPoinList(ctx, req); err != nil {
		return nil, err
	}
	if len(deformData.List) < 2 {
		s.Log.Info("数据量太少，无法计算形变速度")
		return resp, nil
	}
	var velocityList []dto.DeformationVelocityItem
	if velocityList, err = s.CalcaDeformationSpeedStat(deformData.List); err != nil {
		return nil, fmt.Errorf("服务器繁忙")
	}

	resp.LastTime = deformData.LastTime
	resp.List = velocityList
	return resp, nil
}

// 特点：
// 按 整数索引步长 均匀划分数据区间。
// 区间内计算最大值和最小值。
// 选择 离区间均值更远的极值 作为采样点。
// 严格控制采样点数不超过 maxPoints。
// ✅ 优点：
// 输出点数固定，不会超过 maxPoints，适合图表显示或存储限制。
// 保留趋势和波动较大的点。
// ⚠️ 注意：
// 每个区间只选择 一个点（最大值或最小值），如果需要同时保留极值信息，这种方法可能会丢掉另一个极值。
func (s DeformationPoint) sampleDeformationDataAvg(data []collections.DeformationPointModel, maxPoints int) []collections.DeformationPointModel {
	if len(data) <= maxPoints {
		return data
	}

	step := len(data) / maxPoints
	if step < 1 {
		step = 1
	}

	sampled := make([]collections.DeformationPointModel, 0, maxPoints)
	for i := 0; i < len(data); i += step {
		end := i + step
		if end > len(data) {
			end = len(data)
		}

		// 在区间内找极值
		maxVal := data[i]
		minVal := data[i]
		for j := i; j < end; j++ {
			if data[j].Deformation > maxVal.Deformation {
				maxVal = data[j]
			}
			if data[j].Deformation < minVal.Deformation {
				minVal = data[j]
			}
		}

		// 选波动幅度较大的点（离均值更远）
		mid := (data[i].Deformation + data[end-1].Deformation) / 2
		if s.abs(maxVal.Deformation-mid) > s.abs(minVal.Deformation-mid) {
			sampled = append(sampled, maxVal)
		} else {
			sampled = append(sampled, minVal)
		}

		// 控制上限
		if len(sampled) >= maxPoints {
			break
		}
	}

	return sampled
}

// 对形变数据进行整数索引采样
// 将数据分为 maxPoints 个区间，保留每个区间的最大和最小值，趋势和极值都不会丢失。
// 返回两个数组：maxValues 和 minValues
// 对形变数据进行整数索引采样
// 将数据分为 maxPoints 个区间，保留每个区间的最大、最小值和平均值
func (s DeformationPoint) sampleDeformationData(data []collections.DeformationPointModel, maxPoints int) []dto.DeformationDataItem {
	n := len(data)
	if n == 0 || maxPoints <= 0 {
		return nil
	}

	values := make([]dto.DeformationDataItem, 0, maxPoints)
	// 如果数据量小于等于 maxPoints，则每个数组直接返回原数据
	if n <= maxPoints {
		for _, item := range data {
			values = append(values, dto.DeformationDataItem{
				Time:           item.SvrTime,
				DeformationMax: item.Deformation,
				DeformationMin: item.Deformation,
				DeformationAvg: item.Deformation,
				Distance:       item.Distance,
			})
		}
		return values
	}

	// 每个区间的起止索引
	for i := 0; i < maxPoints; i++ {
		start := i * n / maxPoints
		end := (i + 1) * n / maxPoints
		if end > n {
			end = n
		}
		if start >= end {
			break
		}

		maxVal := data[start]
		minVal := data[start]
		sum := 0
		count := 0

		for j := start; j < end; j++ {
			val := data[j].Deformation
			sum += val
			count++
			if val > maxVal.Deformation {
				maxVal = data[j]
			}
			if val < minVal.Deformation {
				minVal = data[j]
			}
		}

		avgVal := sum / count

		values = append(values, dto.DeformationDataItem{
			Time:           maxVal.SvrTime,
			DeformationMax: maxVal.Deformation,
			DeformationMin: minVal.Deformation,
			DeformationAvg: avgVal,
			Distance:       maxVal.Distance,
		})
	}

	// 按时间排序，保证时间序列正确
	sort.SliceStable(values, func(i, j int) bool {
		return values[i].Time.Before(values[j].Time)
	})

	return values
}

func (s DeformationPoint) abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// 根据时间范围计算最大采样点数
func (s DeformationPoint) getMaxPointsForRange(hours int64) int {
	// 用 map 定义时间范围对应的最大采样点数
	var maxPointsMap = map[int64]int{
		1:       3600, //原始数据
		3:       200,
		6:       250,
		12:      300,
		24:      400,
		3 * 24:  500,
		7 * 24:  500,
		30 * 24: 600,
		90 * 24: 600,
	}
	if val, ok := maxPointsMap[hours]; ok {
		return val
	}
	return 500 // 默认值
}

// 获得时间格式转化
func (s DeformationPoint) getTimeFormatByType(t time.Time, timeType string) (string, time.Time) {
	var key string
	var bucket time.Time
	switch timeType {
	case "seconds":
		bucket = t.Truncate(time.Second)
		key = bucket.Format("2006-01-02 15:04:05")
	case "minutes":
		bucket = t.Truncate(time.Minute)
		key = bucket.Format("2006-01-02 15:04")
	case "hours":
		bucket = t.Truncate(time.Hour)
		key = bucket.Format("2006-01-02 15")
	case "days":
		bucket = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		key = bucket.Format("2006-01-02")
	default:
		bucket = t.Truncate(time.Second)
		key = bucket.Format("2006-01-02 15:04:05")
	}
	// fmt.Println("key:", key)
	// fmt.Println("bucket:", bucket.String())
	return key, bucket
}

// 按颗粒度聚合
// func (s DeformationPoint) aggregateByTimeType(data []dto.DeformationDataItem, granularity string) []dto.DeformationDataItem {
// 	if len(data) == 0 {
// 		return nil
// 	}
// 	group := make(map[string][]dto.DeformationDataItem)
// 	timeMap := make(map[string]time.Time, 0)

// 	for _, d := range data {
// 		var bucket time.Time
// 		key, bucket := s.getTimeFormatByType(d.Time, granularity)
// 		group[key] = append(group[key], d)
// 		timeMap[key] = bucket
// 	}

// 	results := make([]dto.DeformationDataItem, 0, len(group))
// 	for k, vals := range group {
// 		sumDeformMax := 0
// 		sumDeformMin := 0
// 		sumDeformAvg := 0
// 		sumDist := 0
// 		for _, v := range vals {
// 			sumDeformMax += v.DeformationMax
// 			sumDeformMin += v.DeformationMin
// 			sumDeformAvg += v.DeformationAvg
// 			sumDist += v.Distance
// 		}
// 		count := len(vals)
// 		results = append(results, dto.DeformationDataItem{
// 			Time:           timeMap[k],
// 			DeformationMax: sumDeformMax / count, // 平均
// 			DeformationMin: sumDeformMin / count, // 平均
// 			DeformationAvg: sumDeformAvg / count, // 平均
// 			Distance:       sumDist / count,      // 平均
// 		})
// 	}

// 	// 按时间排序
// 	sort.Slice(results, func(i, j int) bool {
// 		return results[i].Time.Before(results[j].Time)
// 	})

// 	return results
// }

// 按颗粒度聚合，保留极值，平均值计算平均
func (s DeformationPoint) aggregateByTimeType(data []dto.DeformationDataItem, granularity string) []dto.DeformationDataItem {
	if len(data) == 0 {
		return nil
	}

	group := make(map[string][]dto.DeformationDataItem)
	timeMap := make(map[string]time.Time)

	// 按时间颗粒度分组
	for _, d := range data {
		key, bucket := s.getTimeFormatByType(d.Time, granularity)
		group[key] = append(group[key], d)
		timeMap[key] = bucket
	}

	results := make([]dto.DeformationDataItem, 0, len(group))
	for k, vals := range group {
		if len(vals) == 0 {
			continue
		}

		// 初始化极值
		maxVal := vals[0].DeformationMax
		minVal := vals[0].DeformationMin
		sumAvg := 0
		sumDist := 0

		for _, v := range vals {
			if v.DeformationMax > maxVal {
				maxVal = v.DeformationMax
			}
			if v.DeformationMin < minVal {
				minVal = v.DeformationMin
			}
			sumAvg += v.DeformationAvg
			sumDist += v.Distance
		}

		count := len(vals)
		avgVal := sumAvg / count
		distVal := sumDist / count

		// 保证 Min <= Avg <= Max
		if avgVal < minVal {
			avgVal = minVal
		}
		if avgVal > maxVal {
			avgVal = maxVal
		}

		results = append(results, dto.DeformationDataItem{
			Time:           timeMap[k],
			DeformationMax: maxVal,
			DeformationMin: minVal,
			DeformationAvg: avgVal,
			Distance:       distVal,
		})
	}

	// 按时间升序排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Time.Before(results[j].Time)
	})

	return results
}

// 计算形变速度统计（单位：mm/s）
// deformData: 采样并聚合后的形变数据
func (s DeformationPoint) CalcaDeformationSpeedStat(deformData []dto.DeformationDataItem) ([]dto.DeformationVelocityItem, error) {
	if len(deformData) < 2 {
		return nil, fmt.Errorf("数据量太少，无法计算速度")
	}

	stats := make([]dto.DeformationVelocityItem, 0, len(deformData)-1)

	// 打印输入数据，便于调试
	fmt.Printf("deformData:%+v \n", deformData)

	for i := 1; i < len(deformData); i++ {
		prev := deformData[i-1]
		curr := deformData[i]

		// 打印每个数据点的最大值、平均值、最小值
		fmt.Println("打印max:", curr.DeformationMax)
		fmt.Println("打印avg:", curr.DeformationAvg)
		fmt.Println("打印min:", curr.DeformationMin)
		fmt.Println(" ")

		// 时间差（秒）
		deltaTime := curr.Time.Sub(prev.Time).Seconds()
		if deltaTime <= 0 {
			// 时间差为0时跳过，避免除以0
			continue
		}

		// 最大速度 = 最大值变化 / 100(mm) / deltaTime(秒)
		speedMax := (float64(curr.DeformationMax) - float64(prev.DeformationMax)) / 100.0 / deltaTime
		// 最小速度 = 最小值变化 / 100(mm) / deltaTime(秒)
		speedMin := (float64(curr.DeformationMin) - float64(prev.DeformationMin)) / 100.0 / deltaTime
		// 平均速度 = 平均值变化 / 100(mm) / deltaTime(秒)
		speedAvg := (float64(curr.DeformationAvg) - float64(prev.DeformationAvg)) / 100.0 / deltaTime

		// 保留四位小数，避免小幅度速度被截断为0
		speedMax = math.Round(speedMax*10000) / 10000
		speedMin = math.Round(speedMin*10000) / 10000
		speedAvg = math.Round(speedAvg*10000) / 10000

		stats = append(stats, dto.DeformationVelocityItem{
			Time: curr.Time,
			Max:  speedMax,
			Min:  speedMin,
			Avg:  speedAvg,
		})
	}

	// 打印最终速度统计，便于调试
	fmt.Printf("stats:%+v \n", stats)

	return stats, nil
}
