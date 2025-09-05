package apis

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/app/monsvr/mongosvr"
	"go-admin/common/actions"
)

type RadarPoint struct {
	api.Api
}

// GetPage 获取监测点管理列表
// @Summary 获取监测点管理列表
// @Description 获取监测点管理列表
// @Tags 监测点管理
// @Param id query int false "PointID"
// @Param pointName query string false "监测点名称"
// @Param pointKey query string false "监测点编号"
// @Param radarId query int64 false "雷达ID"
// @Param aStatus query string false "激活状态"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.RadarPoint}} "{"code": 200, "data": [...]}"
// @Router /api/v1/radar-point [get]
// @Security Bearer
func (e RadarPoint) GetPage(c *gin.Context) {
	req := dto.RadarPointGetPageReq{}
	s := service.RadarPoint{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.RadarPoint, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取监测点管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetDeptPage 获取某个部门下的所有监测点列表
// @Summary 获取某个部门下的所有监测点列表
// @Description 获取某个部门下的所有监测点列表
// @Tags 获取某个部门下的所有监测点列表
// @Param id query int false "PointID"
// @Param pointName query string false "监测点名称"
// @Param pointKey query string false "监测点编号"
// @Param radarId query int64 false "雷达ID"
// @Param aStatus query string false "激活状态"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.RadarPoint}} "{"code": 200, "data": [...]}"
// @Router /api/v1/radar-point/dept [get]
// @Security Bearer
func (e RadarPoint) GetDeptPage(c *gin.Context) {
	req := dto.RadarPointGetPageReq{}
	s := service.RadarPoint{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.RadarPoint, 0)
	var count int64

	sc := SysCommon{}
	admin, _, err := sc.IsSuperAdmin(c)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if admin {
		err = s.GetPage(&req, p, &list, &count)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("获取监测点管理失败，\r\n失败信息 %s", err.Error()))
			return
		}
	} else {
		// data, _ := c.Get(jwtauth.JwtPayloadKey)
		// fmt.Println(data)
		// v := data.(jwtauth.MapClaims)
		// deptid, err := v.Int("deptId")
		deptid := user.GetDeptId(c)
		err = s.GetDeptPage(&req, deptid, p, &list, &count)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("获取监测点管理失败，\r\n失败信息 %s", err.Error()))
			return
		}
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取监测点管理
// @Summary 获取监测点管理
// @Description 获取监测点管理
// @Tags 监测点管理
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.RadarPoint} "{"code": 200, "data": [...]}"
// @Router /api/v1/radar-point/{id} [get]
// @Security Bearer
func (e RadarPoint) Get(c *gin.Context) {
	req := dto.RadarPointGetReq{}
	s := service.RadarPoint{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models.RadarPoint

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取监测点管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建监测点管理
// @Summary 创建监测点管理
// @Description 创建监测点管理
// @Tags 监测点管理
// @Accept application/json
// @Product application/json
// @Param data body dto.RadarPointInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/radar-point [post]
// @Security Bearer
func (e RadarPoint) Insert(c *gin.Context) {
	req := dto.RadarPointInsertReq{}
	s := service.RadarPoint{}
	var err error
	if err := e.MakeContext(c).MakeOrm().Bind(&req).MakeService(&s.Service).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	fmt.Printf("创建监测点管理参数：%+v \n", req)
	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, "创建监测点管理失败")
		return
	}
	// var points []int64
	// if points, err = s.GetPointsByRadarId(req.RadarId); err != nil {
	// 	e.Error(500, err, "获得监测点列表出错")
	// 	return
	// }
	// param := []dto.RadarPointIndex{}

	param := []dto.RadarPointIndex{{Index: req.PointIndex}}
	err = mongosvr.InsertCommandData(&mongosvr.CommandData{
		RadarId:     req.RadarId,
		CommandCode: mongosvr.CMD_RD_ADDPOINT,
		Message:     "add monitor point",
		TimeStamp:   time.Now().Unix(),
		Parameters:  map[string]interface{}{"polygon": param},
	})
	if err != nil {
		e.Error(500, err, "创建监测点管理失败")
		return
	}

	e.OK(req.GetId(), "创建成功")
}

