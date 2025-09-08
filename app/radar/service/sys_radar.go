package service

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/monsvr/mongosvr"
	"go-admin/app/radar/models"
	"go-admin/app/radar/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type Radar struct {
	service.Service
}

// GetList 获取Radar列表
func (e *Radar) GetList(c *dto.RadarGetPageReq, p *actions.DataPermission, list *[]models.Radar, count *int64) error {
	var err error
	var data models.Radar

	err = e.Orm.Preload("Dept").Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).Count(count).Error

	if err != nil {
		e.Log.Errorf("RadarService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取Radar对象
func (e *Radar) Get(d *dto.RadarGetReq, p *actions.DataPermission, model *models.Radar) error {
	var data models.Radar

	o := e.Orm.Model(&data)
	if p != nil {
		o = o.Scopes(
			actions.Permission(data.TableName(), p),
		)
	}
	err := o.First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetRadar error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Get 获取Radar对象
func (e *Radar) GetByKey(d *dto.RadarKeyGetReq, p *actions.DataPermission, model *models.Radar) error {
	var data models.Radar
	var err error
	db := e.Orm.Model(&data)
	if p != nil {
		db = db.Scopes(
			actions.Permission(data.TableName(), p),
		)
	}
	if err = db.Where("radar_key = ?", d.GetRadarKey()).First(model).Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.Log.Errorf("Service GetRadar error:%s \n", errors.New("查看对象不存在或无权查看"))
		}
		return err
	}
	return nil
}

// Insert 创建Radar对象
func (e *Radar) Insert(c *dto.RadarInsertReq) error {
	var err error
	var data models.Radar
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("RadarService Insert error:%s \r\n", err)
		return err
	}
	c.RadarId = data.RadarId
	return nil
}

// Update 修改Radar对象
func (e *Radar) Update(c *dto.RadarUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.Radar{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("RadarService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除Radar
func (e *Radar) Remove(d *dto.RadarDeleteReq, p *actions.DataPermission) error {
	var data models.Radar

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveRadar error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

// Get 获取SysDept对象
func (e *Radar) GetById(id int, model *models.Radar) error {
	// if e.Orm == nil {
	// 	e.Orm = getFirstOrm()
	// }
	d := dto.RadarGetReq{}
	d.RadarId = int64(id)
	return e.Get(&d, nil, model)
}

// Get 获取Radar对象
func (e *Radar) GetByRadarKey(key string, model *models.Radar) error {
	// if e.Orm == nil {
	// 	e.Orm = getFirstOrm()
	// }
	d := dto.RadarKeyGetReq{}
	d.RadarKey = key
	return e.GetByKey(&d, nil, model)
}

// func (e *Radar) ConfirmProjetInfo(proId uint32, devId uint32, lng float64, lat float64) error {
// 	ds := SysDept{}
// 	ds.Service = e.Service
// 	dept := models.SysDept{}
// 	if ds.GetOrCreateById(int(proId), &dept) != nil { //如果机构不存在 ，则挂在根节点之下等待去用户修改
// 		proId = 1 //根节点为1
// 	}
// 	flng := strconv.FormatFloat(lng, 'f', -1, 64)
// 	flat := strconv.FormatFloat(lat, 'f', -1, 64)
// 	radar := models.Radar{}
// 	if e.GetById(int(devId), &radar) != nil { //如果雷达不存在 ，则新增
// 		var data models.Radar
// 		data.Lng = flng
// 		data.Lat = flat
// 		data.RadarId = int64(devId)
// 		data.DeptId = int64(proId)
// 		data.FromProject = 1
// 		err := e.Orm.Create(&data).Error
// 		if err != nil {
// 			e.Log.Errorf("ConfirmProjetInfo Insert error:%s \r\n", err)
// 			return err
// 		}
// 	} else if radar.DeptId != int64(proId) || radar.Lng != flng || radar.Lat != flat { //存在了且所在单位不一致，则更新
// 		radar.DeptId = int64(proId)
// 		radar.Lng = flng
// 		radar.Lat = flat
// 		db := e.Orm.Save(&radar)
// 		if err := db.Error; err != nil {
// 			e.Log.Errorf("ConfirmProjetInfo Save error:%s \r\n", err)
// 			return err
// 		}
// 	}
// 	return nil
// }

// Get 获取Radar对象
// func (e *Radar) GetImage(d *dto.RadarGetImageReq) (*mongosvr.DistanceData, error) {

// 	return mongosvr.GetLatestDistanceData(d.RadarId)
// }

// Get 获取Radar对象
func (e *Radar) GetImageV2(d *dto.RadarGetImageReq) (*mongosvr.DistanceDataV2, error) {

	return mongosvr.GetLatestDistanceDataV2(d.RadarId)
}
