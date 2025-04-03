package client

import "mononoke-go/net/packets"

const AuthClientResultID = 10000
const (
	LoginFlagEulaAccepted        = 1
	LoginFlagAccountBlockWarning = 2
)

type AuthClientResult struct {
	Header           packets.Message
	RequestMessageID uint16
	Result           uint16
	LoginFlag        int32
}
