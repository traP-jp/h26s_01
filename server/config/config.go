package config

import (
	"net"
	"strconv"

	"github.com/alecthomas/kong"
	"github.com/go-sql-driver/mysql"
)

type Config struct {
	AppAddr    string `env:"APP_ADDR" default:":8000"`
	DBUser     string `env:"DB_USER" default:"root"`
	DBPassword string `env:"DB_PASSWORD" default:"password"`
	DBHost     string `env:"DB_HOST" default:"localhost"`
	DBPort     int    `env:"DB_PORT" default:"3306"`
	DBName     string `env:"DB_NAME" default:"h26s_01"`
}

func (c *Config) MySQLConfig() *mysql.Config {
	mySQLConfig := mysql.NewConfig()

	mySQLConfig.User = c.DBUser
	mySQLConfig.Passwd = c.DBPassword
	mySQLConfig.Net = "tcp"
	mySQLConfig.Addr = net.JoinHostPort(c.DBHost, strconv.Itoa(c.DBPort))
	mySQLConfig.DBName = c.DBName
	mySQLConfig.Collation = "utf8mb4_general_ci"
	mySQLConfig.ParseTime = true

	return mySQLConfig
}

func Load() *Config {
	var cfg Config

	kong.Parse(&cfg)

	return &cfg
}
