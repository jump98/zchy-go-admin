package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

	"go-admin/app/radar/apis"
	"go-admin/common/actions"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerRadarPointRouter)
}

// registerRadarPointRouter
func registerRadarPointRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.RadarPoint{}
	r := v1.Group("/radar_point").Use(authMiddleware.MiddlewareFunc())
	//形变点管理
	{
		r.GET("getList", actions.PermissionAction(), api.GetRadarPointList)                     //获取检测点列表
		r.GET("/getById/:id", actions.PermissionAction(), api.GetRadarPointById)                //获得监测点ById
		r.POST("/add", api.InsertRadarPoint)                                                    //插入监测点
		r.PUT("/update/:id", actions.PermissionAction(), api.UpdateRadarPoint)                  //修改监测点
		r.DELETE("/delete", api.DeleteRadarPoint)                                               //删除监测点
		r.GET("/getPointListByDeptId", actions.PermissionAction(), api.GetRadarPointListDeptId) //获取某个部门下的所有监测点列表
	}
	//形变数据
	{
		r.POST("/getDeformationData", actions.PermissionAction(), api.GetDeformationData)                 //获取形变数据
		r.POST("/getDeformationVelocity", actions.PermissionAction(), api.GetDeformationVelocity)         //获取形变速度
		r.POST("/getDeformationAcceleration", actions.PermissionAction(), api.GetDeformationAcceleration) //获取形变加速度
	}
}
