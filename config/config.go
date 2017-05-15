package config

import (
	"gopkg.in/yaml.v2"
	"beebe/log"
	"io/ioutil"
	"os"
)

type Config struct {
	DbConfig struct{
		Dialect		string		`yaml:"db_dialect"`
		Host 		string		`yaml:"db_host"`
		UserName	string		`yaml:"db_user"`
		DbName		string		`yaml:"db_name"`
		Password	string		`yaml:"db_password"`
		ConfigStr 	string		`yaml:"db_config"`
	}							`yaml:"dbconfig"`
	Upload struct{
		Default 	string		`yaml:"default"`
		UserPath	string 		`yaml:"user"`
		ProjectPath	string		`yaml:"project"`
	}							`yaml:"upload"`
}

var config *Config

func init() {
	config = new(Config)
	data, filerr := ioutil.ReadFile("config/config.yaml")
	if filerr != nil {
		log.Logger().Fatal("error: read file")
		os.Exit(-1)
	}
	err := yaml.Unmarshal(data, config)
	if err != nil {
		log.Logger().Fatal("error: yaml file to config")
		os.Exit(-1)
	}
}

func GetConfig() *Config {
	return config
}