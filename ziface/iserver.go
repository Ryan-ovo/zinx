package ziface

type IServer interface {
	Start()
	Stop()
	Serve()
	// AddRouter 给服务器注册路由
	AddRouter(msgID uint32, router IRouter)
	GetConnManager() IConnManager
	SetAfterConnStart(func(conn IConnection)) // 设置连接启动之后调用的钩子函数
	SetBeforeConnStop(func(conn IConnection)) // 设置连接关闭之前的钩子函数
	CallAfterConnStart(conn IConnection)      // 调用连接启动之后的钩子函数
	CallBeforeConnStop(conn IConnection)      // 调用连接关闭之前的钩子函数
}
