package znet

import (
	"fmt"
	"log"
	"net"
	"zinx/ziface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func (s *Server) Start() {
	log.Printf("[Start] Server Listener at IP: %s, Port: %d, is starting\n", s.IP, s.Port)
	// 1. 获取一个tcp的addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		log.Printf("resolve tcp addr err = [%+v]\n", err)
		return
	}
	// 2. 监听服务器的地址 listen
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		log.Printf("listen %s, err = [%+v]\n", s.IPVersion, err)
		return
	}
	// 监听成功
	log.Println("Zinx Server start success, listening...")
	// 3. 阻塞等待客户端连接，处理客户端连接的业务
	for {
		// 3.1 阻塞等待客户端建立连接请求
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Accept err = [%+v]\n", err)
			continue
		}
		// 对于每一个客户端连接，开协程去处理业务
		go func() {
			for {
				// 客户端输入的回显任务
				buf := make([]byte, 512)
				cnt, err := conn.Read(buf)
				if err != nil {
					log.Printf("receive buf err = [%+v]\n", err)
					continue
				}
				fmt.Println(cnt)
				// 回显
				if _, err := conn.Write(buf[:cnt]); err != nil {
					log.Printf("write back buf err = [%+v]\n", err)
					continue
				}
			}
		}()
	}
}

func (s *Server) Stop() {
	// todo: 将服务器的资源，状态或者一些已经开辟的链接信息进行回收或者停止
}

func (s *Server) Serve() {
	// 监听客户端连接的操作单独用协程去承载，避免主协程因为没有客户端连接而阻塞
	go s.Start()
	// todo: 启动服务后的额外业务

	// 阻塞主协程
	select {}
}

/*
	初始化Server的方法
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
