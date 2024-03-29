package login_service_v1

import (
	"context"
	"log"
	"time"

	common "cyan.com/pear-common"
	"cyan.com/pear-common/errs"
	"cyan.com/pear-user/pkg/dao"
	"cyan.com/pear-user/pkg/model"
	"cyan.com/pear-user/pkg/repo"
	"cyan.com/pear-user/pkg/util"
)

type LoginService struct {
	UnimplementedLoginServiceServer
	Cache repo.Cache
}

func New() *LoginService {
	return &LoginService {
		Cache: dao.RC,
	}
}

func (s *LoginService) GetCaptcha(ctx context.Context, req *CaptchaReq) (*CaptchResp, error) {
	mobile := req.Mobile
	if !common.VerifyMobile(mobile) {
		return nil, errs.GrpcError(model.NoLegalMobile)
	}  

	code := util.RandCode(6, util.TYPE_MIXED)

	go func (){
		time.Sleep(2 * time.Second)
		c, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		err := s.Cache.Put(c, "REGISTER_" + mobile, code, 15 * time.Minute)

		if err != nil {
			log.Println(err, "redis put error")
		}
	}()

return &CaptchResp{Code: code}, nil
}