package config

import (
	"github.com/BwezB/Wikno-backend/pkg/common/configuration"
)

type Server struct {
	Host string `yaml:"host" validate:"required,hostname"`
	Port string `yaml:"port" validate:"required,min=1,max=65535"`
}

func (s *Server) SetDefaults() {
	s.Host = defaultServerHost
	s.Port = defaultServerPort
}

func (s *Server) AddFromEnv() {
	configuration.SetEnvValue(&s.Host, envServerHost)
	configuration.SetEnvValue(&s.Port, envServerPort)
}

func (s *Server) AddFromFlags() {
	configuration.SetFlagValue(&s.Host, flagServerHost)
	configuration.SetFlagValue(&s.Port, flagServerPort)
}

func (s *Server) GetAddress() string {
	return s.Host + ":" + s.Port
}
