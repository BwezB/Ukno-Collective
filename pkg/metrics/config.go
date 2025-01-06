package metrics

import (
	"strconv"
	c "github.com/BwezB/Wikno-backend/pkg/configs"
)

type MetricsServerConfig struct {
	Host string `yaml:"host" validate:"required"`
	Port int `yaml:"port" validate:"required,min=1,max=65535"`
	Path string `yaml:"path" validate:"required"`
}

func (msc *MetricsServerConfig) SetDefaults() {
	msc.Host = "localhost"
	msc.Port = 2112  
	msc.Path = "/metrics"
}

func (msc *MetricsServerConfig) AddFromEnv() {
	c.SetEnvValue(&msc.Host, "METRICS_HOST")
	c.SetEnvValue(&msc.Port, "METRICS_PORT")
	c.SetEnvValue(&msc.Path, "METRICS_PATH")
}

var (
	flagMetricsHost = c.NewFlag("metrics-host", "", "Metrics server host")
	flagMetricsPort = c.NewFlag("metrics-port", "", "Metrics server port")
	flagMetricsPath = c.NewFlag("metrics-path", "", "Metrics server path")
)
func (msc *MetricsServerConfig) AddFromFlags() {
	c.SetFlagValue(&msc.Host, flagMetricsHost)
	c.SetFlagValue(&msc.Port, flagMetricsPort)
	c.SetFlagValue(&msc.Path, flagMetricsPath)
}

func (msc *MetricsServerConfig) GetAddress() string {
	return msc.Host + ":" + strconv.Itoa(msc.Port)
}
