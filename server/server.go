package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)


type Server struct {
	Host     string
	Port     string
	instance net.Listener
	clients  map[string]*Client
	leaving  chan *Client
	entering chan net.Conn
	message  chan Message
	system   *Client
}


func Default() *Server {
	_server := &Server{
		Host: "0.0.0.0",
		Port: "8888",
		entering: make(chan net.Conn),
		leaving: make(chan *Client),
		message: make(chan Message, 100),
		clients: make(map[string]*Client),
		system:  &Client{
			Name: "system",
		},
	}
	return _server
}


func (s *Server) Run() {
	address := s.Host + ":" + s.Port
	log.Println("listen on ", address)
	listener, _ := net.Listen("tcp", address)
	defer listener.Close()
	go s.dispatch()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panic(err)
		}
		s.entering <- conn
	}
}


func (s *Server) HandleEnter(conn net.Conn) {
	ip := conn.RemoteAddr().String()
	if client, ok := s.clients[ip]; ok {
		msg := constructMsg("welcome " + client.Name + " again", client)
		s.message <- msg
		return
	}
	log.Println(ip, "enter")
	_, _ = fmt.Fprintln(conn, "hey guy, please tell me your name")
	_, _ = fmt.Fprint(conn, "$ ")

	scan := bufio.NewScanner(conn)
	if scan.Scan() {
		name := scan.Text()
		// todo name must be unique
		client := Client{
			IP: ip,
			Name: name,
			Conn: conn,
		}
		s.clients[ip] = &client
		log.Println(s.clients)

		welcomeStr := "welcome " + client.Name
		message := constructMsg(welcomeStr, s.system)
		s.message <- message
		go s.ReceiveMsg(&client)
	}
}


func (s *Server) ReceiveMsg(client *Client) {
	scan := bufio.NewScanner(client.Conn)
	for scan.Scan() {
		text := scan.Text()
		message := constructMsg(text, client)
		s.message <- message
	}

	s.leaving <- client
}


func (s *Server) dispatch() {
	for {
		select {
			case conn := <- s.entering:
				go s.HandleEnter(conn)
			case client := <- s.leaving:
				go s.handleLeave(client)
			case msg := <- s.message:
				go s.broadcast(msg)
			}
	}
}

func (s *Server) handleLeave(client *Client) {
	message := constructMsg(client.Name + " is leave" , s.system)
	s.message <- message
	delete(s.clients, client.IP)
}

func (s *Server) broadcast(message Message) {
	log.Println("broadcast")
	for _, client := range s.clients {
		s.sendMsg(client, message)
	}
}


func (s *Server) sendMsg(client *Client, message Message){
	_, err := fmt.Fprintf(client.Conn, "\n%s\n$ ", message.msg)
	if err != nil {
		s.leaving <- client
	}
}


func constructMsg(msg string, own *Client) Message {
	name := fmt.Sprintf("[%s]:", own.Name)
	t    := time.Now().Format("2006-01-02 15:04:05")
	return Message{
		msg:       fmt.Sprintf("%s %10s %s", t, name, msg),
		owner: 	   own,
		createdAt: time.Now(),
	}
}