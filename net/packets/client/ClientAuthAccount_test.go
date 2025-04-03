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

func TestClientAuthAccountDefault(t *testing.T) {
	pkt := accountPktDefault{
		Header:      packets.Message{HeaderMessageSize: 58, HeaderMessageId: 10000, HeaderMessageChecksum: 0},
		AccountName: [19]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64},
		Password:    [32]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64},
	}
	result, err := utils.Marshal(binary.LittleEndian, pkt, 0x040101)
	if err != nil {
		t.Error(err.Error())
	}
	reader := bytes.NewBuffer(result)
	newPkt := client.ClientAuthAccount{}
	err = utils.Unmarshal(reader, binary.LittleEndian, &newPkt, 0x040101)
	if err != nil {
		t.Error(err.Error())
	}
	if utils.CToGoString(newPkt.Account) != "Hello World" {
		t.Errorf("Invalid account. Expected 'Hello World', received %s %v", utils.CToGoString(newPkt.Account), newPkt.Account)
	}
	if utils.CToGoString(newPkt.Password) != "Hello World" {
		t.Errorf("Invalid PW. Expected 'Hello World', received %s %v", utils.CToGoString(newPkt.Password), newPkt.Password)
	}
}

func TestClientAuthAccount050200(t *testing.T) {
	pkt := accountPkt050200{
		Header:      packets.Message{HeaderMessageSize: 58, HeaderMessageId: 10000, HeaderMessageChecksum: 0},
		AccountName: [61]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64},
		Password:    [61]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64},
	}
	result, err := utils.Marshal(binary.LittleEndian, pkt, 0x050200)
	if err != nil {
		t.Error(err.Error())
	}
	reader := bytes.NewBuffer(result)
	newPkt := client.ClientAuthAccount{}
	err = utils.Unmarshal(reader, binary.LittleEndian, &newPkt, 0x050200)
	if err != nil {
		t.Error(err.Error())
	}
	if utils.CToGoString(newPkt.Account) != "Hello World" {
		t.Errorf("Invalid account. Expected 'Hello World', received %s %v", utils.CToGoString(newPkt.Account), newPkt.Account)
	}
	if utils.CToGoString(newPkt.Password) != "Hello World" {
		t.Errorf("Invalid PW. Expected 'Hello World', received %s %v", utils.CToGoString(newPkt.Password), newPkt.Password)
	}
}

func TestClientAuthAccount080101(t *testing.T) {
	pkt := accountPkt080101{
		Header:       packets.Message{HeaderMessageSize: 58, HeaderMessageId: 10000, HeaderMessageChecksum: 0},
		AccountName:  [61]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64},
		PasswordSize: 11,
		Password:     [77]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64},
	}
	result, err := utils.Marshal(binary.LittleEndian, pkt, 0x080101)
	if err != nil {
		t.Error(err.Error())
	}
	reader := bytes.NewBuffer(result)
	newPkt := client.ClientAuthAccount{}
	err = utils.Unmarshal(reader, binary.LittleEndian, &newPkt, 0x080101)
	if err != nil {
		t.Error(err.Error())
	}
	if utils.CToGoString(newPkt.Account) != "Hello World" {
		t.Errorf("Invalid account. Expected 'Hello World', received %s %v", utils.CToGoString(newPkt.Account), newPkt.Account)
	}
	if utils.CToGoString(newPkt.Password) != "Hello World" {
		t.Errorf("Invalid PW. Expected 'Hello World', received %s %v", utils.CToGoString(newPkt.Password), newPkt.Password)
	}
	if len(utils.CToGoString(newPkt.Password)) != int(newPkt.PasswordSize) {
		t.Errorf("Invalid len. Expected %d, got %d", len(utils.CToGoString(newPkt.Password)), newPkt.PasswordSize)
	}
}

