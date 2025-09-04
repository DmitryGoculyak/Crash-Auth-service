package routers

import "time"

type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
}
