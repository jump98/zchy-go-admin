package monsvr

import (
	"fmt"
	"go-admin/app/monsvr/mongosvr"
)

func InitMonSvr() {
	fmt.Println("================>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	go mongosvr.Init()
	//go tcpsvr.InitTcpServer()
}
