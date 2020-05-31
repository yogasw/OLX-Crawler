package main

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"os"
	"time"
)

func sendMessage(remoteJid string, message string, imageUrl string, wac *whatsapp.Conn) {

	msgId, err := wac.Send(getMessage(remoteJid, message, imageUrl))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending message: %v", err)
		//os.Exit(1)
		<-time.After(5 * time.Second)
		wac.Send(getMessage(remoteJid, message, imageUrl))
	} else {
		fmt.Println("Message Sent -> ID : " + msgId)
	}
}

func getMessage(remoteJid string, message string, imageUrl string) whatsapp.ImageMessage {
	image := GetImageFromUrl(imageUrl)
	imagePath := SaveImage(image)
	imageRender := GetReaderFromImage(imagePath)
	thumb, err := getThumbnail(image)
	failOnError(err,"error get message")

	/*previousMessage := "😘"
	quotedMessage := proto.Message{
		Conversation: &previousMessage,
	}

	ContextInfo := whatsapp.ContextInfo{
		QuotedMessage:   &quotedMessage,
		QuotedMessageID: "",
		Participant:     "", //Whot sent the original message
	}

	textMessage := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: remoteJid,
		},
		ContextInfo: ContextInfo,
		Text:        message,
	}*/
	imageMessage := whatsapp.ImageMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: remoteJid,
		},
		Type:      "image/jpeg",
		Caption:   message,
		Content:   imageRender,
		Thumbnail: thumb,
	}

	return imageMessage
}
