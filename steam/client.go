package steam

import (
	"bufio"
	"net"
	"fmt"
	"encoding/binary"
	"bytes"
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
			t	uint32
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

		if err = binary.Read(bytes.NewReader(buffer), binary.LittleEndian, &t); err != nil {
			break
		}

		fmt.Println(t)
		fmt.Println(buffer)
	}

	c.Fatalf("Error reading from the connection: %v", err)
}
