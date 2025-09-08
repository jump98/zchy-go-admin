package apis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	adminApi "go-admin/app/admin/apis"
	adminModel "go-admin/app/admin/models"
	"go-admin/app/monsvr/mongosvr"
	"go-admin/app/radar/models"
	"go-admin/app/radar/service"
	"go-admin/app/radar/service/dto"
	"go-admin/common/actions"
)

type Radar struct {
	api.Api
}

// GetList 获取雷达管理列表
// @Summary 获取雷达管理列表
// @Description 获取雷达管理列表
// @Tags 雷达管理
// @Param radarId query int64 false "RadarID"
// @Param radarName query string false "雷达名称"
// @Param radarKey query string false "雷达编号"
// @Param specialKey query string false "雷达特殊编号"
// @Param deptId query int64 false "部门"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Radar}} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-radar [get]
// @Security Bearer
func (e Radar) GetList(c *gin.Context) {
	req := dto.RadarGetPageReq{}
	s := service.Radar{}
	err := e.MakeContext(c).MakeOrm().Bind(&req).MakeService(&s.Service).Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)
	list := make([]models.Radar, 0)
	var count int64
	sc := adminApi.SysCommon{}
	parentID := 0
	var admin bool
	var u *adminModel.SysUser
	if admin, u, err = sc.IsSuperAdmin(c); err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if !admin && u != nil {
		parentID = u.DeptId
	}
	if parentID != 0 && req.DeptJoin.DeptId == "" {
		req.DeptJoin.DeptId = strconv.FormatInt(int64(parentID), 10)
	}
	if err = s.GetList(&req, p, &list, &count); err != nil {
		e.Error(500, err, fmt.Sprintf("获取雷达管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取雷达管理
// @Summary 获取雷达管理
// @Description 获取雷达管理
// @Tags 雷达管理
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Radar} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-radar/{id} [get]
// @Security Bearer
func (e Radar) Get(c *gin.Context) {
	req := dto.RadarGetReq{}
	s := service.Radar{}
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
	var object models.Radar

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取雷达管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建雷达管理
// @Summary 创建雷达管理
// @Description 创建雷达管理
// @Tags 雷达管理
// @Accept application/json
// @Product application/json
// @Param data body dto.RadarInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/sys-radar [post]
// @Security Bearer
func (e Radar) Insert(c *gin.Context) {
	req := dto.RadarInsertReq{}
	s := service.Radar{}
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
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建雷达管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改雷达管理
// @Summary 修改雷达管理
// @Description 修改雷达管理
// @Tags 雷达管理
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.RadarUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/sys-radar/{id} [put]
// @Security Bearer
func (e Radar) Update(c *gin.Context) {
	req := dto.RadarUpdateReq{}
	s := service.Radar{}
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
		e.Error(500, err, fmt.Sprintf("修改雷达管理失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除雷达管理
// @Summary 删除雷达管理
// @Description 删除雷达管理
// @Tags 雷达管理
// @Param data body dto.RadarDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/sys-radar [delete]
// @Security Bearer
func (e Radar) Delete(c *gin.Context) {
	s := service.Radar{}
	req := dto.RadarDeleteReq{}
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
	spt := service.RadarPoint{}
	err = e.MakeContext(c).
		MakeOrm().
		MakeService(&spt.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	for i := 0; i < len(req.Ids); i++ {
		err = spt.RemoveRadarPoint(int64(req.Ids[i]), p)
	}

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除雷达管理失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

// // Confirm 确认雷达管理
// // @Summary 确认雷达管理
// // @Description 确认雷达管理
// // @Tags 雷达管理
// // @Accept application/json
// // @Product application/json
// // @Param id path int true "id"
// // @Param data body dto.RadarConfirmReq true "body"
// // @Success 200 {object} response.Response	"{"code": 200, "message": "确认成功"}"
// // @Router /api/v1/sys-radar/confirm/{id} [put]
// // @Security Bearer
// func (e Radar) Confirm(c *gin.Context) {
// 	req := dto.RadarConfirmReq{}
// 	s := service.Radar{}
// 	err := e.MakeContext(c).
// 		MakeOrm().
// 		Bind(&req).
// 		MakeService(&s.Service).
// 		Errors
// 	if err != nil {
// 		e.Logger.Error(err)
// 		e.Error(500, err, err.Error())
// 		return
// 	}
// 	p := actions.GetPermissionFromContext(c)

// 	err = s.Update(&req, p)
// 	if err != nil {
// 		e.Error(500, err, fmt.Sprintf("修改雷达管理失败，\r\n失败信息 %s", err.Error()))
// 		return
// 	}
// 	e.OK(req.GetId(), "修改成功")
// }

// GetRadarImage 获取雷达影像
// @Summary 获取雷达影像
// @Description 获取雷达影像
// @Tags 雷达管理
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.Radar} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-radar/radarimage/{id} [get]
// @Security Bearer
func (e Radar) GetRadarImage(c *gin.Context) {
	req := dto.RadarGetImageReq{}
	s := service.Radar{}
	var err error
	if err := e.MakeContext(c).MakeOrm().Bind(&req).MakeService(&s.Service).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	//插入获取图像命令
	err = mongosvr.InsertCommandData(&mongosvr.CommandData{
		RadarId:     req.RadarId,
		CommandCode: mongosvr.CMD_RD_GETRAWDATA,
		Message:     "get raw data",
		TimeStamp:   time.Now().Unix(),
		Parameters:  map[string]interface{}{},
	})
	err = mongosvr.InsertCommandData(&mongosvr.CommandData{
		RadarId:     req.RadarId,
		CommandCode: mongosvr.CMD_RD_GETSTATEINFO,
		Message:     "get state info",
		TimeStamp:   time.Now().Unix(),
		Parameters:  map[string]interface{}{},
	})
	err = mongosvr.InsertCommandData(&mongosvr.CommandData{
		RadarId:     req.RadarId,
		CommandCode: mongosvr.CMD_RD_GETDEVINFO,
		Message:     "get dev info",
		TimeStamp:   time.Now().Unix(),
		Parameters:  map[string]interface{}{},
	})
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除监测点管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	//直接获取是上次的，如果想获取即时的，应当sleep一下
	//p := actions.GetPermissionFromContext(c)
	image, err := s.GetImageV2(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取雷达管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(image, "查询成功")
}

// GetAlarms 获取雷达报警列表
// @Summary 获取雷达报警列表
// @Description 获取雷达报警列表
// @Tags 雷达管理
// @Param startTime query string false "StartTime"
// @Success 200 {object} response.Response{data=[]mongosvr.AlarmData} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-radar/get_alarms [post]
// @Security Bearer
func (e Radar) GetAlarms(c *gin.Context) {
	req := dto.RadarGetPageReq{}
	s := service.Radar{}
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
	list := make([]models.Radar, 0)
	var count int64
	sc := adminApi.SysCommon{}
	parentID := 0
	admin, u, err := sc.IsSuperAdmin(c)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}
	if !admin && u != nil {
		parentID = u.DeptId
	}
	if parentID != 0 && req.DeptJoin.DeptId == "" {
		req.DeptJoin.DeptId = strconv.FormatInt(int64(parentID), 10)
	}
	req.PageIndex = 1
	req.PageSize = 1000
	err = s.GetList(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取雷达失败，\r\n失败信息 %s", err.Error()))
		return
	}
	//找到所有的雷达id
	radarIDs := make([]int64, 0)
	for i := 0; i < len(list); i++ {
		radarIDs = append(radarIDs, list[i].RadarId)
	}

	// 查询每个雷达ID的最后一条告警记录
	alarms, err := mongosvr.QueryLastAlarmForRadarIDs(radarIDs)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取告警记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(alarms, "查询成功")
}

// GetAlarmsOfIds 通过ID列表获取雷达报警列表
// @Summary 通过ID列表获取雷达报警列表
// @Description 通过ID列表获取雷达报警列表
// @Tags 雷达管理
// @Accept application/json
// @Product application/json
// @Param data body dto.RadarGetAlarmsOfIdsReq true "雷达ID列表"
// @Success 200 {object} response.Response{data=[]mongosvr.AlarmData} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-radar/get_alarmsofids [post]
// @Security Bearer
func (e Radar) GetAlarmsOfIds(c *gin.Context) {
	req := dto.RadarGetAlarmsOfIdsReq{}
	s := service.Radar{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		e.Error(400, err, "请求参数错误")
		return
	}

	if len(req.GetIds()) == 0 {
		e.Error(400, err, "请求参数错误")
		return
	}

	// 查询每个雷达ID的最后一条告警记录
	alarms, err := mongosvr.QueryLastAlarmForRadarIDs(req.GetIds())
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取告警记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(alarms, "查询成功")
}

// GetDevInfo 获取雷达最新设备信息
// @Summary 获取雷达最新设备信息
// @Description 获取雷达最新设备信息
// @Tags 雷达管理
// @Accept application/json
// @Product application/json
// @Param data body dto.RadarGetReq true "data"
// @Success 200 {object} response.Response{data=mongosvr.RadarDevInfo} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-radar/get_dev_info [post]
// @Security Bearer
func (e Radar) GetDevInfo(c *gin.Context) {
	s := service.Radar{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req := dto.RadarGetReq{}
	if err = c.ShouldBindJSON(&req); err != nil {
		e.Error(400, err, "请求参数错误")
		return
	}

	p := actions.GetPermissionFromContext(c)
	var object models.Radar

	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取雷达管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	// 获取雷达最新设备信息
	devInfo, err := mongosvr.GetLatestRadarDevInfo(object.RadarId)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取雷达最新设备信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(devInfo, "查询成功")
}

// GetStateInfo 获取雷达最新状态信息
// @Summary 获取雷达最新状态信息
// @Description 获取雷达最新状态信息
// @Tags 雷达管理
// @Accept application/json
// @Product application/json
// @Param data body dto.RadarGetReq true "data"
// @Success 200 {object} response.Response{data=mongosvr.RadarStatus} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-radar/get_state_info [post]
// @Security Bearer
func (e Radar) GetStateInfo(c *gin.Context) {
	s := service.Radar{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	req := dto.RadarGetReq{}
	if err = c.ShouldBindJSON(&req); err != nil {
		e.Error(400, err, "请求参数错误")
		return
	}

	p := actions.GetPermissionFromContext(c)
	var object models.Radar

	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取雷达管理失败，\r\n失败信息 %s", err.Error()))
		return
	}

	// 获取雷达最新状态信息
	statusInfo, err := mongosvr.GetLatestRadarStatus(object.RadarId)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取雷达最新状态信息失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(statusInfo, "查询成功")
}

// GetAlarmsBefore 获取指定时间之前的雷达报警列表
// @Summary 获取指定时间之前的雷达报警列表
// @Description 获取指定时间之前的雷达报警列表
// @Tags 雷达管理
// @Accept application/json
// @Product application/json
// @Param data body dto.RadarGetAlarmsBeforeReq true "雷达ID、时间和数量"
// @Success 200 {object} response.Response{data=[]mongosvr.AlarmData} "成功"
// @Router /api/v1/sys-radar/get_alarms_before [post]
// @Security Bearer
func (e Radar) GetAlarmsBefore(c *gin.Context) {
	req := dto.RadarGetAlarmsBeforeReq{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	// 解析时间参数
	startTime, err := time.Parse("2006-01-02 15:04:05", req.GetTime())
	if err != nil {
		e.Error(500, err, "时间格式错误")
		return
	}

	// 调用queryAlarmsDataTimeBefore查询警告数据
	alarms, err := mongosvr.QueryAlarmsDataTimeBefore(req.GetRadarId(), startTime, req.GetNum())
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取告警记录失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(alarms, "查询成功")
}
