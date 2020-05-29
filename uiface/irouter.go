package uiface

type IRouter interface {
	//在处理conn业务之前的钩子方法
	PreHandle(request IRequest)
	//处理conn业务的方法
	Handle(request IRequest)
	//处理conn业务之后的钩子方法
	PostHandle(request IRequest)
}
