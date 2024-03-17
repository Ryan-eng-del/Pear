package router

import (
	"github.com/gin-gonic/gin"
)

// Register Item/Objective
type Router interface {
	Register(r *gin.Engine) 
}


// 注册 center 方式1 在自身初始化
// 这种方式缺点是 得需要手动 _ 引入，方式2直接依赖
var Routers []Router
func Register(router ...Router)  {
	Routers = append(Routers, router...)
}

func InitRouter(r *gin.Engine) {
	for _, router := range Routers {
		router.Register(r)
	}
}


// 注册 center 方式 2, 在 center 初始化
// type RegisterRouter struct {}
// func New() *RegisterRouter {
// 	return &RegisterRouter{}
// }

// func (*RegisterRouter) Register(r *gin.Engine) {
// 	userRouter := &user.RouterUser{}
// 	userRouter.Register(r)
// }