//nolint:revive // It has to be that way to stay original
package client

import "mononoke-go/net/packets"

const ClientAuthSelectServerID = 10023

type ClientAuthSelectServer struct {
	Header    packets.Message
	ServerIdx uint32 `subtype:"9" subversion:"0x000000.0x090605"`
}
