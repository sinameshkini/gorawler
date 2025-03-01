package repository

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func NewDBWithConf(conf *Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(conf.dsn()), &gorm.Config{Logger: getLogger(conf.Debug)})
}

func NewDBWithDsn(dsn string, debug bool) (*gorm.DB, error) {
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{Logger: getLogger(debug)})

}

type Config struct {
	Name   string `json:"name"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	DBName string `json:"db_name" mapstructure:"db_name"`
	Debug  bool   `json:"debug"`
}

func (c *Config) dsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
		c.Host,
		c.User,
		c.Pass,
		c.DBName,
		c.Port,
	)
}

func getLogger(debug bool) logger.Interface {
	logLevel := logger.Error
	if debug {
		logLevel = logger.Info
	}

	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logLevel,    // Log level
			Colorful:      true,        // Disable color
		},
	)
}
