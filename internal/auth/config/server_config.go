package config

import (
	"github.com/BwezB/Wikno-backend/pkg/common/configuration"
)

type Server struct {
	Host string `yaml:"host" validate:"required,hostname"`
	Port string `yaml:"port" validate:"required,min=1,max=65535"`
}

func (s *Server) SetDefaults() {
	s.Host = "localhost"
	s.Port = "50051"
}

func (s *Server) AddFromEnv() {
	configuration.SetEnvValue(&s.Host, "SERVER_HOST")
	configuration.SetEnvValue(&s.Port, "SERVER_PORT")
}

var (
	flagServerHost = configuration.NewFlag("server_host", "", "Server Host")
	flagServerPort = configuration.NewFlag("server_port", "", "Server Port")
)

func (s *Server) AddFromFlags() {
	configuration.SetFlagValue(&s.Host, flagServerHost)
	configuration.SetFlagValue(&s.Port, flagServerPort)
}

func (s *Server) GetAddress() string {
	return s.Host + ":" + s.Port
}