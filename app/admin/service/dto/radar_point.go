package dto

import (
	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type GetRadarPointListDeptIdReq struct {
	dto.Pagination `search:"-"`
	Id             int    `form:"id"  search:"type:exact;column:id;table:radar_point" comment:"PointID"`
	PointName      string `form:"pointName"  search:"type:contains;column:point_name;table:radar_point" comment:"监测点名称"`
	PointKey       string `form:"pointKey"  search:"type:exact;column:point_key;table:radar_point" comment:"监测点编号"`
	RadarId        int64  `form:"radarId"  search:"type:exact;column:radar_id;table:radar_point" comment:"雷达ID"`
	AStatus        string `form:"aStatus"  search:"type:exact;column:a_status;table:radar_point" comment:"激活状态"`
	RadarPointOrder
}

type RadarPointOrder struct {
	Id         string `form:"idOrder"  search:"type:order;column:id;table:radar_point"`
	PointName  string `form:"pointNameOrder"  search:"type:order;column:point_name;table:radar_point"`
	PointKey   string `form:"pointKeyOrder"  search:"type:order;column:point_key;table:radar_point"`
	PointType  string `form:"pointTypeOrder"  search:"type:order;column:point_type;table:radar_point"`
	RadarId    string `form:"radarIdOrder"  search:"type:order;column:radar_id;table:radar_point"`
	Lng        string `form:"lngOrder"  search:"type:order;column:lng;table:radar_point"`
	Lat        string `form:"latOrder"  search:"type:order;column:lat;table:radar_point"`
	Alt        string `form:"altOrder"  search:"type:order;column:alt;table:radar_point"`
	Distance   string `form:"distanceOrder"  search:"type:order;column:distance;table:radar_point"`
	PointIndex int    `form:"pointIndex"  search:"type:order;column:point_index;table:radar_point"`
	Remark     string `form:"remarkOrder"  search:"type:order;column:remark;table:radar_point"`
	AStatus    string `form:"aStatusOrder"  search:"type:order;column:a_status;table:radar_point"`
	XStatus    string `form:"xStatusOrder"  search:"type:order;column:x_status;table:radar_point"`
	MTypeId    string `form:"mTypeIdOrder"  search:"type:order;column:m_type_id;table:radar_point"`
	CreateBy   string `form:"createByOrder"  search:"type:order;column:create_by;table:radar_point"`
	UpdateBy   string `form:"updateByOrder"  search:"type:order;column:update_by;table:radar_point"`
	CreatedAt  string `form:"createdAtOrder"  search:"type:order;column:created_at;table:radar_point"`
	UpdatedAt  string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:radar_point"`
	DeletedAt  string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:radar_point"`
}

func (m *GetRadarPointListDeptIdReq) GetNeedSearch() interface{} {
	return *m
}

type InsertRadarPointReq struct {
	Id         int    `json:"-" comment:"PointID"` // PointID
	PointName  string `json:"pointName" comment:"监测点名称"`
	PointKey   string `json:"pointKey" comment:"监测点编号"`
	PointType  string `json:"pointType" comment:"监测点类型"`
	RadarId    int64  `json:"radarId" comment:"雷达ID"`
	Lng        string `json:"lng" comment:"经度"`
	Lat        string `json:"lat" comment:"纬度"`
	Alt        string `json:"alt" comment:"高度"`
	Distance   string `json:"distance" comment:"距离"`
	PointIndex int    `json:"pointIndex" comment:"下标"`
	Remark     string `json:"remark" comment:"备注"`
	AStatus    string `json:"aStatus" comment:"激活状态"`
	XStatus    string `json:"xStatus" comment:"消警状态"`
	MTypeId    string `json:"mTypeId" comment:"门限类型"`
	common.ControlBy
}

func (s *InsertRadarPointReq) Generate(model *models.RadarPoint) {
	if s.Id == 0 {
		model.Id = s.Id
	}
	model.PointName = s.PointName
	model.PointKey = s.PointKey
	model.PointType = s.PointType
	model.RadarId = s.RadarId
	model.Lng = s.Lng
	model.Lat = s.Lat
	model.Alt = s.Alt
	model.Distance = s.Distance
	model.PointIndex = s.PointIndex
	model.Remark = s.Remark
	model.AStatus = s.AStatus
	model.XStatus = s.XStatus
	model.MTypeId = s.MTypeId
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *InsertRadarPointReq) GetId() interface{} {
	return s.Id
}

type UpdateRadarPointReq struct {
	Id         int    `uri:"id" comment:"PointID"` // PointID
	PointName  string `json:"pointName" comment:"监测点名称"`
	PointKey   string `json:"pointKey" comment:"监测点编号"`
	PointType  string `json:"pointType" comment:"监测点类型"`
	RadarId    int64  `json:"radarId" comment:"雷达ID"`
	Lng        string `json:"lng" comment:"经度"`
	Lat        string `json:"lat" comment:"纬度"`
	Alt        string `json:"alt" comment:"高度"`
	Distance   string `json:"distance" comment:"距离"`
	PointIndex int    `json:"pointIndex" comment:"下标"`
	Remark     string `json:"remark" comment:"备注"`
	AStatus    string `json:"aStatus" comment:"激活状态"`
	XStatus    string `json:"xStatus" comment:"消警状态"`
	MTypeId    string `json:"mTypeId" comment:"门限类型"`
	common.ControlBy
}

// aStatus
// :
// "0"
// lat
// :
// "60"
// lng
// :
// "50"
// pointIndex
// :
// 4760
// pointKey
// :
// "005"
// pointName
// :
// "005"

func (s *UpdateRadarPointReq) Generate(model *models.RadarPoint) {
	if s.Id == 0 {
		model.Id = s.Id
	}
	model.PointName = s.PointName
	model.PointKey = s.PointKey
	model.PointType = s.PointType
	model.RadarId = s.RadarId
	model.Lng = s.Lng
	model.Lat = s.Lat
	model.Alt = s.Alt
	model.Distance = s.Distance
	model.PointIndex = s.PointIndex
	model.Remark = s.Remark
	model.AStatus = s.AStatus
	model.XStatus = s.XStatus
	model.MTypeId = s.MTypeId
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *UpdateRadarPointReq) GetId() interface{} {
	return s.Id
	//return s.id
}

// GetRadarPointByIdReq 功能获取请求参数
type GetRadarPointByIdReq struct {
	Id int `uri:"id"`
}

func (s *GetRadarPointByIdReq) GetId() interface{} {
	return s.Id
	//return s.id
}

// DeleteRadarPointReq 功能删除请求参数
type DeleteRadarPointReq struct {
	Ids []int `json:"ids"`
}

func (s *DeleteRadarPointReq) GetId() interface{} {
	return s.Ids
}

type RadarPointIndex struct {
	Index int `json:"index" comment:"监测点下标"`
}
