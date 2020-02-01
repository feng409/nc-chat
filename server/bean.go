package server

import (
	"net"
	"time"
)

type Client struct {
	Name string
	IP   string
	Conn net.Conn
}


type Message struct {
	msg       string
	owner     *Client
	createdAt time.Time
}

