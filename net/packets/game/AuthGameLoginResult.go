package game

import "mononoke-go/net/packets"

const AuthGameLoginResultID = 20002

type AuthGameLoginResult struct {
	Header packets.Message
	Result uint16
}
