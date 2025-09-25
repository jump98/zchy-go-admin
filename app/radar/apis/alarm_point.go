package apis

import (
	"errors"
	"fmt"
	"go-admin/app/radar/models"
	"go-admin/app/radar/service"
	"go-admin/app/radar/service/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
)

type AlarmPoint struct {
	api.Api
}

// GetAlarmRules 获取预警规则
// @Summary 获取预警规则
// @Description 获取预警规则
// @Tags 监测点-预警&告警&消警管理
// @Param deptId query int64 true "机构ID"
// @Param radarPointId query int64 false "监测点ID"
// @Success 200 {object} response.Response{data=dto.GetAlarmRulesResp} "成功"
// @Router /api/v1/alarm/getPointAlarmRules [get]
// @Security Bearer
func (e AlarmPoint) GetAlarmRules(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			e.Logger.Error("获取预警规则出错:", err)
		}
	}()

	req := dto.GetAlarmRulesReq{}
	s := service.AlarmPoint{}
	if err = e.MakeContext(c).MakeOrm().Bind(&req).MakeService(&s.Service).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 之后可以在这里验证一下权限
	// deptId := user.GetDeptId(c)
	fmt.Println("deptId:", req.DeptId)
	fmt.Println("radarPointId:", req.RadarPointId)
	var alarmRuleList []models.AlarmPoint
	deptId := req.DeptId
	radarPointId := req.RadarPointId

	if deptId == 0 {
		e.Error(400, err, "没有DeptId")
		return
	}

	if alarmRuleList, err = s.GetAlarmRules(deptId, radarPointId); err != nil {
		e.Error(500, err, "获取预警设定失败")
		return
	}

	resp := dto.GetAlarmRulesResp{
		AlarmRuleList: alarmRuleList,
	}

	e.OK(resp, "success")
}

// UpdateAlarmRule 修改预警规则
// @Summary 修改预警规则
// @Description 修改预警规则，支持传入机构ID、监测点ID、门限类型和规则数组
// @Tags 监测点-预警&告警&消警管理
// @Accept application/json
// @Param request body dto.UpdateAlarmPointReq true "请求参数"
// @Success 200 {object} response.Response{data=dto.AddAlarmPointResp} "成功"
// @Failure 400 {object} response.Response "请求错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/alarm/updatePointAlarmRules [post]
// @Security Bearer
func (e AlarmPoint) UpdateAlarmRule(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			e.Logger.Error("修改预警规则出错:", err)
		}
	}()

	fmt.Println("请求修改预警规则")

	req := dto.UpdateAlarmPointReq{}
	// 绑定 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 访问数据
	fmt.Printf("Items:%+v \n", req.Items)

	s := service.AlarmPoint{}
	if err = e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	items := req.Items
	if len(items) != 3 {
		e.Error(http.StatusBadRequest, errors.New("参数错误"), "参数错误")
		return
	}
	deptId := req.DeptId
	radarPointId := req.RadarPointId
	mode := req.Mode

	// 之后可以在这里验证一下权限
	if err = s.UpdateAlarmRule(items, deptId, radarPointId, mode); err != nil {
		e.Error(500, err, "修改预警规则失败")
		return
	}

	resp := dto.AddAlarmPointResp{
		Success: true,
	}
	e.OK(resp, "success")
}

// @Param beginTime query string false "开始时间：2025-01-01 00:00:00"
// @Param endTime query string false "结束时间：2025-01-01 00:00:00"
// @Param alarmType query int false "预警类型 (100=累计水平位移, 101=水平位移速度, 102=水平位移加速度)" Enums(100,101,102)
// @Param alarmLevel query int false "报警等级 (0=蓝,1=黄,2=橙,3=红)" Enums(0,1,2,3)

// GetAlarmPointLogsPage 获得监测点告警日志
// @Summary 获得监测点告警日志
// @Description 分页查询监测点的告警日志，可以根据机构ID、雷达ID、监测点ID、预警类型、报警等级和时间范围过滤
// @Tags 监测点-预警&告警&消警管理
// @Accept application/json
// @Param pageIndex query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param deptId query int64 true "机构ID"
// @Param radarId query int64 false "雷达ID"
// @Param alarmType query int false "预警类型 (100=累计水平位移, 101=水平位移速度, 102=水平位移加速度)"
// @Param alarmLevel query int false "报警等级 (1=蓝,1=黄,2=橙,3=红)"
// @Param radarPointId query int64 false "监测点ID"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.AlarmPointLogs}} "成功"
// @Router /api/v1/alarm/getAlarmPointLogsPage [get]
// @Security Bearer
func (e AlarmPoint) GetAlarmPointLogsPage(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			e.Logger.Error("获得监测点告警日志出错:", err)
		}
	}()
	fmt.Println("获得监测点告警日志")

	req := dto.GetAlarmPointLogsPageReq{}
	s := service.AlarmPoint{}
	if err = e.MakeContext(c).MakeOrm().Bind(&req, binding.Form).MakeService(&s.Service).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	var list []models.AlarmPointLogs
	var count int64
	if err = s.GetAlarmPointLogsPage(req, &list, &count); err != nil {
		e.Error(500, err, err.Error())
		return
	}
	e.PageOK(list, int(count), req.PageIndex, req.PageIndex, "success")
}

// CloseAlarmPointById 关闭告警
// @Summary 关闭告警
// @Description 关闭告警，根据监测点ID条件
// @Tags 监测点-预警&告警&消警管理
// @Accept application/json
// @Param req body dto.CloseAlarmPointByIdReq true "请求参数"
// @Success 200 {object} response.Response{data=dto.CloseAlarmPointByIdResp} "成功"
// @Router /api/v1/alarm/closeAlarmPointById [post]
// @Security Bearer
func (e AlarmPoint) CloseAlarmPointById(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			e.Logger.Error("关闭告警出错:", err)
		}
	}()
	fmt.Println("关闭告警")
	req := dto.CloseAlarmPointByIdReq{}
	s := service.AlarmPoint{}
	if err = e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).MakeService(&s.Service).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	fmt.Printf("打印请求参数：%+v \n", req)

	radarPointId := req.RadarPointId
	remark := req.ProcessRemark
	userId := user.GetUserId(c)

	var ids []int64
	if ids, err = s.CloseAlarmPointById(radarPointId, int64(userId), remark); err != nil {
		e.Error(500, err, err.Error())
		return
	}

	resp := dto.CloseAlarmPointByIdResp{
		Ids: ids,
	}
	e.OK(resp, "success")
}
