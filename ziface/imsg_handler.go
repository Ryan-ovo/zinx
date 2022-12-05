package ziface

type IMsgHandler interface {
	// DoMsgHandler 以非阻塞的方式处理消息，执行对应的路由业务方法
	DoMsgHandler(request IRequest)
	// AddRouter 为消息绑定相应的处理逻辑handler
	AddRouter(msgID uint32, router IRouter)
	// StartWorkerPool 启动worker工作池
	StartWorkerPool()
	// SendMsgToTaskQueue 将消息交给对应的worker进行处理
	SendMsgToTaskQueue(request IRequest)
}
