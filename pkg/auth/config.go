package auth

import (
	c "github.com/BwezB/Wikno-backend/pkg/configs"
	"strconv"
)

type AuthConfig struct {
	Host string `yaml:"host" validate:"required,hostname"`
	Port int `yaml:"port" validate:"required,min=1,max=65535"`
}

func (a *AuthConfig) SetDefaults() {
	a.Host = "localhost"
	a.Port = 50051
}

func (a *AuthConfig) AddFromEnv() {
	c.SetEnvValue(&a.Host, "AUTH_HOST")
	c.SetEnvValue(&a.Port, "AUTH_PORT")
}

var (
	flagAuthHost = c.NewFlag("auth_host", "", "Auth Host")
	flagAuthPort = c.NewFlag("auth_port", "", "Auth Port")
)
func (a *AuthConfig) AddFromFlags() {
	c.SetFlagValue(&a.Host, flagAuthHost)
	c.SetFlagValue(&a.Port, flagAuthPort)
}


// HELPER FUNCTIONS

func (a *AuthConfig) GetAddress() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}