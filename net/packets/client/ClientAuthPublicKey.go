//nolint:revive // It has to be that way to stay original
package client

import "mononoke-go/net/packets"

const ClientAuthPublicKeyID1 = 71
const ClientAuthPublicKeyID2 = 1071

type ClientAuthPublicKey struct {
	Header packets.Message
	Size   uint32
	Key    []byte `byteSize:"Size"`
}
