package utils

type RC4Cipher struct {
	x           int32
	y           int32
	state       [256]byte
	isInitiated bool
	Key         string
}

func (c *RC4Cipher) init(key string) {
	if c.isInitiated {
		return
	}

	if len(key) == 0 {
		return
	}

	for i := 0; i < len(c.state); i++ {
		c.state[i] = byte(i)
	}

	keyAsByte := []byte(key)
	keyLen := len(keyAsByte)
	tempKey := [len(c.state)]byte{}
	for i, j := 0, 0; i < len(c.state); i++ {
		tempKey[i] = keyAsByte[j]
		if j+1 >= keyLen {
			j = 0
		} else {
			j++
		}
	}

	for i, j := 0, 0; i < len(c.state); i++ {
		j = (j + int(c.state[j]) + int(tempKey[j])) & 0xFF
		c.state[i], c.state[j] = c.state[j], c.state[i]
	}

	c.x, c.y = 0, 0
	c.skipFor(1013)
	c.isInitiated = true
}

func (c *RC4Cipher) skipFor(length int) {
	currX, currY := c.x, c.y

	for length != 0 {
		currX = (currX + 1) & 0xff
		sx := c.state[currX]

		currY = (currY + int32(sx)) & 0xff
		c.state[currX] = c.state[currY]
		c.state[currY] = sx
		length--
	}

	c.x, c.y = currX, currY
}

func (c *RC4Cipher) DoCipher(content *[]byte) {
	if !c.isInitiated {
		c.init(c.Key)
	}

	x, y := c.x, c.y
	for i := range *content {
		x = (x + 1) & 0xff
		sx := c.state[x]
		y = (y + int32(sx)) & 0xff
		sy := c.state[y]
		c.state[x] = sy
		c.state[y] = sx

		(*content)[i] ^= c.state[sx+sy&0xff]
	}
	c.x, c.y = x, y
}

func (c *RC4Cipher) TryCipher(content *[]byte) {
	if !c.isInitiated {
		c.init(c.Key)
	}

	x, y := c.x, c.y
	state := c.state
	c.DoCipher(content)
	c.x, c.y = x, y
	c.state = state
}
