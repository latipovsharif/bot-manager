package platforms

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
)

type Whatsapp struct {
	wac *whatsapp.Conn
}

func (w *Whatsapp) HandleError(err error) {
	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := w.wac.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}

func (w *Whatsapp) HandleTextMessage(message whatsapp.TextMessage) {
	fmt.Printf("%v %v %v %v\n\t%v\n", message.Info.Timestamp, message.Info.Id, message.Info.RemoteJid, message.ContextInfo.QuotedMessageID, message.Text)
}

func (w *Whatsapp) Run() error {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		return fmt.Errorf("error creating connection: %v\n", err)
	}

	w.wac = wac

	w.wac.SetClientVersion(0, 4, 2080)
	err = w.login()
	if err != nil {
		return fmt.Errorf("error logging in: %v\n", err)
	}

	<-time.After(5 * time.Second)

	var chat whatsapp.Chat

	for _, c := range w.wac.Store.Chats {
		if c.Name == "Тест" {
			chat = c
			log.Println("found test chat which we will listen for messages")
		}
	}

	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: chat.Jid,
		},
		Text: "Message sent by github.com/Rhymen/go-whatsapp",
	}

	_, err = w.wac.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}

func (w *Whatsapp) login() error {
	//load saved session
	session, err := w.readSession()
	if err == nil {
		//restore session
		session, err = w.wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v\n", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = w.wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
	}

	//save session
	err = w.writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func (w *Whatsapp) readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func (w *Whatsapp) writeSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}
