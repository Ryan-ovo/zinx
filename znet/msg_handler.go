package znet

import (
	"log"
	"zinx/ziface"
)

type MsgHandler struct {
	APIs map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		APIs: make(map[uint32]ziface.IRouter),
	}
}

func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	msgID := request.GetMsgID()
	handler, ok := m.APIs[msgID]
	if !ok {
		log.Println("api not found, msg_id = ", msgID)
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := m.APIs[msgID]; ok {
		log.Println("router has been registered")
		return
	}
	m.APIs[msgID] = router
	log.Println("Add api success, msg_id = ", msgID)
}
