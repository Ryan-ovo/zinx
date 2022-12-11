package znet

import (
	"fmt"
	"log"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	Name               string                        // 服务器名
	IPVersion          string                        // tcp网络名称
	IP                 string                        // ip地址
	Port               int                           // 端口号
	msgHandler         ziface.IMsgHandler            // 消息管理模块
	ConnMgr            ziface.IConnManager           // 连接管理器
	AfterConnStartFunc func(conn ziface.IConnection) // 连接启动之后调用的钩子函数
	BeforeConnStopFunc func(conn ziface.IConnection) // 连接关闭之后调用的钩子函数
}

func (s *Server) Start() {
	// 启动工作池
	s.msgHandler.StartWorkerPool()
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
		// 3.2 最大连接数控制
		log.Println(">>>>>> ConnMgr Len = ", s.ConnMgr.Len(), ", MAX = ", utils.GlobalObject.MaxConn)
		if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
			log.Printf("conn num exceed limit, MaxConnSize = [%d]\n", utils.GlobalObject.MaxConn)
			if err = conn.Close(); err != nil {
				log.Printf("conn close error = [%+v]\n", err)
			}
			continue
		}
		zinxConn := NewConnection(s, conn, connID, s.msgHandler)
		connID++
		go zinxConn.Start()
	}
}

func (s *Server) Stop() {
	log.Println("[STOP] Zinx server , name = ", s.Name)
	// todo: 将服务器的资源，状态或者一些已经开辟的链接信息进行回收或者停止
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	// 监听客户端连接的操作单独用协程去承载，避免主协程因为没有客户端连接而阻塞
	go s.Start()
	// todo: 启动服务后的额外业务

	// 阻塞主协程
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
	log.Println("Add Router Success!, msg_id = ", msgID)
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnMgr
}

/*
	初始化Server的方法
*/
func NewServer() ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.Port,
		msgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

func (s *Server) SetAfterConnStart(do func(conn ziface.IConnection)) {
	s.AfterConnStartFunc = do
}

func (s *Server) SetBeforeConnStop(do func(conn ziface.IConnection)) {
	s.BeforeConnStopFunc = do
}

func (s *Server) CallAfterConnStart(conn ziface.IConnection) {
	if s.AfterConnStartFunc != nil {
		s.AfterConnStartFunc(conn)
	}
}

func (s *Server) CallBeforeConnStop(conn ziface.IConnection) {
	if s.BeforeConnStopFunc != nil {
		s.BeforeConnStopFunc(conn)
	}
}
