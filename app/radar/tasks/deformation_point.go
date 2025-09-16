package tasks

import (
	"context"
	"go-admin/app/monsvr/mongosvr"
	"sync"
	"time"

	"github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
)

type DeformationPointTask struct {
	ctx    context.Context
	cancel context.CancelFunc
	Logger *logger.Helper
	mu     sync.Mutex
	minute time.Duration //按多少分钟聚合
}

// InitDeformationPointTask 初始化监测点预警任务
func InitDeformationPointTask(parentCtx context.Context) *DeformationPointTask {
	ctx, cancel := context.WithCancel(parentCtx)
	t := &DeformationPointTask{
		Logger: logger.NewHelper(sdk.Runtime.GetLogger()).WithFields(map[string]interface{}{}),
		minute: 1,
		ctx:    ctx,
		cancel: cancel,
	}
	go t.startTask()

	return t
}

func (t *DeformationPointTask) Stop() {
	t.cancel() // 通知任务退出
}

func (t *DeformationPointTask) startTask() {
	ticker := time.NewTicker(t.minute * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-t.ctx.Done():
			t.Logger.Info("DeformationPointTask 停止")
			return
		case <-ticker.C:
			t.mu.Lock()
			t.Logger.Info("DeformationPointTask 开始")
			t.aggregateLastMinute()
			t.mu.Unlock()
		}
	}
}

func (t *DeformationPointTask) aggregateLastMinute() {
	ctx := context.Background()
	now := time.Now()
	startTime := now.Add(-(t.minute) * time.Minute).Truncate(time.Minute)
	endTime := now.Truncate(time.Minute)
	// 打印本地时间
	//fmt.Println("查询开始时间:", startTime.Format("2006-01-02 15:04:05"))
	//fmt.Println("查询结束时间:", endTime.Format("2006-01-02 15:04:05"))

	var err error
	if err = mongosvr.DeformationPointMinuteService.SaveDeformMinuteByTime(ctx, startTime, endTime); err != nil {
		t.Logger.Error("保存形变分钟表出错:", err)
		return
	}
}
