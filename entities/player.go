package entities

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log/slog"
	"math/big"
	"mononoke-go/config"
	"mononoke-go/database"
	"mononoke-go/net"
	"mononoke-go/net/packets"
	"mononoke-go/net/packets/client"
	"mononoke-go/utils"
	"strconv"
	"sync"
)

type Player struct {
	Client          *net.Client
	AccountID       uint32
	AccountName     string
	Age             uint8
	IsBlocked       bool
	LastServerIndex uint32
	IsInGame        bool
	KickNextLogin   bool
	OneTimeKey      uint64
	GameIndex       uint32
	Permission      uint32
}

type PlayerList struct {
	Players map[string]*Player
	mutex   sync.Mutex
}

func (pl *PlayerList) AddPlayer(player *Player) {
	pl.mutex.Lock()
	pl.Players[player.AccountName] = player
	pl.mutex.Unlock()
}

func (pl *PlayerList) RemovePlayer(player *Player) {
	pl.mutex.Lock()
	delete(pl.Players, player.AccountName)
	pl.mutex.Unlock()
}

func (pl *PlayerList) GetPlayer(key string) *Player {
	pl.mutex.Lock()
	player := pl.Players[key]
	pl.mutex.Unlock()
	return player
}

type AuthHandler struct {
	GameSrvs *GameList
	Players  *PlayerList
	DESKey   [8]byte
	DB       *database.GormDatabase
	Config   *config.Configuration
	Log      *slog.Logger
}

func (a *AuthHandler) InitServer(server *net.Server) {
	server.OnNewMessage(func(c *net.Client, header packets.Message, message []byte) {
		go a.HandleMessage(c, header, message)
	})
	server.OnClientConnectionClosed(func(c *net.Client, err error) {
		if player := a.Players.GetPlayer(c.PlayerIdentifier); player != nil {
			if !player.IsInGame {
				a.Players.RemovePlayer(player)
			}
		}
		if err != nil {
			a.Log.Debug("Player disconnected",
				"AutHandler", "OnClientConnectionClosed",
				"identifier", c.PlayerIdentifier,
				"error", err.Error())
		}
	})
	server.OnNewClient(func(c *net.Client) {
		a.Log.Debug("New Player connected!",
			"AutHandler", "OnNewClient",
			"remoteEndpoint", c.GetEndpoint())
	})
}

func (a *AuthHandler) parseMessage(c *net.Client, msg []byte, packetName string, result interface{}) error {
	a.Log.Debug(fmt.Sprintf("%s Packet received!", packetName), "data", fmt.Sprintf("%v", msg))
	reader := bytes.NewBuffer(msg)
	err := utils.Unmarshal(reader, binary.LittleEndian, result, int(c.SupportedVersion))
	if err != nil {
		a.Log.Error("Error while decoding packet",
			"function", "AuthHandler::HandleMessage",
			"packet", packetName,
			"error", err.Error())
		return err
	}
	return nil
}

func (a *AuthHandler) HandleMessage(c *net.Client, header packets.Message, msg []byte) {
	a.setSupportedVersionByPacketID(c, header.HeaderMessageId, header.HeaderMessageSize)
	switch header.HeaderMessageId {
	case client.ClientAuthVersionID:
		versionPkt := client.ClientAuthVersion{}
		if err := a.parseMessage(c, msg, "ClientAuthVersion", &versionPkt); err == nil {
			a.HandleVersion(c, versionPkt)
		}
	case client.ClientAuthAccountID:
		accountPkt := client.ClientAuthAccount{}
		if err := a.parseMessage(c, msg, "ClientAuthAccount", &accountPkt); err == nil {
			a.HandleAccountLogin(c, accountPkt)
		}
	case client.ClientAuthServerListID:
		a.Log.Debug("ClientAuthServerList Packet received!", "data", fmt.Sprintf("%v", msg))
		// doesn't need to be decoded
		a.HandleServerList(c)
	case client.ClientAuthSelectServerID:
		serverSelectPkt := client.ClientAuthSelectServer{}
		if err := a.parseMessage(c, msg, "ClientAuthSelectServer", &serverSelectPkt); err == nil {
			a.HandleServerSelection(c, serverSelectPkt)
		}
	case client.ClientAuthPublicKeyID1, client.ClientAuthPublicKeyID2:
		pubKeyPkt := client.ClientAuthPublicKey{}
		if err := a.parseMessage(c, msg, "ClientAuthPublicKey", &pubKeyPkt); err == nil {
			a.HandlePublicKey(c, pubKeyPkt)
		}
	case 9999:
	default:
		a.Log.Warn("Unknown packet",
			"function", "AuthHandler::HandleMessage",
			"id", header.HeaderMessageId,
			"data", fmt.Sprintf("%v", msg))
	}
}

func (a *AuthHandler) setSupportedVersionByPacketID(c *net.Client, packetID uint16, packetSize uint32) {
	switch packetID {
	case client.ClientAuthAccountID:
		if packetSize > 58 && c.SupportedVersion < packets.Version520 {
			c.SupportedVersion = packets.Version520
		}
	case client.ClientAuthPublicKeyID1: // < 0x090603
		c.SupportedVersion = packets.Version963
	case client.ClientAuthPublicKeyID2: // >= 0x090603
		c.SupportedVersion = packets.Version811
	default:
		return
	}
}

