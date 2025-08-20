package apis

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/app/monsvr/mongosvr"

	mongodto "go-admin/app/monsvr/mongosvr/dto"

	"github.com/golang-jwt/jwt/v4"
)

type RadarAuthenticateRequest struct {
	RadarKey  string `json:"radarkey" binding:"required"`
	Vender    string `json:"vender" binding:"required"`
	Secret    string `json:"secret" binding:"required"`
	Status    int    `json:"status"`
	Timestamp int64  `json:"timestamp" binding:"required"`
}

type RadarAuthenticateResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

// RadarAlarmRequest 雷达告警信息请求
type RadarAlarmRequest struct {
	RadarId     int64 `json:"-"` // 从token中获取
	Time        int64 `json:"time"`
	Voltage     int   `json:"voltage"`
	Temperature int   `json:"temperature"`
	Battery     int   `json:"battery"`
	SolarPanel  int   `json:"solar_panel"`
	RadarData   int   `json:"radar_data"`
}

// 获得token中的radarId
func (e *SysRadar) GetTokenRadarId(c *gin.Context) (int64, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, errors.New("authorization header missing")
	}

	var err error
	var tokenClaims *TokenClaims
	if tokenClaims, err = e.GetParseClaimsToken(authHeader); err != nil {
		fmt.Println("请求解析token出错:", err)
		return 0, fmt.Errorf("请求解析token出错: %v", err)
	}
	// nowTime := time.Now().Unix()
	// if tokenClaims.Exp < nowTime {
	// 	e.Logger.Error("token已过期")
	// 	return 0, errors.New("token已过期")
	// }
	return tokenClaims.RadarId, nil
}

type TokenClaims struct {
	jwt.RegisteredClaims
	RadarId int64 `json:"radarId"`
	Exp     int64
}

// 解密token
func (e *SysRadar) GetParseClaimsToken(tokenStr string) (*TokenClaims, error) {
	tokenStr2 := strings.TrimPrefix(tokenStr, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenStr2, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JwtConfig.Secret), nil
	})
	if err != nil {
		e.Logger.Error("解析token出错:", err)
		return nil, err
	}
	// 获取claims
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		e.Logger.Error("解析token.claims出错:", err)
		return nil, err
	}
	return claims, nil
}

// Authenticate 雷达设备认证
// @Summary 雷达设备认证
// @Description 雷达设备登录认证
// @Tags 雷达管理-终端接口
// @Accept  application/json
// @Product application/json
// @Param data body RadarAuthenticateRequest true "认证信息"
// @Success 200 {object} RadarAuthenticateResponse "{\"code\": 0, \"message\": \"Status received\", \"token\": \"...\", \"expires_in\": 3600}"
// @Router /api/v1/radar/authenticate [post]
func (e SysRadar) Authenticate(c *gin.Context) {
	var err error
	s := service.SysRadar{}
	if err = e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.Logger.Error(err)
		e.Error(400, err, err.Error())
		return
	}

	var req RadarAuthenticateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		e.Error(400, err, "请求参数错误")
		return
	}

	if req.RadarKey == "" || req.Vender == "" || req.Secret == "" {
		e.Error(400, nil, "认证失败: 无效的雷达ID、厂商或密钥")
		return
	}

	fmt.Printf("请求登录的雷达信息:%+v \n", req)

	radarId := int64(0)
	// 生成JWT token
	// 需要先导入 jwt 包，此处假设使用的是 github.com/golang-jwt/jwt/v5
	radar := &models.SysRadar{}
	if err = s.GetByRadarKey(req.RadarKey, radar); err != nil {
		//没有找到雷达，则新建
		reqNew := dto.SysRadarInsertReq{}
		reqNew.RadarName = req.RadarKey
		reqNew.RadarKey = req.RadarKey
		reqNew.SpecialKey = ""
		reqNew.DeptId = 1 //默认智存合一
		reqNew.Lng = "0"
		reqNew.Lat = "0"
		reqNew.Alt = "0"
		reqNew.Remark = ""
		reqNew.Status = "0"
		reqNew.Vender = req.Vender
		reqNew.Secret = req.Secret
		err = s.Insert(&reqNew)
		if err != nil {
			e.Error(500, err, err.Error())
			return
		}
		radarId = reqNew.GetId().(int64)
	} else {
		radarId = radar.RadarId
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"radarId": radarId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// 签名token
	tokenString, err := token.SignedString([]byte(config.JwtConfig.Secret)) // 请替换为实际的密钥
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, "Token生成失败")
		return
	}

	resp := RadarAuthenticateResponse{
		Code:      0,
		Message:   "认证成功",
		Token:     tokenString,
		ExpiresIn: 86400, // 24小时
	}

	e.OK(resp, "认证成功")
}

