package user

import (
	"context"
	"net/http"
	"time"

	login_service_v2 "cyan.com/pear-api/pkg/service/login.service.v1"
	common "cyan.com/pear-common"
	"cyan.com/pear-user/pkg/dao"
	"cyan.com/pear-user/pkg/repo"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
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
	result := &common.Result{}
	params := &GetCaptchaInput{}
	err := c.ShouldBind(params)

	if err != nil {
		c.JSON(http.StatusOK, result.Fail(common.BusinessCode(3001), "mobile must be input"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	rsp, err := LoginServiceClient.GetCaptcha(ctx, &login_service_v2.CaptchaReq{
		Mobile: params.Mobile,
	})

	if err != nil {
		err, _ := status.FromError(err)
		c.JSON(http.StatusOK, result.Fail(common.BusinessCode(err.Code()), err.Message()))
		return
	}

	c.JSON(http.StatusOK, result.Success(rsp.Code))
}

