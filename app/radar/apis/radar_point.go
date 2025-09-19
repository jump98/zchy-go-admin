package apis

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/monsvr/mongosvr"
	"go-admin/app/radar/models"
	"go-admin/app/radar/service"
	"go-admin/app/radar/service/dto"
	"go-admin/common/actions"
)

type RadarPoint struct {
	api.Api
}

// GetRadarPointList 获取监测点管理列表
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
// @Router /api/v1/radar_point [get]
// @Security Bearer
func (e RadarPoint) GetRadarPointList(c *gin.Context) {
	req := dto.GetRadarPointListReq{}
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

// GetRadarPointListDeptId 获取某个部门下的所有监测点列表
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
// @Router /api/v1/radar_point/dept [get]
// @Security Bearer
func (e RadarPoint) GetRadarPointListDeptId(c *gin.Context) {
	req := dto.GetRadarPointListDeptIdReq{}
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

	deptId := req.DeptId
	fmt.Println("根据deptId查询监测点信息：", deptId)
	p := actions.GetPermissionFromContext(c)
	list := make([]models.RadarPoint, 0)
	var count int64

	//sc := adminApi.SysCommon{}
	//admin, _, err := sc.IsSuperAdmin(c)
	//if err != nil {
	//	e.Error(500, err, "查询失败")
	//	return
	//}
	//fmt.Println("admin:", admin)
	if err = s.GetPointListByDeptId(&req, deptId, p, &list, &count); err != nil {
		e.Error(500, err, fmt.Sprintf("获取监测点管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	//if admin {
	//	err = s.GetPage(&req, p, &list, &count)
	//	if err != nil {
	//		e.Error(500, err, fmt.Sprintf("获取监测点管理失败，\r\n失败信息 %s", err.Error()))
	//		return
	//	}
	//} else {
	//	// data, _ := c.Get(jwtauth.JwtPayloadKey)
	//	// fmt.Println(data)
	//	// v := data.(jwtauth.MapClaims)
	//	// deptid, err := v.Int("deptId")
	//	deptid := user.GetDeptId(c)
	//	err = s.GetPointListByDeptId(&req, deptid, p, &list, &count)
	//	if err != nil {
	//		e.Error(500, err, fmt.Sprintf("获取监测点管理失败，\r\n失败信息 %s", err.Error()))
	//		return
	//	}
	//}
	fmt.Println("list:", list)
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// GetRadarPointById 获取监测点管理
// @Summary 获取监测点管理
// @Description 获取监测点管理
// @Tags 监测点管理
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.RadarPoint} "{"code": 200, "data": [...]}"
// @Router /api/v1/radar_point/{id} [get]
// @Security Bearer
func (e RadarPoint) GetRadarPointById(c *gin.Context) {
	req := dto.GetRadarPointByIdReq{}
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
	var radarPointItem models.RadarPoint
	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &radarPointItem)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取监测点管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(radarPointItem, "查询成功")
}

// InsertRadarPoint 创建监测点管理
// @Summary 创建监测点管理
// @Description 创建监测点管理
// @Tags 监测点管理
// @Accept application/json
// @Product application/json
// @Param data body dto.InsertRadarPoint true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/radar_point [post]
// @Security Bearer
func (e RadarPoint) InsertRadarPoint(c *gin.Context) {
	req := dto.InsertRadarPointReq{}
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

	param := []dto.GetRadarPointsIndex{{Position: req.PointIndex, PhaseDepth: req.PhaseDepth, PoseDepth: req.PoseDepth}}
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

//func (e RadarPoint) convertToMapSlice(points []models.RadarPoint) []dto.RadarPointIndex {
//	var mapSlice []dto.RadarPointIndex
//	for _, point := range points {
//		mapSlice = append(mapSlice, dto.RadarPointIndex{
//			Index: point.PointIndex,
//		})
//	}
//	return mapSlice
//}

// UpdateRadarPoint 修改监测点管理
// @Summary 修改监测点管理
// @Description 修改监测点管理
// @Tags 监测点管理
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.UpdateRadarPointReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/radar_point/{id} [put]
// @Security Bearer
func (e RadarPoint) UpdateRadarPoint(c *gin.Context) {
	req := dto.UpdateRadarPointReq{}
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

// DeleteRadarPoint 删除监测点管理
// @Summary 删除监测点管理
// @Description 删除监测点管理
// @Tags 监测点管理
// @Param data body dto.DeleteRadarPointReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/radar_point [delete]
// @Security Bearer
func (e RadarPoint) DeleteRadarPoint(c *gin.Context) {
	s := service.RadarPoint{}
	req := dto.DeleteRadarPointReq{}
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
	if err = s.Remove(&req, p); err != nil {
		e.Error(500, err, fmt.Sprintf("删除监测点管理失败，\r\n失败信息 %s", err.Error()))
		return
	}
	radarId, points, _ := e.getRadarIDandPoints(req.Ids, &s, p)
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

func (e RadarPoint) getRadarIDandPoints(ids []int, s *service.RadarPoint, p *actions.DataPermission) (int64, []dto.GetRadarPointsIndex, error) {
	var points []dto.GetRadarPointsIndex
	var radarId int64 = 0
	for _, id := range ids {
		if radarId == 0 {
			radarId, _ = s.GetRadarIdByPointId(id)
		}
		req := dto.GetRadarPointByIdReq{Id: id}
		object := &models.RadarPoint{}
		err := s.Get(&req, p, object)
		if err == nil {
			points = append(points, dto.GetRadarPointsIndex{
				Position:   object.PointIndex,
				PhaseDepth: object.PhaseDepth,
				PoseDepth:  object.PoseDepth,
			})
		}
	}
	return radarId, points, nil
}

// GetDeformationData 获取形变点数据
// @Summary 获取形变点数据
// @Description 根据设备ID、索引和时间范围获取采样后的形变点数据
// @Tags 监测点管理
// @Accept application/json
// @Product application/json
// @Param data body dto.GetDeformationDataReq true "形变点数据查询参数"
// @Success 200 {object} response.Response{data=[]dto.GetDeformationDataResp} "成功"
// @Router /api/v1/radar_point/deformation_data [post]
// @Security Bearer
func (e RadarPoint) GetDeformationData(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			e.Logger.Error("GetDeformationData panic:", r)
		}
	}()

	req := dto.GetDeformationDataReq{}
	var err error
	var s = service.DeformationPoint{}
	if err = e.MakeContext(c).MakeOrm().MakeService(&s.Service).Bind(&req).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 参数验证
	if req.RadarId <= 0 || req.Index < 0 || req.StartTime == "" || req.EndTime == "" {
		e.Error(400, nil, "参数不完整")
		return
	}

	// 获取数据
	var resp *dto.GetDeformationDataResp
	if resp, err = s.GetDeformationPoinList(c, req); err != nil {
		e.Error(500, err, fmt.Sprintf("获取形变点数据失败: %s", err.Error()))
		return
	}
	// 返回采样后的数据
	e.OK(resp, "查询成功")
}

// GetDeformationVelocity 获取形变速度数据
// @Summary 获取形变速度数据
// @Description 根据设备ID、索引和时间范围获取采样后的形变点速度数据
// @Tags 监测点管理
// @Accept application/json
// @Product application/json
// @Param data body dto.GetDeformationVelocityReq true "形变点数据查询参数"
// @Success 200 {object} response.Response{data=[]mongosvr.DeformationPointData} "成功"
// @Router /api/v1/radar_point/deformation_data [post]
// @Security Bearer
func (e RadarPoint) GetDeformationVelocity(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			e.Logger.Error("GetDeformationData panic:", r)
		}
	}()

	req := dto.GetDeformationDataReq{}
	var err error
	var s = service.DeformationPoint{}
	if err = e.MakeContext(c).MakeOrm().MakeService(&s.Service).Bind(&req).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 参数验证
	if req.RadarId <= 0 || req.Index < 0 || req.StartTime == "" || req.EndTime == "" {
		e.Error(400, nil, "参数不完整")
		return
	}

	// 获取数据
	var resp *dto.GetDeformationVelocityResp
	if resp, err = s.GetDeformationVelocity(c, req); err != nil {
		e.Error(500, err, fmt.Sprintf("获取形变速度统计失败: %s", err.Error()))
		return
	}
	// 返回采样后的数据
	e.OK(resp, "查询成功")
}

// GetDeformationAcceleration 获取形变加速度数据
// @Summary 获取形变加速度数据
// @Description 根据设备ID、索引和时间范围获取采样后的形变点加速度数据
// @Tags 监测点管理
// @Accept application/json
// @Product application/json
// @Param data body dto.GetDeformationVelocityReq true "形变点数据查询参数"
// @Success 200 {object} response.Response{data=[]mongosvr.DeformationPointData} "成功"
// @Router /api/v1/radar_point/deformation_data [post]
// @Security Bearer
func (e RadarPoint) GetDeformationAcceleration(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			e.Logger.Error("GetDeformationData panic:", r)
		}
	}()

	req := dto.GetDeformationDataReq{}
	var err error
	var s = service.DeformationPoint{}
	if err = e.MakeContext(c).MakeOrm().MakeService(&s.Service).Bind(&req).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 参数验证
	if req.RadarId <= 0 || req.Index < 0 || req.StartTime == "" || req.EndTime == "" {
		e.Error(400, nil, "参数不完整")
		return
	}

	// 获取数据
	var resp *dto.GetDeformationAccelerationResp
	if resp, err = s.GetDeformationAcceleration(c, req); err != nil {
		e.Error(500, err, fmt.Sprintf("获取形变加速度统计失败: %s", err.Error()))
		return
	}
	e.OK(resp, "查询成功")
}

// GetDeformCurveList 获得形变曲线数据（形变曲线、速度、加速度）
func (e RadarPoint) GetDeformCurveList(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			e.Logger.Error("GetDeformCurveList panic:", r)
		}
	}()

	req := dto.GetDeformCurveListReq{}
	var err error
	var s = service.DeformationPointV2{}
	if err = e.MakeContext(c).MakeOrm().MakeService(&s.Service).Bind(&req).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 参数验证
	if req.RadarId <= 0 || len(req.Index) == 0 || req.StartTime == "" || req.EndTime == "" {
		e.Error(400, nil, "参数不完整")
		return
	}

	// 获取数据
	var resp *dto.GetDeformCurveListResp
	if resp, err = s.GetDeformCurveList(c, req); err != nil {
		e.Error(500, err, "获取形变曲线失败")
		return
	}
	// 返回采样后的数据
	e.OK(resp, "查询成功")
}
