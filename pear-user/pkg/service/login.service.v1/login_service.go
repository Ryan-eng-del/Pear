package login_service_v1

import (
	"context"
	"errors"
	"log"
	"time"

	common "cyan.com/pear-common"
	"cyan.com/pear-user/pkg/dao"
	"cyan.com/pear-user/pkg/repo"
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
		return nil, errors.New("code error")
	}


code := "123456"
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