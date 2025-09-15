package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/radar/apis"
	"go-admin/common/actions"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerRadarRouter)
	routerNoCheckRole = append(routerNoCheckRole, registerRadarNotAuthRouter)
}

// registerRadarRouter
func registerRadarRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Radar{}
	r := v1.Group("/radar_info").Use(authMiddleware.MiddlewareFunc())
	// r := v1.Group("/sys-radar").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), api.GetList)
		r.GET("/:radarId", actions.PermissionAction(), api.Get)
		r.GET("/radarimage/:radarId", actions.PermissionAction(), api.GetRadarImage)
		r.POST("", api.Insert)
		r.PUT("/:radarId", actions.PermissionAction(), api.Update)
		r.DELETE("", api.Delete)
		r.POST("/getRadarListByDeptId", api.GetRadarListByDeptId)

		//TODO:需要删除
		r.POST("/get_alarms", api.GetAlarms)
		r.POST("/get_alarmsofids", api.GetAlarmsOfIds)
		r.POST("/get_dev_info", api.GetDevInfo)
		r.POST("/get_state_info", api.GetStateInfo)
		r.POST("/get_alarms_before", api.GetAlarmsBefore)
	}
}

// registerRadarRouter
func registerRadarNotAuthRouter(v1 *gin.RouterGroup) {
	api := apis.Radar{}
	// 雷达设备认证路由
	noAuth := v1.Group("/radar")
	{
		noAuth.POST("/authenticate", api.Authenticate) // 雷达设备认证
	}
	// 雷达设备告警信息路由（需要认证）
	radarAuth := v1.Group("/radar")
	{
		radarAuth.POST("/put_alarm", api.PutAlarm)             //雷达设备上传告警信息
		radarAuth.POST("/put_deformation", api.PutDeformation) //雷达设备上传形变数据
		radarAuth.POST("/get_commands", api.GetCommands)       //雷达设备获取下发命令
		radarAuth.POST("/raw_data", api.PutRawData)            //雷达设备上传距离像数据
		radarAuth.POST("/puttestcommand", api.PutTestCommand)  //测试新增一个雷达设备测试指令
		radarAuth.POST("/dev_reboot", api.PutRebootCommand)    //雷达设备重启指令
		radarAuth.POST("/status", api.PutStatus)               //雷达设备状态上报
		radarAuth.POST("/dev_info", api.PutDevInfo)            //雷达设备信息上报
		radarAuth.GET("/get_radar_points", api.GetRadarPoints) //获取指定雷达的监测点列表
		radarAuth.POST("/put_commands", api.PutCommands)       //发送命令到服务器
	}
}
