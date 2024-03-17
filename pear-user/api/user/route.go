package user

import (
	router "cyan.com/pear-user/route"
	"github.com/gin-gonic/gin"
)

type RouterUser struct {}

func init() {
	router.Register(&RouterUser{})
}

func (*RouterUser) Register(r *gin.Engine) {
	c := New()
	r.GET("/ping", c.GetCaptcha)
}
