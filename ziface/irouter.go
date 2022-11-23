package ziface

type IRouter interface {
	// PreHandle 处理业务之前的hook
	PreHandle(request IRequest)
	// Handle 处理业务主方法的hook
	Handle(request IRequest)
	// PostHandle 处理业务之后的hook
	PostHandle(request IRequest)
}
