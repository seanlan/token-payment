package config

type Config struct {
	Debug  bool           `yaml:"debug" json:"debug"`
	App    string         `yaml:"app" json:"app"`
	Secret string         `yaml:"secret" json:"secret"`
	Web    WebConfig      `yaml:"web" json:"web"`
	DB     DatabaseConfig `yaml:"db" json:"db"`
	Redis  RedisConfig    `yaml:"redis" json:"redis"`
}

type WebConfig struct {
	Host   string `yaml:"host" json:"host"`
	Secret string `yaml:"secret" json:"secret"`
}

type DatabaseConfig struct {
	Uri string `yaml:"uri" json:"uri"`
}

type RedisConfig struct {
	Prefix   string `yaml:"prefix" json:"prefix"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
	Hosts    string `yaml:"hosts" json:"Hosts"`
}
