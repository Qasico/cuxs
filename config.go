package cuxs

import (
	"os"
	"strconv"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/qasico/cuxs/log"
)

const (
	VERSION = "0.0.1"
	DEFAULT_RUNMODE = "dev"
)

type (
	ServerConfig struct {
		Graceful      bool
		ServerTimeOut int
		HTTPAddr      string
		EnableHTTPS   bool
		HTTPSCertFile string
		HTTPSKeyFile  string
	}

	DatabaseConfig struct {
		Engine     string
		ServerHost string
		ServerPort int
		DBName     string
		DBUser     string
		DBPassword string
		IdleMax    int
		ConnMax    int
	}

	RedisConfig struct {
		Network  string
		Address  string
		Password string
	}

	AppConfig struct {
		AppPath          string
		WorkPath         string
		Runmode          string
		ServerName       string
		ResponseType     string
		JwtHash          string
		RecoverPanic     bool
		CopyRequestBody  bool
		EnableErrorsShow bool
		EnableGzip       bool
		MaxMemory        int
		DatabaseConfig   DatabaseConfig
		ServerConfig     ServerConfig
		RedisConfig      RedisConfig
	}
)

var Config *AppConfig

func init() {
	Config = &AppConfig{}

	workPath, _ := os.Getwd()
	Config.AppPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	Config.WorkPath, _ = filepath.Abs(workPath)

	Config.LoadConfig()

	Config.Runmode = Config.getString("APP_RUNMODE", DEFAULT_RUNMODE)
	Config.ServerName = Config.getString("APP_NAME", "cuxs " + VERSION)
	Config.ResponseType = Config.getString("APP_RESPONSE_TYPE", "json")
	Config.JwtHash = Config.getString("APP_JWT_SECRET", "123rty890")
	Config.RecoverPanic = Config.getBool("APP_RECOVER", true)
	Config.CopyRequestBody = Config.getBool("APP_CBODY", true)
	Config.EnableErrorsShow = Config.getBool("APP_DEBUG", false)
	Config.EnableGzip = Config.getBool("APP_GZIP", true)
	Config.MaxMemory = Config.getInt("APP_MMEMORY", 1 << 26)

	Config.DatabaseConfig.Engine = Config.getString("DB_ENGINE", "postgres")
	Config.DatabaseConfig.ServerHost = Config.getString("DB_HOST", "127.0.0.1")
	Config.DatabaseConfig.ServerPort = Config.getInt("DB_PORT", 5432)
	Config.DatabaseConfig.DBName = Config.getString("DB_NAME", "foobar")
	Config.DatabaseConfig.DBUser = Config.getString("DB_USER", "root")
	Config.DatabaseConfig.DBPassword = Config.getString("DB_PASS", "")
	Config.DatabaseConfig.IdleMax = Config.getInt("DB_IDLEMAX", 0)
	Config.DatabaseConfig.ConnMax = Config.getInt("DB_CONNMAX", 20)

	Config.ServerConfig.Graceful = Config.getBool("SERVER_GRACEFUL", true)
	Config.ServerConfig.ServerTimeOut = Config.getInt("SERVER_TIMEOUT", 0)
	Config.ServerConfig.HTTPAddr = Config.getString("SERVER_HOST", "0.0.0.0:8088")
	Config.ServerConfig.EnableHTTPS = Config.getBool("SERVER_SSL", false)
	Config.ServerConfig.HTTPSCertFile = Config.getString("SERVER_CERT", "")
	Config.ServerConfig.HTTPSKeyFile = Config.getString("SERVER_KEY", "")

	Config.RedisConfig.Network = Config.getString("REDIS_NETWORK", "")
	Config.RedisConfig.Address = Config.getString("REDIS_ADDRESS", "")
	Config.RedisConfig.Password = Config.getString("REDIS_PASS", "")
}

// Load .env file in app directory to use for config param
// it will stubed into global env variable
func (c *AppConfig) LoadConfig() {
	if Config.WorkPath != Config.AppPath {
		os.Chdir(Config.AppPath)
	}

	configFile := filepath.Join(Config.AppPath, ".env")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Warnf("No Config File Loaded, Using Default Config ...")
	} else {
		godotenv.Load(configFile)
	}
}

// Read env variable with default values as type string
func (c *AppConfig) getString(key string, defaultValue string) (val string) {
	if val = os.Getenv(key); val != "" {
		return val
	} else {
		os.Setenv(key, defaultValue)
	}

	return defaultValue
}

// Read env variable with default value as type int
func (c *AppConfig) getInt(key string, defaultValue int) (val int) {
	p, _ := strconv.ParseInt(os.Getenv(key), 10, 32)
	if val = int(p); val != 0 {
		return val
	} else {
		os.Setenv(key, strconv.Itoa(defaultValue))
	}

	return defaultValue
}

// Read env variable with default value as type bool
func (c *AppConfig) getBool(key string, defaultValue bool) (val bool) {
	if v := os.Getenv(key); v != "" {
		if v == "true" {
			return true
		} else {
			return false
		}
	} else {
		os.Setenv(key, strconv.FormatBool(defaultValue))
	}

	return defaultValue
}