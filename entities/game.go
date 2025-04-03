package entities

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log/slog"
	"mononoke-go/database"
	"mononoke-go/net"
	"mononoke-go/net/packets"
	"mononoke-go/net/packets/game"
	"mononoke-go/utils"
	"sync"
)

type Game struct {
	Client              *net.Client
	ServerIdx           uint32
	ServerName          string
	ServerScreenshotURL string
	IsAdultServer       byte
	ServerIP            string
	ServerPort          int32
}

type GameList struct {
	Games map[uint32]*Game
	mutex sync.Mutex
}

func (gl *GameList) AddGame(game *Game) {
	gl.mutex.Lock()
	gl.Games[game.ServerIdx] = game
	gl.mutex.Unlock()
}

func (gl *GameList) RemoveGame(game *Game) {
	gl.mutex.Lock()
	delete(gl.Games, game.ServerIdx)
	gl.mutex.Unlock()
}

func (gl *GameList) GetGame(key uint32) (*Game, bool) {
	game, ok := gl.Games[key]
	return game, ok
}

type GameHandler struct {
	List       *GameList
	PlayerList *PlayerList
	DB         *database.GormDatabase
	Log        *slog.Logger
}

func (a *GameHandler) InitServer(server *net.Server) {
	server.OnNewMessage(func(c *net.Client, header packets.Message, message []byte) {
		go a.HandleMessage(c, header, message)
	})
	server.OnClientConnectionClosed(func(c *net.Client, err error) {
		if game, exists := a.List.GetGame(c.GameIdentifier); exists && game != nil {
			a.List.RemoveGame(game)
		}
		a.Log.Info("Game disconnected",
			"function", "GameHandler::OnClientConnectionClosed",
			"identifier", c.GameIdentifier,
			"error", err.Error())
	})
	server.OnNewClient(func(c *net.Client) {
		a.Log.Debug("New Gameserver connected!",
			"function", "GameHandler::OnNewClient",
			"remoteEndpoint", c.GetEndpoint())
	})
}

func (a *GameHandler) parseMessage(_ *net.Client, msg []byte, packetName string, result interface{}) error {
	a.Log.Debug(fmt.Sprintf("%s Packet received!", packetName),
		"function", "GameHandler::HandleMessage",
		"data", fmt.Sprintf("%v", msg))
	reader := bytes.NewBuffer(msg)
	err := utils.Unmarshal(reader, binary.LittleEndian, result, 0x000000)
	if err != nil {
		a.Log.Error("Error while decoding packet",
			"function", "GameHandler::HandleMessage",
			"packet", "AuthGameGameLogin",
			"error", err.Error())
		return err
	}
	return nil
}

func (a *GameHandler) HandleMessage(c *net.Client, header packets.Message, msg []byte) {
	switch header.HeaderMessageId {
	case game.AuthGameKickClientID:
		clientKickFailedPkt := game.GameAuthClientKickFailed{}
		if err := a.parseMessage(c, msg, "GameAuthClientKickFailed", &clientKickFailedPkt); err == nil {
			a.HandleClientKickFailed(c, clientKickFailedPkt)
		}
	case game.GameAuthClientLoginID:
		clientLoginPkt := game.GameAuthClientLogin{}
		if err := a.parseMessage(c, msg, "GameAuthClientLogin", &clientLoginPkt); err == nil {
			a.HandleClientLogin(c, clientLoginPkt)
		}
	case game.GameAuthClientLogoutID:
		clientLogoutPkt := game.GameAuthClientLogout{}
		if err := a.parseMessage(c, msg, "GameAuthClientLogout", &clientLogoutPkt); err == nil {
			a.HandleClientLogout(c, clientLogoutPkt)
		}
	case game.GameAuthLoginID:
		loginPkt := game.GameAuthLogin{}
		if err := a.parseMessage(c, msg, "GameAuthLogin", &loginPkt); err == nil {
			a.HandleGameServerLogin(c, loginPkt)
		}
	case game.GameAuthSecurityNoCheckID:
		securityNoPkt := game.GameAuthSecurityNoCheck{}
		if err := a.parseMessage(c, msg, "GameAuthSecurityNoCheck", &securityNoPkt); err == nil {
			a.HandleSecurityNoCheck(c, securityNoPkt)
		}
	default:
		a.Log.Warn("Unknown packet",
			"function", "GameHandler::HandleMessage",
			"id", header.HeaderMessageId,
			"data", fmt.Sprintf("%v", msg))
	}
}

func (a *GameHandler) HandleClientKickFailed(_ *net.Client, clientKickFailedPkt game.GameAuthClientKickFailed) {
	playerName := utils.CToGoString(clientKickFailedPkt.Account[:])
	a.removePlayerFromGame(playerName)
}

