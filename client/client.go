package client

import (
	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
	"github.com/sinchop/smpp/message"
)

type Client struct {
	ServerAddr string
	UserName   string
	Password   string
	tx         *smpp.Transmitter
}

func NewClient(serverAddr string, userName string, password string) *Client {
	return &Client{
		ServerAddr: serverAddr,
		UserName:   userName,
		Password:   password,
	}

}

func (client *Client) Bind() {
	client.tx = &smpp.Transmitter{
		Addr:   client.ServerAddr,
		User:   client.UserName,
		Passwd: client.Password,
	}
	conn := client.tx.Bind()
	var status smpp.ConnStatus
	if status = <-conn; status.Error() != nil {
		log.Fatalln("Unable to connect, aborting:", status.Error())
	}
	log.Println("Connection completed, status:", status.Status().String())
	// connection checker goroutine
	go func() {
		for c := range conn {
			log.Println("SMPP connection status:", c.Status())
		}
	}()
}

func (client *Client) SendSubmitSM(sm *message.ShortMessage) (*message.ShortMessageResp, error) {
	s := &smpp.ShortMessage{
		Src:           sm.Src,
		SourceAddrNPI: sm.SourceAddrNPI,
		SourceAddrTON: sm.SourceAddrTON,
		Dst:           sm.Dst,
		DestAddrNPI:   sm.DestAddrNPI,
		DestAddrTON:   sm.DestAddrTON,
		Register:      pdufield.DeliverySetting(sm.Register),
	}

	switch sm.DataCoding {
	case message.DefaultType:
		s.Text = pdutext.GSM7(sm.Text)
		break
	case message.UCS2Type:
		s.Text = pdutext.UCS2(sm.Text)
		break
	}
	s, err := client.tx.Submit(s)
	if err != nil {
		return nil, err
	}
	rawResp := s.Resp()
	return &message.ShortMessageResp{
		MessageID: rawResp.Fields()[pdufield.MessageID].String(),
		Status:    message.Status(rawResp.Header().Status),
	}, nil

}
