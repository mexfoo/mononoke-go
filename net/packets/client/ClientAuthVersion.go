//nolint:revive // It has to be that way to stay original
package client

import "mononoke-go/net/packets"

const ClientAuthVersionID = 10001

type ClientAuthVersion struct {
	Header  packets.Message
	Version [20]byte
}
