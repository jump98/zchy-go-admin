package tasks

import (
	"context"
	"database/sql"
	"fmt"
	"go-admin/app/monsvr/mongosvr"
	"go-admin/app/monsvr/mongosvr/collections"
	"go-admin/app/radar/models"
	"go-admin/app/radar/service"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"gorm.io/gorm"
)

//我需要执行定时任务，用来执行监控预警。
//我有配置表：
//radar_point:是雷达监测点表
//alarm_point:是监测点配置表,配置了每个radar_point的间隔时间，时间单位是分
//我需要启动定时任务，轮询radar_point表，得到alarm_point配置，然后去做逻辑处理，判断是否有预警信息

type AlarmPointTask struct {
	ctx    context.Context
	cancel context.CancelFunc
	Logger *logger.Helper
	DB     *gorm.DB
	Config []AlarmPointConifg //预警配置
	mu     sync.Mutex
}

// AlarmPointConifg 监测点预警配置
type AlarmPointConifg struct {
	AlarmPointId  int64               `json:"alarmPointId"` //监测点ID
	LastAlarmTime time.Time           `json:"last_time"`    //最近一次检测预警的时间
	RadarId       int64               `json:"radarId"`      //雷达Id
	PointIndex    int64               `json:"rointIndex"`   //监测点index
	DeptId        int64               `json:"deptId"`       //机构ID
	AlarmPoint    []models.AlarmPoint `json:"alarmPoint"`   //预警规则
	RadarPoint    models.RadarPoint   `json:"radarPoint"`   //监测点信息
}

// InitAlarmPointTask 初始化监测点预警任务
func InitAlarmPointTask(parentCtx context.Context) *AlarmPointTask {
	ctx, cancel := context.WithCancel(parentCtx)
	t := &AlarmPointTask{

		Logger: logger.NewHelper(sdk.Runtime.GetLogger()).WithFields(map[string]interface{}{}),
		DB:     sdk.Runtime.GetDbByKey("*"),
		Config: []AlarmPointConifg{},
		ctx:    ctx,
		cancel: cancel,
	}
	go t.startTask()
	return t
}

