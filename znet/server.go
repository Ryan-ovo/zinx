package znet

import (
	"fmt"
	"log"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name      string // 服务器名
	IPVersion string // tcp网络名称
	IP        string // ip地址
	Port      int    // 端口号
	Router    ziface.IRouter
}

func (s *Server) Start() {
	log.Printf("[Zinx] Server name = [%+v], listen at ip = [%+v], port = [%+v]\n", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.Port)
	log.Printf("[Zinx] Version = [%+v], MaxConn = [%+v], MaxPkgSize = [%+v]\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPkgSize)
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
	var connID uint32 = 0
	for {
		// 3.1 阻塞等待客户端建立连接请求
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Accept err = [%+v]\n", err)
			continue
		}
		zinxConn := NewConnection(conn, connID, s.Router)
		connID++
		go zinxConn.Start()
	}
}

func (s *Server) Stop() {
	log.Println("[STOP] Zinx server , name = ", s.Name)
	// todo: 将服务器的资源，状态或者一些已经开辟的链接信息进行回收或者停止
}

func (s *Server) Serve() {
	// 监听客户端连接的操作单独用协程去承载，避免主协程因为没有客户端连接而阻塞
	go s.Start()
	// todo: 启动服务后的额外业务

	// 阻塞主协程
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	log.Println("Add Router Success!")
	s.Router = router
}

/*
	初始化Server的方法
*/
func NewServer() ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.Port,
		Router:    nil,
	}
	return s
}
