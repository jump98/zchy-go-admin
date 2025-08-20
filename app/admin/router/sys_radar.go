package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/admin/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysRadarRouter)
}

// registerSysRadarRouter
func registerSysRadarRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SysRadar{}
	r := v1.Group("/sys-radar").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.GET("", actions.PermissionAction(), api.GetList)
		r.GET("/:radarId", actions.PermissionAction(), api.Get)
		r.GET("/radarimage/:radarId", actions.PermissionAction(), api.GetRadarImage)
		r.POST("", api.Insert)
		r.PUT("/:radarId", actions.PermissionAction(), api.Update)
		r.DELETE("", api.Delete)
		//r.PUT("/confirm/:radarId", actions.PermissionAction(), api.Confirm)
		r.POST("/get_alarms", api.GetAlarms)
		r.POST("/get_alarmsofids", api.GetAlarmsOfIds)
		r.POST("/get_dev_info", api.GetDevInfo)
		r.POST("/get_state_info", api.GetStateInfo)
		r.POST("/get_alarms_before", api.GetAlarmsBefore)
	}

	// 雷达设备认证路由
	noAuth := v1.Group("/radar")
	{
		noAuth.POST("/authenticate", api.Authenticate)
	}

	// 雷达设备告警信息路由（需要认证）
	radarAuth := v1.Group("/radar").Use(authMiddleware.MiddlewareFunc()) //.Use(middleware.AuthCheckRole())
	{
		radarAuth.POST("/put_alarm", api.PutAlarm)
		radarAuth.POST("/put_deformation", api.PutDeformation)
		radarAuth.POST("/get_commands", api.GetCommands)
		radarAuth.POST("/raw_data", api.PutRawData)
		//radarAuth.POST("/put_rebootcommand", api.PutRebootCommand)
		radarAuth.POST("/puttestcommand", api.PutTestCommand)
		radarAuth.POST("/dev_reboot", api.PutRebootCommand)
		radarAuth.POST("/status", api.PutStatus)
		// 在radarAuth路由组中添加
		radarAuth.POST("/dev_info", api.PutDevInfo)
	}
}