func (a *AuthHandler) HandleVersion(c *net.Client, versionPkt client.ClientAuthVersion) {
	version, err := strconv.Atoi(utils.CToGoString(versionPkt.Version[:]))
	if err != nil {
		a.Log.Error("Cannot parse version string as integer",
			"function", "AuthHandler::HandleVersion",
			"accountName", utils.CToGoString(versionPkt.Version[:]))
		c.Close()
		return
	}
	switch version {
	case packets.DateVersion200:
		c.SupportedVersion = packets.Version200
	case packets.DateVersion410:
		c.SupportedVersion = packets.Version410
	case packets.DateVersion920:
		c.SupportedVersion = packets.Version920
	case packets.DateVersion967:
		c.SupportedVersion = packets.Version967
	default:
		c.SupportedVersion = packets.Version200
	}
}

func (a *AuthHandler) HandleAccountLogin(c *net.Client, accountPkt client.ClientAuthAccount) {
	if c.IsAuthenticated || c.PlayerIdentifier != "" {
		a.Log.Error("User already authenticated",
			"function", "AuthHandler::HandleAccountLogin",
			"accountName", utils.CToGoString(accountPkt.Account))
		c.Close()
		return
	}
	player := new(Player)
	player.AccountName = utils.CToGoString(accountPkt.Account)
	// DES attempt
	var password string
	if len(c.AESKey) == 0 {
		block, _ := des.NewCipher(a.DESKey[:])
		var decryptedPassword []byte
		encryptedBlock := make([]byte, 8)
		for i := range 4 {
			block.Decrypt(encryptedBlock, accountPkt.Password[i*8:(i+1)*8])
			decryptedPassword = append(decryptedPassword, encryptedBlock...)
		}
		password = utils.CToGoString(decryptedPassword)
	} else {
		var decryptedPassword []byte
		decryptedBlock := make([]byte, 16)
		block, err := aes.NewCipher(c.AESKey[:16])
		if err != nil {
			a.Log.Error("Cannot decrypt AES password",
				"function", "AuthHandler::HandleAccountLogin",
				"accountName", player.AccountName,
				"error", err.Error())
		}
		mode := cipher.NewCBCDecrypter(block, c.AESKey[16:])

		for bytesRead := 0; bytesRead+15 < int(accountPkt.PasswordSize); bytesRead += 16 {
			mode.CryptBlocks(decryptedBlock, accountPkt.Password[bytesRead:bytesRead+16])
			decryptedPassword = append(decryptedPassword, decryptedBlock...)
		}
		decryptedPassword = utils.PKCS5Trimming(decryptedPassword)
		password = utils.CToGoString(decryptedPassword)
	}

	account, found := a.DB.GetUserByNameAndPW(player.AccountName,
		fmt.Sprintf("%s%s", a.Config.Database.DefaultSalt, password), a.Config)

	if !found {
		resultPkt := client.AuthClientResult{
			RequestMessageID: 10010,
			Result:           packets.ResultNotExist,
			LoginFlag:        client.LoginFlagEulaAccepted,
		}
		a.Log.Debug("Failed login attempt",
			"function", "AuthHandler::HandleAccountLogin",
			"accountName", player.AccountName)
		c.Send(resultPkt, client.AuthClientResultID)
		return
	}

	player.Age = account.Age
	player.AccountID = account.AccountID
	player.IsBlocked = account.Blocked
	player.LastServerIndex = account.LastLoginServerIdx
	player.Permission = account.Permission

	var resultPkt client.AuthClientResult
	if player.IsBlocked {
		resultPkt = client.AuthClientResult{
			RequestMessageID: 10010,
			Result:           packets.ResultAccessDenied,
			LoginFlag:        client.LoginFlagAccountBlockWarning,
		}
		c.Send(resultPkt, client.AuthClientResultID)
		return
	}

	c.IsAuthenticated = true
	c.PlayerIdentifier = player.AccountName
	a.Players.AddPlayer(player)
	resultPkt = client.AuthClientResult{
		RequestMessageID: 10010,
		Result:           packets.ResultSuccess,
		LoginFlag:        client.LoginFlagEulaAccepted,
	}
	c.Send(resultPkt, client.AuthClientResultID)
}

