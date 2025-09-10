package models

import "go-admin/common/models"

// 权限范围值
//     value: '1',
//     label: '全部数据权限'
//     value: '2',
//     label: '自定数据权限'
//     value: '3',
//     label: '本机构数据权限'
//     value: '4',
//     label: '本机构及以下数据权限'
//     value: '5',
//     label: '仅本人数据权限'

type SysRole struct {
	RoleId    int        `json:"roleId" gorm:"primaryKey;autoIncrement"` // 角色编码
	RoleName  string     `json:"roleName" gorm:"size:128;"`              // 角色名称
	Status    string     `json:"status" gorm:"size:4;"`                  // 状态 1禁用 2正常
	RoleKey   string     `json:"roleKey" gorm:"size:128;"`               //角色代码
	RoleSort  int        `json:"roleSort" gorm:""`                       //角色排序
	Flag      string     `json:"flag" gorm:"size:128;"`                  //
	Remark    string     `json:"remark" gorm:"size:255;"`                //备注
	Admin     bool       `json:"admin" gorm:"size:4;"`
	DataScope string     `json:"dataScope" gorm:"size:128;"` //权限范围
	Params    string     `json:"params" gorm:"-"`
	MenuIds   []int      `json:"menuIds" gorm:"-"`
	DeptIds   []int      `json:"deptIds" gorm:"-"`
	SysDept   []SysDept  `json:"sysDept" gorm:"many2many:sys_role_dept;foreignKey:RoleId;joinForeignKey:role_id;references:DeptId;joinReferences:dept_id;"`
	SysMenu   *[]SysMenu `json:"sysMenu" gorm:"many2many:sys_role_menu;foreignKey:RoleId;joinForeignKey:role_id;references:MenuId;joinReferences:menu_id;"`
	models.ControlBy
	models.ModelTime
}

func (*SysRole) TableName() string {
	return "sys_role"
}

func (e *SysRole) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysRole) GetId() interface{} {
	return e.RoleId
}
