package service

import (
	"errors"
	"fmt"
	"go-admin/app/radar/models"
	"go-admin/app/radar/service/dto"
	cDto "go-admin/common/dto"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"
)

type AlarmPoint struct {
	service.Service
}

// GetAlarmRules 获取所有的预警规则
func (e *AlarmPoint) GetAlarmRules(deptId int64, radarPointId int64) ([]models.AlarmPoint, error) {
	var err error
	alarmPointItems := make([]models.AlarmPoint, 0)
	if err = e.Orm.Model(&models.AlarmPoint{}).Where("dept_id = ? and radar_point_id = ?", deptId, radarPointId).Find(&alarmPointItems).Error; err != nil {
		return nil, err
	}
	if len(alarmPointItems) == 0 {
		fmt.Println("deptId:", deptId)
		fmt.Println("radarPointId:", radarPointId)
		alarmPointItems = e.GetDefaultRadarPointConfig(deptId, 0, radarPointId)
		if radarPointId == 0 {
			//创建全局预警规则
			//TODO:应该是在定时监测的时候创建
			if err = e.Orm.Model(&models.AlarmPoint{}).Create(&alarmPointItems).Error; err != nil {
				return nil, err
			}
		}
	}
	return alarmPointItems, nil
}

// AddAlarmRule 增加预警规则
//func (e *AlarmPoint) AddAlarmRule(req dto.AddAlarmPointReq) error {
//	var errTx error
//	tx := e.Orm.Begin()
//	defer func() {
//		if errTx != nil {
//			tx.Commit()
//		} else {
//			tx.Rollback()
//		}
//	}()
//
//	alarmPointItems := make([]*models.AlarmPoint, 4)
//	for i, item := range req.Items {
//		alarmPointItems[i].DeptId = item.DeptId
//		alarmPointItems[i].AlarmCheckType = item.AlarmCheckType
//		alarmPointItems[i].AlarmName = item.AlarmName
//		//alarmPointItems[i].RadarId = item.RadarId
//		alarmPointItems[i].RadarPointId = item.RadarPointId
//		alarmPointItems[i].AlarmType = item.AlarmType
//		alarmPointItems[i].RedOption = item.RedOption
//		alarmPointItems[i].OrangeOption = item.OrangeOption
//		alarmPointItems[i].YellowOption = item.YellowOption
//		alarmPointItems[i].BlueOption = item.BlueOption
//		alarmPointItems[i].Interval = item.Interval
//		alarmPointItems[i].Duration = item.Duration
//	}
//	if errTx = tx.Create(alarmPointItems).Error; errTx != nil {
//		return errTx
//	}
//	for _, item := range alarmPointItems {
//		fmt.Println("打印预警LevelID:", item.Id)
//	}
//	return errTx
//}

// UpdateAlarmRule 修改预警规则
func (e *AlarmPoint) UpdateAlarmRule(items []dto.AlarmPointItem, deptId, radarPointId int64, mode models.RadarPointMType) error {
	var err error
	db := e.Orm

	//设置监测点为全局门限
	if radarPointId != 0 && mode == models.RadarPointMTypeGlobal {
		radarPointItem := &models.RadarPoint{}
		if err = db.Model(&models.RadarPoint{}).Where("id = ?", radarPointId).First(radarPointItem).Error; err != nil {
			return err
		}
		radarPointItem.MTypeId = mode
		if err = db.Save(radarPointItem).Error; err != nil {
			return err
		}
		return nil
	}

	var errTx error
	tx := e.Orm.Begin()
	defer func() {
		if errTx != nil {
			tx.Rollback()
			e.Log.Error("出错：事务回滚:", err)
		} else {
			tx.Commit()
		}
	}()

	//设置全局门限 || 设置监测点为独立门限
	for _, item := range items {
		var alarmPoint models.AlarmPoint
		fmt.Println("item.AlarmType:", item.AlarmType)
		fmt.Println("item.RadarPointId:", item.RadarPointId)
		fmt.Println("item.DeptId:", item.DeptId)

		errTx = tx.Where("dept_id = ? AND radar_point_id = ? AND alarm_type = ?", deptId, radarPointId, item.AlarmType).First(&alarmPoint).Error

		if errors.Is(errTx, gorm.ErrRecordNotFound) {
			errTx = nil
			// 不存在则插入
			alarmPoint = models.AlarmPoint{
				DeptId:         deptId,
				RadarPointId:   radarPointId,
				AlarmCheckType: item.AlarmCheckType,
				AlarmName:      item.AlarmName,
				AlarmType:      item.AlarmType,
				RedOption:      item.RedOption,
				OrangeOption:   item.OrangeOption,
				YellowOption:   item.YellowOption,
				BlueOption:     item.BlueOption,
				Interval:       item.Interval,
				Duration:       item.Duration,
			}
			if errTx = tx.Create(&alarmPoint).Error; errTx != nil {
				return errTx
			}
		} else if errTx == nil {
			// 已存在则更新
			alarmPoint.AlarmName = item.AlarmName
			alarmPoint.RedOption = item.RedOption
			alarmPoint.OrangeOption = item.OrangeOption
			alarmPoint.YellowOption = item.YellowOption
			alarmPoint.BlueOption = item.BlueOption
			alarmPoint.Interval = item.Interval
			alarmPoint.Duration = item.Duration
			if errTx = tx.Save(&alarmPoint).Error; errTx != nil {
				return errTx
			}
		}
	}

	//设置监测点为独立门限
	if radarPointId != 0 {
		radarPointItem := &models.RadarPoint{}
		if errTx = tx.Model(&models.RadarPoint{}).Where("id = ?", radarPointId).First(radarPointItem).Error; errTx != nil {
			return errTx
		}
		radarPointItem.MTypeId = mode
		if errTx = tx.Save(radarPointItem).Error; errTx != nil {
			return errTx
		}
	}

	return errTx
}

