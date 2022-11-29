package ziface

/*
	封包，拆包模块
	直接面向TCP连接中的数据流，处理TCP粘包问题
*/

type IDataPack interface {
	// GetHeadLen 获取数据包头部长度
	GetHeadLen() uint32
	// Pack 封包
	Pack(msg IMessage) ([]byte, error)
	// Unpack 拆包
	Unpack([]byte) (IMessage, error)
}
