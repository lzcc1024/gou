package unet

import (
	"fmt"
	"strconv"

	"github.com/lzcc1024/gou/uiface"
	"github.com/lzcc1024/gou/utils"
)

type Msghandler struct {
	//  存放每个MsgId 所对应的处理方法的map属性
	Apis map[uint32]uiface.IRouter
	//业务工作Worker池的数量
	WorkerPoolSize uint32
	//Worker负责取任务的消息队列
	TaskQueue []chan uiface.IRequest
}

func NewMsgHandle() *Msghandler {
	return &Msghandler{
		Apis:           make(map[uint32]uiface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan uiface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//马上以非阻塞方式处理消息
func (m *Msghandler) DoMsgHandler(request uiface.IRequest) {
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (m *Msghandler) AddRouter(msgId uint32, router uiface.IRouter) {
	if _, ok := m.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	m.Apis[msgId] = router
}

//启动一个Worker工作流程
func (m *Msghandler) StartOneWorker(workerID int, taskQueue chan uiface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

//启动worker工作池
func (m *Msghandler) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		m.TaskQueue[i] = make(chan uiface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

//将消息交给TaskQueue,由worker进行处理
func (m *Msghandler) SendMsgToTaskQueue(request uiface.IRequest) {

	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMsgID(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	m.TaskQueue[workerID] <- request
}
