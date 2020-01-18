package steam

import (
	"bufio"
	"net"
	"fmt"
	"encoding/binary"
	"bytes"

	"github.com/lucasmbaia/goskins/steam/protocol/steamland"
)

const (
	tcpConnectionMagic uint32 = 0x31305456 // "VT01"
)

type Client struct {
	connection  *Connection
	servers	    steamServers
	auth	    *Auth
	events	    chan interface{}
}

type Connection struct {
	write	*bufio.Writer
	read	*bufio.Reader
	conn	net.Conn
}

type Auth struct {
	Username    string
	Password    string
	Key	    string
	Code	    string
}

func NewClient() (c *Client, err error) {
	c = &Client{
		connection: &Connection{},
		auth:	    &Auth{},
		events:	    make(chan interface{}),
	}

	return
}

func (c *Client) Connect() (err error) {
	if c.servers, err = GetServersSteam(); err != nil {
		return
	}

	if c.connection.conn, err = net.Dial("tcp", c.servers.GetRandomServer()); err != nil {
	//if c.connection.conn, err = net.Dial("tcp", "162.254.193.7:27019"); err != nil {
		return
	}

	c.connection.write = bufio.NewWriter(c.connection.conn)
	c.connection.read = bufio.NewReader(c.connection.conn)

	go c.read()

	return
}

func (c *Client) handlePackets(packet *steamland.Packet) {
	switch packet.EMsg {
	case steamland.EMsg_ChannelEncryptRequest:
	case steamland.EMsg_ChannelEncryptResult:
	case steamland.EMsg_Multi:
	case steamland.EMsg_ClientCMList:
	}
}

func (c *Client) handleChannelEncryptRequest(packet *steamland.Packet) {
	var (
		body	*steamland.MsgChannelEncryptRequest
	)

	body = NewMsgChannelEncryptRequest()
	packet.ReadMsg(body)
}

func (c *Client) handleEvents() {
}

func (c *Client) Fatalf(f string, e ...interface{}) {
	fmt.Println(e...)
}

func (c *Client) read() {
	var (
		err	error
	)

	for {
		var (
			buffer	[]byte
			pl	uint32
			pm	uint32
			packet	*Packet
		)

		if err = binary.Read(c.connection.conn, binary.LittleEndian, &pl); err != nil {
			break
		}

		if err = binary.Read(c.connection.conn, binary.LittleEndian, &pm); err != nil {
			break
		}

		if pm != tcpConnectionMagic {
			err = fmt.Errorf("Invalid connection magic! Expected %d, got %d!", tcpConnectionMagic, pm)
			break
		}

		buffer = make([]byte, pl, pl)
		if _, err = c.connection.read.Read(buffer); err != nil {
			break
		}

		if packet,  err = steamland.NewPacket(buffer); err != nil {
			break
		}

		fmt.Println(buffer)
		c.handlePackets(packet)
	}

	c.Fatalf("Error reading from the connection: %v", err)
}
