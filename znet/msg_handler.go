package znet

import (
	"log"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandler struct {
	APIs           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		APIs:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

func (m *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan ziface.IRequest)
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

func (m *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	log.Printf("Worker[%+v] is started...\n", workerID)
	for {
		select {
		case req := <-taskQueue:
			m.DoMsgHandler(req)
		}
	}
}

func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 根据RequestID分配当前消息应该由哪个worker处理
	workerID := request.GetRequestID() % m.WorkerPoolSize
	log.Printf("Send Request to worker[%d], MsgID = [%d], RequestID = [%d]", workerID, request.GetMsgID(), request.GetRequestID())
	// 将请求发送给消息队列
	m.TaskQueue[workerID] <- request
}
