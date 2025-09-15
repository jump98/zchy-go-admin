package apis

import (
	"errors"
	"fmt"
	"go-admin/app/radar/models"
	"go-admin/app/radar/service"
	"go-admin/app/radar/service/dto"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/utils"
)

type AlarmPoint struct {
	api.Api
}

// GetAlarmRules 获取预警规则
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
	var alarmRuleList []*models.AlarmPoint
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

// AddAlarmRule 增加预警规则
func (e AlarmPoint) AddAlarmRule(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			e.Logger.Error("增加预警规则出错:", err)
		}
	}()

	req := dto.AddAlarmPointReq{}
	s := service.AlarmPoint{}
	if err = e.MakeContext(c).MakeOrm().Bind(&req).MakeService(&s.Service).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	levelItems := req.Items
	if len(levelItems) != 3 {
		// e.Error(400, errors.New("参数错误"), "参数错误")
		utils.ParameterError("参数错误")
		return
	}
	//var levels = make([]models.AlarmLevel, 4)
	//for i, item := range levelItems {
	//	levels[i] = item.AlarmLevel
	//}
	//if !checkAlarmLevel(levels) {
	//	utils.ParameterError("预警等级参数错误")
	//	return
	//}

	// 之后可以在这里验证一下权限
	// deptId := user.GetDeptId(c)

	// var alarmRuleList []models.AlarmPoint
	// var alarmRuleLevelList []models.AlarmRuleLevel
	// deptId := req.DeptId
	if err = s.AddAlarmRule(req); err != nil {
		e.Error(500, err, "获取警报规则失败")
		return
	}

	resp := dto.AddAlarmPointResp{
		// AlarmRuleList:      alarmRuleList,
		// AlarmRuleLevelList: alarmRuleLevelList,
	}

	e.OK(resp, "success")
}

// UpdateAlarmRule 修改预警规则
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

// DeleteAlarmRule 删除预警规则
func (e AlarmPoint) DeleteAlarmRule(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			e.Logger.Error("删除预警规则出错:", err)
		}
	}()

	req := dto.DeleteAlarmPointReq{}
	s := service.AlarmPoint{}
	if err = e.MakeContext(c).MakeOrm().Bind(&req).MakeService(&s.Service).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return

	}

	alarmRuleId := req.AlarmRuleId

	// 之后可以在这里验证一下权限
	if err = s.DeleteAlarmRule(alarmRuleId); err != nil {
		e.Error(500, err, "获取警报规则失败")
		return
	}

	resp := dto.DeleteAlarmPointResp{}
	e.OK(resp, "success")
}

func checkAlarmLevel(levels []models.AlarmLevel) bool {
	var rightLevels = []models.AlarmLevel{1, 2, 3, 4}
	return !slices.Equal(levels, rightLevels)
}
