package config

import (
	"time"

	"github.com/spf13/viper"
)

// Provider defines a set of read-only methods for accessing the application
// configuration params as defined in one of the config files.
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
}

var defaultConfig *viper.Viper

func Config() Provider {
	return defaultConfig
}

func LoadConfigProvider(appName string) Provider {
	return readViperConfig(appName)
}

func init() {
	defaultConfig = readViperConfig("IOTHUB")
}

func readViperConfig(appName string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	// global defaults
	v.SetDefault("json_logs", false)
	v.SetDefault("loglevel", "debug")
	v.SetDefault("mode", "debug") // debug, release, test
	v.SetDefault("listen_address", ":5000")
	v.SetDefault("secret", "887yff9898yfhuiew3489fy3hewfuig239f8ghew32yfh")

	// HTTP Server Config
	v.SetDefault("secure", false)
	v.SetDefault("read_timeout", "0m10s")
	v.SetDefault("write_timeout", "0m10s")
	v.SetDefault("max_header_bytes", 1048576)

	// TLS Config
	v.SetDefault("cert_file", "ssl/server.crt")
	v.SetDefault("key_file", "ssl/server.key")

	return v
}
