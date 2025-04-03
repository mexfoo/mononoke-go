package client_test

import (
	"bytes"
	"encoding/binary"
	"mononoke-go/net/packets"
	"mononoke-go/net/packets/client"
	"mononoke-go/utils"
	"reflect"
	"testing"
)

func TestClientAuthPublicKey(t *testing.T) {
	accountStruct := client.ClientAuthPublicKey{
		Header: packets.Message{
			HeaderMessageId:       10000,
			HeaderMessageSize:     0,
			HeaderMessageChecksum: 0,
		},
		Size: 5,
		Key:  []byte{0x1, 0x2, 0x3, 0x4, 0x5},
	}
	result, err := utils.Marshal(binary.LittleEndian, accountStruct, 0x100000)
	if err != nil {
		t.Error(err.Error())
	}
	if len(result) != 16 {
		t.Error("len <> 16")
	}

	reader := bytes.NewBuffer(result)

	var message client.ClientAuthPublicKey

	err = utils.Unmarshal(reader, binary.LittleEndian, &message, 0x010000)
	if err != nil {
		t.Error(err.Error())
	}
	if accountStruct.Size != message.Size {
		t.Error("invalid size")
	}
	t.Log("Testing DeepEqual for key...")
	if !reflect.DeepEqual(accountStruct.Key, message.Key) {
		t.Error("invalid key")
	}
}
