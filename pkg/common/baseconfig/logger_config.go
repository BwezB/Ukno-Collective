package baseconfig

type Logger struct {
	Level string `yaml:"level" validate:"required,oneof=debug info warning error"`
}

func (l *Logger) SetDefaults() {
	l.Level = defaultLogLevel
}

func (l *Logger) AddFromEnv() {
	SetEnvValue(&l.Level, envLogLevel)
}

func (l *Logger) AddFromFlags() {
	SetFlagValue(&l.Level, flagLogLevel)
}
