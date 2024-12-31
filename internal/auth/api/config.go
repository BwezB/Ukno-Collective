package api

import (
	c "github.com/BwezB/Wikno-backend/pkg/configs"
)

type ServerConfig struct {
	Host string `yaml:"host" validate:"required,hostname"`
	Port string `yaml:"port" validate:"required,min=1,max=65535"`
}

func (s *ServerConfig) SetDefaults() {
	s.Host = defaultServerHost
	s.Port = defaultServerPort
}

func (s *ServerConfig) AddFromEnv() {
	c.SetEnvValue(&s.Host, envServerHost)
	c.SetEnvValue(&s.Port, envServerPort)
}

func (s *ServerConfig) AddFromFlags() {
	c.SetFlagValue(&s.Host, flagServerHost)
	c.SetFlagValue(&s.Port, flagServerPort)
}

func (s *ServerConfig) GetAddress() string {
	return s.Host + ":" + s.Port
}

// DEFAULTS
const (
	// defaultServerHost is the default host for the server.
	defaultServerHost = "localhost"
	// defaultServerPort is the default port for the server.
	defaultServerPort = "50051"
)

// ENV
const (
	// envServerHost is the environment variable for the server host.
	envServerHost = "SERVER_HOST"
	// envServerPort is the environment variable for the server port.
	envServerPort = "SERVER_PORT"
)

// FLAGS
var (
	flagServerHost = c.NewFlag("server_host", "", "Server Host")
	flagServerPort = c.NewFlag("server_port", "", "Server Port")
)
