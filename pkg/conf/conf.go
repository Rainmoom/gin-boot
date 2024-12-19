package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	Cfg = &ConfigYaml{
		Mode: "debug",
		Http: &HttpConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Log: &LogConfig{
			Level:         "info",
			MaxSize:       100, // megabytes
			MaxBackups:    5,
			MaxAge:        15, // 15 days
			Compress:      true,
			Path:          "./log/app.log",
			ConsoleEnable: true,
		},
		Postgres: &PostgresConfig{
			Ip:   "127.0.0.1",
			Port: 5432,
			Ssl:  "disable",
		},
	}
)

type HttpConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Db       int    `yaml:"db"`
	Password string `yaml:"password"`
	MaxIdle  int    `yaml:"maxIdle"`
	PoolSize int    `yaml:"poolSize"`
}

type MysqlConfig struct {
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type PostgresConfig struct {
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Ssl      string `yaml:"ssl"`
}

type KV struct {
	Key   string
	Value string
}

type LogConfig struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
	// MaxSize max size of single file, unit is MB
	MaxSize int `yaml:"maxSize"`
	// MaxBackups max number of backup files
	MaxBackups int `yaml:"maxBackups"`
	// MaxAge max days of backup files, unit is day
	MaxAge int `yaml:"maxAge"`
	// Compress whether compress backup file
	Compress bool `yaml:"compress"`
	// Format
	Format string `yaml:"format"`
	// Console output
	ConsoleEnable bool `yaml:"consoleEnable"`
}

type Config interface {
}

type ConfigYaml struct {
	Config
	Mode      string          `yaml:"mode"`
	Http      *HttpConfig     `yaml:"http"`
	Log       *LogConfig      `yaml:"log"`
	Redis     *RedisConfig    `yaml:"redis"`
	Mysql     *MysqlConfig    `yaml:"mysql"`
	Postgres  *PostgresConfig `yaml:"postgres"`
	SecretKey string          `json:"secretKey"`
}

func ParseConfig(conf any, filePath string) {
	viper.SetConfigFile(filePath)

	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("config file: %s", err))
	}

	if err = viper.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("parse config from config.yaml failed:%s", err))
	}
}
