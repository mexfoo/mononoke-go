//nolint:revive // It has to be that way to stay original
package game

import "mononoke-go/net/packets"

const GameAuthClientLogoutID = 20012

type GameAuthClientLogout struct {
	Header             packets.Message
	Account            [61]byte
	ContinuousPlayTime uint32
}
