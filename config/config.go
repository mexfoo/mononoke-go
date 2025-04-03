package config

import (
	"log/slog"
	"os"

	"github.com/jinzhu/configor"
)

type Configuration struct {
	Database struct {
		Dialect     string `default:"sqlite3"`
		Connection  string `default:"data/mononoke-go.db"`
		DefaultSalt string `default:""`
	}
	DefaultUser struct {
		Name     string `default:"test"`
		Password string `default:"test"`
	}
	Server struct {
		AuthClient struct {
			ListenIP      string `default:"127.0.0.1"`
			ListenPort    int32  `default:"4500"`
			UseEncryption bool   `default:"true"`
			EncryptionKey string `default:""`
		}
		AuthGame struct {
			ListenIP      string `default:"127.0.0.1"`
			ListenPort    int32  `default:"4502"`
			UseEncryption bool   `default:"false"`
			EncryptionKey string `default:""`
		}
		DefaultDESKey  string `default:""`
		AgeRestriction uint8  `default:"18"`
	}
	LoggerLevel string `default:"Info"`
	LoggerType  string `default:"Text"`
}

// Get returns the configuration extracted from env variables or config file.
func Get() *Configuration {
	conf := new(Configuration)
	err := configor.New(&configor.Config{ENVPrefix: "MONONOKE", Silent: true}).Load(conf, "config.yml")
	if err != nil {
		panic(err)
	}
	return conf
}

func InitLogger(level, loggerType string) *slog.Logger {
	var programLevel = new(slog.LevelVar)
	switch level {
	case "Debug":
		programLevel.Set(slog.LevelDebug)
	case "Error":
		programLevel.Set(slog.LevelError)
	case "Info":
		programLevel.Set(slog.LevelInfo)
	case "Warn":
		programLevel.Set(slog.LevelWarn)
	default:
		programLevel.Set(slog.LevelInfo)
	}
	var handler slog.Handler
	if loggerType == "JSON" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	}
	return slog.New(handler)
}
