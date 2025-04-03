package game

import "mononoke-go/net/packets"

const AuthGameClientLoginID = 20011

type AuthGameClientLogin struct {
	Header               packets.Message
	Account              [61]byte
	AccountID            uint32
	Result               uint16
	Permission           uint32
	PCBangUser           uint8
	EventCode            uint32
	Age                  uint32
	ContinuousPlayTime   uint32
	ContinuousLogoutTime uint32
}
