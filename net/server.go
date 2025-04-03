package net

import (
	"fmt"
	"log/slog"
	"mononoke-go/net/packets"
	"net"
)

// TCP Server.
type Server struct {
	address                  string // Address to open connection: localhost:9999.
	Log                      *slog.Logger
	encryptClient            bool
	encryptionKey            string
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, header packets.Message, message []byte)
}

// Called right after Server starts listening new client.
func (s *Server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed.
func (s *Server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message.
func (s *Server) OnNewMessage(callback func(c *Client, header packets.Message, message []byte)) {
	s.onNewMessage = callback
}

// Listen starts network Server.
func (s *Server) Listen() error {
	var listener net.Listener
	var err error
	listener, err = net.Listen("tcp", s.address)
	if err != nil {
		s.Log.Error(fmt.Sprintf("Error starting TCP Server: %s", err))
		return err
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		client := Client{
			conn:   conn,
			Server: s,
			Log:    s.Log,
		}
		go client.listen(s.encryptionKey)
	}
}

// Creates new tcp Server instance.
func NewTCPServer(address string, encrypt bool, key string, log *slog.Logger) *Server {
	log.Info(fmt.Sprintf("Creating Server with address %s", address))
	serverInstance := &Server{
		address:       address,
		encryptClient: encrypt,
		Log:           log,
		encryptionKey: key,
	}

	serverInstance.OnNewClient(func(_ *Client) {})
	serverInstance.OnNewMessage(func(_ *Client, _ packets.Message, _ []byte) {})
	serverInstance.OnClientConnectionClosed(func(_ *Client, _ error) {})

	return serverInstance
}