func (a *GameHandler) HandleClientLogin(c *net.Client, clientLoginPkt game.GameAuthClientLogin) {
	if !a.gameServerAuthenticated(c, "HandleClientLogin") {
		c.Close()
		return
	}

	playerName := utils.CToGoString(clientLoginPkt.Account[:])

	loginResultPkt := game.AuthGameClientLogin{
		Result:  packets.ResultAccessDenied,
		Account: clientLoginPkt.Account,
	}

	currGame, found := a.List.GetGame(c.GameIdentifier)
	if !found {
		a.Log.Error("Gameserver not in list, Clientlogin failed!",
			"function", "GameHandler::HandleClientLogin",
			"GameIdentifier", c.GameIdentifier)
		c.Close()
		return
	}

	a.Log.Debug(fmt.Sprintf("%p", a.PlayerList))
	player := a.PlayerList.GetPlayer(playerName)
	if player == nil {
		a.Log.Error("Login attempt with non-existant player!",
			"function", "GameHandler::HandleClientLogin",
			"accountName", playerName)
		c.Send(loginResultPkt, game.AuthGameClientLoginID)
		return
	}

	if player.OneTimeKey != clientLoginPkt.OneTimeKey {
		a.Log.Error("Client tried to login with wrong key",
			"function", "GameHandler::HandleClientLogin",
			"accountID", player.AccountID,
			"accountName", player.AccountName,
			"expectedKey", player.OneTimeKey,
			"receivedKey", clientLoginPkt.OneTimeKey)
		c.Send(loginResultPkt, game.AuthGameClientLoginID)
		return
	}

	player.GameIndex = currGame.ServerIdx
	player.IsInGame = true
	loginResultPkt.Result = packets.ResultSuccess
	loginResultPkt.AccountID = player.AccountID
	loginResultPkt.Permission = player.Permission
	loginResultPkt.PCBangUser = 0
	loginResultPkt.Age = uint32(player.Age)
	loginResultPkt.EventCode = 0
	loginResultPkt.ContinuousPlayTime = 0
	loginResultPkt.ContinuousLogoutTime = 0
	if succ, err := a.DB.UpdateLastLoginServerIdx(player.AccountID, player.GameIndex); !succ {
		a.Log.Error("Cannot update Last Login ServerIdx",
			"function", "GameHandler::HandleClientLogin",
			"accountID", player.AccountID,
			"accountName", player.AccountName,
			"error", err.Error())
	}

	c.Send(loginResultPkt, game.AuthGameClientLoginID)
}

func (a *GameHandler) HandleClientLogout(c *net.Client, clientLogoutPkt game.GameAuthClientLogout) {
	if !a.gameServerAuthenticated(c, "HandleClientLogout") {
		c.Close()
		return
	}

	playerName := utils.CToGoString(clientLogoutPkt.Account[:])
	a.removePlayerFromGame(playerName)
}

func (a *GameHandler) HandleSecurityNoCheck(c *net.Client, _ game.GameAuthSecurityNoCheck) {
	if !a.gameServerAuthenticated(c, "HandleSecurityNoCheck") {
		c.Close()
		return
	}
}

func (a *GameHandler) HandleGameServerLogin(c *net.Client, loginPkt game.GameAuthLogin) {
	srv := Game{
		Client:              c,
		ServerIdx:           uint32(loginPkt.ServerIdx),
		ServerName:          utils.CToGoString(loginPkt.ServerName[:]),
		ServerScreenshotURL: utils.CToGoString(loginPkt.ServerScreenshotURL[:]),
		IsAdultServer:       loginPkt.IsAdultServer,
		ServerIP:            utils.CToGoString(loginPkt.ServerIP[:]),
		ServerPort:          loginPkt.ServerPort,
	}
	if _, exists := a.List.GetGame(srv.ServerIdx); exists {
		a.Log.Error("Gameserver already registered",
			"function", "GameHandler::HandleGameServerLogin",
			"serverIdx", srv.ServerIdx,
			"serverName", srv.ServerName,
			"serverIP", srv.ServerIP,
			"serverPort", srv.ServerPort,
			"serverScreenshotURL", srv.ServerScreenshotURL,
			"isAdultServer", srv.IsAdultServer)
		resultPkt := game.AuthGameLoginResult{
			Result: packets.ResultAccessDenied,
		}
		c.Send(resultPkt, game.AuthGameLoginResultID)
		c.Close()
		return
	}

	a.List.AddGame(&srv)
	a.Log.Info("Gameserver registered",
		"function", "GameHandler::HandleGameServerLogin",
		"serverIdx", srv.ServerIdx,
		"serverName", srv.ServerName,
		"serverIP", srv.ServerIP,
		"serverPort", srv.ServerPort,
		"serverScreenshotURL", srv.ServerScreenshotURL,
		"isAdultServer", srv.IsAdultServer)

	c.GameIdentifier = srv.ServerIdx
	c.IsAuthenticated = true

	resultPkt := game.AuthGameLoginResult{
		Result: packets.ResultSuccess,
	}

	c.Send(resultPkt, game.AuthGameLoginResultID)
}

func (a *GameHandler) KickPlayer(playerName string, srv *Game) {
	kickPkt := game.AuthGameKickClient{}
	copy(kickPkt.Account[:], []byte(playerName))
	srv.Client.Send(kickPkt, game.AuthGameKickClientID)
}

func (a *GameHandler) removePlayerFromGame(playerName string) {
	player := a.PlayerList.GetPlayer(playerName)
	if player == nil {
		a.Log.Error("Cannot remove player from game, player does not exist",
			"function", "GameHandler::removePlayerFromGame",
			"accountName", playerName)
		return
	}
	a.PlayerList.RemovePlayer(player)
}

func (a *GameHandler) gameServerAuthenticated(c *net.Client, funcName string) bool {
	if !c.IsAuthenticated {
		a.Log.Error("Gameserver not authenticated",
			"GameHandler", funcName,
			"remoteAddr", c.GetEndpoint())
		return false
	}

	if _, exist := a.List.GetGame(c.GameIdentifier); !exist {
		a.Log.Error("Gameserver not in game list",
			"GameHandler", funcName,
			"remoteAddr", c.GetEndpoint(),
			"serverIdx", c.GameIdentifier)
		return false
	}

	return true
}
