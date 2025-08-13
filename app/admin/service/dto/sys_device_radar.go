package dto

// import (
// 	"go-admin/app/admin/models"

// 	"go-admin/common/dto"
// 	common "go-admin/common/models"
// )

// type SysDeviceRadarGetPageReq struct {
// 	dto.Pagination `search:"-"`
// 	RadarId        int    `form:"radarId" search:"type:exact;column:radar_id;table:sys_radar" comment:"雷达ID"`
// 	RadarName      string `form:"radarName" search:"type:contains;column:radar_name;table:sys_radar" comment:"雷达名称"`
// 	RadarKey       string `form:"radarKey" search:"type:contains;column:radar_key;table:sys_radar" comment:"雷达编号"`
// 	SpecialKey     string `form:"specialkey" search:"type:contains;column:special_key;table:sys_radar" comment:"特殊编号"`
// 	DeptJoin       `search:"type:left;on:dept_id:dept_id;table:sys_radar;join:sys_dept"`
// 	SysDeviceRadarOrder
// }

// type SysDeviceRadarOrder struct {
// 	RadarIdOrder   string `search:"type:order;column:radar_id;table:sys_radar" form:"radarIdOrder"`
// 	RadarNameOrder string `search:"type:order;column:radar_name;table:sys_radar" form:"radarNameOrder"`
// 	StatusOrder    string `search:"type:order;column:status;table:sys_radar" form:"statusOrder"`
// 	CreatedAtOrder string `search:"type:order;column:created_at;table:sys_radar" form:"createdAtOrder"`
// }

// func (m *SysDeviceRadarGetPageReq) GetNeedSearch() interface{} {
// 	return *m
// }

// type UpdateSysDeviceRadarStatusReq struct {
// 	RadarId int    `json:"radarId" comment:"RadarID" vd:"$>0"` // 用户ID
// 	Status  string `json:"status" comment:"状态" vd:"len($)>0"`
// 	common.ControlBy
// }

// func (s *UpdateSysDeviceRadarStatusReq) GetId() interface{} {
// 	return s.RadarId
// }

// func (s *UpdateSysDeviceRadarStatusReq) Generate(model *models.SysDeviceRadar) {
// 	if s.RadarId != 0 {
// 		model.RadarId = s.RadarId
// 	}
// 	model.Status = s.Status
// }

// type SysDeviceRadarInsertReq struct {
// 	RadarId    int    `json:"radarId" comment:"雷达ID"` // 雷达ID
// 	RadarName  string `json:"radarName" comment:"雷达名称" vd:"len($)>0"`
// 	RadarKey   string `json:"radarKey" comment:"雷达编号" vd:"len($)>0"`
// 	SpecialKey string `json:"specialKey" comment:"雷达特殊编号"`
// 	DeptId     int    `json:"deptId" comment:"所在单位" vd:"$>0"`
// 	Lng        string `json:"lng" comment:"经度坐标"`
// 	Lat        string `json:"lat" comment:"纬度坐标"`
// 	Alt        string `json:"alt" comment:"坐标高度" vd:"$>0"`
// 	Remark     string `json:"remark" comment:"备注"`
// 	Status     string `json:"status" comment:"状态" vd:"len($)>0" default:"1"`
// 	common.ControlBy
// }

// func (s *SysDeviceRadarInsertReq) Generate(model *models.SysDeviceRadar) {
// 	if s.RadarId != 0 {
// 		model.RadarId = s.RadarId
// 	}
// 	model.RadarName = s.RadarName
// 	model.RadarKey = s.RadarKey
// 	model.SpecialKey = s.SpecialKey
// 	model.DeptId = s.DeptId
// 	model.Lng = s.Lng
// 	model.Lat = s.Lat
// 	model.Alt = s.Alt
// 	model.Remark = s.Remark
// 	model.Status = s.Status
// 	model.CreateBy = s.CreateBy
// }

// func (s *SysDeviceRadarInsertReq) GetId() interface{} {
// 	return s.RadarId
// }

// // type SysDeviceRadarGetReq struct {
// // 	Id int `uri:"id"`
// // }

// // func (s *SysDeviceRadarGetReq) GetId() interface{} {
// // 	return s.Id
// // }

// type SysDeviceRadarUpdateReq struct {
// 	RadarId    int    `json:"radarId" comment:"雷达ID"` // 雷达ID
// 	RadarName  string `json:"radarName" comment:"雷达名称" vd:"len($)>0"`
// 	RadarKey   string `json:"radarKey" comment:"雷达编号" vd:"len($)>0"`
// 	SpecialKey string `json:"specialKey" comment:"雷达特殊编号"`
// 	DeptId     int    `json:"deptId" comment:"所在单位" vd:"$>0"`
// 	Lng        string `json:"lng" comment:"经度坐标"`
// 	Lat        string `json:"lat" comment:"纬度坐标"`
// 	Alt        string `json:"alt" comment:"坐标高度" vd:"$>0"`
// 	Remark     string `json:"remark" comment:"备注"`
// 	Status     string `json:"status" comment:"状态" vd:"len($)>0" default:"1"`
// 	common.ControlBy
// }

// func (s *SysDeviceRadarUpdateReq) Generate(model *models.SysDeviceRadar) {
// 	if s.RadarId != 0 {
// 		model.RadarId = s.RadarId
// 	}
// 	model.RadarName = s.RadarName
// 	model.RadarKey = s.RadarKey
// 	model.SpecialKey = s.SpecialKey
// 	model.DeptId = s.DeptId
// 	model.Lng = s.Lng
// 	model.Lat = s.Lat
// 	model.Alt = s.Alt
// 	model.Remark = s.Remark
// 	model.Status = s.Status
// }

// func (s *SysDeviceRadarUpdateReq) GetId() interface{} {
// 	return s.RadarId
// }

// type SysDeviceRadarById struct {
// 	dto.ObjectById
// 	common.ControlBy
// }

// func (s *SysDeviceRadarById) GetId() interface{} {
// 	if len(s.Ids) > 0 {
// 		s.Ids = append(s.Ids, s.Id)
// 		return s.Ids
// 	}
// 	return s.Id
// }

// func (s *SysDeviceRadarById) GenerateM() (common.ActiveRecord, error) {
// 	return &models.SysDeviceRadar{}, nil
// }

// type SysDeviceRadarDeleteReq struct {
// 	Ids []int `json:"ids"`
// }

// func (s *SysDeviceRadarDeleteReq) GetId() interface{} {
// 	return s.Ids
// }
