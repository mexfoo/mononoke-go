//nolint:gochecknoglobals // We only have version as global.
package main

import (
	"fmt"
	"mononoke-go/config"
	"mononoke-go/database"
	"mononoke-go/engine"
	"mononoke-go/utils"
	"os"
)

var (
	// Version the version of mononoke-go.
	Version = "unknown"
	// Commit the git commit hash of this version.
	Commit = "unknown"
	// BuildDate the date on which this binary was build.
	BuildDate = "unknown"
)

func main() {
	conf := config.Get()
	logger := config.InitLogger(conf.LoggerLevel, conf.LoggerType)

	passw, err := utils.HashPassword(fmt.Sprintf("%s%s", conf.Database.DefaultSalt, conf.DefaultUser.Password))
	if err != nil {
		logger.Error("Cannot hash password for default user!",
			"function", "main::main",
			"error", err.Error())
		os.Exit(1)
	}

	db, err := database.New(
		conf.Database.Dialect,
		conf.Database.Connection,
		conf.DefaultUser.Name,
		passw,
		true)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	logger.Info(fmt.Sprintf("Starting mononoke-go version %s:%s@%s", Version, Commit, BuildDate))

	if err = engine.Create(db, conf, logger); err != nil {
		logger.Error("Server error!",
			"function", "main::main",
			"error", err.Error())
	}
}
