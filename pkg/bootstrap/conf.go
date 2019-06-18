package bootstrap

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/labstack/gommon/log"
	"time"
)

var (
	Conf              config // holds the global app config.
	defaultConfigFile = "conf/conf.toml"
)

type config struct {
	ReleaseMode bool   `toml:"release_mode"`
	LogLevel    string `toml:"log_level"`

	SessionStore string `toml:"session_store"`
	CacheStore   string `toml:"cache_store"`

	// 应用配置
	App app

	// 模板
	Tmpl tmpl

	Server server

	// 静态资源
	Static static

	// Opentracing
	Opentracing opentracing

	// Session
	Session session

	// image
	Image image

	System system
}

type app struct {
	Name string `toml:"name"`
	Mode string
}

type server struct {
	Addr string `toml:"addr"`

	ReadTimeout    Duration `toml:"readTimeout"`
	WriteTimeout   Duration `toml:"writeTimeout"`
	MaxHeaderBytes int
	DomainApi      string `toml:"domain_api"`
	DomainWeb      string `toml:"domain_web"`
}

type system struct {
	Pid      int
	Hostname string
}

type image struct {
	Domain string
	Space  string
	Path   string
}

type session struct {
	Addr string
	Pwd  string
}

type static struct {
	Type string `toml:"type"`
}

type tmpl struct {
	Type   string `toml:"type"`   // PONGO2,TEMPLATE(TEMPLATE Default)
	Data   string `toml:"data"`   // BINDATA,FILE(FILE Default)
	Dir    string `toml:"dir"`    // PONGO2(template/pongo2),TEMPLATE(template)
	Suffix string `toml:"suffix"` // .html,.tpl
}

type opentracing struct {
	Disable     bool   `toml:"disable"`
	Type        string `toml:"type"`
	ServiceName string `toml:"service_name"`
	Address     string `toml:"address"`
}

// initConfig initializes the app configuration by first setting defaults,
// then overriding settings from the app config file, then overriding
// It returns an error if any.
func InitConfig(configFile string) (err error) {
	if configFile == "" {
		configFile = defaultConfigFile
	}

	// Set defaults.
	Conf = config{
		ReleaseMode: false,
		LogLevel:    "DEBUG",
	}

	if _, err = os.Stat(configFile); err != nil {
		return errors.New("config file err:" + err.Error())
	} else {
		log.Infof("load config from file:" + configFile)
		configBytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			return errors.New("config load err:" + err.Error())
		}
		_, err = toml.Decode(string(configBytes), &Conf)
		if err != nil {
			return errors.New("config decode err:" + err.Error())
		}
	}

	Conf.System.Pid = os.Getpid()
	Conf.System.Hostname, err = os.Hostname()
	if err != nil {
		log.Errorf("config hostname err: %v", err)
	}

	// @TODO 配置检查
	log.Infof("config data:%v", Conf)

	return nil
}

func GetLogLvl() log.Lvl {
	//DEBUG INFO WARN ERROR OFF
	switch Conf.LogLevel {
	case "DEBUG":
		return log.DEBUG
	case "INFO":
		return log.INFO
	case "WARN":
		return log.WARN
	case "ERROR":
		return log.ERROR
	case "OF":
		return log.OFF
	}

	return log.DEBUG
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

const (
	// Template Type
	PONGO2   = "PONGO2"
	TEMPLATE = "TEMPLATE"

	// Bindata
	BINDATA = "BINDATA"

	// File
	FILE = "FILE"

	// Redis
	REDIS = "REDIS"

	// Memcached
	MEMCACHED = "MEMCACHED"

	// Cookie
	COOKIE = "COOKIE"

	// In Memory
	IN_MEMORY = "IN_MEMARY"

	// 账号规则
	RegexpAccount = `^[a-zA-Z][a-zA-z0-9\.]{2,50}$`
)
