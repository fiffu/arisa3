package engine

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/lib/functional"

	dgo "github.com/bwmarrin/discordgo"
)

var (
	errNotCommand        = errors.New("not a command")
	errNoHandler         = errors.New("no handler")
	errDuplicatedRequest = errors.New("duplicated request")
)

type CommandsRegistry struct {
	cmds        map[string]types.ICommand
	clock       func() time.Time
	idempotency *idempotency
}

func NewCommandRegistry() *CommandsRegistry {
	cmds := make(map[string]types.ICommand)
	return &CommandsRegistry{cmds, time.Now, newIdempotencyChecker()}
}

// Register creates an ApplicationCommand with the given ICommands.
func (r *CommandsRegistry) Register(ctx context.Context, sess *dgo.Session, cmds ...types.ICommand) error {
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Internal("CommandsRegistry.Register"))
	defer span.End()

	for _, cmd := range cmds {

		appID := sess.State.User.ID
		data := cmd.Data()
		log.Infof(context.Background(), "Binding command /%s", cmd.Name())

		_, span := instrumentation.SpanInContext(ctx, instrumentation.Vendor(sess.ApplicationCommandCreate))
		_, err := sess.ApplicationCommandCreate(appID, "", data, dgo.WithContext(ctx))
		span.End()
		if err != nil {
			span.RecordError(err)
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
	ctx, err := r.registryHandler(s, i)
	if err == errDuplicatedRequest {
		log.Warnf(ctx, "Ignoring duplicated request")
		return
	}
	if err != nil {
		log.Errorf(ctx, err, "Error handling interaction")

		ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Vendor(s.InteractionRespond))
		defer span.End()

		if err := s.InteractionRespond(
			i.Interaction,
			types.NewResponse().Content("Hmm, seems like something went wrong. Try again later?").Data(),
			dgo.WithContext(ctx),
		); err != nil {
			log.Errorf(ctx, err, "Error sending response, maybe interaction already acknowledged?")
		}
	}
}

// registryHandler routes the InteractionCreate event to the appropriate command's handler.
func (r *CommandsRegistry) registryHandler(s *dgo.Session, i *dgo.InteractionCreate) (ctx context.Context, err error) {
	ctx = context.Background()
	startTime := r.clock()

	if i.Interaction.Data.Type() != dgo.InteractionApplicationCommand {
		err = errNotCommand
		return
	}

	commandName := i.ApplicationCommandData().Name
	cmd, ok := r.cmds[commandName]
	if !ok {
		return
	}

	id := i.ID
	if ok := r.idempotency.Check(id); !ok {
		err = errDuplicatedRequest
		return
	}

	// Code before this line executes for all commands; be careful to avoid excess logging.

	// Setup context for handler
	traceID := log.Hash(id)[:10]
	ctx = log.Put(ctx, log.TraceID, traceID)

	// Extract arguments and stuff
	who := i.User
	if who == nil && i.Member != nil {
		who = i.Member.User
	}
	ctx = log.Put(ctx, log.User, fmt.Sprintf("%s#%s:%s", who.Username, who.Discriminator, who.ID))

	if i.GuildID != "" {
		ctx = log.Put(ctx, log.Guild, i.GuildID)
	}

	opts := make(map[string]interface{})
	for _, o := range i.ApplicationCommandData().Options {
		opts[o.Name] = o.Value
	}
	log.Infof(ctx, "Interaction incoming <<< user=%s options=%+v", who, opts)

	// Instrumentation for the command handler
	ctx, span := instrumentation.SpanInContext(ctx, instrumentation.Command(cmd.Name()))
	span.SetAttributes(
		instrumentation.KV.CommandName(cmd.Name()),
		instrumentation.KV.TraceID(traceID),
		instrumentation.KV.User(who.String()),
		instrumentation.KV.Params(opts),
	)
	defer span.End()

	// Invoke handler
	handler := cmd.HandlerFunc()
	if handler == nil {
		return ctx, r.fallbackHandler(ctx, s, i, cmd)
	}
	args := parseArgs(ctx, cmd, i.ApplicationCommandData().Options)
	err = mustHandleCommand(ctx, cmd, handler, args, s, i)
	if err != nil {
		log.Errorf(ctx, err, "Handler errored")
	}

	endTime := r.clock()
	elapsed := endTime.Sub(startTime)
	log.Infof(ctx, "Interaction served in %d millisecs", elapsed.Milliseconds())

	return ctx, err
}

// fallbackHandler is invoked if a command has no associated handler.
func (r *CommandsRegistry) fallbackHandler(ctx context.Context, s *dgo.Session, i *dgo.InteractionCreate, cmd types.ICommand) error {
	log.Warnf(ctx, "No interaction handler registered for command: %s")
	return fmt.Errorf("%w: %s", errNoHandler, cmd.Name())
}

// parseArgs wraps user-supplied options in the InteractionCreate payload inside IArgs.
func parseArgs(ctx context.Context, cmd types.ICommand, args []*dgo.ApplicationCommandInteractionDataOption) types.IArgs {
	mapping := make(map[types.IOption]*dgo.ApplicationCommandInteractionDataOption)
	for _, arg := range args {
		if opt, ok := cmd.FindOption(arg.Name); ok {
			mapping[opt] = arg
		}
	}
	log.Infof(ctx, "Parsed options for command %s: %v", cmd.Name(), functional.Deref(args))
	return types.NewArgs(cmd, mapping)
}
