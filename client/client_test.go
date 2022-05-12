package client

import (
	"github.com/sinchop/smpp/message"
	"github.com/sinchop/smpp/server"

	"testing"
)

func TestClient(t *testing.T) {

	// Create SMPP server with simple handler function
	s := server.NewServer("test", 3601,
		func(c server.Conn, submit *message.ShortMessage) (*message.ShortMessageResp, error) {
			return &message.ShortMessageResp{
				MessageID: "1234",
				Status:    message.Status_OK,
			}, nil
		})

	// Add client account to server
	s.AddAccount(&server.Account{
		UserName: "client",
		Password: "pw",
	})

	// Start server
	s.Start()

	// Create SMPP client
	client := NewClient("127.0.0.1:3601", "client", "pw")
	// Bind client
	client.Bind()

	// Send SubmitSM
	submitSmResp, err := client.SendSubmitSM(&message.ShortMessage{
		Src:        "467019191695",
		Dst:        "467373737373",
		Text:       []byte("Hello world"),
		DataCoding: message.DefaultType,
	})
	if err != nil {
		t.Fatal(err)
	}

	if submitSmResp.Status != message.Status_OK {
		t.Fatalf("Status is not OK")
	}
	if submitSmResp.MessageID != "1234" {
		t.Fatalf("Wrong messageID received")
	}

}