// PutAlarm 雷达设备上传告警信息
// @Summary 雷达设备上传告警信息
// @Description 雷达设备上传告警信息
// @Tags 雷达管理-终端接口
// @Accept  application/json
// @Product application/json
// @Param data body RadarAlarmRequest true "告警信息"
// @Success 200 {object} response.Response "{\"code\": 0, \"message\": \"告警信息接收成功\"}"
// @Router /api/v1/radar/put_alarm [post]
// @Security Bearer
func (e SysRadar) PutAlarm(c *gin.Context) {
	e.MakeContext(c)

	var radarId int64
	var err error
	if radarId, err = e.GetTokenRadarId(c); err != nil {
		e.Error(400, err, err.Error())
		return
	}

	var req RadarAlarmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		e.Error(400, err, "请求参数错误")
		return
	}

	//time = time.Unix(req.Time, 0)

	// 设置雷达ID
	req.RadarId = radarId

	// 处理告警信息逻辑
	e.Logger.Infof("接收到雷达 %s 的告警信息: %+v", radarId, req)

	// 存储告警信息到MongoDB
	alarmData := &mongosvr.AlarmData{
		RadarId:     req.RadarId,
		TimeStamp:   req.Time,
		Voltage:     req.Voltage,
		Temperature: req.Temperature,
		Battery:     req.Battery,
		SolarPanel:  req.SolarPanel,
		RadarData:   req.RadarData,
	}

	if err := mongosvr.InsertAlarmData(alarmData); err != nil {
		e.Logger.Errorf("存储告警信息失败: %v", err)
		// 这里可以根据需求决定是否返回错误
		e.Error(500, err, "存储告警信息失败")
		return
	}

	e.OK(nil, "告警信息接收成功")
}

// PutRebootCommand 雷达设备重启指令
// @Summary 雷达设备重启指令
// @Description 雷达设备重启指令
// @Tags 雷达管理-终端接口
// @Accept  application/json
// @Product application/json
// @Param data body object false "测试信息"
// @Success 200 {object} response.Response "{\"code\": 0, \"message\": \"信息接收成功\"}"
// @Router /api/v1/radar/dev_reboot [post]
// @Security Bearer
func (e SysRadar) PutRebootCommand(c *gin.Context) {
	e.MakeContext(c)

	var radarId int64
	var err error
	if radarId, err = e.GetTokenRadarId(c); err != nil {
		e.Error(400, err, err.Error())
		return
	}

	//测试命令
	mongosvr.InsertCommandData(&mongosvr.CommandData{
		RadarId:     radarId,
		TimeStamp:   time.Now().Unix(),
		CommandCode: mongosvr.CMD_RD_REBOOT,
		Message:     "reboot",
		Parameters:  map[string]interface{}{},
	})

	e.OK(nil, "测试信息接收成功")
}

// PutTestCommand 测试新增一个雷达设备测试指令
// @Summary 测试新增一个雷达设备测试指令
// @Description 测试新增一个雷达设备测试指令
// @Tags 雷达管理-终端接口
// @Accept  application/json
// @Product application/json
// @Param data body RadarTestCommandRequest true "测试指令"
// @Success 200 {object} response.Response "{\"code\": 0, \"message\": \"信息接收成功\"}"
// @Router /api/v1/radar/put_testcommand [post]
// @Security Bearer
func (e SysRadar) PutTestCommand(c *gin.Context) {
	e.MakeContext(c)
	var radarId int64
	var err error
	if radarId, err = e.GetTokenRadarId(c); err != nil {
		e.Error(400, err, err.Error())
		return
	}

	var req RadarTestCommandRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		e.Error(400, err, "请求参数错误")
		return
	}

	//测试命令
	mongosvr.InsertCommandData(&mongosvr.CommandData{
		RadarId:     radarId,
		TimeStamp:   time.Now().Unix(),
		CommandCode: req.CommandCode,
		Message:     req.Message,
		Parameters:  req.Parameters,
	})

	e.OK(nil, "测试信息接收成功")
}

