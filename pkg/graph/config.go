package graph

import (
	c "github.com/BwezB/Wikno-backend/pkg/configs"
)

type GraphConfig struct {
	Host string
	Port string
}

func (gc *GraphConfig) SetDefaults() {
	gc.Host = "localhost"
	gc.Port = "50052"
}

func (gc *GraphConfig) AddFromEnv() {
	c.SetEnvValue(&gc.Host, "GRAPH_HOST")
	c.SetEnvValue(&gc.Port, "GRAPH_PORT")
}

var (
	flagGraphHost = c.NewFlag("graph_host", "", "Graph Host")
	flagGraphPort = c.NewFlag("graph_port", "", "Graph Port")
)
func (gc *GraphConfig) AddFromFlags() {
	c.SetFlagValue(&gc.Host, flagGraphHost)
	c.SetFlagValue(&gc.Port, flagGraphPort)
}


// HELPER FUNCTIONS

func (gc *GraphConfig) GetAddress() string {
	return gc.Host + ":" + gc.Port
}
