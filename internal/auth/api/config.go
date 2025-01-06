package api

import (
	c "github.com/BwezB/Wikno-backend/pkg/configs"
	m "github.com/BwezB/Wikno-backend/pkg/metrics"
	"strconv"
)

type ServerConfig struct {
	Metrics m.MetricsServerConfig `yaml:"metrics"`
	Host string `yaml:"host" validate:"required,hostname"`
	Port int `yaml:"port" validate:"required,min=1,max=65535"`
}


// DEFAULTS

func (s *ServerConfig) SetDefaults() {
	s.Host = "localhost"
	s.Port = 50051
	s.Metrics.SetDefaults()
}


// ENV

func (s *ServerConfig) AddFromEnv() {
	c.SetEnvValue(&s.Host, "SERVER_HOST")
	c.SetEnvValue(&s.Port, "SERVER_PORT")
	s.Metrics.AddFromEnv()
}


// FLAGS

var (
	flagServerHost = c.NewFlag("server_host", "", "Server Host")
	flagServerPort = c.NewFlag("server_port", "", "Server Port")
)
func (s *ServerConfig) AddFromFlags() {
	c.SetFlagValue(&s.Host, flagServerHost)
	c.SetFlagValue(&s.Port, flagServerPort)
	s.Metrics.AddFromFlags()
}


func (s *ServerConfig) GetAddress() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}
