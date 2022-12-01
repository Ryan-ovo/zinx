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
	log.Println("Call Router Handle...")
	log.Println("receive msg from client: msgID =", request.GetMsgID(), "data = ", request.GetData())
	if err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping...")); err != nil {
		log.Println("err")
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()
}
