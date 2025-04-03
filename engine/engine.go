package engine

import (
	"errors"
	"fmt"
	"log/slog"
	"mononoke-go/config"
	"mononoke-go/database"
	"mononoke-go/entities"
	"mononoke-go/net"
	"mononoke-go/utils"
	"os"
	"os/signal"
	"syscall"
)

func Create(db *database.GormDatabase, conf *config.Configuration, log *slog.Logger) error {
	shutdown := make(chan error)
	go doShutdownOnSignal(shutdown)
	playerList := new(entities.PlayerList)
	playerList.Players = make(map[string]*entities.Player)

	gameList := &entities.GameList{
		Games: make(map[uint32]*entities.Game),
	}

	authClient := net.NewTCPServer(
		fmt.Sprintf("%s:%d", conf.Server.AuthClient.ListenIP, conf.Server.AuthClient.ListenPort),
		conf.Server.AuthClient.UseEncryption,
		conf.Server.AuthClient.EncryptionKey,
		log)

	if authClient == nil {
		return errors.New("error starting AuthClient, stopping")
	}
	authHandler := entities.AuthHandler{
		GameSrvs: gameList,
		Players:  playerList,
		DESKey:   utils.InitDESKey(conf.Server.DefaultDESKey),
		DB:       db,
		Config:   conf,
		Log:      log,
	}
	authHandler.InitServer(authClient)

	gameClient := net.NewTCPServer(
		fmt.Sprintf("%s:%d", conf.Server.AuthGame.ListenIP, conf.Server.AuthGame.ListenPort),
		conf.Server.AuthGame.UseEncryption,
		conf.Server.AuthClient.EncryptionKey,
		log)

	if gameClient == nil {
		return errors.New("error starting AuthClient, stopping")
	}
	gameHandler := entities.GameHandler{
		List:       gameList,
		PlayerList: playerList,
		DB:         db,
		Log:        log,
	}
	gameHandler.InitServer(gameClient)

	go func() {
		err := gameClient.Listen()
		doShutdown(shutdown, err)
	}()
	go func() {
		err := authClient.Listen()
		doShutdown(shutdown, err)
	}()

	err := <-shutdown
	log.Error("Shutting down",
		"function", "Engine::Create",
		"error", err.Error())
	return err
}

func doShutdownOnSignal(shutdown chan<- error) {
	onSignal := make(chan os.Signal, 1)
	signal.Notify(onSignal, os.Interrupt, syscall.SIGTERM)
	sig := <-onSignal
	doShutdown(shutdown, fmt.Errorf("received signal %s", sig))
}

func doShutdown(shutdown chan<- error, err error) {
	select {
	case shutdown <- err:
	default:
		// If there is no one listening on the shutdown channel, then the
		// shutdown is already initiated and we can ignore these errors.
	}
}
