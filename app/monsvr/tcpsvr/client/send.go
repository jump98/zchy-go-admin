package client

// import (
// 	"encoding/binary"
// 	"fmt"
// )

// func (clt *DevClient) sendData(pkg *DevPackage, dat interface{}) {

// 	pkg.Head = 0x7E7E9C9C
// 	b, _ := MergeAddCRC32Sum(pkg, dat)

// 	err := binary.Write(clt.conn, binary.BigEndian, b)
// 	if err != nil {
// 		fmt.Println("Error sending message:", err)
// 		return
// 	}

// }
