package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug   bool
	Discord struct {
		Token string `yaml:"token"`
	}
	Db struct {
		Kind string `yaml:"kind"`
		Path string `yaml:"path"`
	}
	Voices struct {
		Watson struct {
			Token string `yaml:"token"`
			Api   string `yaml:"api"`
		}
	}
	Guild Guild
}

type Guild struct {
	Prefix  string
	Lang    string
	MaxChar int
	Voice   struct {
		Source string
		Type   string
	}
	ReadBots   bool
	Policy     string
	PolicyList map[string]string `yaml:",flow"`
	Replace    map[string]string `yaml:",flow"`
}

const configFile = "./config.yaml"

var CurrentConfig Config

func init() {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Config load failed:", err)
	}
	err = yaml.Unmarshal(file, &CurrentConfig)
	if err != nil {
		log.Fatal("Config parse failed:", err)
	}
}
