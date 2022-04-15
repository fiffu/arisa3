package types

import (
	dgo "github.com/bwmarrin/discordgo"
)

type ICommand interface {
	Name() string
	Data() *dgo.ApplicationCommand
	ForChat() *Command
	ForUser() *Command
	ForMessage() *Command
	Desc(string) *Command
	Handler(hdlr Handler) *Command
	GetHandler() Handler
}

type Handler func(*dgo.Session, *dgo.InteractionCreate, ICommand, IArgs) error

type Command struct {
	name    string
	data    *dgo.ApplicationCommand
	handler Handler
}

func NewCommand(name string) *Command {
	d := dgo.ApplicationCommand{Name: name}
	return &Command{name: name, data: &d}
}

// Name is command name
func (c *Command) Name() string { return c.name }

// Data returns the underlying command definition.
func (c *Command) Data() *dgo.ApplicationCommand { return c.data }

// Desc sets this command description.
func (c *Command) Desc(description string) *Command { c.data.Description = description; return c }

// ForChat sets command type to Chat, the default command type.
// These are slash commands (i.e. called directly from the chat).
func (c *Command) ForChat() *Command { c.data.Type = dgo.ChatApplicationCommand; return c }

// ForUser sets command type to User, adds command to user context menu.
func (c *Command) ForUser() *Command { c.data.Type = dgo.UserApplicationCommand; return c }

// ForMessage sets command type to Message, adds command to message context menu.
func (c *Command) ForMessage() *Command { c.data.Type = dgo.MessageApplicationCommand; return c }

// Handler assigns a callback to this command.
func (c *Command) Handler(hdlr Handler) *Command { c.handler = hdlr; return c }

// GetHandler returns the callback assigned to this command.
func (c *Command) GetHandler() Handler { return c.handler }

// Option defines options accepted by this command.
func (c *Command) Options(opts ...IOption) *Command {
	for _, opt := range opts {
		c.data.Options = append(c.data.Options, opt.Data())
	}
	return c
}
