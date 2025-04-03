//nolint:revive // It has to be that way to stay original
package client

import "mononoke-go/net/packets"

const ClientAuthServerListID = 10021

type ClientAuthServerList struct {
	Header packets.Message
}
