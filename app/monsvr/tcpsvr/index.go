package tcpsvr

// import (
// 	"fmt"
// 	"go-admin/app/monsvr/tcpsvr/client"
// 	"net"
// 	"os"
// )

// func InitTcpServer() {
// 	port := "10881"
// 	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error listening: %s\n", err)
// 		os.Exit(1)
// 	}
// 	defer listener.Close()
// 	fmt.Println("Monitor Server is running on port " + port)

// 	for {
// 		// 接受连接
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "Error accepting: %s\n", err)
// 			continue
// 		}
// 		// 为每个连接创建一个goroutine进行处理
// 		go client.HandleConnection(conn)
// 	}
// }
