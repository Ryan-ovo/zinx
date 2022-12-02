package znet

import (
	"errors"
	"io"
	"log"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn       *net.TCPConn       // 当前连接的 socket
	ConnID     uint32             // 连接ID
	isClosed   bool               // 当前的连接状态
	ExitChan   chan bool          // 通知当前连接停止的 channel
	MsgHandler ziface.IMsgHandler // 消息管理模块
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgHandler: msgHandler,
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
		// v0.5 集成拆包过程
		dp := NewDataPack()
		// 读取客户端的消息头
		msgHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, msgHead); err != nil {
			log.Printf("Read Message Head Error = [%+v]\n", err)
			c.ExitChan <- true
			continue
		}
		// 拆包，把消息id和消息长度读取到dp对象中
		msg, err := dp.Unpack(msgHead)
		if err != nil {
			log.Printf("unpack error = [%+v]\n", err)
			c.ExitChan <- true
			continue
		}
		// 根据消息长度读取消息体，把消息体读取到dp对象中
		var msgBody []byte
		if msg.GetMsgLen() > 0 {
			msgBody = make([]byte, msg.GetMsgLen())
			if _, err = io.ReadFull(c.Conn, msgBody); err != nil {
				log.Printf("Read Message Data Error = [%+v]\n", err)
				c.ExitChan <- true
				continue
			}
		}
		msg.SetData(msgBody)
		req := &Request{
			conn: c,
			msg:  msg,
		}
		go c.MsgHandler.DoMsgHandler(req)
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

func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection already closed")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		log.Printf("pack message error = [%+v]\n", err)
	}
	if _, err = c.Conn.Write(msg); err != nil {
		log.Printf("Write message error = [%+v], msg_id = [%d]", err, msgID)
		c.ExitChan <- true
		return errors.New("conn write error")
	}
	return nil
}
