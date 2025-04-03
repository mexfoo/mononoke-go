package client

import "mononoke-go/net/packets"

const AuthClientSelectServerID = 10024

type AuthClientSelectServer struct {
	Header        packets.Message
	Result        uint16
	OneTimeKey    uint64   `version:"0x000000.0x080100"`
	EncryptedSize int32    `version:"0x080101.0x999999"`
	EncryptedData [24]byte `version:"0x080101.0x999999"`
	PendingTime   uint32
}
