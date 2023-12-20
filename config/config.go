package config

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	log "zerologix/logger"

	"gopkg.in/yaml.v2"
)

type config struct {
	Env        string            `yaml:"env"`
	App        *appConfig        `yaml:"app"`
	Log        *logConfig        `yaml:"log"`
	HttpServer *httpServerConfig `yaml:"httpServer"`
	Mysql      *MysqlConfig      `yaml:"mysql"`
}

type appConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Debug   bool   `yaml:"debug"`
}

type logConfig struct {
	Level int `yaml:"level"`
}

type httpServerConfig struct {
	Addr        string `yaml:"addr"`
	MetricsAddr string `yaml:"metrics_addr"`
}

type MysqlConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Server   string `yaml:"server"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Debug    bool
	// TODO else
}

// All config
var (
	Conf       = &config{}
	configPath string
)

func Init(ctx context.Context) error {
	flag.Parse()
	// check if there is a config.yaml from configMap.
	if _, err := os.Stat("./config/config.yaml"); err == nil {
		configPath = "./config/config.yaml"
	}
	err := initGlobalConf(ctx, configPath)
	if err != nil {
		return err
	}
	log.Log(ctx).Infof("Config init successfully")
	return nil
}

func initGlobalConf(ctx context.Context, configPath string) (err error) {
	if len(configPath) == 0 {
		return errors.New("config file path is empty")
	}

	yamlFile, err := filepath.Abs(configPath)
	if err != nil {
		return
	}
	// Read file as binary
	b, err := os.ReadFile(yamlFile)
	if err != nil {
		return
	}
	log.Log(ctx).Infof("00000 %s", b)
	newConfig := &config{}
	// Yaml
	if err = yaml.Unmarshal(b, newConfig); err != nil {
		return
	}

	// Overwrite by env
	readFromEnv(newConfig)

	Conf = newConfig
	if Conf.Mysql != nil {
		Conf.Mysql.Debug = Conf.App.Debug
	}

	bf, _ := json.MarshalIndent(Conf, "", "	")
	fmt.Printf("Config:\n%s\n", string(bf))
	return nil
}

func readFromEnv(conf *config) {
	// TODO something from secret
}
