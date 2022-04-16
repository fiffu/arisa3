package engine

import (
	"arisa3/app/types"
	"errors"
	"fmt"
	"runtime/debug"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var (
	errNotCommand              = errors.New("not a command")
	errNoHandler               = errors.New("no handler")
	errUnrecognisedInteraction = errors.New("unrecognised interaction name")
	errPanic                   = errors.New("panic while executing command handler")
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
		if _, err := s.ApplicationCommandCreate(appID, "", data); err != nil {
			return err
		}
		registryLog(log.Info()).Msgf("Binding command '%s' -> %+v", cmd.Name(), cmd)
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
	registryLog(log.Info()).Msgf("Incoming interaction from '%s': %+v", i.User, i.Data)
	if err := r.registryHandler(s, i); err != nil {
		registryLog(log.Error()).Err(err).Msgf("Error handling interaction")
		err = s.InteractionRespond(
			i.Interaction,
			types.NewResponse().Content("Hmm, seems like something went wrong. Try again later?").Data(),
		)
		if err != nil {
			registryLog(log.Error()).Err(err).Msgf("Error sending response")
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
		err = fmt.Errorf("%w: %s", errUnrecognisedInteraction, commandName)
		return
	}
	handler := cmd.GetHandler()
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
	registryLog(log.Debug()).Msgf("Handler executing")
	err = handler(s, i, cmd, args)
	if err == nil {
		registryLog(log.Debug()).Msgf("Handler completed")
	} else {
		registryLog(log.Error()).Err(err).Msgf("Handler errored")
	}
	return
}

// fallbackHandler is invoked in lieu of runHandler if a command has no associated handler.
func (r *CommandsRegistry) fallbackHandler(s *dgo.Session, i *dgo.InteractionCreate, cmd types.ICommand) error {
	log.Error().Str("registry::command", cmd.Name()).Msgf("Missing interaction handler")
	return fmt.Errorf("%w: %s", errNoHandler, cmd.Name())
}

// parseArgs wraps user-supplied options in the InteractionCreate payload inside IArgs.
func parseArgs(cmd types.ICommand, opts []*dgo.ApplicationCommandInteractionDataOption) types.IArgs {
	args := make(map[string]*dgo.ApplicationCommandInteractionDataOption)
	for _, opt := range opts {
		args[opt.Name] = opt
	}
	log.Debug().Str("registry::command", cmd.Name()).Msgf("Parsed options: %v", args)
	return types.NewArgs(args)
}
