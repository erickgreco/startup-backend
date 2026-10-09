package config

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string
	DBURL    string
	GinMode  string
	LogLevel string

	DBMaxConns              int32
	DBMinConns              int32
	DBMaxConnLifeTime       time.Duration
	DBMaxConnIdleTime       time.Duration
	DBHealthCheckPeriod     time.Duration
	DBMaxConnLifeTimeJitter time.Duration
}

func Load() *Config {
	viper.SetDefault("PORT", "8081")
	viper.SetDefault("DB_URL", "")
	viper.SetDefault("GIN_MODE", "debug")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("DB_MAX_CONNS", 10)
	viper.SetDefault("DB_MIN_CONNS", 5)
	viper.SetDefault("MAX_CONN_LIFETIME", time.Hour)
	viper.SetDefault("MAX_CONN_IDLE_TIME", 10*time.Minute)
	viper.SetDefault("HEALTH_CHECK_PERIOD", time.Minute)
	viper.SetDefault("MAX_CONN_LIFETIME_JITTER", 30*time.Second)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// TODO: Cuando viper publique v2 con FileLookupError migrar a
	// errors.As(err, new(viper.FileLookupError))
	if err := viper.ReadInConfig(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		panic(fmt.Errorf("error reading .env: %w", err))
	}

	return &Config{
		Port:                    ":" + viper.GetString("PORT"),
		DBURL:                   viper.GetString("DB_URL"),
		GinMode:                 viper.GetString("GIN_MODE"),
		LogLevel:                viper.GetString("LOG_LEVEL"),
		DBMaxConns:              viper.GetInt32("DB_MAX_CONNS"),
		DBMinConns:              viper.GetInt32("DB_MIN_CONNS"),
		DBMaxConnLifeTime:       viper.GetDuration("MAX_CONN_LIFETIME"),
		DBMaxConnIdleTime:       viper.GetDuration("MAX_CONN_IDLE_TIME"),
		DBHealthCheckPeriod:     viper.GetDuration("HEALTH_CHECK_PERIOD"),
		DBMaxConnLifeTimeJitter: viper.GetDuration("MAX_CONN_LIFETIME_JITTER"),
	}
}
