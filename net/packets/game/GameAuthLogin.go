//nolint:revive // It has to be that way to stay original
package game

import "mononoke-go/net/packets"

const GameAuthLoginID = 20001

type GameAuthLogin struct {
	Header              packets.Message
	ServerIdx           uint16
	ServerName          [21]byte
	ServerScreenshotURL [256]byte
	IsAdultServer       byte
	ServerIP            [16]byte
	ServerPort          int32
}
