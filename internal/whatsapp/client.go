package whatsapp

import (
	"context"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var (
	client *whatsmeow.Client
	ctx    = context.Background()
)

func ConnectToWhatsApp() error {
	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New(ctx, "sqlite3", fmt.Sprint("file:", os.Getenv("DB"), "?_foreign_keys=on"), dbLog)
	if err != nil {
		return err
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return err
	}

	clientLog := waLog.Stdout("Client", "ERROR", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)

	client.AddEventHandler(eventHandler)

	if client.Store.ID != nil {
		err := client.Connect()
		if err != nil {
			fmt.Println("error trying to connect")
			return err
		}

		fmt.Println("connected!")
		return nil
	}

	qrChan, _ := client.GetQRChannel(context.Background())
	err = client.Connect()
	if err != nil {
		return err
	}

	for evt := range qrChan {
		if evt.Event != "code" {
			continue
		}

		qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
	}

	return nil
}

func SendMessage(recipientJID string, message string) {
	jid, err := types.ParseJID(recipientJID)
	if err != nil {
		fmt.Printf("error parsing JID: %v\n", err)
		return
	}

	msg := &waE2E.Message{
		Conversation: &message,
	}

	_, err = client.SendMessage(ctx, jid, msg)
	if err != nil {
		fmt.Printf("error sending message: %v\n", err)
	} else {
		fmt.Println("message sent successfully!")
	}
}

func GetContacts() (map[types.JID]types.ContactInfo, error) {
	contacts, err := client.Store.Contacts.GetAllContacts(ctx)
	return contacts, err
}

func eventHandler(evt any) {
	switch v := evt.(type) {
	case *events.Connected:
		fmt.Println("connected!")
	case *events.Message:
		if v.Info.IsGroup {
			break
		}

		sender := v.Info.Sender.String()

		messageText := ""

		if imgMsg := v.Message.GetImageMessage(); imgMsg != nil {
			break
		}

		if msg := v.Message.GetConversation(); msg != "" {
			messageText = msg
		} else if extMsg := v.Message.GetExtendedTextMessage(); extMsg != nil {
			messageText = extMsg.GetText()
		} else {
			messageText = "[Unsupported message type]"
		}

		if strings.ToLower(messageText) == "ping" {
			SendMessage(sender, "pong")
		}

	case *events.Receipt:

	case *events.Disconnected:
		fmt.Println("disconnected from WhatsApp")
	}
}
