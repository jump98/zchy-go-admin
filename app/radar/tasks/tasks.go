package tasks

import (
	"context"
	"fmt"
)

var deformationTask *DeformationPointTask
var alarmPointTask *AlarmPointTask

func Init(ctx context.Context) {
	alarmPointTask = InitAlarmPointTask(ctx)
	deformationTask = InitDeformationPointTask(ctx)
}

func Stop() {
	fmt.Println("停止定时任务")
	deformationTask.Stop()
	alarmPointTask.Stop()
}
