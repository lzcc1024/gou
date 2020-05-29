package unet

import "github.com/lzcc1024/gou/uiface"

type BaseRouter struct{}

//在处理conn业务之前的钩子方法
func (r *BaseRouter) PreHandle(request uiface.IRequest) {}

//处理conn业务的方法
func (r *BaseRouter) Handle(request uiface.IRequest) {}

//处理conn业务之后的钩子方法
func (r *BaseRouter) PostHandle(request uiface.IRequest) {}
