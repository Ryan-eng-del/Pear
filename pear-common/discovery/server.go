package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/resolver"
)

type Server struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
	Version string `json:"version"`
	Weight int64 `json:"weight"`
}

func BuildPrefix (info Server) string {
	if info.Version == "" {
		return fmt.Sprintf("/%s/", info.Name)
	}
	return fmt.Sprintf("/%s/%s", info.Name, info.Version)
}


func BuildRegPath (info Server) string {
	return fmt.Sprintf("%s%s", BuildPrefix(info), info.Addr)
}

func ParseValue(value []byte) (Server, error) {
	info := Server{}
	if err := json.Unmarshal(value, &info); err != nil {
		return info, err
	}
	return info, nil
}

func SpitPath(path string) (Server, error ) {
	info := Server{}
	strs := strings.Split(path, "/")
	if len(strs) == 0 {
		return info, errors.New("invalid path")
	}
	info.Addr = strs[len(strs)-1]
	return info, nil
}

func Exists(l []resolver.Address, addr resolver.Address) bool {
	for i := range l {
		if l[i].Addr == addr.Addr {
			return true
		}
	}
	return false
}

func Remove(s []resolver.Address, addr resolver.Address) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr.Addr {
			s[i] = s[len(s) - 1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}

func BuildResolveUrl(app string) string {
	return schema + ":///" + app
}
