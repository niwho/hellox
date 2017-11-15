package config

// SectionCore is sub section of config.
type SectionCore struct {
	BindAddr string `yaml:"http_bind"`
	Mode     string `yaml:"mode"`      // 关闭gin的日志
	LogLevel int    `yaml:"log_level"` // 关闭gin的日志
}

type ReportSend struct {
	Mail     []string `yaml:"mail"`
	DingDing string   `yaml:"dingding"`
	Db       struct {
		Name   string `yaml:"name"`
		Server string `yaml:"server"`
		User   string `yaml:"user"`
		Passwd string `yaml:"passwd"`
	} `yaml:"db"`
}
