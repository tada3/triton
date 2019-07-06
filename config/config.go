package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/koding/multiconfig"
	"github.com/tada3/triton/logging"
)

const (
	defaultConfigFile = "config.yaml"
)

var (
	homeDir        string
	configInEffect *Config
	log            *logging.Entry
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
	hd, err := getHomeDir()
	if err != nil {
		panic(err)
	}

	homeDir = hd

	configLogging()
	log = logging.NewEntry("config")

	err = parseConfig()
	if err != nil {
		panic(err)
	}
}

func GetHomeDir() string {
	return homeDir
}

func GetConfig() *Config {
	return configInEffect
}

func getHomeDir() (string, error) {
	dir := os.Getenv("TRITON_HOME")
	if dir != "" {
		if exists(dir) {
			return dir, nil
		}
		return "", fmt.Errorf("Invalid $TRITON_HOME: %s", dir)
	}
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("ERROR! Failed to get the current working directory.", err)
		wd = "???"
	}
	fmt.Printf("$TRITON_HOME is not defined. Use current dir (%s) as home dir.\n", wd)
	return ".", nil
}

func parseConfig() error {

	configFile := os.Getenv("TRITON_CONFIG_FILE")
	if configFile == "" {
		configFile = filepath.Join(homeDir, defaultConfigFile)
	}

	var mc *multiconfig.DefaultLoader
	if exists(configFile) {
		//fmt.Printf("Parsing %s..\n", configFile)
		log.Info("Parsing %s..", configFile)
		mc = multiconfig.NewWithPath(configFile)

	} else {
		//fmt.Printf("WARN %s is not found.\n", configFile)
		log.Warn("%s is not found.", configFile)
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

func configLogging() {
	logging.SetLevel(logging.DEBUG)

	conf1 := logging.OutputConfig{
		OutputType: logging.STDOUT,
	}
	conf2 := logging.FileOutputConfig{
		OutputConfig: logging.OutputConfig{
			OutputType: logging.FILE,
		},
		Filename: filepath.Join(homeDir, "log", "triton.log"),
	}
	configs := []interface{}{conf1, conf2}

	err := logging.SetOutputByOutputConfig(configs)
	if err != nil {
		fmt.Printf("ERROR! Failed to configure logging: %+v\n", err)
	}
}
