package config

import (
	"github.com/BwezB/Wikno-backend/pkg/common/configs"
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
	configs.SetEnvValue(&s.Host, envServerHost)
	configs.SetEnvValue(&s.Port, envServerPort)
}

func (s *Server) AddFromFlags() {
	configs.SetFlagValue(&s.Host, flagServerHost)
	configs.SetFlagValue(&s.Port, flagServerPort)
}

func (s *Server) GetAddress() string {
	return s.Host + ":" + s.Port
}
