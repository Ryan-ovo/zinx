package main

import (
	"log"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) PreHandle(request ziface.IRequest) {
	log.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping\n"))
	if err != nil {
		log.Println("call back before ping error")
	}
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	log.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("main ping\n"))
	if err != nil {
		log.Println("call back before ping error")
	}
}

func (p *PingRouter) PostHandle(request ziface.IRequest) {
	log.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping\n"))
	if err != nil {
		log.Println("call back before ping error")
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()
}
