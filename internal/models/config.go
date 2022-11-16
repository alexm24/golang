package models

type HTTPServerConfig struct {
	Port string `yaml:"port"`
	Path string `yaml:"path"`
}

type CentrifugoConfig struct {
	Url             string `yaml:"url"`
	TokenHmacSecret string `yaml:"token_hmac_secret"`
	APIKey          string `yaml:"api_key"`
}

type RedisConfig struct {
	Url  string `yaml:"url"`
	Port string `yaml:"port"`
}

type Config struct {
	HTTPServerConfig `yaml:"http_server_config"`
	CentrifugoConfig `yaml:"centrifugo_config"`
	RedisConfig      `yaml:"redis_config"`
	DBConfig         string `yaml:"db_config"`
}
