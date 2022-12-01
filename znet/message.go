package znet

type Message struct {
	ID      uint32 // 消息id
	DataLen uint32 // 消息长度
	Data    []byte // 消息内容
}

func NewMessage(msgID uint32, data []byte) *Message {
	return &Message{
		ID:      msgID,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (m *Message) GetMsgID() uint32 {
	return m.ID
}

func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}

func (m *Message) SetMsgLen(len uint32) {
	m.DataLen = len
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}
