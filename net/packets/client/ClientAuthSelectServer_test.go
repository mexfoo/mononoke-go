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

func TestClientAuthSelectServerU16(t *testing.T) {
	accountStruct := client.ClientAuthSelectServer{
		Header: packets.Message{
			HeaderMessageId:       10000,
			HeaderMessageSize:     0,
			HeaderMessageChecksum: 0,
		},
		ServerIdx: 20,
	}
	result, err := utils.Marshal(binary.LittleEndian, accountStruct, 0x040000)
	if err != nil {
		t.Error(err.Error())
	}
	if len(result) != 9 {
		t.Errorf("invalid length. Expected 9, received %d", len(result))
	}

	reader := bytes.NewBuffer(result)
	var message client.ClientAuthSelectServer

	err = utils.Unmarshal(reader, binary.LittleEndian, &message, 0x040000)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(accountStruct, message) {
		t.Error("invalid DeepEqual!")
	}
}

func TestClientAuthSelectServerU32(t *testing.T) {
	accountStruct := client.ClientAuthSelectServer{
		Header: packets.Message{
			HeaderMessageId:       10000,
			HeaderMessageSize:     0,
			HeaderMessageChecksum: 0,
		},
		ServerIdx: 20,
	}
	result, err := utils.Marshal(binary.LittleEndian, accountStruct, 0x090606)
	if err != nil {
		t.Error(err.Error())
	}
	if len(result) != 11 {
		t.Errorf("invalid length. Expected 11, received %d", len(result))
	}

	reader := bytes.NewBuffer(result)
	var message client.ClientAuthSelectServer

	err = utils.Unmarshal(reader, binary.LittleEndian, &message, 0x090606)
	if err != nil {
		t.Error(err.Error())
	}
	if !reflect.DeepEqual(accountStruct, message) {
		t.Error("invalid DeepEqual!")
	}
}