func (a *AuthHandler) HandleServerList(c *net.Client) {
	if !a.IsLoggedIn(c, "HandleServerList") {
		return
	}

	player := a.Players.GetPlayer(c.PlayerIdentifier)

	a.GameSrvs.mutex.Lock()
	serverPkt := client.AuthClientServerList{
		LastLoginServerIdx: player.LastServerIndex,
	}
	if len(a.GameSrvs.Games) < 0xFFFF {
		serverPkt.Servers = uint32(len(a.GameSrvs.Games)) //nolint:gosec // We have a check for overflow
	}
	for _, value := range a.GameSrvs.Games {
		val := client.ServerInfo{
			ServerIdx:     value.ServerIdx,
			IsAdultServer: value.IsAdultServer,
			ServerPort:    value.ServerPort,
			UserRatio:     0,
		}
		copy(val.ServerIP[:], []byte(value.ServerIP))
		copy(val.ServerScreenshotURL[:], []byte(value.ServerScreenshotURL))
		copy(val.ServerName[:], []byte(value.ServerName))
		serverPkt.ServerInfo = append(serverPkt.ServerInfo, val)
	}
	a.GameSrvs.mutex.Unlock()

	c.Send(serverPkt, client.AuthClientServerListID)
}

func (a *AuthHandler) HandlePublicKey(c *net.Client, pubKeyPkt client.ClientAuthPublicKey) {
	key, err := utils.BytesToPublicKey(pubKeyPkt.Key)
	if err != nil {
		a.Log.Error("Cannot parse public key",
			"function", "AuthHandler::HandlePublicKey",
			"error", err.Error())
		c.Close()
		return
	}

	aesKey := make([]byte, 32)
	for i := range aesKey {
		otk, randErr := rand.Int(rand.Reader, big.NewInt(0xFFFFFFFF))
		if randErr != nil {
			a.Log.Error("Error generating randon key.",
				"function", "AuthHandler::HandlePublicKey",
				"error", randErr.Error())
		}
		aesKey[i] = byte(otk.Int64() & 0xFF)
	}

	encryptedAES, err := utils.EncryptWithPublicKey(aesKey, key)
	if err != nil {
		a.Log.Error("Cannot encrypt AES key with public key",
			"function", "AuthHandler::HandlePublicKey",
			"error", err.Error())
		c.Close()
		return
	}

	resultPkt := client.ClientAuthPublicKey{
		Size: uint32(len(encryptedAES)), //nolint:gosec // this is fine
		Key:  encryptedAES,
	}
	c.AESKey = aesKey
	if pubKeyPkt.Header.HeaderMessageId == client.ClientAuthPublicKeyID1 {
		c.Send(resultPkt, client.AuthClientAESKeyID1)
	} else {
		c.Send(resultPkt, client.AuthClientAESKeyID2)
	}
}

func (a *AuthHandler) HandleServerSelection(c *net.Client, serverSelectPkt client.ClientAuthSelectServer) {
	if !a.IsLoggedIn(c, "HandleServerSelection") {
		return
	}

	player := a.Players.GetPlayer(c.PlayerIdentifier)
	if player.IsInGame {
		a.Log.Error("Player already in game!",
			"function", "AuthHandler::HandleServerSelection",
			"requestedServerIDX", serverSelectPkt.ServerIdx,
			"existingServerIdx", player.GameIndex)
		resultPkt := client.AuthClientSelectServer{Result: packets.ResultAccessDenied}
		c.Send(resultPkt, client.AuthClientSelectServerID)
		return
	}

	resultPkt := client.AuthClientSelectServer{Result: packets.ResultAccessDenied}

	srv, ok := a.GameSrvs.GetGame(serverSelectPkt.ServerIdx)
	if !ok {
		a.Log.Error("Invalid server selection!",
			"function", "AuthHandler::HandleServerSelection",
			"serverIdx", serverSelectPkt.ServerIdx)
		c.Send(resultPkt, client.AuthClientSelectServerID)
		return
	}

	if (srv.IsAdultServer == 1) && (player.Age < a.Config.Server.AgeRestriction) {
		resultPkt.Result = packets.ResultTooYoung
		a.Log.Debug("Player too young to join adult server!",
			"function", "AuthHandler::HandleServerSelection",
			"serverIdx", serverSelectPkt.ServerIdx,
			"playerAge", player.Age)
		c.Send(resultPkt, client.AuthClientSelectServerID)
		return
	}

	otk, err := rand.Int(rand.Reader, big.NewInt(0x7FFFFFFFFFFFFFFF))
	if err != nil {
		a.Log.Error("Error generating randon one-time-key.",
			"function", "AuthHandler::HandleServerSelection",
			"error", err.Error())
	}

	player.IsInGame = true
	player.GameIndex = serverSelectPkt.ServerIdx
	player.OneTimeKey = otk.Uint64()
	resultPkt.Result = packets.ResultSuccess
	resultPkt.OneTimeKey = player.OneTimeKey
	resultPkt.PendingTime = 0
	c.Send(resultPkt, client.AuthClientSelectServerID)
}

func (a *AuthHandler) IsLoggedIn(c *net.Client, funcName string) bool {
	if !c.IsAuthenticated || c.PlayerIdentifier == "" {
		a.Log.Error("User is not authenticated!",
			"function", fmt.Sprintf("AuthHandler::%s", funcName))
		c.Close()
		return false
	}

	player := a.Players.GetPlayer(c.PlayerIdentifier)
	if player == nil {
		a.Log.Error("Player not found!",
			"function", fmt.Sprintf("AuthHandler::%s", funcName),
			"accountName", c.PlayerIdentifier)
		c.Close()
		return false
	}
	return true
}
