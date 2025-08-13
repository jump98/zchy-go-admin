package service

import (
	"errors"
	"strconv"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/app/monsvr/mongosvr"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type SysRadar struct {
	service.Service
}

// GetPage 获取SysRadar列表
func (e *SysRadar) GetPage(c *dto.SysRadarGetPageReq, p *actions.DataPermission, list *[]models.SysRadar, count *int64) error {
	var err error
	var data models.SysRadar

	err = e.Orm.Preload("Dept").Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("SysRadarService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取SysRadar对象
func (e *SysRadar) Get(d *dto.SysRadarGetReq, p *actions.DataPermission, model *models.SysRadar) error {
	var data models.SysRadar

	o := e.Orm.Model(&data)
	if p != nil {
		o = o.Scopes(
			actions.Permission(data.TableName(), p),
		)
	}
	err := o.First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSysRadar error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Get 获取SysRadar对象
func (e *SysRadar) GetByKey(d *dto.SysRadarKeyGetReq, p *actions.DataPermission, model *models.SysRadar) error {
	var data models.SysRadar

	o := e.Orm.Model(&data)
	if p != nil {
		o = o.Scopes(
			actions.Permission(data.TableName(), p),
		)
	}
	err := o.Where("radar_key = ?", d.GetKey()).First(model).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSysRadar error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建SysRadar对象
func (e *SysRadar) Insert(c *dto.SysRadarInsertReq) error {
	var err error
	var data models.SysRadar
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysRadarService Insert error:%s \r\n", err)
		return err
	}
	c.RadarId = data.RadarId
	return nil
}

// Update 修改SysRadar对象
func (e *SysRadar) Update(c *dto.SysRadarUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.SysRadar{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("SysRadarService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除SysRadar
func (e *SysRadar) Remove(d *dto.SysRadarDeleteReq, p *actions.DataPermission) error {
	var data models.SysRadar

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveSysRadar error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// Get 获取SysDept对象
func (e *SysRadar) GetById(id int, model *models.SysRadar) error {
	if e.Orm == nil {
		e.Orm = getFirstOrm()
	}
	d := dto.SysRadarGetReq{}
	d.RadarId = int64(id)
	return e.Get(&d, nil, model)
}

// Get 获取SysDept对象
func (e *SysRadar) GetByRadarKey(key string, model *models.SysRadar) error {
	if e.Orm == nil {
		e.Orm = getFirstOrm()
	}
	d := dto.SysRadarKeyGetReq{}
	d.RadarKey = key
	return e.GetByKey(&d, nil, model)
}

func (e *SysRadar) ConfirmProjetInfo(proId uint32, devId uint32, lng float64, lat float64) error {
	ds := SysDept{}
	ds.Service = e.Service
	dept := models.SysDept{}
	if ds.GetOrCreateById(int(proId), &dept) != nil { //如果机构不存在 ，则挂在根节点之下等待去用户修改
		proId = 1 //根节点为1
	}
	flng := strconv.FormatFloat(lng, 'f', -1, 64)
	flat := strconv.FormatFloat(lat, 'f', -1, 64)
	radar := models.SysRadar{}
	if e.GetById(int(devId), &radar) != nil { //如果雷达不存在 ，则新增
		var data models.SysRadar
		data.Lng = flng
		data.Lat = flat
		data.RadarId = int64(devId)
		data.DeptId = int64(proId)
		data.FromProject = 1
		err := e.Orm.Create(&data).Error
		if err != nil {
			e.Log.Errorf("ConfirmProjetInfo Insert error:%s \r\n", err)
			return err
		}
	} else if radar.DeptId != int64(proId) || radar.Lng != flng || radar.Lat != flat { //存在了且所在单位不一致，则更新
		radar.DeptId = int64(proId)
		radar.Lng = flng
		radar.Lat = flat
		db := e.Orm.Save(&radar)
		if err := db.Error; err != nil {
			e.Log.Errorf("ConfirmProjetInfo Save error:%s \r\n", err)
			return err
		}
	}
	return nil
}

// Get 获取SysRadar对象
func (e *SysRadar) GetImage(d *dto.SysRadarGetImageReq) (*mongosvr.DistanceData, error) {

	return mongosvr.GetLatestDistanceData(d.RadarId)
}

// Get 获取SysRadar对象
func (e *SysRadar) GetImageV2(d *dto.SysRadarGetImageReq) (*mongosvr.DistanceDataV2, error) {

	return mongosvr.GetLatestDistanceDataV2(d.RadarId)
}
