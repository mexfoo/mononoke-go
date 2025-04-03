package client

import (
	"mononoke-go/net/packets"
)

const AuthClientResultWithStringID = 10002

type AuthClientResultWithString struct {
	Header      packets.Message
	Body        AuthClientResult
	MessageSize uint32
	Message     []byte `byteSize:"MessageSize"`
}
