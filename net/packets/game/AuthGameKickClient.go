package game

import "mononoke-go/net/packets"

const AuthGameKickClientID = 20013

type AuthGameKickClient struct {
	Header  packets.Message
	Account [61]byte
}
