//nolint:revive // It has to be that way to stay original
package game

import "mononoke-go/net/packets"

const GameAuthSecurityNoCheckID = 40001

type GameAuthSecurityNoCheck struct {
	Header   packets.Message
	Account  [61]byte
	Security [19]byte
}
