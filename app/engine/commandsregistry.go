package engine

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/fiffu/arisa3/app/types"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var (
	errNotCommand = errors.New("not a command")
	errNoHandler  = errors.New("no handler")
	errPanic      = errors.New("panic while executing command handler")
)

type CommandsRegistry struct {
	cmds map[string]types.ICommand
}

func NewCommandRegistry() *CommandsRegistry {
	cmds := make(map[string]types.ICommand)
	return &CommandsRegistry{cmds}
}

// Register creates an ApplicationCommand with the given ICommands.
func (r *CommandsRegistry) Register(s *dgo.Session, cmds ...types.ICommand) error {
	for _, cmd := range cmds {
		appID := s.State.User.ID
		data := cmd.Data()
		registryLog(log.Info()).Msgf("Binding command /%s", cmd.Name())
		if _, err := s.ApplicationCommandCreate(appID, "", data); err != nil {
			return err
		}
		r.cmds[cmd.Name()] = cmd
	}
	return nil
}

// BindCallbacks binds InteractionCreate event to the registry's onInteractionCreate handler.
func (r *CommandsRegistry) BindCallbacks(s *dgo.Session) {
	s.AddHandler(func(sess *dgo.Session, i *dgo.InteractionCreate) {
		r.onInteractionCreate(sess, i)
	})
}

// onInteractionCreate logs errors from registryHandler.
func (r *CommandsRegistry) onInteractionCreate(s *dgo.Session, i *dgo.InteractionCreate) {
	if err := r.registryHandler(s, i); err != nil {
		registryLog(log.Error()).Err(err).Msgf("Error handling interaction")
		err = s.InteractionRespond(
			i.Interaction,
			types.NewResponse().Content("Hmm, seems like something went wrong. Try again later?").Data(),
		)
		if err != nil {
			registryLog(log.Error()).Err(err).Msgf("Error sending response, maybe interaction already acknowledged?")
		}
	}
}

// registryHandler routes the InteractionCreate event to the appropriate command's handler.
func (r *CommandsRegistry) registryHandler(s *dgo.Session, i *dgo.InteractionCreate) (err error) {
	if i.Interaction.Data.Type() != dgo.InteractionApplicationCommand {
		err = errNotCommand
		return
	}
	commandName := i.ApplicationCommandData().Name
	cmd, ok := r.cmds[commandName]
	if !ok {
		return
	}

	// Code before this line executes for all commands; be careful to avoid excess logging.

	// Logging
	who := i.User
	if who == nil && i.Member != nil {
		who = i.Member.User
	}
	opts := make(map[string]interface{})
	for _, o := range i.ApplicationCommandData().Options {
		opts[o.Name] = o.Value
	}
	registryLog(log.Info()).
		Str(types.CtxCommand, i.ApplicationCommandData().Name).
		Str(types.CtxInteraction, i.ID).
		Msgf(
			"Interaction incoming <<< user=%s options=%+v",
			who, opts,
		)

	// Invoke handler
	handler := cmd.HandlerFunc()
	if handler == nil {
		return r.fallbackHandler(s, i, cmd)
	}
	args := parseArgs(cmd, i.ApplicationCommandData().Options)
	return r.runHandler(s, i, cmd, handler, args)
}

// runHandler executes a command's handler, trapping and logging any panics/errors.
func (r *CommandsRegistry) runHandler(
	s *dgo.Session, i *dgo.InteractionCreate,
	cmd types.ICommand, handler types.Handler, args types.IArgs) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errPanic
			fmt.Printf("%+v:\n%s\n", r, string(debug.Stack()))
		}
	}()
	registryLog(log.Debug()).Str(types.CtxCommand, cmd.Name()).Msgf("Handler executing")
	err = handler(types.NewCommandEvent(s, i, cmd, args))
	if err == nil {
		registryLog(log.Debug()).Msgf("Handler completed")
	} else {
		registryLog(log.Error()).Err(err).Msgf("Handler errored")
	}
	return
}

// fallbackHandler is invoked in lieu of runHandler if a command has no associated handler.
func (r *CommandsRegistry) fallbackHandler(s *dgo.Session, i *dgo.InteractionCreate, cmd types.ICommand) error {
	registryLog(log.Error()).Str(types.CtxCommand, cmd.Name()).Msgf("Missing interaction handler")
	return fmt.Errorf("%w: %s", errNoHandler, cmd.Name())
}

// parseArgs wraps user-supplied options in the InteractionCreate payload inside IArgs.
func parseArgs(cmd types.ICommand, args []*dgo.ApplicationCommandInteractionDataOption) types.IArgs {
	mapping := make(map[types.IOption]*dgo.ApplicationCommandInteractionDataOption)
	for _, arg := range args {
		if opt, ok := cmd.FindOption(arg.Name); ok {
			mapping[opt] = arg
		}
	}
	registryLog(log.Info()).Str(types.CtxCommand, cmd.Name()).Msgf("Parsed options: %v", args)
	return types.NewArgs(cmd, mapping)
}
