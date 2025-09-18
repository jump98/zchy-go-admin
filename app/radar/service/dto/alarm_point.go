package dto

import (
	"go-admin/app/radar/models"
	"go-admin/common/dto"
	"time"
)

// GetAlarmRulesReq  获取预警规则列表
type GetAlarmRulesReq struct {
	DeptId       int64 `form:"deptId"`       //机构ID
	RadarPointId int64 `form:"radarPointId"` //机构ID
}

type GetAlarmRulesResp struct {
	AlarmRuleList []models.AlarmPoint `json:"alarmRule"` //机构ID
}

// AddAlarmPointReq 增加预警规则列表
type AddAlarmPointReq struct {
	Items []AlarmPointItem `json:"items"`
}
type AddAlarmPointResp struct {
	Success bool `json:"success"` //成功
}

// AlarmPointItem  预警规则等级表
type AlarmPointItem struct {
	DeptId         int64                 `json:"deptId"`         //机构ID
	RadarPointId   int64                 `json:"radarPointId"`   //监测点ID
	AlarmCheckType models.AlarmCheckType `json:"alarmCheckType"` //监测类型
	AlarmName      string                `json:"alarmName"`      //预警规则名称
	Remark         string                `json:"remark"`         //预警规则介绍
	AlarmType      models.AlarmType      `json:"alarmType"`      //预警类型
	RedOption      string                `json:"redOption"`      //预警条件
	OrangeOption   string                `json:"orangeOption"`   //预警条件
	YellowOption   string                `json:"yellowOption"`   //预警条件
	BlueOption     string                `json:"blueOption"`     //预警条件
	Interval       uint64                `json:"interval"`       //预警间隔时间（分）
	Duration       uint64                `json:"duration"`       //连续预警次数
}

// UpdateAlarmPointReq 修改预警规则列表
type UpdateAlarmPointReq struct {
	Mode         models.RadarPointMType `form:"mode"`         //门限类型
	RadarPointId int64                  `form:"radarPointId"` //监测点ID
	DeptId       int64                  `form:"deptId"`
	Items        []AlarmPointItem       `form:"items"`
}
type UpdateAlarmPointesp struct {
}

// DeleteAlarmPointReq  删除预警设定列表
type DeleteAlarmPointReq struct {
	AlarmRuleId int64 `json:"alarmRuleId"` //预警规矩ID
}
type DeleteAlarmPointResp struct {
}

type GetAlarmPointLogsPageReq struct {
	dto.Pagination `search:"-"`
	DeptId         int64             `form:"deptId" search:"type:exact;column:dept_id;table:alarm_point_logs"`              //机构ID
	RadarId        int64             `form:"radarId" search:"type:exact;column:radar_id;table:alarm_point_logs"`            //雷达ID
	radarPointId   int64             `form:"radarPointId" search:"type:exact;column:radar_point_id;table:alarm_point_logs"` //雷达ID
	AlarmType      models.AlarmType  `form:"alarmType" search:"type:exact;column:alarm_type;table:alarm_point_logs"`        //预警类型
	AlarmLevel     models.AlarmLevel `form:"alarmLevel" search:"type:exact;column:alarm_level;table:alarm_point_logs"`      //报警等级
	StartTime      time.Time         `form:"start" search:"type:gte;column:created_at;table:alarm_point_logs"`              //开始时间
	EndTime        time.Time         `form:"end" search:"type:lte;column:created_at;table:alarm_point_logs"`                //结束时间
}

func (m *GetAlarmPointLogsPageReq) GetNeedSearch() interface{} {
	return *m
}

type GetAlarmPointLogsPageResp struct {
	List []models.AlarmPointLogs `json:"list"`
}
