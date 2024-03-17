package user

import (
	"log"
	"time"

	common "cyan.com/pear-common"
	"cyan.com/pear-user/pkg/dao"
	"cyan.com/pear-user/pkg/model"
	"cyan.com/pear-user/pkg/repo"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type ControllerUser struct {
	Cache repo.Cache
}

func New() *ControllerUser {
	return &ControllerUser{
		Cache: dao.RC,
	}
}

func (h *ControllerUser) GetCaptcha(c *gin.Context) {
	rsp := &common.Result{}
	mobile := c.PostForm("mobile")

	if !common.VerifyMobile(mobile) {
		c.JSON(200, rsp.Fail(model.NoLegalMobile, "invalid mobile"))
		return
	}

	code := "123456"

	go func (){
		time.Sleep(2 * time.Second)
		c, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		err := h.Cache.Put(c, "REGISTER_" + mobile, code, 15 * time.Minute)

		if err != nil {
			log.Println(err, "redis put error")
		}
	}()

	c.JSON(200, rsp.Success(code))
}

