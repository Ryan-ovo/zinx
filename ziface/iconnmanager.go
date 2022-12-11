package ziface

type IConnManager interface {
	Add(conn IConnection)                   // 添加连接
	Remove(conn IConnection)                // 删除连接
	Get(connID uint32) (IConnection, error) // 获取连接
	Len() int                               // 获取连接数
	ClearConn()                             // 清空连接
	ClearOneConn(connID uint32)             // 清空指定连接
}
