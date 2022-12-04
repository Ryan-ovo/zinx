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
	log.Println("receive msg from client: msgID =", request.GetMsgID(), "data = ", request.GetData())
	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping...")); err != nil {
		log.Println("err")
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

func (h *HelloRouter) Handle(request ziface.IRequest) {
	log.Println("Call hello Router Handle...")
	log.Println("receive msg from client: msgID =", request.GetMsgID(), "data = ", request.GetData())
	if err := request.GetConnection().SendMsg(2, []byte("Hello Zinx Router V0.6")); err != nil {
		log.Println("err")
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(1, &PingRouter{})
	s.AddRouter(2, &HelloRouter{})
	s.Serve()
}
