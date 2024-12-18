package config

import (
	"flag"
)

// Defult vales for the database
const (
	defaultServerAddress = ":50051"
)

// Environment variables for the database
const (
	envServerAddress = "SERVER_ADDRESS"
)

// Flags for the database
var (
	flagServerAddress = flag.String("server-address", "", "Server address")
)

type ServerConfig struct {
	Address string `yaml:"address" validate:"required" json:"-"`
}

func NewServerConfig(fileServerConfig *ServerConfig) (*ServerConfig, error) {
	if !flag.Parsed() {
		flag.Parse() // Parse the flags if they have not been parsed
	}

	return &ServerConfig{
		Address: getConfigValue(*flagServerAddress, envServerAddress, fileServerConfig.Address, defaultServerAddress),
	}, nil

}