func (e RadarPoint) convertToMapSlice(points []models.RadarPoint) []dto.RadarPointIndex {
	var mapSlice []dto.RadarPointIndex
	for _, point := range points {
		mapSlice = append(mapSlice, dto.RadarPointIndex{
			Index: point.PointIndex,
		})
	}
	return mapSlice
}

// Update 修改监测点管理
// @Summary 修改监测点管理
// @Description 修改监测点管理
// @Tags 监测点管理
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.RadarPointUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/radar-point/{id} [put]
// @Security Bearer
func (e RadarPoint) Update(c *gin.Context) {
	req := dto.RadarPointUpdateReq{}
	s := service.RadarPoint{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改监测点管理失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除监测点管理
// @Summary 删除监测点管理
// @Description 删除监测点管理
// @Tags 监测点管理
// @Param data body dto.RadarPointDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/radar-point [delete]
// @Security Bearer
func (e RadarPoint) Delete(c *gin.Context) {
	s := service.RadarPoint{}
	req := dto.RadarPointDeleteReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	radarId, points, _ := e.getRadarIDandPoints(req.Ids, &s, p)
	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除监测点管理失败，\r\n失败信息 %s", err.Error()))
		return
	}
	err = mongosvr.InsertCommandData(&mongosvr.CommandData{
		RadarId:     radarId,
		CommandCode: mongosvr.CMD_RD_DELETEPOINT,
		Message:     "remove monitor point",
		TimeStamp:   time.Now().Unix(),
		Parameters:  map[string]interface{}{"polygon": points},
	})
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除监测点管理失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

func (e RadarPoint) getRadarIDandPoints(ids []int, s *service.RadarPoint, p *actions.DataPermission) (int64, []dto.RadarPointIndex, error) {
	var points []dto.RadarPointIndex
	var radarId int64 = 0
	for _, id := range ids {
		if radarId == 0 {
			radarId, _ = s.GetRadarIdByPointId(id, p)
		}
		req := dto.RadarPointGetReq{Id: id}
		object := &models.RadarPoint{}
		err := s.Get(&req, p, object)
		if err == nil {
			points = append(points, dto.RadarPointIndex{
				Index: object.PointIndex,
			})
		}
	}
	return radarId, points, nil
}

// GetDeformationData 获取变形点数据
// @Summary 获取变形点数据
// @Description 根据设备ID、索引和时间范围获取采样后的变形点数据
// @Tags 监测点管理
// @Accept application/json
// @Product application/json
// @Param data body dto.DeformationPointQueryReq true "变形点数据查询参数"
// @Success 200 {object} response.Response{data=[]mongosvr.DeformationPointData} "成功"
// @Router /api/v1/radar-point/deformation-data [post]
// @Security Bearer
func (e RadarPoint) GetDeformationData(c *gin.Context) {
	req := dto.DeformationPointQueryReq{}
	var err error
	if err = e.MakeContext(c).MakeOrm().Bind(&req).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// 参数验证
	if req.Devid <= 0 || req.Index < 0 || req.StartTime == "" || req.EndTime == "" {
		e.Error(400, nil, "参数不完整")
		return
	}
	hours := req.Hours
	radarId := req.Devid
	pointIndex := req.Index
	startTime := req.StartTime
	endTime := req.EndTime
	timeType := req.TimeType
	if hours < 0 {
		hours = 0
	}

	// 获取数据
	var lastTime time.Time
	var data []mongosvr.DeformationPointData
	if data, lastTime, err = mongosvr.QueryDeformationPointData(radarId, pointIndex, startTime, endTime, hours, timeType); err != nil {
		e.Error(500, err, fmt.Sprintf("获取变形点数据失败: %s", err.Error()))
		return
	}

	resp := dto.DeformationPointQueryResp{}
	resp.LastTime = lastTime
	resp.List = data

	// 返回采样后的数据
	e.OK(resp, "查询成功")
}
