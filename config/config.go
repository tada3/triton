package config

import (
	"fmt"
	"os"

	"github.com/koding/multiconfig"
)

const (
	defaultConfigFile = "config.yaml"
)

var (
	configInEffect *Config
)

type Config struct {
	MySQLHost     string `default:"localhsot"`
	MySQLPort     int    `default:"3306"`
	MySQLDatabase string `default:"triton"`
	MySQLUser     string `default:"triton"`
	MySQLPasswd   string `default:triton`

	RedisHost   string `default:"localhost"`
	RedisPort   int    `default:6379`
	RedisPasswd string

	SoundFileBaseUrl string
}

func init() {
	err := parseConfig()
	if err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	return configInEffect
}

func parseConfig() error {

	configFile := os.Getenv("TRITON_CONFIG_FILE")
	if configFile == "" {
		configFile = defaultConfigFile
	}

	var mc *multiconfig.DefaultLoader
	if exists(configFile) {
		fmt.Printf("Parsing %s..\n", configFile)
		mc = multiconfig.NewWithPath(configFile)

	} else {
		fmt.Printf("WARN %s is not found.\n", configFile)
		mc = multiconfig.New()
	}

	cfg := &Config{}
	err := mc.Load(cfg)
	if err != nil {
		return err
	}
	configInEffect = cfg
	return nil
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
