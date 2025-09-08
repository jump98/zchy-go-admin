package service

import (
	"errors"
	"fmt"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/radar/models"
	"go-admin/app/radar/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type RadarPoint struct {
	service.Service
}

// GetPage 获取RadarPoint列表
func (e *RadarPoint) GetPage(c *dto.GetRadarPointListDeptIdReq, p *actions.DataPermission, list *[]models.RadarPoint, count *int64) error {
	var err error
	var data models.RadarPoint

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("RadarPointService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// GetPage 获取RadarPoint列表
func (e *RadarPoint) GetDeptPage(c *dto.GetRadarPointListDeptIdReq, deptid int, p *actions.DataPermission, list *[]models.RadarPoint, count *int64) error {
	var err error
	var data models.RadarPoint

	err = e.Orm.Debug().Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).Where("radar_id in(select radar_id from sys_radar where dept_id=?)", deptid).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("RadarPointService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取RadarPoint对象
func (e *RadarPoint) Get(d *dto.GetRadarPointByIdReq, p *actions.DataPermission, model *models.RadarPoint) error {
	var data models.RadarPoint

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetRadarPoint error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建RadarPoint对象
func (e *RadarPoint) Insert(c *dto.InsertRadarPointReq) error {
	var err error
	var data models.RadarPoint
	c.Generate(&data)
	fmt.Printf("data:%+v \n", data)
	fmt.Println("data.Lat", data.Lat)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("RadarPointService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改RadarPoint对象
func (e *RadarPoint) Update(c *dto.UpdateRadarPointReq, p *actions.DataPermission) error {
	var err error
	var data = models.RadarPoint{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("RadarPointService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除RadarPoint
func (e *RadarPoint) Remove(d *dto.DeleteRadarPointReq, p *actions.DataPermission) error {
	var data models.RadarPoint

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveRadarPoint error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// RemoveRadarPoint 删除指定雷达ID的监测点
func (e *RadarPoint) RemoveRadarPoint(radarId int64, p *actions.DataPermission) error {
	var data models.RadarPoint

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Where("radar_id = ?", radarId).
		Delete(&data)
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveRadarPoint error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据或未找到相关监测点")
	}
	return nil
}

// GetPointByRadarId 根据雷达ID获取监测点列表
func (e *RadarPoint) GetPointByRadarId(radarId int64, p *actions.DataPermission) ([]models.RadarPoint, error) {

	var data models.RadarPoint
	var list []models.RadarPoint

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		Where("radar_id = ?", radarId).
		Find(&list).Error
	if err != nil {
		e.Log.Errorf("Service GetPointByRadarId error:%s \r\n", err)
		return nil, err
	}
	return list, nil
}

// GetRadarIdByPointId 根据监测点ID获取对应的雷达ID
func (e *RadarPoint) GetRadarIdByPointId(pointId int, p *actions.DataPermission) (int64, error) {
	var radarId int64
	if err := e.Orm.Model(&models.RadarPoint{}).Select("radar_id").Where("id = ?", pointId).First(&radarId).Error; err != nil {
		return 0, err
	}
	return radarId, nil
}

// 获得指定雷达是所有监测点
func (e *RadarPoint) GetPointsByRadarId(radarId int64) ([]int64, error) {
	pointIndexs := make([]int64, 0)
	var err error
	if err = e.Orm.Model(&models.RadarPoint{}).Select("point_index").Where("radar_id = ?", radarId).Scan(&pointIndexs).Error; err != nil {
		e.Log.Errorf("Service GetPointsByRadarId error:%s \r\n", err)
	}
	e.Log.Info("查询雷达点位列表:", pointIndexs)
	return pointIndexs, err
}
