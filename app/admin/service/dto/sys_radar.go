package dto

import (
	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysRadarGetPageReq struct {
	dto.Pagination `search:"-"`
	RadarId        int64  `form:"radarId"  search:"type:exact;column:radar_id;table:sys_radar" comment:"RadarID"`
	RadarName      string `form:"radarName"  search:"type:exact;column:radar_name;table:sys_radar" comment:"雷达名称"`
	RadarKey       string `form:"radarKey"  search:"type:exact;column:radar_key;table:sys_radar" comment:"雷达编号"`
	SpecialKey     string `form:"specialKey"  search:"type:exact;column:special_key;table:sys_radar" comment:"雷达特殊编号"`
	//DeptId         int64  `form:"deptId"  search:"type:exact;column:dept_id;table:sys_radar" comment:"部门"`
	DeptJoin `search:"type:left;on:dept_id:dept_id;table:sys_radar;join:sys_dept"`
	SysRadarOrder
}

type SysRadarOrder struct {
	RadarId    string `form:"radarIdOrder"  search:"type:order;column:radar_id;table:sys_radar"`
	RadarName  string `form:"radarNameOrder"  search:"type:order;column:radar_name;table:sys_radar"`
	RadarKey   string `form:"radarKeyOrder"  search:"type:order;column:radar_key;table:sys_radar"`
	SpecialKey string `form:"specialKeyOrder"  search:"type:order;column:special_key;table:sys_radar"`
	DeptId     string `form:"deptIdOrder"  search:"type:order;column:dept_id;table:sys_radar"`
	Lng        string `form:"lngOrder"  search:"type:order;column:lng;table:sys_radar"`
	Lat        string `form:"latOrder"  search:"type:order;column:lat;table:sys_radar"`
	Alt        string `form:"altOrder"  search:"type:order;column:alt;table:sys_radar"`
	Remark     string `form:"remarkOrder"  search:"type:order;column:remark;table:sys_radar"`
	Status     string `form:"statusOrder"  search:"type:order;column:status;table:sys_radar"`
	CreateBy   string `form:"createByOrder"  search:"type:order;column:create_by;table:sys_radar"`
	UpdateBy   string `form:"updateByOrder"  search:"type:order;column:update_by;table:sys_radar"`
	CreatedAt  string `form:"createdAtOrder"  search:"type:order;column:created_at;table:sys_radar"`
	UpdatedAt  string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:sys_radar"`
	DeletedAt  string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:sys_radar"`
}

func (m *SysRadarGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SysRadarInsertReq struct {
	RadarId     int64  `json:"-" comment:"RadarID"` // RadarID
	RadarName   string `json:"radarName" comment:"雷达名称"`
	RadarKey    string `json:"radarKey" comment:"雷达编号"`
	SpecialKey  string `json:"specialKey" comment:"雷达特殊编号"`
	DeptId      int    `json:"deptId" comment:"部门" vd:"$>0"`
	Lng         string `json:"lng" comment:"经度"`
	Lat         string `json:"lat" comment:"纬度"`
	Alt         string `json:"alt" comment:"高度"`
	Remark      string `json:"remark" comment:"备注"`
	Status      string `json:"status" comment:"状态"`
	Vender      string `json:"vender" comment:"设备厂家名"`
	Secret      string `json:"secret" comment:"密钥"`
	FromProject int    `json:"fromProject" comment:"来自项目"`
	common.ControlBy
}

func (s *SysRadarInsertReq) Generate(model *models.SysRadar) {
	if s.RadarId == 0 {
		//model.Model = common.Model{Id: s.RadarId}
		model.RadarId = s.RadarId
	}
	model.RadarName = s.RadarName
	model.RadarKey = s.RadarKey
	model.SpecialKey = s.SpecialKey
	model.DeptId = int64(s.DeptId)
	model.Lng = s.Lng
	model.Lat = s.Lat
	model.Alt = s.Alt
	model.Remark = s.Remark
	model.Status = s.Status
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.FromProject = s.FromProject
	model.Vender = s.Vender
	model.Secret = s.Secret
}

func (s *SysRadarInsertReq) GetId() interface{} {
	return s.RadarId
}

type SysRadarUpdateReq struct {
	RadarId     int64  `uri:"radarId" comment:"RadarID"` // RadarID
	RadarName   string `json:"radarName" comment:"雷达名称"`
	RadarKey    string `json:"radarKey" comment:"雷达编号"`
	SpecialKey  string `json:"specialKey" comment:"雷达特殊编号"`
	DeptId      int    `json:"deptId" comment:"部门" vd:"$>0"`
	Lng         string `json:"lng" comment:"经度"`
	Lat         string `json:"lat" comment:"纬度"`
	Alt         string `json:"alt" comment:"高度"`
	Remark      string `json:"remark" comment:"备注"`
	Status      string `json:"status" comment:"状态"`
	FromProject int    `json:"fromProject" comment:"来自项目"`
	common.ControlBy
}

func (s *SysRadarUpdateReq) Generate(model *models.SysRadar) {
	if s.RadarId == 0 {
		//model.Model = common.Model{Id: s.RadarId}
		model.RadarId = s.RadarId
	}
	model.RadarName = s.RadarName
	model.RadarKey = s.RadarKey
	model.SpecialKey = s.SpecialKey
	model.DeptId = int64(s.DeptId)
	model.Lng = s.Lng
	model.Lat = s.Lat
	model.Alt = s.Alt
	model.Remark = s.Remark
	model.Status = s.Status
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.FromProject = s.FromProject
}

func (s *SysRadarUpdateReq) GetId() interface{} {
	return s.RadarId
}

// SysRadarGetReq 功能获取请求参数
type SysRadarGetReq struct {
	RadarId int64 `uri:"radarId"`
}

func (s *SysRadarGetReq) GetId() interface{} {
	return s.RadarId
}

// SysRadarGetReq 功能获取请求参数
type SysRadarKeyGetReq struct {
	RadarKey string `uri:"radarkey"`
}

func (s *SysRadarKeyGetReq) GetRadarKey() interface{} {
	return s.RadarKey
}

// SysRadarGetAlarmsOfIdsReq 获取告警列表请求参数
type SysRadarGetAlarmsOfIdsReq struct {
	Ids []int64 `json:"ids"`
}

func (s *SysRadarGetAlarmsOfIdsReq) GetIds() []int64 {
	return s.Ids
}

// SysRadarGetImageReq 功能获取请求参数
type SysRadarGetImageReq struct {
	RadarId int64 `uri:"radarId"`
}

func (s *SysRadarGetImageReq) GetId() interface{} {
	return s.RadarId
}

// SysRadarDeleteReq 功能删除请求参数
type SysRadarDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *SysRadarDeleteReq) GetId() interface{} {
	return s.Ids
}

// SysRadarGetAlarmsBeforeReq 获取指定时间之前的告警列表请求参数
type SysRadarGetAlarmsBeforeReq struct {
	RadarId int64  `json:"radarId"`
	Time    string `json:"time"`
	Num     int    `json:"num"`
}

func (s *SysRadarGetAlarmsBeforeReq) GetRadarId() int64 {
	return s.RadarId
}

func (s *SysRadarGetAlarmsBeforeReq) GetTime() string {
	return s.Time
}

func (s *SysRadarGetAlarmsBeforeReq) GetNum() int {
	return s.Num
}

// // SysRadarConfirmReq 功能删除请求参数
// type SysRadarConfirmReq struct {
// 	RadarId int64 `uri:"radarId" comment:"RadarID"` // RadarID
// }

// func (s *SysRadarConfirmReq) GetId() interface{} {
// 	return s.RadarId
// }
