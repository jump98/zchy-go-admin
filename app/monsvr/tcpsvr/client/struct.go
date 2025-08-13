package client

// type DevPackage struct {
// 	Head uint32 //固定0x7E7E9C9C
// 	DevPackageNoHeader
// }
// type DevPackageNoHeader struct {
// 	DevType   uint32 //设备类型
// 	Cmd       uint32 //命令码
// 	CmdExtend uint32 //扩展码
// 	DataLen   uint32 //数据长度
// }

// type DevProjectInfo struct {
// 	ProID  uint32  //项目编号
// 	DevID  uint32  //设备编号，雷达。。。
// 	DevLng float64 //经度
// 	DevLat float64 //纬度
// 	CRC32  uint32
// }

// // 用于返回的信息，不带CRC，因为CRC在发送时计算
// type RetDevProjectInfo struct {
// 	Code uint8 //返回值
// }

// type RadarDeformatQueryInfo struct {
// 	DefUpdateFreq         uint16  //形变更新率uint16
// 	DefSendInterval       uint16  //形变发送间隔uint16
// 	FreqReadyTime         uint16  //频率准备时间 uint16
// 	FreqSendInterval      uint16  //频率发送间隔uint16
// 	ADI                   float32 //ADI float32
// 	ADIFactor             float32 //ADI系数 float32
// 	SmoothThreshold       uint16  //平滑门限uint16
// 	FreqComputeRange      uint8   //频率计算范围uint8
// 	SmoothTimes           float32 //平滑倍数float32
// 	SmoothLevel           float32 //平滑程度float32
// 	AntiObstruction       uint8   //防遮挡开启uint8
// 	AtmosphericCorrection uint8   //大气校正开启uint8
// 	PeakTracking          uint8   //峰值跟踪开启uint8
// 	Reserve               uint8   //备用uint8
// 	WorkMode              uint8   //工作模式uint8
// 	Band                  uint16  //带宽uint16
// 	CRC32                 uint32
// }

// type RadarAddCtrlPointInfo struct {
// 	RetValue uint8 //返回值
// 	CRC32    uint32
// }

// type RadarHeartBreakInfo struct {
// 	Option     uint8  //选项
// 	TimeStamp  int32  //时间戳
// 	DetectNum  int32  //检测点数量
// 	TestTimes  int32  //测量次数
// 	FrameNO    uint16 //帧数
// 	FrameTimes uint16 //帧数循环次数
// }
