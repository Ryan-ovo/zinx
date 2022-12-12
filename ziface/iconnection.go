package ziface

import "net"

type IConnection interface {
	// Start 启动连接
	Start()
	// Stop 停止连接
	Stop()
	// GetTCPConnection 获取当前连接绑定的socket
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取连接ID
	GetConnID() uint32
	// RemoteAddr 获取远程客户端的连接类型，IP和端口
	RemoteAddr() net.Addr
	// SendMsg 发送封包后的数据
	SendMsg(msgID uint32, data []byte) error
	// SetProperty 设置连接属性
	SetProperty(key string, value interface{})
	// GetProperty 获取连接属性
	GetProperty(key string) (interface{}, error)
	// RemoveProperty 移除连接属性
	RemoveProperty(key string)
}

/*
	HandleFunc 抽象方法
	1. socket连接
	2. 客户端请求的数据
	3. 客户端请求的数据长度
*/
// 实现router类后，这个抽象方法不需要了
//type HandleFunc func(*net.TCPConn, []byte, int) error
