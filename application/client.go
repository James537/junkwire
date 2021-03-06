package application

import (
	"net"

	"github.com/sodapanda/junkwire/connection"
	"github.com/sodapanda/junkwire/misc"
)

//AppClient client
type AppClient struct {
	conn       *net.UDPConn
	connAddr   *net.UDPAddr
	clientConn *connection.ClientConn
	rcv        int
}

//NewAppClient new client
func NewAppClient(listenPort string) *AppClient {
	ac := new(AppClient)
	addr, err := net.ResolveUDPAddr("udp4", ":"+listenPort)
	misc.CheckErr(err)
	conn, err := net.ListenUDP("udp4", addr)
	misc.CheckErr(err)
	ac.conn = conn
	return ac
}

//Start start
func (ac *AppClient) Start() {
	go ac.socketToDevice()
}

//SetClientConn set conn
func (ac *AppClient) SetClientConn(clientConn *connection.ClientConn) {
	ac.clientConn = clientConn
	if clientConn != nil {
		ac.clientConn.AddHandler(clientHandler{ac: ac})
	}
}

func (ac *AppClient) socketToDevice() {
	buffer := make([]byte, 2000)
	for {
		length, addr, err := ac.conn.ReadFromUDP(buffer)
		misc.CheckErr(err)
		ac.connAddr = addr
		data := buffer[:length]
		if ac.clientConn != nil {
			ac.clientConn.Write(data, false)
		}
	}
}

type clientHandler struct {
	ac *AppClient
}

func (ch clientHandler) OnData(data []byte) {
	ch.ac.rcv++
	_, err := ch.ac.conn.WriteToUDP(data, ch.ac.connAddr)
	if err != nil {
		misc.PLog(err.Error())
	}
}
func (ch clientHandler) OnDisconnect(cc *connection.ClientConn) {}
func (ch clientHandler) OnConnect(cc *connection.ClientConn)    {}

//IClient client
type IClient interface {
	Start()
	SetClientConn(clientConn *connection.ClientConn)
}
