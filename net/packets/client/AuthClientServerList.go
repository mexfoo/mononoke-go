package client

import (
	"mononoke-go/net/packets"
)

const AuthClientServerListID = 10022

type ServerInfo struct {
	ServerIdx           uint32 `subtype:"9" subversion:"0x000000.0x090604"`
	ServerName          [21]byte
	IsAdultServer       uint8     `version:"0x040100.0x999999"`
	ServerScreenshotURL [256]byte `version:"0x040100.0x999999"`
	ServerIP            [16]byte
	ServerPort          int32
	UserRatio           uint16
}

type AuthClientServerList struct {
	Header             packets.Message
	LastLoginServerIdx uint32       `subtype:"9" subversion:"0x000000.0x090604"`
	Servers            uint32       `subtype:"9" subversion:"0x000000.0x090604"`
	ServerInfo         []ServerInfo `byteSize:"Servers"`
}
