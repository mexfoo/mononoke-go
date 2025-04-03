//nolint:revive // It has to be that way to stay original
package client

import "mononoke-go/net/packets"

const ClientAuthAccountID = 10010

type ClientAuthAccount struct {
	Header       packets.Message
	Account      []byte  `loop:"3" len0:"19" version0:"0x000000.0x050199" len1:"61" version1:"0x050200.0x090605" len2:"56" version2:"0x090606.0x999999"` //nolint:lll // Has to be.
	MacStamp     [8]byte `version:"0x090606.0x999999"`
	PasswordSize uint32  `version:"0x080101.0x999999"`
	Password     []byte  `loop:"5" len0:"32" version0:"0x000000.0x050199" len1:"61" version1:"0x050200.0x080100" len2:"77" version2:"0x080101.0x090605" len3:"516" version3:"0x090606.0x090606" len4:"PasswordSize" version4:"0x090707.0x999999"` //nolint:lll // Has to be.
}
