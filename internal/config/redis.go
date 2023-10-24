package config

import (
	"token-payment/internal/dao"
	"token-payment/pkg/xlredis"
)

func initRedis(c Redis) {
	r, err := xlredis.NewClient(c.Hosts, c.Username, c.Password, c.Prefix, c.DB)
	if err != nil {
		panic(err)
	}
	dao.Redis = r
}
