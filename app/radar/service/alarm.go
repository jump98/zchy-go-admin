package service

import (
	"errors"
	"fmt"
	"go-admin/app/radar/models"
	"go-admin/app/radar/service/dto"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"github.com/samber/lo"
)

type Alarm struct {
	service.Service
}

// 获取所有的预警规则
func (e *Alarm) GetAlarmRules(deptId int64) ([]models.AlarmRule, []models.AlarmRuleLevel, error) {
	var err error
	alarmRuleList := make([]models.AlarmRule, 0)
	alarmRuleLevelList := make([]models.AlarmRuleLevel, 0)
	if err = e.Orm.Model(&models.AlarmRule{}).Where("dept_id = ?", deptId).Find(&alarmRuleList).Error; err != nil {
		return nil, nil, err
	}
	if err = e.Orm.Model(&models.AlarmRuleLevel{}).Where("dept_id = ?", deptId).Find(&alarmRuleLevelList).Error; err != nil {
		return nil, nil, err
	}
	return alarmRuleList, alarmRuleLevelList, nil
}

// 增加预警规则
func (e *Alarm) AddAlarmRule(req dto.AddAlarmRuleReq) error {
	var errTx error
	tx := e.Orm.Begin()
	defer func() {
		if errTx != nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	deptId := req.DeptId
	alarmRuleItem := &models.AlarmRule{
		DeptId: deptId,
	}
	alarmRuleLevelItems := make([]*models.AlarmRuleLevel, 4)
	for i, item := range req.AlarmRuleLevelItem {
		alarmRuleLevelItems[i].DeptId = deptId
		alarmRuleLevelItems[i].AlarmRuleId = alarmRuleItem.Id
		alarmRuleLevelItems[i].AlarmLevel = item.AlarmLevel
		alarmRuleLevelItems[i].Option.Data = item.Option
		alarmRuleLevelItems[i].OptionMode = item.OptionMode
		alarmRuleLevelItems[i].Suggestion = item.Suggestion
		alarmRuleLevelItems[i].Horn = item.Horn
	}

	if errTx = tx.Create(alarmRuleItem).Error; errTx != nil {
		return errTx
	}
	if errTx = tx.Create(alarmRuleLevelItems).Error; errTx != nil {
		return errTx
	}
	fmt.Println("打印预警ID:", alarmRuleItem.Id)
	for _, item := range alarmRuleLevelItems {
		fmt.Println("打印预警LevelID:", item.Id)

	}
	return errTx
}

// 修改预警规则
func (e *Alarm) UpdateAlarmRule(req dto.UpdateAlarmRuleReq) error {
	ruleId := req.AlarmRuleId
	alarmRuleItem := &models.AlarmRule{
		Id: req.AlarmRuleId,
	}

	alarmRuleLevelItems := make([]*models.AlarmRuleLevel, 4)
	var err error
	db := e.Orm
	if err = db.First(alarmRuleItem).Error; err != nil {
		return err
	}
	if err = db.Where("rule_id = ?", ruleId).Find(alarmRuleLevelItems).Error; err != nil {
		return err
	}

	alarmRuleItem.AlarmName = req.AlarmName
	alarmRuleItem.Remark = req.Remark
	for i, item := range alarmRuleLevelItems {
		newItem, ok := lo.Find(req.AlarmRuleLevelItem, func(u dto.AlarmRuleLevelItem) bool {
			return u.AlarmLevel == item.AlarmLevel
		})
		if !ok {
			return errors.New("参数错误,未找到预警级别")
		}
		alarmRuleLevelItems[i].Option.Data = newItem.Option
		alarmRuleLevelItems[i].OptionMode = newItem.OptionMode
		alarmRuleLevelItems[i].Suggestion = newItem.Suggestion
		alarmRuleLevelItems[i].Horn = newItem.Horn

	}

	var errTx error
	tx := e.Orm.Begin()
	defer func() {
		if errTx != nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if errTx = tx.Save(alarmRuleItem).Error; errTx != nil {
		return errTx
	}
	if errTx = tx.Save(alarmRuleLevelItems).Error; errTx != nil {
		return errTx
	}
	return errTx
}

// 删除预警规则
func (e *Alarm) DeleteAlarmRule(alarmRuleId int64) error {
	alarmRuleLevelItems := make([]*models.AlarmRuleLevel, 4)
	var err error
	db := e.Orm
	alarmRuleItem := &models.AlarmRule{
		Id: alarmRuleId,
	}
	if err = db.First(alarmRuleItem).Error; err != nil {
		return err
	}
	if err = db.Where("rule_id = ?", alarmRuleId).Find(alarmRuleLevelItems).Error; err != nil {
		return err
	}
	var errTx error
	tx := e.Orm.Begin()
	defer func() {
		if errTx != nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	//1.删除预警规则
	if errTx = tx.Model(&models.AlarmRule{}).Where("id = ?", alarmRuleId).Delete(alarmRuleItem).Error; errTx != nil {
		return errTx
	}
	//2.删除预警规矩等级
	if errTx = tx.Model(&models.AlarmRuleLevel{}).Where("alarm_rule_id = ?", alarmRuleId).Delete(alarmRuleLevelItems).Error; errTx != nil {
		return errTx
	}
	//TODO:
	//3.删除预警告警配置
	//4.删除预警联系人组
	//5.删除预警联系人员
	return errTx
}
