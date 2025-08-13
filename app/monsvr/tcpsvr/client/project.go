package client

// import (
// 	"encoding/binary"
// 	"fmt"
// 	"go-admin/app/admin/service"
// )

// // 收到盒子信息，加入数据库
// func (clt *DevClient) processProjectInfo(pkg *DevPackage) {

// 	var info DevProjectInfo

// 	// 读取结构体的字节数据
// 	err := binary.Read(clt.conn, binary.LittleEndian, &info)
// 	if err != nil {
// 		fmt.Println("Error reading data:", err)
// 		return
// 	}
// 	prj := service.SysRadar{}
// 	prj.ConfirmProjetInfo(info.ProID, info.DevID, info.DevLng, info.DevLat)
// 	//不需要读取CRC了，已经在结构体中
// 	clt.devid = info.DevID
// 	retData := RetDevProjectInfo{}
// 	retData.Code = 0
// 	pkg.DataLen = 1
// 	clt.sendData(pkg, retData)
// }
