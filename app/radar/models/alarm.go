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
	AlarmTypeRadarPointDeformation  AlarmType = 100 //监测点-累计水平位移，单位 ： mm
	AlarmTypeRadarPointVelocity     AlarmType = 101 //监测点-水平位移速度（瞬时位移量）预警阈值，单位：mm/h
	AlarmTypeRadarPointAcceleration AlarmType = 102 //监测点-水平位移加速度预警阈值，单位： mm/h2
)

const (
	AlarmLevelNone   AlarmLevel = iota
	AlarmLevelBlue              // 蓝色
	AlarmLevelYellow            // 黄色
	AlarmLevelOrange            // 橙色
	AlarmLevelRed               // 红色
)
