package config

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gopkg.in/yaml.v3"
	"math/rand"
	"time"
)

var C *Config

func setup(c *Config) {
	initLogging(c.Debug, c.App)
	initDB(c.Mysql, c.Debug)
	initRedis(c.Redis)
}

func init() {
	var err error
	// 初始化随机种子
	rand.Seed(time.Now().Unix())
	// 初始化时区 设置时区为中国时区
	local, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = local
	// 初始化配置
	path := "conf.yaml"
	c := config.New(
		config.WithSource(
			file.NewSource(path),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}))
	err = c.Load()
	if err != nil {
		panic(err)
	}
	var s Config
	if err = c.Scan(&s); err != nil {
		panic(err)
	}
	C = &s
	setup(C)
}
