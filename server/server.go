package smppserver

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
)

// HandlerFunc is the signature of a function passed to Server instances,
// that is called when client PDU messages arrive.
type HandlerFunc func(c Conn, m pdu.Body)

type Account struct {
	UserName string `valid: "required"`
	Password string `valid: "required"`
}

// Server is an SMPP server.
type Server struct {
	Accounts map[string]Account
	SystemID string `valid: "required"`
	Port     int    `valid: "required"`
	TLS      *tls.Config
	Handler  HandlerFunc

	conns []Conn
	mu    sync.Mutex
	l     net.Listener
}

func NewServer(systemID string, port int) *Server {
	return &Server{
		SystemID: systemID,
		Port:     port,
		Accounts: make(map[string]Account),
	}
}

func newLocalListener(port int) net.Listener {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err == nil {
		return l
	}
	if l, err = net.Listen("tcp6", "[::1]:"+strconv.Itoa(port)); err != nil {
		panic(fmt.Sprintf("smpp server: failed to listen on a port: %v", err))
	}
	return l
}

// Start starts the server.
func (srv *Server) Start() {
	_, err := govalidator.ValidateStruct(srv)
	if err != nil {
		panic(fmt.Sprintf("Server missing mandatory configuration: %v", err))

	}
	srv.l = newLocalListener(srv.Port)
	go srv.Serve()
}

// Addr returns the local address of the server, or an empty string
// if the server hasn't been started yet.
func (srv *Server) Addr() string {
	if srv.l == nil {
		return ""
	}
	return srv.l.Addr().String()
}

// Close stops the server, causing the accept loop to break out.
func (srv *Server) Close() {
	if srv.l == nil {
		panic("smpptest: server is not started")
	}
	srv.l.Close()
}

// Serve accepts new clients and handle them by authenticating the
// first PDU, expected to be a Bind PDU. Other PDUs will be handled by the Handler function defined for the server
func (srv *Server) Serve() {
	for {
		cli, err := srv.l.Accept()
		log.Println("Got conn")
		if err != nil {
			break // on srv.l.Close
		}

		c := newConn(cli)
		srv.conns = append(srv.conns, c)
		go srv.handle(c)
	}
}

// handle new clients.
func (srv *Server) handle(c *conn) {
	defer c.Close()
	if err := srv.auth(c); err != nil {
		if err != io.EOF {
			log.Println("smpp server: server auth failed:", err)
		}
		return
	}
	for {
		p, err := c.Read()
		if err != nil {
			if err != io.EOF {
				log.Println("smpp server: read failed:", err)
			}
			break
		}
		srv.Handler(c, p)
	}
}

// auth authenticate new clients.
func (srv *Server) auth(c *conn) error {
	p, err := c.Read()
	if err != nil {
		return err
	}
	var resp pdu.Body
	switch p.Header().ID {
	case pdu.BindTransmitterID:
		resp = pdu.NewBindTransmitterResp()
	case pdu.BindReceiverID:
		resp = pdu.NewBindReceiverResp()
	case pdu.BindTransceiverID:
		resp = pdu.NewBindTransceiverResp()
	default:
		return errors.New("unexpected pdu, want bind")
	}
	f := p.Fields()
	user := f[pdufield.SystemID]
	passwd := f[pdufield.Password]
	if user == nil || passwd == nil {
		return errors.New("malformed pdu, missing system_id/password")
	}
	account, ok := srv.Accounts[user.String()]
	if !ok {
		return errors.New("invalid user")
	}

	if passwd.String() != account.Password {
		return errors.New("invalid passwd")
	}
	resp.Fields().Set(pdufield.SystemID, srv.SystemID)

	return c.Write(resp)
}
