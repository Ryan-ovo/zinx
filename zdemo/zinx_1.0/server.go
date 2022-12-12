package main

import (
	"log"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	log.Println("Call Ping Router Handle...")
	log.Println("receive msg from client: msgID =", request.GetMsgID(), "data = ", string(request.GetData()))
	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping...")); err != nil {
		log.Println("err")
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

func (h *HelloRouter) Handle(request ziface.IRequest) {
	log.Println("Call hello Router Handle...")
	log.Println("receive msg from client: msgID =", request.GetMsgID(), "data = ", string(request.GetData()))
	if err := request.GetConnection().SendMsg(2, []byte("Hello Zinx Router V0.9")); err != nil {
		log.Println("err")
	}
}

func DoConnBegin(conn ziface.IConnection) {
	log.Println("---> DoConnBegin is called...")
	if err := conn.SendMsg(250, []byte("hello world")); err != nil {
		log.Println(err)
	}
	conn.SetProperty("name", "Ryan")
	conn.SetProperty("age", "18")
}

func DoConnEnd(conn ziface.IConnection) {
	log.Println("--> DoConnEnd is called...")
	log.Printf("connID = [%d] is offline", conn.GetConnID())

	// 获取连接属性
	if v, err := conn.GetProperty("name"); err == nil {
		log.Printf("key = [%s], value = [%+v]\n", "name", v)
	}
	if v, err := conn.GetProperty("age"); err == nil {
		log.Printf("key = [%s], value = [%+v]]\n", "age", v)
	}
}

func main() {
	s := znet.NewServer()

	// 注册连接的钩子函数
	s.SetAfterConnStart(DoConnBegin)
	s.SetBeforeConnStop(DoConnEnd)

	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})

	s.Serve()
}
