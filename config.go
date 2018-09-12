package main

import (
	"encoding/json"
	"io/ioutil"
)

//Config The config for the entire application
type Config struct {
	HTTPConfig        *HTTPConfig        `json:"http"`
	ApplicationConfig *ApplicationConfig `json:"application"`
}

//HTTPConfig The config for the server/request/responses
type HTTPConfig struct {
	Port int `json:"port"`
}

//ApplicationConfig The config for the Actual Application
type ApplicationConfig struct {
	OutputWidth int    `json:"output_width"`
	CharSet     string `json:"char_set"`
}

//InitConfig Initializes all the configuration by reading from config.json
func InitConfig() *Config {
	configString, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	config := &Config{}
	err = json.Unmarshal(configString, config)
	if err != nil {
		panic(err)
	}
	return config
}
