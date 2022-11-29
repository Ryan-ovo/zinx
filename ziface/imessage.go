package ziface

type IMessage interface {
	/*
		getter
	*/
	GetMsgID() uint32
	GetMsgLen() uint32
	GetData() []byte

	/*
		setter
	*/
	SetMsgID(id uint32)
	SetMsgLen(len uint32)
	SetData(data []byte)
}
