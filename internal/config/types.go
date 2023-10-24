package config

type Config struct {
	Debug bool   `yaml:"debug" json:"debug"`
	App   string `yaml:"app" json:"app"`
	Web   Web    `yaml:"web" json:"web"`
	Mysql Mysql  `yaml:"mysql" json:"mysql"`
	Redis Redis  `yaml:"redis" json:"redis"`
}

type Web struct {
	Host   string `yaml:"host" json:"host"`
	Secret string `yaml:"secret" json:"secret"`
}

type Mysql struct {
	Uri string `yaml:"uri" json:"uri"`
}

type Redis struct {
	Prefix   string `yaml:"prefix" json:"prefix"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
	Hosts    string `yaml:"hosts" json:"Hosts"`
}
