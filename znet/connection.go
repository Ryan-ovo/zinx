package znet

import (
	"log"
	"net"
	"zinx/ziface"
)

type Connection struct {
	// 当前连接的 socket
	Conn *net.TCPConn
	// 连接ID
	ConnID uint32
	// 当前的连接状态
	isClosed bool
	// 通知当前连接停止的 channel
	ExitChan chan bool
	// 处理该连接的router
	Router ziface.IRouter
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router:   router,
	}
}

func (c *Connection) Start() {
	log.Println("Connection start, ConnID = ", c.ConnID)
	// 启动从当前连接读取数据的业务
	go c.StartReader()
	// todo: 启动从当前连接写入数据的业务
}

func (c *Connection) StartReader() {
	log.Println("Reader Goroutine Start is running...")
	defer log.Println("connID = ", c.ConnID, ", Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			log.Println("receive buf error = ", err)
			continue
		}
		req := &Request{
			conn: c,
			data: buf[:cnt],
		}
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(req)
		}(req)
	}
}

func (c *Connection) Stop() {
	log.Println("Connection stop, connID = ", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	// 关闭socket连接
	_ = c.Conn.Close()
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	//TODO implement me
	panic("implement me")
}
