package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
	"fmt"
)

type Config struct {
	DbConfig struct{
		Dialect		string		`yaml:"db_dialect"`
		Host 		string		`yaml:"db_host"`
		UserName	string		`yaml:"db_user"`
		DbName		string		`yaml:"db_name"`
		Password	string		`yaml:"db_password"`
		ConfigStr 	string		`yaml:"db_config"`
		MaxOpen		int			`yaml:"db_maxopen"`
		MaxIdle		int			`yaml:"db_maxidle"`
		MaxLifeTime	time.Duration	`yaml:"db_maxlifetime"`
	}							`yaml:"dbconfig"`
	Upload struct{
		Base		string		`yaml:"base"`
		Default 	string		`yaml:"default"`
		UserPath	string 		`yaml:"user"`
		ProjectPath	string		`yaml:"project"`
		TeamPath	string		`yaml:"team"`
	}							`yaml:"upload"`
	Web struct{
		Port 		string		`yaml:"port"`
		IsCors	 	bool		`yaml:"is_cors"`
	}							`yaml:"web"`
	Cors struct{
		AllowOrigin		string		`yaml:"allow_origin"`
		AllowMethods	string		`yaml:"allow_methods"`
		AllowHeaders	string		`yaml:"allow_headers"`
		AllowCred		bool		`yaml:"allow_credentials"`
	}						`yaml:"cors"`
	Log struct{
		IsDebug			bool 		`yaml:"is_debug"`
		LogPath 		string 		`yaml:"log_path"`
		ErrorPath		string		`yaml:"error_path"`
		MacaronPath		string		`yaml:"macaron_path"`
	}						`yaml:"log"`
}

var config *Config

func init() {
	config = new(Config)
	data, filerr := ioutil.ReadFile("config/config.yaml")
	if filerr != nil {
		fmt.Println("log: read file")
		os.Exit(-1)
	}
	err := yaml.Unmarshal(data, config)
	if err != nil {
		fmt.Println("log: yaml file to config")
		os.Exit(-1)
	}
}

func GetConfig() *Config {
	return config
}