package znet

import "zinx/ziface"

type Request struct {
	conn      ziface.IConnection
	msg       ziface.IMessage
	requestID uint32
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}

func (r *Request) GetRequestID() uint32 {
	return r.requestID
}
