package service

import "github.com/alexandrevicenzi/go-sse"

var SseNotification *sse.Server

func InitSseServe(sse *sse.Server) {
	SseNotification = sse
}

func SendNotification(message string, userId string) {
	SseNotification.SendMessage("/notifications/"+userId, sse.SimpleMessage(message))
}
