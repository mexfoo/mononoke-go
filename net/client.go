package net

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log/slog"
	"mononoke-go/net/packets"
	"mononoke-go/utils"
	"net"
)

// Client holds info about connection.
type Client struct {
	conn             net.Conn
	Server           *Server
	Log              *slog.Logger
	PlayerIdentifier string
	GameIdentifier   uint32
	IsAuthenticated  bool
	encryptCipher    utils.RC4Cipher
	decryptCipher    utils.RC4Cipher
	AESKey           []byte
	SupportedVersion int32
}

// Get endpoint.
func (c *Client) GetEndpoint() string {
	return c.conn.RemoteAddr().String()
}

// Read client data from channel.
func (c *Client) listen(key string) {
	if c.Server.encryptClient {
		c.encryptCipher, c.decryptCipher = utils.RC4Cipher{Key: key}, utils.RC4Cipher{Key: key}
	}
	c.Server.onNewClientCallback(c)
	reader := bufio.NewReader(c.conn)
	for {
		header, err := reader.Peek(7)
		if err != nil {
			c.Log.Error(fmt.Sprintf("[net.Client] Cannot read message! %s", err.Error()))
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}

		miscHeader := make([]byte, 7)
		copy(miscHeader, header)
		parsedHeader := packets.Message{}
		if !c.readHeader(miscHeader, &parsedHeader) {
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}

		message := make([]byte, parsedHeader.HeaderMessageSize)
		_, err = reader.Read(message)

		if err != nil {
			c.Log.Error("[net.Client] Cannot parse TS_MESSAGE!")
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}

		if c.Server.encryptClient {
			c.decryptCipher.DoCipher(&message)
		}
		c.Server.onNewMessage(c, parsedHeader, message)
	}
}

func (c *Client) readHeader(header []byte, msg *packets.Message) bool {
	if c.Server.encryptClient {
		c.decryptCipher.TryCipher(&header)
	}
	if _, err := binary.Decode(header, binary.LittleEndian, msg); err != nil {
		return false
	}
	if msg.HeaderMessageChecksum == msg.GetHeaderChecksum() {
		return true
	}
	return false
}

// Send bytes to client.
func (c *Client) Send(content any, packetID uint16) {
	packet, err := utils.Marshal(binary.LittleEndian, content, int(c.SupportedVersion))
	if err != nil {
		c.Log.Error(fmt.Sprintf("[Client] Cannot write packet: %s", err.Error()))
		c.conn.Close()
		c.Server.onClientConnectionClosed(c, err)
		return
	}

	binary.LittleEndian.PutUint32(packet, uint32(len(packet))) //nolint:gosec // this is fine
	binary.LittleEndian.PutUint16(packet[4:], packetID)
	packet[6] = packets.SetHeaderChecksum(len(packet), int(packetID))

	if c.Server.encryptClient {
		c.encryptCipher.DoCipher(&packet)
	}

	_, err = c.conn.Write(packet)
	if err != nil {
		c.conn.Close()
		c.Server.onClientConnectionClosed(c, err)
	}
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	c.Server.onClientConnectionClosed(c, nil)
	return c.conn.Close()
}
