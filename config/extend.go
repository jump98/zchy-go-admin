package config

var ExtConfig Extend

// Extend 扩展配置
//
//	extend:
//	  demo:
//	    name: demo-name
//
// 使用方法： config.ExtConfig......即可！！
type Extend struct {
	MongoDB MongoDB `json:"mongodb"`
}

// type AMap struct {
// 	Key string
// }

type MongoDB struct {
	//连接地址
	Source string `json:"source"`
	//连接地址
	RadarDBName string `json:"radar_db_name"`
}
