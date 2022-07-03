package types

import dgo "github.com/bwmarrin/discordgo"

const (
	MessageCreateEvent = "MessageCreateEvent"
)

type IEvent interface {
	Name() string
	Session() *dgo.Session
}

type event struct {
	name string
	sess *dgo.Session
}

func NewEvent(sess *dgo.Session, name string) IEvent {
	return &event{name, sess}
}

func (e *event) Name() string {
	return e.name
}

func (e *event) Session() *dgo.Session {
	return e.sess
}
