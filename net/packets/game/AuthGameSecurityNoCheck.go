package game

import "mononoke-go/net/packets"

type AuthGameSecurityNoCheck struct {
	Header  packets.Message
	Account [61]byte
	// Mode    uint32
	Result uint32
}
