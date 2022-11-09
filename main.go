package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/mjehanno/gchat/service/chat"
	"google.golang.org/grpc"
)

func main() {
	a := app.New()
	w := a.NewWindow("gchat")
	rand.Seed(time.Now().UnixNano())
	length := 9

	ran_str := make([]byte, length)

	// Generating Random string
	for i := 0; i < length; i++ {
		ran_str[i] = byte(65 + rand.Intn(25))
	}
	username := string(ran_str)

	data := binding.BindStringList(&[]string{})

	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("client cannot connect to server : %s", err.Error())
	}

	client := chat.NewChatServiceClient(conn)

	ctx := context.Background()

	messageStream, err := client.ExchangeMessage(ctx)
	if err != nil {
		log.Fatal("failed to connect to grpc server")
	}

	go func() {
		for {
			message, err := messageStream.Recv()
			if err == io.EOF {
				return
			}
			var mu sync.Mutex
			if message.Author != username {
				mu.Lock()
				err = data.Append(fmt.Sprintf("%s send : %s", message.Author, message.Content))
				if err != nil {
					log.Println("error while adding message to list : ", err.Error())
				}
				mu.Unlock()
			}
		}
	}()

	c := container.NewGridWithRows(2, messageDisplayer(data), messageBox(ctx, messageStream, data, username))

	w.SetContent(c)

	w.Resize(fyne.NewSize(500, 300))
	w.ShowAndRun()
}

func messageDisplayer(data binding.ExternalStringList) *fyne.Container {
	list := widget.NewListWithData(data, func() fyne.CanvasObject {
		return widget.NewLabel("message:")
	}, func(i binding.DataItem, o fyne.CanvasObject) {
		o.(*widget.Label).Bind(i.(binding.String))
	})

	c := container.NewMax(list)
	return c
}

func messageBox(ctx context.Context, client chat.ChatService_ExchangeMessageClient, data binding.ExternalStringList, user string) *fyne.Container {
	text := widget.NewMultiLineEntry()
	c := container.NewVBox(text, widget.NewButton("Send", func() {
		err := client.Send(&chat.Message{Author: user, Content: text.Text})
		if err != nil {
			log.Println("error while sending message : ", err.Error())
		}
		var mu sync.Mutex
		mu.Lock()
		data.Append(fmt.Sprintf("I send : %s", text.Text))
		mu.Unlock()
		text.SetText("")
	}))

	return c
}