// RadarTestCommandRequest 测试指令请求
type RadarTestCommandRequest struct {
	CommandCode int                    `json:"command_code"`
	Message     string                 `json:"message"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// DeformationDataPoint 形变数据点
type DeformationDataPoint struct {
	Index       int     `json:"index"`
	Deformation float32 `json:"deformation"`
	Distance    float32 `json:"distance"`
}

// DeformationRequest 形变数据请求
type DeformationRequest struct {
	RadarKey  string                 `json:"radarkey"` // 从token中获取
	Timestamp int64                  `json:"timestamp"`
	Interval  int                    `json:"interval"`
	Data      []DeformationDataPoint `json:"data"`
}

func (e SysRadar) InitDeformationData(dd *mongosvr.DeformationData, dr *DeformationRequest) {
	dd.RadarKey = dr.RadarKey
	dd.TimeStamp = time.Unix(dr.Timestamp, 0)
	dd.Interval = dr.Interval
	for _, dp := range dr.Data {
		dd.DefData = append(dd.DefData, mongosvr.DeformationDefData{
			Index:       dp.Index,
			Deformation: dp.Deformation,
			Distance:    dp.Distance,
		})
	}
}

// GetCommandsRequest 获取命令请求
type GetCommandsRequest struct {
	RadarKey  string `json:"radarkey"`
	Timestamp int64  `json:"timestamp"`
}

// GetCommandsResponse 获取命令响应
type GetCommandsResponse struct {
	Code     int                       `json:"code"`
	Commands []mongodto.CommandDataDto `json:"commands"`
}

// GetCommands 雷达设备获取下发命令
// @Summary 雷达设备获取下发命令
// @Description 雷达设备获取下发命令
// @Tags 雷达管理-终端接口
// @Accept  application/json
// @Product application/json
// @Param data body GetCommandsRequest true "获取命令请求"
// @Success 200 {object} GetCommandsResponse "{\"code\": 0, \"commands\": [{\"command_code\": 100, \"message\": \"reboot\", \"parameters\": {}}]}"
// @Router /api/v1/radar/get_commands [post]
// @Security Bearer
func (e SysRadar) GetCommands(c *gin.Context) {
	e.MakeContext(c)
	var radarId int64
	var err error
	if radarId, err = e.GetTokenRadarId(c); err != nil {
		e.Error(400, err, err.Error())
		return
	}

	var req GetCommandsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		e.Error(400, err, "请求参数错误")
		return
	}

	cmds, err := mongosvr.QueryRadarComandData(radarId)
	if err != nil {
		cmds = make([]mongodto.CommandDataDto, 0)
	}
	// 这里应该是从数据库或其他存储中获取该雷达的待执行命令
	// 示例代码返回两个示例命令
	response := GetCommandsResponse{
		Code:     0,
		Commands: cmds,
	}

	fmt.Println("获得到命令：", len(cmds))
	e.OK(response, "获取命令成功")
}

// RawDataRequest 距离像数据请求
type RawDataRequest struct {
	RadarKey    string    `json:"radarkey"` // 从token中获取
	Timestamp   int64     `json:"timestamp"`
	CommandCode int       `json:"command_code"`
	Data        []float32 `json:"data"`
}

// PutRawData 雷达设备上传距离像数据
// @Summary 雷达设备上传距离像数据
// @Description 雷达设备上传距离像数据
// @Tags 雷达管理-终端接口
// @Accept  application/json
// @Product application/json
// @Param data body RawDataRequest true "距离像数据"
// @Success 200 {object} response.Response "{\"code\": 0, \"message\": \"距离像数据接收成功\"}"
// @Router /api/v1/radar/raw_data [post]
// @Security Bearer
func (e SysRadar) PutRawData(c *gin.Context) {
	e.MakeContext(c)
	var radarId int64
	var err error
	if radarId, err = e.GetTokenRadarId(c); err != nil {
		e.Error(400, err, err.Error())
		return
	}

	var req RawDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		e.Error(400, err, "请求参数错误")
		return
	}

	// 处理距离像数据逻辑
	e.Logger.Infof("接收到雷达 %s 的距离像数据: %+v", radarId, req)
	d := mongosvr.DistanceDataV2{}
	d.RadarID = radarId
	d.CommandCode = req.CommandCode
	d.TimeStamp = req.Timestamp
	d.Data = req.Data
	err = mongosvr.InsertDistanceDataV2(&d)
	if err != nil {
		e.Error(500, err, "距离像数据存储失败")
		return
	}

	e.OK(nil, "距离像数据接收成功")
}

// PutDeformation 雷达设备上传形变数据
// @Summary 雷达设备上传形变数据
// @Description 雷达设备上传形变数据
// @Tags 雷达管理-终端接口
// @Accept  application/json
// @Product application/json
// @Param data body DeformationRequest true "形变数据"
// @Success 200 {object} response.Response "{\"code\": 0, \"message\": \"形变数据接收成功\"}"
// @Router /api/v1/radar/put_deformation [post]
// @Security Bearer
func (e SysRadar) PutDeformation(c *gin.Context) {
	e.MakeContext(c)

	var radarId int64
	var err error
	if radarId, err = e.GetTokenRadarId(c); err != nil {
		e.Error(400, err, err.Error())
		return
	}

	var req DeformationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		e.Error(400, err, "请求参数错误")
		return
	}

	defData := mongosvr.DeformationData{}
	// 设置雷达ID
	defData.RadarId = radarId
	e.InitDeformationData(&defData, &req)
	// 调用mongosvr插入形变数据
	if err := mongosvr.InsertDeformationData(&defData); err != nil {
		e.Logger.Errorf("插入形变数据失败: %v", err)
		e.Error(500, err, "形变数据存储失败")
		return
	}

	e.Logger.Infof("接收到雷达 %s 的形变数据: %+v", radarId, req)

	e.OK(nil, "形变数据接收成功")
}

// PutStatus 雷达设备状态上报
// @Summary 雷达设备状态上报
// @Description 雷达设备状态上报
// @Tags 雷达管理-终端接口
// @Accept application/json
// @Product application/json
// @Param data body mongosvr.RadarStatusRequest true "设备状态数据"
// @Success 200 {object} response.Response "{\"code\": 0, \"message\": \"状态接收成功\"}"
// @Router /api/v1/radar/status [post]
// @Security Bearer
func (e SysRadar) PutStatus(c *gin.Context) {
	e.MakeContext(c)
	fmt.Println("雷达设备状态上报")
	var radarId int64
	var err error
	if radarId, err = e.GetTokenRadarId(c); err != nil {
		e.Error(400, err, err.Error())
		return
	}

	var req mongosvr.RadarStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		e.Logger.Error("终端状态信息上报参数错误：", err)
		e.Error(400, err, "请求参数错误")
		return
	}

	var status mongosvr.RadarStatus
	status.RadarStatusRequest = req
	status.RadarId = radarId

	if err := mongosvr.InsertRadarStatus(&status); err != nil {
		e.Logger.Errorf("插入状态数据失败: %v", err)
		e.Error(500, err, "状态数据存储失败")
		return
	}
	// 处理状态信息逻辑
	e.Logger.Infof("接收到雷达 %s 的状态信息: %+v", req.RadarKey, req)

	e.OK(nil, "状态接收成功")
}

// PutDevInfo 雷达设备信息上报
// @Summary 雷达设备信息上报
// @Description 雷达设备信息上报
// @Tags 雷达管理-终端接口
// @Accept application/json
// @Product application/json
// @Param data body mongosvr.RadarDevInfoRequest true "设备信息数据"
// @Success 200 {object} response.Response "{\"code\": 0, \"message\": \"设备信息接收成功\"}"
// @Router /api/v1/radar/dev_info [post]
// @Security Bearer
func (e SysRadar) PutDevInfo(c *gin.Context) {
	e.MakeContext(c)

	var radarId int64
	var err error
	if radarId, err = e.GetTokenRadarId(c); err != nil {
		e.Error(400, err, err.Error())
		return
	}

	var req mongosvr.RadarDevInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		e.Error(400, err, "请求参数错误")
		return
	}

	var status mongosvr.RadarDevInfo
	status.RadarDevInfoRequest = req
	status.RadarId = radarId

	if err := mongosvr.InsertRadarDevInfo(&status); err != nil {
		e.Logger.Errorf("插入设备信息数据失败: %v", err)
		e.Error(500, err, "设备信息数据存储失败")
		return
	}
	// 处理设备信息逻辑
	e.Logger.Infof("接收到雷达 %s 的设备信息: %+v", req.RadarKey, req)

	e.OK(nil, "设备信息接收成功")
}
