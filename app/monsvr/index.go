package monsvr

import (
	"go-admin/app/monsvr/mongosvr"
)

func InitMonSvr() {
	go mongosvr.Init()
	//go tcpsvr.InitTcpServer()
}
