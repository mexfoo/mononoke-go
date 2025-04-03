package client

import "mononoke-go/net/packets"

const (
	AuthClientAESKeyID1 = 72
	AuthClientAESKeyID2 = 1072
)

type AuthClientAESKey struct {
	Header  packets.Message
	KeySize uint32
	Key     []byte `byteSize:"KeySize"`
}
