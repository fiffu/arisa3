package types

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	dgo "github.com/bwmarrin/discordgo"
)

var (
	// Application command naming requirements
	// This is a Golang adaptation of the given Regexp pattern from the docs:
	//      /^[-_\p{L}\p{N}\p{sc=Deva}\p{sc=Thai}]{1,32}$/u
	// Ref: https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-naming
	commandNamePattern = regexp.MustCompile(`^[a-z-_\p{L}\p{N}\p{Devanagari}\p{Thai}]{1,32}$`)
)

type ICommand interface {
	Name() string
	Data() *dgo.ApplicationCommand
	ForChat() *Command
	ForUser() *Command
	ForMessage() *Command
	Desc(string) *Command
	Handler(hdlr CommandHandler) *Command
	HandlerFunc() CommandHandler
	FindOption(string) (IOption, bool)
}

type CommandHandler func(context.Context, ICommandEvent) error

type Command struct {
	name    string
	data    *dgo.ApplicationCommand
	opts    map[string]IOption
	handler CommandHandler
}

func NewCommand(name string) *Command {
	data := dgo.ApplicationCommand{Name: name}
	opts := make(map[string]IOption)
	cmd := &Command{
		name: name,
		data: &data,
		opts: opts,
	}
	cmd.mustValidate()
	return cmd
}

func (c *Command) mustValidate() {
	if !commandNamePattern.MatchString(c.name) {
		msg := fmt.Sprintf("invalid command name (regexp mismatch), got: %s, expected: %s", c.name, commandNamePattern.String())
		panic(msg)
	}
	if c.name != strings.ToLower(c.name) {
		msg := fmt.Sprintf("invalid command name (must use lowercase), got: %s, expected: %s", c.name, strings.ToLower(c.name))
		panic(msg)
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
func (c *Command) Handler(hdlr CommandHandler) *Command { c.handler = hdlr; return c }

// HandlerFunc returns the callback assigned to this command.
func (c *Command) HandlerFunc() CommandHandler { return c.handler }

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