func TestClientAuthAccount090606(t *testing.T) {
	pkt := accountPkt090606{
		Header:       packets.Message{HeaderMessageSize: 58, HeaderMessageId: 10000, HeaderMessageChecksum: 0},
		AccountName:  [56]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64},
		MAC:          [8]byte{0, 1, 2, 3, 4, 5, 6, 7},
		PasswordSize: 11,
		Password:     [516]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64},
	}
	result, err := utils.Marshal(binary.LittleEndian, pkt, 0x090606)
	if err != nil {
		t.Error(err.Error())
	}
	reader := bytes.NewBuffer(result)
	newPkt := client.ClientAuthAccount{}
	err = utils.Unmarshal(reader, binary.LittleEndian, &newPkt, 0x090606)
	if err != nil {
		t.Error(err.Error())
	}
	if utils.CToGoString(newPkt.Account) != "Hello World" {
		t.Errorf("Invalid account. Expected 'Hello World', received %s %v", utils.CToGoString(newPkt.Account), newPkt.Account)
	}
	if utils.CToGoString(newPkt.Password) != "Hello World" {
		t.Errorf("Invalid PW. Expected 'Hello World', received %s %v", utils.CToGoString(newPkt.Password), newPkt.Password)
	}
	if len(utils.CToGoString(newPkt.Password)) != int(newPkt.PasswordSize) {
		t.Errorf("Invalid length. Expected %d, got %d", len(utils.CToGoString(newPkt.Password)), newPkt.PasswordSize)
	}
	if !reflect.DeepEqual(pkt.MAC, newPkt.MacStamp) {
		t.Errorf("Invalid MAC. Expected %v, got %v", pkt.MAC, newPkt.MacStamp)
	}
}

func TestClientAuthAccount090607(t *testing.T) {
	pkt := accountPkt090607{
		Header:       packets.Message{HeaderMessageSize: 58, HeaderMessageId: 10000, HeaderMessageChecksum: 0},
		AccountName:  [56]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64},
		MAC:          [8]byte{0, 1, 2, 3, 4, 5, 6, 7},
		PasswordSize: 11,
		Password:     []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64},
	}
	result, err := utils.Marshal(binary.LittleEndian, pkt, 0x090607)
	if err != nil {
		t.Error(err.Error())
	}
	reader := bytes.NewBuffer(result)
	newPkt := client.ClientAuthAccount{}
	err = utils.Unmarshal(reader, binary.LittleEndian, &newPkt, 0x090607)
	if err != nil {
		t.Error(err.Error())
	}
	if utils.CToGoString(newPkt.Account) != "Hello World" {
		t.Errorf("Invalid account. Expected 'Hello World', received %s %v", utils.CToGoString(newPkt.Account), newPkt.Account)
	}
	if utils.CToGoString(newPkt.Password) != "Hello World" {
		t.Errorf("Invalid PW. Expected 'Hello World', received %s %v", utils.CToGoString(newPkt.Password), newPkt.Password)
	}
	if len(utils.CToGoString(newPkt.Password)) != int(newPkt.PasswordSize) {
		t.Errorf("Invalid length. Expected %d, got %d", len(utils.CToGoString(newPkt.Password)), newPkt.PasswordSize)
	}
	if !reflect.DeepEqual(pkt.MAC, newPkt.MacStamp) {
		t.Errorf("Invalid MAC. Expected %v, got %v", pkt.MAC, newPkt.MacStamp)
	}
}

func TestAESDecryption(t *testing.T) {
	msg := []byte{149, 0, 0, 0, 26, 39, 214, 116, 101, 115, 116, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 16, 0, 0, 0, 112, 135, 125, 94, 113, 36, 24, 52, 254,
		49, 255, 234, 1, 172, 17, 149, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0}
	reader := bytes.NewBuffer(msg)
	newPkt := client.ClientAuthAccount{}
	err := utils.Unmarshal(reader, binary.LittleEndian, &newPkt, 0x080101)
	if err != nil {
		t.Error(err.Error())
	}
	if utils.CToGoString(newPkt.Account) != "test" {
		t.Errorf("Invalid account, expected %s, received %s", "test", utils.CToGoString(newPkt.Account))
	}
}

type accountPktDefault struct {
	Header      packets.Message
	AccountName [19]byte
	Password    [32]byte
}

type accountPkt050200 struct {
	Header      packets.Message
	AccountName [61]byte
	Password    [61]byte
}

type accountPkt080101 struct {
	Header       packets.Message
	AccountName  [61]byte
	PasswordSize uint32
	Password     [77]byte
}

type accountPkt090606 struct {
	Header       packets.Message
	AccountName  [56]byte
	MAC          [8]byte
	PasswordSize uint32
	Password     [516]byte
}

type accountPkt090607 struct {
	Header       packets.Message
	AccountName  [56]byte
	MAC          [8]byte
	PasswordSize uint32
	Password     []byte
}
