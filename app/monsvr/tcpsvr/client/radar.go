package client

// import (
// 	"encoding/binary"
// 	"fmt"
// 	"go-admin/app/monsvr/mongosvr"
// )

// func (clt *DevClient) processRadarMsg(pkg *DevPackage) {
// 	switch pkg.Cmd {
// 	case D_CMD_HEARTBREAK:
// 		clt.processRadarHeartBreakInfo(pkg)
// 	case D_CMD_RADAR_QUERYOUTPARAM:
// 		clt.processRadarQueryInfo(pkg)
// 	case D_CMD_RADAR_ADDPOINT_IMAGE:
// 		if pkg.CmdExtend == D_CMD_RADAR_EXT_ADDPOINT {
// 			clt.processRadarAddCtrlPoint(pkg)
// 		} else {
// 			clt.processRadarDistanceImage(pkg)
// 		}
// 	}

// }

// // 查询雷达信息
// func (clt *DevClient) SendRadarQueryInfo() {
// 	pkg := DevPackage{}
// 	pkg.DevType = D_TYPE_RADAR
// 	pkg.Cmd = D_CMD_RADAR_QUERYOUTPARAM
// 	pkg.CmdExtend = 0
// 	pkg.DataLen = 0
// 	clt.sendData(&pkg, nil)
// }

// // 添加选点雷达信息
// func (clt *DevClient) SendRadarAddCtrl() {
// 	pkg := DevPackage{}
// 	pkg.DevType = D_TYPE_RADAR
// 	pkg.Cmd = D_CMD_RADAR_ADDPOINT_IMAGE
// 	pkg.CmdExtend = D_CMD_RADAR_EXT_ADDPOINT
// 	pkg.DataLen = 0
// 	clt.sendData(&pkg, nil)
// }

// // 收到雷达查询信息
// func (clt *DevClient) processRadarQueryInfo(pkg *DevPackage) {

// 	var info RadarDeformatQueryInfo

// 	// 读取结构体的字节数据
// 	err := binary.Read(clt.conn, binary.LittleEndian, &info)
// 	if err != nil {
// 		fmt.Println("Error reading data:", err)
// 		return
// 	}
// 	//不需要读取CRC了，已经在结构体中
// }

// // 收到雷达添加点信息
// func (clt *DevClient) processRadarAddCtrlPoint(pkg *DevPackage) {

// 	var info RadarAddCtrlPointInfo

// 	// 读取结构体的字节数据
// 	err := binary.Read(clt.conn, binary.LittleEndian, &info)
// 	if err != nil {
// 		fmt.Println("Error reading data:", err)
// 		return
// 	}
// 	//不需要读取CRC了，已经在结构体中
// }

// // 收到雷达距离影像信息
// func (clt *DevClient) processRadarDistanceImage(pkg *DevPackage) {

// 	var info mongosvr.DistanceData

// 	initToDistanceImageData(&info, pkg, clt)
// 	// 读取结构体的字节数据
// 	err := binary.Read(clt.conn, binary.LittleEndian, &info.EchoData)
// 	if err != nil {
// 		fmt.Println("Error reading data:", err)
// 		return
// 	}
// 	//需要读取CRC
// 	clt.readCrc32()
// 	mongosvr.InsertDistanceData(&info)
// }

// // 收到雷达形变数据信息
// func (clt *DevClient) processRadarHeartBreakInfo(pkg *DevPackage) {

// 	clt.testCount++
// 	//fmt.Println("processRadarHeartBreakInfo:", clt.testCount)

// 	var info RadarHeartBreakInfo

// 	// 读取结构体的字节数据
// 	err := binary.Read(clt.conn, binary.LittleEndian, &info)
// 	if err != nil {
// 		fmt.Println("Error reading data:", err)
// 		return
// 	}
// 	size := info.DetectNum * info.TestTimes * 2 //*2是因为float32距离(米) +float32形变(毫米)
// 	var defArray = make([]float32, size)
// 	// 读取float数组
// 	err = binary.Read(clt.conn, binary.LittleEndian, &defArray)
// 	if err != nil {
// 		fmt.Println("Error reading data:", err)
// 		return
// 	}
// 	//需要读取CRC
// 	clt.readCrc32()
// 	var data mongosvr.DeformationData
// 	initToDeformationData(&data, pkg, clt, &info, defArray)
// 	mongosvr.InsertDeformationData(&data)
// }
