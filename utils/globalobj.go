package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

// 全局对象
var GlobalObject *GlobalObj

type GlobalObj struct {
	// server相关
	TcpServer ziface.IServer // 全局server对象
	Host      string         // ip
	Port      int            // 端口号
	Name      string         // 服务器名
	// zinx相关
	Version           string // zinx版本号
	MaxConn           int    // 最大连接数
	MaxPkgSize        uint32 //数据包的最大值
	WorkerPoolSize    uint32 // 工作池的最大线程数
	MaxWorkerTaskSize uint32 // 每个协程允许排队的最大任务数
}

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("zdemo/zinx_0.8/conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Host:              "0.0.0.0",
		Port:              8999,
		Name:              "ZinxServerApp",
		Version:           "v0.4",
		MaxConn:           1000,
		MaxPkgSize:        4096,
		WorkerPoolSize:    10,
		MaxWorkerTaskSize: 1024,
	}
	GlobalObject.Reload()
}