// 1分钟执行一次任务
func (t *AlarmPointTask) startTask() {
	defer func() {
		if err := recover(); err != nil {
			t.Logger.Error(err)
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	//ticker := time.NewTicker(time.Minute * 1)
	defer ticker.Stop()

	for {
		select {
		case <-t.ctx.Done():
			t.Logger.Info("AlarmPointTask 停止")
			return
		case <-ticker.C:
			t.mu.Lock()
			t.Logger.Info("AlarmPointTask 开始")
			t.monitor()
			t.mu.Unlock()
		}
	}
}

func (t *AlarmPointTask) Stop() {
	t.cancel()
}

func (t *AlarmPointTask) monitor() {
	ctx := context.Background()
	var err error
	t.Config = make([]AlarmPointConifg, 0)

	radarPointItems := make([]models.RadarPoint, 0)
	if err = t.DB.Model(models.RadarPoint{}).Find(&radarPointItems).Error; err != nil {
		t.Logger.Error(err)
		return
	}
	t.Logger.Info("打印监测点数量：", len(radarPointItems))

	radarItems := make([]models.Radar, 0)
	if err = t.DB.Model(models.Radar{}).Find(&radarItems).Error; err != nil {
		return
	}
	radarToDeptMap := map[int64]int64{}
	for _, item := range radarItems {
		radarToDeptMap[item.RadarId] = item.DeptId
	}
	alarmPointItems := make([]models.AlarmPoint, 0)
	if err = t.DB.Model(models.AlarmPoint{}).Find(&alarmPointItems).Error; err != nil {
		return
	}

	alarmMap := map[string][]models.AlarmPoint{}
	for _, item := range alarmPointItems {
		key := fmt.Sprintf("%d_%d", item.DeptId, item.RadarPointId)
		alarmMap[key] = append(alarmMap[key], item)
	}

	for _, item := range radarPointItems {
		var deptId int64
		var ok bool
		if deptId, ok = radarToDeptMap[item.RadarId]; !ok {
			t.Logger.Error("监测点的雷达信息不存在,radarPointId:", item.Id)
			continue
		}
		var key string
		if item.MTypeId == models.RadarPointMTypeAlone {
			key = fmt.Sprintf("%d_%d", deptId, item.Id)
		}
		if item.MTypeId == models.RadarPointMTypeGlobal {
			key = fmt.Sprintf("%d_%d", deptId, 0)
			if _, ok := alarmMap[key]; !ok {
				s := service.AlarmPoint{}
				defaultM := s.GetDefaultRadarPointConfig(deptId, 0, 0)
				alarmMap[key] = defaultM
			}
		}
		if aps, ok := alarmMap[key]; ok {
			cfg := AlarmPointConifg{
				AlarmPointId:  item.Id,
				PointIndex:    item.PointIndex,
				LastAlarmTime: item.LastAlarmTime.Time,
				RadarId:       item.RadarId,
				DeptId:        deptId,
				AlarmPoint:    aps,
				RadarPoint:    item,
			}
			t.Config = append(t.Config, cfg)
		}
	}

	fmt.Println("打印配置数量:", len(t.Config))

	for _, item := range t.Config {
		data := t.findDeformMintueList(ctx, item)
		if len(data) == 0 {
			return
		}

		velocityData := t.getVelocityDataList(item.LastAlarmTime, data)

		t.Logger.Info("加速度的数据:", len(velocityData))
		for _, item := range velocityData {
			t.Logger.Info("时间：", item.Time.Local().Format("2006-01-02 15:04:05"))
		}

		deptId := item.DeptId
		//启动预警监测
		for _, alarmItem := range item.AlarmPoint {
			switch alarmItem.AlarmType {
			case models.AlarmTypeRadarPointDeformation:
				t.monitorDeform(deptId, &item.RadarPoint, alarmItem, data)
			case models.AlarmTypeRadarPointVelocity:
				t.monitorVelocity(deptId, &item.RadarPoint, alarmItem, velocityData)
			case models.AlarmTypeRadarPointAcceleration:
				t.monitorAcceleration(deptId, &item.RadarPoint, alarmItem, velocityData)
			}
		}
		//保存最近一次检测预警的时间
		t.saveRadarPointTime(item.RadarPoint)
	}
}

func (t *AlarmPointTask) saveRadarPointTime(radarPoint models.RadarPoint) {
	var err error
	lastTime := sql.NullTime{
		Time:  time.Now().Local(), // 真实时间
		Valid: true,               // 表示这个时间有效
	}
	radarPoint.LastAlarmTime = lastTime
	if err = t.DB.Save(&radarPoint).Error; err != nil {
		return
	}
}

// 监测形变报警
func (t *AlarmPointTask) findDeformMintueList(ctx context.Context, item AlarmPointConifg) []collections.DeformationPointMinuteModel {
	var err error
	defer func() {
		if err != nil {
			t.Logger.Error("findDeformMintueList.err:", err)
		}
	}()

	interval := item.AlarmPoint[0].Interval
	duration := item.AlarmPoint[0].Duration

	now := time.Now()
	elapsed := now.Sub(item.LastAlarmTime)
	//fmt.Println("elapsed:", elapsed.String())
	if elapsed < time.Duration(interval)*time.Minute {
		t.Logger.Infof("预警间隔时间中:%d minute...", interval)
		return []collections.DeformationPointMinuteModel{}
	}
	radarId := item.RadarId
	pointIndex := item.PointIndex
	// 使用固定格式的字符串
	startTime := now.Add(-time.Hour * time.Duration(duration))
	endTime := now
	// 获取数据
	var data []collections.DeformationPointMinuteModel
	if data, err = mongosvr.DeformationPointMinuteService.FindDeformMinuteByTime(ctx, radarId, pointIndex, startTime, endTime); err != nil {
		t.Logger.Error("查询出错:", err)
		return []collections.DeformationPointMinuteModel{}
	}
	fmt.Println("查询形变数据结果:", len(data))
	if len(data) == 0 {
		return data
	}

	return data
}

// 加速度形变数据 - 只取上一次预警到当前时间的形变列表
func (t *AlarmPointTask) getVelocityDataList(lastAlarmTime time.Time, data []collections.DeformationPointMinuteModel) []collections.DeformationPointMinuteModel {
	var velocityData []collections.DeformationPointMinuteModel
	for _, item := range data {
		if item.Time.After(lastAlarmTime) {
			velocityData = append(velocityData, item)
		}
	}
	return velocityData
}

// 监测累计形变
func (t *AlarmPointTask) monitorDeform(deptId int64, radarPoint *models.RadarPoint, alarmPoint models.AlarmPoint, data []collections.DeformationPointMinuteModel) {
	var err error
	defer func() {
		if err != nil {
			t.Logger.Error("monitorDeform.err:", err)
		}
	}()

	if len(data) == 0 {
		t.Logger.Info("无形变数据，无法计算累计形变")
		return
	}

	// 计算累计形变值（原始值）
	var totalDeformation int64 = 0
	for _, d := range data {
		totalDeformation += d.Deformation
	}
	// 转换为毫米
	totalDeformationMM := totalDeformation / 100
	t.Logger.Info("累计形变位移:", totalDeformationMM)

	// 解析各预警阈值
	var redValue, orangeValue, yellowValue, blueValue int64
	if redValue, err = strconv.ParseInt(alarmPoint.RedOption, 10, 64); err != nil {
		return
	}
	if orangeValue, err = strconv.ParseInt(alarmPoint.OrangeOption, 10, 64); err != nil {
		return
	}
	if yellowValue, err = strconv.ParseInt(alarmPoint.YellowOption, 10, 64); err != nil {
		return
	}
	if blueValue, err = strconv.ParseInt(alarmPoint.BlueOption, 10, 64); err != nil {
		return
	}

	// 使用绝对值判断预警等级
	absDeformation := totalDeformationMM
	if absDeformation < 0 {
		absDeformation = -absDeformation
	}

	var alarmLevel models.AlarmLevel
	var alarmValue int64 // 触发的预警阈值
	switch {
	case absDeformation >= redValue:
		alarmLevel = models.AlarmLevelRed
		alarmValue = redValue
	case absDeformation >= orangeValue:
		alarmLevel = models.AlarmLevelOrange
		alarmValue = orangeValue
	case absDeformation >= yellowValue:
		alarmLevel = models.AlarmLevelYellow
		alarmValue = yellowValue
	case absDeformation >= blueValue:
		alarmLevel = models.AlarmLevelBlue
		alarmValue = blueValue
	default:
		alarmLevel = models.AlarmLevelNone
		alarmValue = 0
		return
	}

	if alarmLevel > radarPoint.AlarmLevel {
		radarPoint.AlarmLevel = alarmLevel
	}

	// 输出日志
	t.Logger.Info(fmt.Sprintf("监测点ID: %d, 累计形变: %d, 预警等级: %d, 触发阈值: %d",
		radarPoint.Id, totalDeformationMM, alarmLevel, alarmValue))

	db := t.DB
	alarmPointLogs := &models.AlarmPointLogs{
		AlarmType:     alarmPoint.AlarmType,
		RadarId:       radarPoint.RadarId,
		RadarPointId:  radarPoint.Id,
		AlarmLevel:    alarmLevel,
		DeptId:        deptId,
		CurrentValue:  strconv.FormatInt(totalDeformationMM, 10), // 原始累计形变
		AlarmValue:    strconv.FormatInt(alarmValue, 10),         // 触发阈值
		Interval:      alarmPoint.Interval,
		Duration:      1,
		ProcessRemark: "",
	}

	if err := db.Create(alarmPointLogs).Error; err != nil {
		t.Logger.Error("保存AlarmPointLogs失败:", err)
		return
	}

	t.Logger.Info("累计形变预警记录已保存到数据库")
}

// 监测形变位移速度（瞬时位移量）预警阈值，单位：mm/m
func (t *AlarmPointTask) monitorVelocity(deptId int64, radarPoint *models.RadarPoint, alarmPoint models.AlarmPoint, data []collections.DeformationPointMinuteModel) {
	if len(data) < 2 {
		t.Logger.Info("数据不足，无法计算瞬时形变速度")
		return
	}

	// 解析各预警阈值
	var err error
	var redValue, orangeValue, yellowValue, blueValue float64
	if redValue, err = strconv.ParseFloat(alarmPoint.RedOption, 64); err != nil {
		return
	}
	if orangeValue, err = strconv.ParseFloat(alarmPoint.OrangeOption, 64); err != nil {
		return
	}
	if yellowValue, err = strconv.ParseFloat(alarmPoint.YellowOption, 64); err != nil {
		return
	}
	if blueValue, err = strconv.ParseFloat(alarmPoint.BlueOption, 64); err != nil {
		return
	}

	var maxAlarmLevel models.AlarmLevel = models.AlarmLevelNone
	var alarmValue float64
	var currentValue float64
	var alarmCount int

	// 遍历每两条数据计算瞬时速度
	for i := 1; i < len(data); i++ {
		prev := data[i-1]
		curr := data[i]

		deltaDeform := float64(curr.Deformation-prev.Deformation) / 100.0 // mm
		deltaTime := curr.Time.Sub(prev.Time).Minutes()                   // 分钟
		if deltaTime <= 0 {
			continue
		}
		speed := deltaDeform / deltaTime
		absSpeed := math.Abs(speed)
		t.Logger.Info("形变A:", curr.Deformation, "  形变B:", prev.Deformation, "  形变差值:", deltaDeform, "  时间:", deltaTime, "  曲线速度:", speed)

		var level models.AlarmLevel
		var curTrigger float64
		switch {
		case absSpeed >= redValue:
			level = models.AlarmLevelRed
			curTrigger = redValue
		case absSpeed >= orangeValue:
			level = models.AlarmLevelOrange
			curTrigger = orangeValue
		case absSpeed >= yellowValue:
			level = models.AlarmLevelYellow
			curTrigger = yellowValue
		case absSpeed >= blueValue:
			level = models.AlarmLevelBlue
			curTrigger = blueValue
		default:
			continue
		}

		// 更新最大等级
		if level > maxAlarmLevel {
			maxAlarmLevel = level
			alarmValue = curTrigger
			currentValue = speed // 保留原始值
		}
		alarmCount++
	}

	if maxAlarmLevel == models.AlarmLevelNone {
		t.Logger.Info("未触发任何瞬时速度预警")
		return
	}
	if maxAlarmLevel > radarPoint.AlarmLevel {
		radarPoint.AlarmLevel = maxAlarmLevel
	}

	t.Logger.Info(fmt.Sprintf("监测点ID: %d, 瞬时速度预警次数: %d, 最大预警等级: %d, 触发阈值: %.2f, 原始速度: %.2f",
		radarPoint.Id, alarmCount, maxAlarmLevel, alarmValue, currentValue))

	//TODO 后续需要查询是否存在同等级预警，如果存在就不上报
	db := t.DB
	alarmPointLogs := &models.AlarmPointLogs{
		AlarmType:     alarmPoint.AlarmType,
		RadarId:       radarPoint.RadarId,
		RadarPointId:  radarPoint.Id,
		AlarmLevel:    maxAlarmLevel,
		DeptId:        deptId,
		CurrentValue:  fmt.Sprintf("%.2f", currentValue), // 瞬时速度原始值
		AlarmValue:    fmt.Sprintf("%.2f", alarmValue),   // 触发阈值
		Interval:      alarmPoint.Interval,
		Duration:      uint64(alarmCount), // 触发次数
		ProcessRemark: "",
	}

	if err := db.Create(alarmPointLogs).Error; err != nil {
		t.Logger.Error("保存瞬时速度AlarmPointLogs失败:", err)
		return
	}

	t.Logger.Info("瞬时速度预警记录已保存到数据库")
}

// 监测形变位移加速度（瞬时加速度），单位：mm/m²
func (t *AlarmPointTask) monitorAcceleration(deptId int64, radarPoint *models.RadarPoint, alarmPoint models.AlarmPoint, data []collections.DeformationPointMinuteModel) {
	if len(data) < 3 { // 至少 3 条数据才能计算两段速度的加速度
		t.Logger.Info("数据不足，无法计算瞬时加速度")
		return
	}

	// 解析各预警阈值
	var err error
	var redValue, orangeValue, yellowValue, blueValue float64
	if redValue, err = strconv.ParseFloat(alarmPoint.RedOption, 64); err != nil {
		return
	}
	if orangeValue, err = strconv.ParseFloat(alarmPoint.OrangeOption, 64); err != nil {
		return
	}
	if yellowValue, err = strconv.ParseFloat(alarmPoint.YellowOption, 64); err != nil {
		return
	}
	if blueValue, err = strconv.ParseFloat(alarmPoint.BlueOption, 64); err != nil {
		return
	}

	var maxAlarmLevel models.AlarmLevel = models.AlarmLevelNone
	var alarmValue float64
	var currentValue float64
	var alarmCount int

	// 先计算每条数据的瞬时速度
	speeds := make([]float64, len(data)-1)
	times := make([]float64, len(data)-1) // 时间间隔分钟
	for i := 1; i < len(data); i++ {
		deltaDeform := float64(data[i].Deformation-data[i-1].Deformation) / 100.0
		deltaTime := data[i].Time.Sub(data[i-1].Time).Minutes()
		if deltaTime <= 0 {
			deltaTime = 1 // 防止除0，默认1分钟
		}
		speeds[i-1] = deltaDeform / deltaTime // mm/m
		times[i-1] = deltaTime
	}

	// 计算加速度 = 速度变化 / 时间变化
	for i := 1; i < len(speeds); i++ {
		deltaV := speeds[i] - speeds[i-1]
		deltaT := (times[i] + times[i-1]) / 2 // 平均时间间隔
		acc := deltaV / deltaT
		absAcc := math.Abs(acc)

		var level models.AlarmLevel
		var curTrigger float64
		switch {
		case absAcc >= redValue:
			level = models.AlarmLevelRed
			curTrigger = redValue
		case absAcc >= orangeValue:
			level = models.AlarmLevelOrange
			curTrigger = orangeValue
		case absAcc >= yellowValue:
			level = models.AlarmLevelYellow
			curTrigger = yellowValue
		case absAcc >= blueValue:
			level = models.AlarmLevelBlue
			curTrigger = blueValue
		default:
			continue
		}

		if level > maxAlarmLevel {
			maxAlarmLevel = level
			alarmValue = curTrigger
			currentValue = acc // 保留原始值
		}
		alarmCount++
	}

	if maxAlarmLevel == models.AlarmLevelNone {
		t.Logger.Info("未触发任何瞬时加速度预警")
		return
	}
	if maxAlarmLevel > radarPoint.AlarmLevel {
		radarPoint.AlarmLevel = maxAlarmLevel
	}

	t.Logger.Info(fmt.Sprintf("监测点ID: %d, 瞬时加速度预警次数: %d, 最大预警等级: %d, 触发阈值: %.2f, 原始加速度: %.2f",
		radarPoint.Id, alarmCount, maxAlarmLevel, alarmValue, currentValue))

	// 保存到数据库
	db := t.DB
	alarmPointLogs := &models.AlarmPointLogs{
		AlarmType:     alarmPoint.AlarmType,
		RadarId:       radarPoint.RadarId,
		RadarPointId:  radarPoint.Id,
		AlarmLevel:    maxAlarmLevel,
		DeptId:        deptId,
		CurrentValue:  fmt.Sprintf("%.2f", currentValue),
		AlarmValue:    fmt.Sprintf("%.2f", alarmValue),
		Interval:      alarmPoint.Interval,
		Duration:      uint64(alarmCount),
		ProcessRemark: "",
	}

	if err := db.Create(alarmPointLogs).Error; err != nil {
		t.Logger.Error("保存瞬时加速度AlarmPointLogs失败:", err)
		return
	}

	t.Logger.Info("瞬时加速度预警记录已保存到数据库")
}
