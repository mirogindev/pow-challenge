package config

import (
	"github.com/mirogindev/pow-challenge/internal/tools"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	Difficulty    int    `yaml:"difficulty"`
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	MaxIterations int64  `yaml:"maxIterations"`
}

var cfg *Config
var configPath string

func init() {
	var err error
	configPath = path.Join(path.Join(tools.GetBasePath(), "../../", "config.yaml"))
	cfg, err = parseConfig(configPath)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Debug("cannot parse config, will be used default values")
		cfg = &Config{Difficulty: 2, Host: "localhost", Port: 8085, MaxIterations: 1000000}
	}
}

func GetConfig() *Config {
	return cfg
}

func parseConfig(p string) (*Config, error) {
	file, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	confContent := []byte(os.ExpandEnv(string(file)))
	config := Config{}
	err = yaml.Unmarshal(confContent, &config)
	if err != nil {
		return nil, err
	}
	return &config, err
}
