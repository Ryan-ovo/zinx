package znet

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"zinx/ziface"
)

type Connection struct {
	TcpServer    ziface.IServer         // 连接隶属的服务器
	Conn         *net.TCPConn           // 当前连接的 socket
	ConnID       uint32                 // 连接ID
	isClosed     bool                   // 当前的连接状态
	ExitChan     chan bool              // 通知当前连接停止的 channel
	MsgHandler   ziface.IMsgHandler     // 消息管理模块
	msgChan      chan []byte            // 读写协程之间的通信
	property     map[string]interface{} // 连接属性
	propertyLock sync.RWMutex           // 保护连接属性读写的锁
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}
	c.TcpServer.GetConnManager().Add(c)
	return c
}

func (c *Connection) Start() {
	log.Println("Connection start, ConnID = ", c.ConnID)
	// 启动从当前连接读取数据的业务
	go c.StartReader()
	// 启动从当前连接写数据到客户端的业务
	go c.StartWriter()
	// 调用钩子函数
	c.TcpServer.CallAfterConnStart(c)
}

func (c *Connection) StartReader() {
	log.Println("Reader Goroutine Start is running...")
	defer log.Println("connID = ", c.ConnID, ", Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	var ReqID uint32
	for {
		// v0.5 集成拆包过程
		dp := NewDataPack()
		// 读取客户端的消息头
		msgHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, msgHead); err != nil {
			log.Printf("Read Message Head Error = [%+v]\n", err)
			break
		}
		// 拆包，把消息id和消息长度读取到dp对象中
		msg, err := dp.Unpack(msgHead)
		if err != nil {
			log.Printf("unpack error = [%+v]\n", err)
			break
		}
		// 根据消息长度读取消息体，把消息体读取到dp对象中
		var msgBody []byte
		if msg.GetMsgLen() > 0 {
			msgBody = make([]byte, msg.GetMsgLen())
			if _, err = io.ReadFull(c.Conn, msgBody); err != nil {
				log.Printf("Read Message Data Error = [%+v]\n", err)
				break
			}
		}
		msg.SetData(msgBody)
		req := &Request{
			conn:      c,
			msg:       msg,
			requestID: ReqID,
		}
		ReqID++
		c.MsgHandler.SendMsgToTaskQueue(req)
	}
}

func (c *Connection) StartWriter() {
	log.Println("Writer Goroutine Start is running...")
	defer log.Println("connID = ", c.ConnID, ", Writer is exit, remote addr is ", c.RemoteAddr().String())

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				log.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
		case <-c.ExitChan:
			// 如果连接已经关闭，直接退出即可
			return
		}
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
	// 通知writer关闭
	c.ExitChan <- true

	// 调用连接关闭前的钩子函数
	c.TcpServer.CallBeforeConnStop(c)

	// 删除连接
	c.TcpServer.GetConnManager().Remove(c)

	// 关闭相关通道
	close(c.ExitChan)
	close(c.msgChan)
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

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection already closed")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		log.Printf("pack message error = [%+v]\n", err)
	}
	c.msgChan <- msg
	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("property not found")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
