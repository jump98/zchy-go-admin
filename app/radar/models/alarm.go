package models

// AlarmCheckType 监测类型
type AlarmCheckType int

// AlarmType 预警类型
type AlarmType int

// AlarmLevel 预警等级枚举
type AlarmLevel int

// AlarmRuleOptionMode 预警满足条件
//type AlarmRuleOptionMode int

//const (
//	AlarmRuleOptionModeAll AlarmRuleOptionMode = iota + 1 //预警满足条件-全部满足
//	AlarmRuleOptionModeOr                                 //预警满足条件-满足其一
//)

const (
	AlarmCheckRadarPoint AlarmCheckType = iota + 1 //预警-检测点
)

const (
	AlarmType_RadarPoint_Deformation  AlarmType = 100 //监测点-形变
	AlarmType_RadarPoint_Velocity     AlarmType = 101 //监测点-速度-状态信息
	AlarmType_RadarPoint_Acceleration AlarmType = 102 //监测点-加速度-状态信息
)

const (
	AlarmLevelRed    AlarmLevel = iota + 1 // 红色
	AlarmLevelOrange                       // 橙色
	AlarmLevelYellow                       // 黄色
	AlarmLevelBlue                         // 蓝色
)
