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

func TestVersionPkt(t *testing.T) {
	ver := client.ClientAuthVersion{
		Header:  packets.Message{HeaderMessageSize: 27, HeaderMessageId: 50, HeaderMessageChecksum: 0},
		Version: [20]byte{0, 1, 2, 3, 4, 5, 6, 7},
	}

	marshall, ok := utils.Marshal(binary.LittleEndian, ver, 0x080100)
	if ok != nil {
		t.Errorf("%s", ok.Error())
	}

	ver2 := client.ClientAuthVersion{}
	reader := bytes.NewBuffer(marshall)
	utils.Unmarshal(reader, binary.LittleEndian, &ver2, 0x00000)

	if !reflect.DeepEqual(ver, ver2) {
		t.Errorf("Invalid result")
	}
}
