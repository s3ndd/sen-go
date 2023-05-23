package log

// Config is the configuration for the logger.
type Config struct {
	LogLevel  string
	LogFormat string
	Program   string
	Env       string
}

// Level returns the log level.
func (c *Config) Level() string {
	return c.LogLevel
}

// Format returns the log format.
func (c *Config) Format() string {
	return c.LogFormat
}

// SourceProgram returns the source program.
func (c *Config) SourceProgram() string {
	return c.Program
}

// IsDev returns true if the environment is dev.
func (c *Config) IsDev() bool {
	return c.Env == "dev"
}
