//nolint:revive // It has to be that way to stay original
package game

import "mononoke-go/net/packets"

const GameAuthClientKickFailedID = 20014

type GameAuthClientKickFailed struct {
	Header  packets.Message
	Account [61]byte
}
