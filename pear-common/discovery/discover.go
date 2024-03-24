package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Register struct {
	EtcdAddrs []string
	DialTimeout int
	closeCh chan struct{}
	leasesID clientv3.LeaseID
	keepAliveCh <- chan *clientv3.LeaseKeepAliveResponse
	srvInfo Server
	srvTTL int64
	cli *clientv3.Client
	logger *zap.Logger
}

func NewRegister (etcdAddrs []string, logger *zap.Logger) *Register {
	return &Register{
		EtcdAddrs: etcdAddrs,
		DialTimeout: 10,
		logger: logger,
	}
}

func (r *Register) Register(srv Server,  ttl int64) (chan <- struct{}, error) {
	var err error
	if strings.Split(srv.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip")
	}


	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints: r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	r.srvInfo = srv
	r.srvTTL = ttl


	if err := r.register(); err != nil {
		return nil, err
	}

	r.closeCh = make(chan struct{})
	go r.keepAlive()
	return r.closeCh, nil
}


func (r Register) unregister() error {
	_ ,err := r.cli.Delete(context.Background(), BuildPrefix(r.srvInfo))
	return err
}


func (r *Register) register() error {
	leaseCtx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout) * time.Second)

	defer cancel()
	leaseResp, err := r.cli.Grant(leaseCtx, r.srvTTL)
	if err != nil {
		return err
	}

	r.leasesID = leaseResp.ID

	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), leaseResp.ID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)

	if err != nil {
		return err
	}

	r.logger.Info(BuildRegPath(r.srvInfo))
	_, err = r.cli.Put(context.Background(), BuildRegPath(r.srvInfo), string(data))
	return err
}

func (r *Register) keepAlive() {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		case <- r.closeCh:
			if err := r.unregister(); err != nil {
				r.logger.Error("unregister failed", zap.Error(err))
			}

			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				r.logger.Error("revoke failed", zap.Error(err))
			}
			return
		case res := <- r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}
		case <- ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}
		}
	}
} 

func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}


func (r *Register) UpdateHandler () http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request){
		wi := req.URL.Query().Get("weight")
		weight, err := strconv.Atoi(wi)

		if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return 
		}

		var update = func () error {
			r.srvInfo.Weight = int64(weight)
			data, err := json.Marshal(r.srvInfo)
			if err != nil {
				
				return err
			}
			_, err = r.cli.Put(context.Background(), BuildRegPath(r.srvInfo), string(data))
			return err
		}
	
		
		if err := update(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return 
		}
		w.Write([]byte("update server weight succes"))
	})
}


func (r *Register) GetServerInfo ()(Server, error) {
	resp, err := r.cli.Get(context.Background(), BuildRegPath(r.srvInfo))

	if err != nil {
		return r.srvInfo, err
	}

	info := Server{}

	if resp.Count >= 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &info); err != nil {
			return info, err
		}
	}

	return info, nil

}



