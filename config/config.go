package config

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

var (
	Build   string
	Version string
)

func New() *Config {
	return &Config{}
}
