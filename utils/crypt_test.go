package utils_test

import (
	"mononoke-go/utils"
	"reflect"
	"testing"
)

func TestInitDESKey(t *testing.T) {
	key := "test"
	want := [8]byte{0x74, 0x65, 0x73, 0x74, 0x00, 0x00, 0x00, 0x00}
	result := utils.InitDESKey(key)
	if !reflect.DeepEqual(want, result) {
		t.Errorf(`InitDESKey("test") = %#v, want match for %#v, nil`, result, want)
	}
}

func TestHashPassword(t *testing.T) {
	password := "helloworld"
	result, err := utils.HashPassword(password)
	if err != nil {
		t.Errorf(`HashPassword("helloworld") threw error %s`, err.Error())
	}
	if !utils.VerifyPassword(password, result) {
		t.Errorf(`HashPassword("helloworld") = %s, can't verify!`, result)
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "$2a$14$YYtz2pCu3YBI8fOVYUSYuOXAgkBzeOZO2k02p/JqUpUFzBYJ8AE9O"
	if !utils.VerifyPassword("helloworld", password) {
		t.Errorf(`Verifypassword("helloworld", "$2a$14$YYtz2pCu3YBI8fOVYUSYuOXAgkBzeOZO2k02p/JqUpUFzBYJ8AE9O") = false`)
	}
}
