package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var Conf Config

// Config configuration
type Config struct {
	Core       SectionCore `yaml:core`
	ReportSend ReportSend  `yaml:"report_send"`

	ConfigFile  string
	LogDir      string
	LogFile     string
	ServiceName string
	PprofPort   string
	RpcConfDir  string
}

// String ..
func (s *Config) String() string {
	b, err := json.Marshal(*s)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// LoadConfig loads configurations from yaml file
func LoadConfig(path string) *Config {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Load sched configuration file error:%#v", err))
	}

	err = yaml.Unmarshal(text, &Conf)
	if err != nil {
		panic(fmt.Sprintf("Load sched configuration yaml error:%#v", err))
	}
	return &Conf
}

func Usage() {
	usage := `
        -Conf  Config file
        -log   log dir
        -svc   svc name
        -port  listen port
        -rpc   rpc Conf dir
		-version commit info
        `
	fmt.Fprintln(os.Stderr, os.Args[0], usage)
	os.Exit(-1)
}

func ParseCommandParams(commitver, date, goversion string) {
	flag.StringVar(&Conf.ConfigFile, "conf", "", "support Config file.")
	flag.StringVar(&Conf.LogDir, "log", "", "support log dir.")
	flag.StringVar(&Conf.ServiceName, "svc", "", "support svc name.")
	flag.StringVar(&Conf.PprofPort, "port", "", "support service port")
	flag.StringVar(&Conf.RpcConfDir, "rpc", "", "support rpc Conf dir")

	versionFlag := false
	flag.BoolVar(&versionFlag, "version", false, "support rpc Conf dir")
	flag.Parse()

	if versionFlag {
		fmt.Printf("compile info: %s %s %s\n", commitver, date, goversion)
		os.Exit(0)
	}
	if Conf.ConfigFile == "" {
		fmt.Fprintf(os.Stderr, "Configfile is empty, use -Conf option ")
		Usage()
	}

	if Conf.LogDir == "" {
		fmt.Fprintf(os.Stderr, "logdir is empty, use -log option")
		Usage()
	}

	if Conf.ServiceName == "" {
		fmt.Fprintf(os.Stderr, "servicename is empty, use -svc option")
		Usage()
	}

	var err error
	Conf.ConfigFile, err = filepath.Abs(Conf.ConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get abs Config file error: %s\n", err)
		os.Exit(-1)
	}

	Conf.LogDir, err = filepath.Abs(Conf.LogDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get abs log dir error: %s\n", err)
		os.Exit(-1)
	}
	Conf.LogFile = filepath.Join(Conf.LogDir, "app", Conf.ServiceName+".log")
	LoadConfig(Conf.ConfigFile)
}

func IsTestEnv() bool {
	return os.Getenv("STRATEGY_ENV") == "test"
}
