package tasks

import (
	"fmt"
	"go-admin/app/radar/models"
	"go-admin/app/radar/service"
	"time"

	"github.com/gin-gonic/gin"
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
	Ticker  *time.Ticker
	Context *gin.Context
	Logger  *logger.Helper
	DB      *gorm.DB
	Config  []AlarmPointConifg //预警配置
}

// AlarmPointConifg 监测点预警配置
type AlarmPointConifg struct {
	AlarmPointId  int64            `json:"alarmPointId"` //监测点ID
	LastAlarmTime time.Time        `json:"last_time"`    //最近一次检测预警的时间
	RadarId       int64            `json:"radarId"`      //雷达Id
	DeptId        int64            `json:"deptId"`       //机构ID
	AlarmType     models.AlarmType `json:"alarmType"`    //预警类型
	Interval      uint64           `json:"interval"`     //预警间隔时间(m)
	Duration      uint64           `json:"duration"`     //形变值的查询事件跨度（h）
	RedOption     string           `json:"redOption"`    //红色预警条件
	OrangeOption  string           `json:"orangeOption"` //橙色预警条件
	YellowOption  string           `json:"yellowOption"` //黄色预警条件
	BlueOption    string           `json:"blueOption"`   //蓝色预警条件
}

// InitAlarmPointTask 初始化监测点预警任务
func InitAlarmPointTask() {
	t := &AlarmPointTask{
		Ticker: time.NewTicker(time.Second * 3),
		Logger: logger.NewHelper(sdk.Runtime.GetLogger()).WithFields(map[string]interface{}{}),
		DB:     sdk.Runtime.GetDbByKey("*"),
		Config: []AlarmPointConifg{},
	}
	t.startTask()
}

func (t *AlarmPointTask) startTask() {
	for {
		<-t.Ticker.C
		t.Logger.Info("执行监测点预警定时任务")
		t.monitor()
	}
}

func (t *AlarmPointTask) Stop() {
	if t.Ticker != nil {
		t.Ticker.Stop()
	}
	t.Logger.Info("监测点预警定时任务已停止")
}

func (t *AlarmPointTask) monitor() {
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
		fmt.Println("设置radarToDeptMap[item.RadarId] ：", item.RadarId, item.DeptId)
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
				for _, item := range defaultM {
					key := fmt.Sprintf("%d_%d", deptId, item.RadarPointId)
					t.Logger.Info("key: ", key)
					alarmMap[key] = append(alarmMap[key], item)
				}
			}
		}
		if aps, ok := alarmMap[key]; ok {
			for _, ap := range aps {
				cfg := AlarmPointConifg{
					AlarmPointId:  ap.Id,
					LastAlarmTime: item.LastAlarmTime.Time,
					RadarId:       item.RadarId,
					DeptId:        deptId,
					Interval:      ap.Interval,
					Duration:      ap.Duration,
					AlarmType:     ap.AlarmType,
					RedOption:     ap.RedOption,
					OrangeOption:  ap.OrangeOption,
					YellowOption:  ap.YellowOption,
					BlueOption:    ap.BlueOption,
				}
				t.Config = append(t.Config, cfg)
			}
		}
	}
	fmt.Println("打印配置:", len(t.Config))

	for _, item := range t.Config {
		switch item.AlarmType {
		case models.AlarmTypeRadarPointDeformation:
			t.monitorDeformation()
		case models.AlarmTypeRadarPointVelocity:
			t.monitorVelocity()
		case models.AlarmTypeRadarPointAcceleration:
			t.monitorAcceleration()
		}
	}
}

// 监测形变报警
func (t *AlarmPointTask) monitorDeformation() {}

// 监测形变报警
func (t *AlarmPointTask) monitorVelocity() {}

// 监测形变报警
func (t *AlarmPointTask) monitorAcceleration() {}
