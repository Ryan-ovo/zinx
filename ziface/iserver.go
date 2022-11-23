package ziface

type IServer interface {
	Start()
	Stop()
	Serve()
	// AddRouter 给服务器注册路由
	AddRouter(router IRouter)
}
