package models

import "time"

// AlarmContactGroupMember 预警联系人组-成员
type AlarmContactGroupMember struct {
	Id                  int64      `gorm:"primaryKey;autoIncrement"`
	DeptId              int64      `json:"deptId"  gorm:"comment:机构"`                                                                                           //机构ID
	AlarmContactGroupId int64      `json:"alarmContactGroupId" gorm:"uniqueIndex:uniq_group_level_userid; not null;comment:联系人分组ID"`                            //联系人分组ID
	UserId              int64      `json:"userId"              gorm:"uniqueIndex:uniq_group_level_userid; not null"`                                            //用户ID
	AlarmLevel          AlarmLevel `json:"alarmLevel"          gorm:"uniqueIndex:uniq_group_level_userid; not null; type:tinyint;comment:预警等级 1=红 2=橙 3=黄 4=蓝"` //预警等级 1=红 2=橙 3=黄 4=蓝
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (AlarmContactGroupMember) TableName() string {
	return "alarm_contact_group_member"
}
