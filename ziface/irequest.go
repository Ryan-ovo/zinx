package ziface

/*
	把客户端请求的连接信息和请求数据绑定到一起，封装为Request
*/
type IRequest interface {
	// GetConnection 获取连接
	GetConnection() IConnection
	// GetData 获取请求的消息数据
	GetData() []byte
	// GetMsgID 获取请求的ID
	GetMsgID() uint32
}
