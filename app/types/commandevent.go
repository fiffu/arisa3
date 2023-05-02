package types

//go:generate mockgen -source=commandevent.go -destination=./commandevent_mock.go -package=types

import (
	"context"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type ICommandEvent interface {
	Session() *dgo.Session
	Interaction() *dgo.InteractionCreate
	User() *dgo.User
	Command() ICommand
	Args() IArgs
	Respond(context.Context, ICommandResponse) error
}

// commandEvent implements ICommandEvent
type commandEvent struct {
	s    *dgo.Session
	i    *dgo.InteractionCreate
	cmd  ICommand
	args IArgs
}

func NewCommandEvent(s *dgo.Session, i *dgo.InteractionCreate, cmd ICommand, args IArgs) ICommandEvent {
	return &commandEvent{s, i, cmd, args}
}

func (evt *commandEvent) Session() *dgo.Session               { return evt.s }
func (evt *commandEvent) Interaction() *dgo.InteractionCreate { return evt.i }
func (evt *commandEvent) Command() ICommand                   { return evt.cmd }
func (evt *commandEvent) Args() IArgs                         { return evt.args }
func (evt *commandEvent) User() *dgo.User {
	user := evt.i.User
	if user == nil && evt.i.Member != nil {
		user = evt.i.Member.User
	}
	return user
}
func (evt *commandEvent) Respond(ctx context.Context, resp ICommandResponse) error {
	itr := evt.i.Interaction
	data := resp.Data()
	log.Info().
		Str(CtxCommand, evt.Command().Name()).
		Str(CtxInteraction, itr.ID).
		Msgf("Interaction response >>> resp: \n| %s", resp.String())
	return evt.s.InteractionRespond(itr, data)
}
