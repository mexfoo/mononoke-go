package client_test

import (
	"encoding/binary"
	"mononoke-go/net/packets"
	"mononoke-go/net/packets/client"
	"mononoke-go/utils"
	"testing"
)

func TestAuthClientServerList(t *testing.T) {
	pkt := client.AuthClientServerList{
		Header: packets.Message{
			HeaderMessageSize:     400,
			HeaderMessageId:       10,
			HeaderMessageChecksum: 3,
		},
		LastLoginServerIdx: 1,
		Servers:            2,
		ServerInfo: []client.ServerInfo{
			{
				ServerIdx:           1,
				ServerName:          [21]byte{1, 2, 3, 4},
				IsAdultServer:       0,
				ServerScreenshotURL: [256]byte{1, 2, 3, 4},
				ServerIP:            [16]byte{1, 2, 3, 4},
				ServerPort:          4500,
				UserRatio:           1,
			},
			{
				ServerIdx:           2,
				ServerName:          [21]byte{1, 2, 3, 4},
				IsAdultServer:       1,
				ServerScreenshotURL: [256]byte{1, 2, 3, 4},
				ServerIP:            [16]byte{1, 2, 3, 4},
				ServerPort:          4502,
				UserRatio:           2,
			},
		},
	}

	result, err := utils.Marshal(binary.LittleEndian, pkt, 0x090603)
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%d", len(result))
}
