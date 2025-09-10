package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/radar/apis"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerAlarmRouter)
}

// 预警路由
func registerAlarmRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.Alarm{}
	r := v1.Group("/alarm").Use(authMiddleware.MiddlewareFunc())
	{
		r.GET("/getAlarmRules", api.GetAlarmRules)     //获取预警规则
		r.GET("/addAlarmRule", api.AddAlarmRule)       //增加预警规则
		r.GET("/updateAlarmRule", api.UpdateAlarmRule) //修改预警规则
		r.GET("/deleteAlarmRule", api.DeleteAlarmRule) //删除预警规则
		// r.POST("/get_alarmsofids", api.GetAlarmsOfIds)
		// r.POST("/get_dev_info", api.GetDevInfo)
		// r.POST("/get_state_info", api.GetStateInfo)
		// r.POST("/get_alarms_before", api.GetAlarmsBefore)
	}
}
