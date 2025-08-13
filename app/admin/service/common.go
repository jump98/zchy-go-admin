package service

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"gorm.io/gorm"
)

func generateGinContext() *gin.Context {
	w := httptest.NewRecorder()

	// 创建一个测试用的 gin.Context
	c, _ := gin.CreateTestContext(w)

	// 手动设置请求对象（可选）
	c.Request = &http.Request{
		Method: "GET",
		Header: http.Header{},
		// 可以设置请求体（如 POST 请求）
		// Body: io.NopCloser(bytes.NewBufferString(`{"name":"test"}`)),
	}

	// 设置路径参数（如 /user/:id）
	c.Params = []gin.Param{{Key: "id", Value: "123"}}

	// 设置查询参数（如 ?name=test）
	c.Request.URL = &url.URL{
		RawQuery: "name=test",
	}
	return c
}

func getFirstOrm() *gorm.DB {

	x := sdk.Runtime.GetDb()
	for _, v := range x {
		return v
	}
	return nil
}
