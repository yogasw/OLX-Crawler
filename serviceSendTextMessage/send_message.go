package main

import (
	"fmt"
	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	"os"
	"time"
)

func sendMessage(remoteJid string, message string, imageUrl string, wac *whatsapp.Conn) {
	if imageUrl == "" {
		msgId, err := wac.Send(getTextMessage(remoteJid, message))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error sending message: %v", err)
			//os.Exit(1)
			<-time.After(5 * time.Second)
			//wac.Send(getImageMessage(remoteJid, message, imageUrl))
		} else {
			fmt.Println("Message Text Sent -> ID : " + msgId)
		}
	} else {
		msgId, err := wac.Send(getImageMessage(remoteJid, message, imageUrl))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error sending message: %v", err)
			//os.Exit(1)
			<-time.After(5 * time.Second)
			//wac.Send(getImageMessage(remoteJid, message, imageUrl))
		} else {
			fmt.Println("Message Image Sent -> ID : " + msgId)
		}
	}
}

func getImageMessage(remoteJid string, message string, imageUrl string) whatsapp.ImageMessage {
	image := GetImageFromUrl(imageUrl)
	imagePath := SaveImage(image)
	imageRender := GetReaderFromImage(imagePath)
	thumb, err := getThumbnail(image)
	failOnError(err, "error get message")
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

func getTextMessage(remoteJid string, message string) whatsapp.TextMessage {
	previousMessage := "😘"
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
	}
	return textMessage
}