// DeleteAlarmRule 删除预警规则
//func (e *AlarmPoint) DeleteAlarmRule(alarmRuleId int64) error {
//	alarmPointItems := make([]*models.AlarmPoint, 3)
//	var err error
//	db := e.Orm
//	alarmRuleItem := &models.AlarmPoint{
//		Id: alarmRuleId,
//	}
//	if err = db.First(alarmRuleItem).Error; err != nil {
//		return err
//	}
//	if err = db.Where("radar_point_id = ?", alarmRuleId).Find(alarmPointItems).Error; err != nil {
//		return err
//	}
//	var errTx error
//	tx := e.Orm.Begin()
//	defer func() {
//		if errTx != nil {
//			tx.Commit()
//		} else {
//			tx.Rollback()
//		}
//	}()
//	//1.删除预警规则
//	if errTx = tx.Model(&models.AlarmPoint{}).Where("id = ?", alarmRuleId).Delete(alarmRuleItem).Error; errTx != nil {
//		return errTx
//	}
//	////2.删除预警规矩等级
//	//if errTx = tx.Model(&models.AlarmRuleLevel{}).Where("alarm_rule_id = ?", alarmRuleId).Delete(alarmPointItems).Error; errTx != nil {
//	//	return errTx
//	//}
//	//TODO:
//	//3.删除预警告警配置
//	//4.删除预警联系人组
//	//5.删除预警联系人员
//	return errTx
//}

// GetDefaultRadarPointConfig 获得默认监测点预警配置
func (e *AlarmPoint) GetDefaultRadarPointConfig(deptId, radarId, radarPointId int64) []models.AlarmPoint {
	fmt.Println("GetDefaultRadarPointConfig:", deptId, radarId, radarPointId)
	items := make([]models.AlarmPoint, 0)
	items = append(items, models.AlarmPoint{
		DeptId:         deptId,
		AlarmCheckType: models.AlarmCheckRadarPoint,
		AlarmName:      "",
		RadarId:        radarId,
		RadarPointId:   radarPointId,
		AlarmType:      models.AlarmTypeRadarPointDeformation,
		RedOption:      "150",
		OrangeOption:   "100",
		YellowOption:   "60",
		BlueOption:     "30",
		Interval:       10,
		Duration:       24,
	})
	//监测点-速度
	items = append(items, models.AlarmPoint{
		DeptId:         deptId,
		AlarmCheckType: models.AlarmCheckRadarPoint,
		AlarmName:      "",
		RadarId:        radarId,
		RadarPointId:   radarPointId,
		AlarmType:      models.AlarmTypeRadarPointVelocity,
		RedOption:      "15",
		OrangeOption:   "10",
		YellowOption:   "6",
		BlueOption:     "3",
		Interval:       10,
		Duration:       24,
	})
	//监测点-加速度
	items = append(items, models.AlarmPoint{
		DeptId:         deptId,
		AlarmCheckType: models.AlarmCheckRadarPoint,
		AlarmName:      "",
		RadarId:        radarId,
		RadarPointId:   radarPointId,
		AlarmType:      models.AlarmTypeRadarPointAcceleration,
		RedOption:      "15",
		OrangeOption:   "10",
		YellowOption:   "6",
		BlueOption:     "3",
		Interval:       10,
		Duration:       24,
	})
	return items
}

// GetAlarmPointLogsPage 获得监测点的告警日志
func (e *AlarmPoint) GetAlarmPointLogsPage(c dto.GetAlarmPointLogsPageReq, list []*models.AlarmPointLogs, count *int64) error {
	var err error
	var data models.AlarmPointLogs

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	return nil
}
