package znet

import (
	"zinx/ziface"
)

/*
	实现router时，先嵌入基类，并且根据需要重写其中的一个或多个方法
	并且基类的所有方法实现都为空，是因为有些router不需要某些hook，继承baseRouter后不需要实现所有的hook也可以被实例化
*/
type BaseRouter struct {
}

func (b *BaseRouter) PreHandle(request ziface.IRequest) {}

func (b *BaseRouter) Handle(request ziface.IRequest) {}

func (b *BaseRouter) PostHandle(request ziface.IRequest) {}
