//nolint:revive // It has to be that way to stay original
package game

import "mononoke-go/net/packets"

const GameAuthClientLoginID = 20010

type GameAuthClientLogin struct {
	Header     packets.Message
	Account    [61]byte
	OneTimeKey uint64
}
