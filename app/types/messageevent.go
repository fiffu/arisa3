package types

//go:generate mockgen -source=messageevent.go -destination=./messageevent_mock.go -package=types

import (
	dgo "github.com/bwmarrin/discordgo"
)

type IMessageEvent interface {
	Event() IEvent
	Message() *dgo.Message
	GuildID() string
	User() *dgo.User
}

func NewMessageEvent(sess *dgo.Session, source *dgo.MessageCreate) IMessageEvent {
	return &msgEvent{
		source: source,
		event:  NewEvent(sess, MessageCreateEvent),
	}
}

type msgEvent struct {
	source *dgo.MessageCreate
	event  IEvent
}

func (m *msgEvent) Event() IEvent {
	return m.event
}

func (m *msgEvent) Message() *dgo.Message {
	return m.source.Message
}

func (m *msgEvent) User() *dgo.User {
	return m.source.Author
}

func (m *msgEvent) GuildID() string {
	return m.source.GuildID
}
