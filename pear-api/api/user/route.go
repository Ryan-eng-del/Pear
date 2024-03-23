package user

import (
	router "cyan.com/pear-api/route"
	"github.com/gin-gonic/gin"
)

type RouterUser struct {}

func init() {
	router.Register(&RouterUser{})
}

func (*RouterUser) Register(r *gin.Engine) {
	InitUserClient()
	c := New()
	r.POST("/getCaptcha", c.GetCaptcha)
}
