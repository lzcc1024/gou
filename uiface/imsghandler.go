package uiface

// 定义 IMsgHandle 接口
type IMsgHandle interface {
	//马上以非阻塞方式处理消息
	DoMsgHandler(IRequest)
	//为消息添加具体的处理逻辑
	AddRouter(uint32, IRouter)
	//启动worker工作池
	StartWorkerPool()
	//将消息交给TaskQueue,由worker进行处理
	SendMsgToTaskQueue(request IRequest)
}
