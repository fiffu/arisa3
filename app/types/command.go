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
	HandlerFunc() Handler
	FindOption(string) (IOption, bool)
}

type ICommandEvent interface {
	Session() *dgo.Session
	Interaction() *dgo.InteractionCreate
	Command() ICommand
	Args() IArgs
	Respond(ICommandResponse) error
}

type Handler func(ICommandEvent) error

type Command struct {
	name    string
	data    *dgo.ApplicationCommand
	opts    map[string]IOption
	handler Handler
}

func NewCommand(name string) *Command {
	data := dgo.ApplicationCommand{Name: name}
	opts := make(map[string]IOption)
	return &Command{
		name: name,
		data: &data,
		opts: opts,
	}
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

// HandlerFunc returns the callback assigned to this command.
func (c *Command) HandlerFunc() Handler { return c.handler }

// Option defines options accepted by this command.
func (c *Command) Options(opts ...IOption) *Command {
	for _, opt := range opts {
		c.opts[opt.Name()] = opt
		c.data.Options = append(c.data.Options, opt.Data())
	}
	return c
}

func (c *Command) FindOption(name string) (opt IOption, ok bool) {
	opt, ok = c.opts[name]
	return
}
