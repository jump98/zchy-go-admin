package main

import (
	"bufio"
	"fmt"
	"go-admin/cmd"
	"os"

	"github.com/gin-gonic/gin"
)

//go:generate swag init --parseDependency --parseDepth=6 --instanceName admin -o ./docs/admin

// @title go-admin API
// @version 2.0.0
// @description 接口文档
// @description radar api
// @license.name MIT

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	fmt.Println("	gin.Mode() :", gin.Mode())
	if gin.Mode() == "debug" {
		fmt.Println("监听回车按键")
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				text := scanner.Text()
				if text == "" {
					fmt.Println() // 按回车输出空行
				}
			}
		}()
	}

	cmd.Execute()

}
