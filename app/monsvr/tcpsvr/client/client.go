package client

// import (
// 	"encoding/binary"
// 	"fmt"
// 	"net"
// )

// type DevClient struct {
// 	conn      net.Conn
// 	devid     uint32 //设备id（radarid）
// 	testCount uint64
// }

// var clientConns = make([]DevClient, 0) // 创建一个长度为5，容量也为5的int类型切片

// func HandleConnection(conn net.Conn) {
// 	clt := DevClient{devid: 0}
// 	clt.testCount = 0
// 	clt.conn = conn
// 	clientConns = append(clientConns, clt)

// 	clt.SendRadarQueryInfo()

// 	clt.loopHandle()
// }

// func (clt *DevClient) loopHandle() {
// 	defer clt.conn.Close() // 确保连接被关闭
// 	//reader := bufio.NewReader(clt.conn)
// 	for {
// 		var pkg DevPackage

// 		// 读取header
// 		err := binary.Read(clt.conn, binary.LittleEndian, &pkg.Head)
// 		if err != nil {
// 			fmt.Println("Error reading data:", err)
// 			break
// 		}
// 		if pkg.Head != 0x7E7E9C9C {
// 			continue
// 		}
// 		// 读取结构体的字节数据
// 		err = binary.Read(clt.conn, binary.LittleEndian, &pkg.DevPackageNoHeader)
// 		if err != nil {
// 			fmt.Println("Error reading data:", err)
// 			break
// 		}

// 		//fmt.Print("Message received:", message)              // 打印接收到的消息
// 		//clt.conn.Write([]byte("Message echoed: " + message)) // 将消息回显给客户端
// 		clt.processPackage(&pkg)
// 	}
// }

// func (clt *DevClient) processPackage(pkg *DevPackage) {
// 	switch pkg.DevType {
// 	case D_TYPE_CBOX:
// 		clt.processCBoxMsg(pkg)
// 	case D_TYPE_RADAR:
// 		clt.processRadarMsg(pkg)
// 	}
// }
// func (clt *DevClient) readCrc32() uint32 {
// 	var CRC32 uint32
// 	err := binary.Read(clt.conn, binary.LittleEndian, &CRC32)
// 	if err != nil {
// 		fmt.Println("Error reading data:", err)
// 		return 0
// 	}
// 	return CRC32
// }
