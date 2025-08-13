package apis

import (
	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
)

type SysCommon struct {
	api.Api
}

func (e SysCommon) IsSuperAdmin(c *gin.Context) (bool, *models.SysUser, error) {
	u := service.SysUser{}
	p := actions.GetPermissionFromContext(c)
	requser := dto.SysUserById{}
	sysUser := models.SysUser{}
	requser.Id = user.GetUserId(c)
	r := service.SysRole{}

	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&u.Service).
		MakeService(&r.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return false, nil, err
	}
	err = u.Get(&requser, p, &sysUser)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return false, &sysUser, nil
	}
	reqr := dto.SysRoleGetReq{}
	ur := models.SysRole{}
	reqr.Id = sysUser.RoleId
	err = r.Get(&reqr, &ur)
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return false, &sysUser, nil
	}
	if (ur.RoleKey == "admin" && ur.RoleName == "系统管理员") ||
		(ur.RoleKey == "超级管理员" && ur.RoleName == "超级管理员") {
		return true, &sysUser, nil
	}
	return false, &sysUser, nil
}
