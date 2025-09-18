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
	api := apis.AlarmPoint{}
	r := v1.Group("/alarm").Use(authMiddleware.MiddlewareFunc())
	{
		r.GET("/getPointAlarmRules", api.GetAlarmRules) //获取监测点预警规则
		//	r.GET("/addPointAlarmRule", api.AddAlarmRule)        //增加预警规则
		r.POST("/updatePointAlarmRules", api.UpdateAlarmRule) //修改预警规则
		//r.POST("/deletePointAlarmRule", api.DeleteAlarmRule)  //删除预警规则
		r.POST("/getAlarmPointLogsPage", api.GetAlarmPointLogsPage) //获得监测点告警日志

		// r.POST("/get_alarmsofids", api.GetAlarmsOfIds)
		// r.POST("/get_dev_info", api.GetDevInfo)
		// r.POST("/get_state_info", api.GetStateInfo)
		// r.POST("/get_alarms_before", api.GetAlarmsBefore)
	}
}
