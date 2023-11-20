package colours

// interfaces.go describes the interfaces of the Colour Roles domain.
// The domain logic should operate on these primitives.

//go:generate mockgen -source=interfaces.go -destination=./interfaces_mock.go -package=colours

import (
	"context"
	"time"
)

// IColoursDomain describes the colour roles domain
type IColoursDomain interface {
	// Get a member's last freeze time. If Never is returned, member's colour is not frozen.
	GetLastFrozen(context.Context, IDomainMember) (time.Time, error)
	// Get a member's last mutate time.
	GetLastMutate(context.Context, IDomainMember) (time.Time, bool, error)
	// Get a member's last reroll time.
	GetLastReroll(context.Context, IDomainMember) (time.Time, bool, error)

	// Get the reroll cooldown end time for member.
	GetRerollCooldownEndTime(context.Context, IDomainMember) (time.Time, error)
	// Get history of colours associated with member.
	GetHistory(context.Context, IDomainMember) (*History, error)

	// Apply a mutation on member's colour role. If frozen, mutation is not allowed.
	Mutate(context.Context, IDomainSession, IDomainMember) (*Colour, error)
	// Reroll the colour for a member's colour role.
	Reroll(context.Context, IDomainSession, IDomainMember) (*Colour, error)

	// Freeze a member's colour role, i.e. disable mutations.
	Freeze(context.Context, IDomainMember) error
	// Unfreeze a member's colour role, i.e. enable mutations.
	Unfreeze(context.Context, IDomainMember) error

	// Returns whether member has colour role.
	HasColourRole(context.Context, IDomainMember) bool
	// Search member's roles for one that looks like a colour role, based on name.
	GetColourRole(context.Context, IDomainMember) IDomainRole
	// Derive a role name based on the member's properties, like username.
	GetColourRoleName(context.Context, IDomainMember) string
	// Generate a colour role name based on the member's nickname or username.
	CreateColourRole(context.Context, IDomainSession, IDomainMember, *Colour) (IDomainRole, error)
	// Get height that colour roles should be at, based on position of role with maxHeightRoleName
	GetColourRoleHeight(context.Context, IDomainSession, IDomainGuild) (int, error)
	// Set the height of a role
	SetRoleHeight(context.Context, IDomainSession, IDomainGuild, string, int) error
	// Create and assign a colour role to a member.
	AssignColourRole(context.Context, IDomainSession, IDomainMember, IDomainRole) error
}

// IDomainSession wraps methods of discordgo.Session that IColoursDomain will use.
type IDomainSession interface {
	GuildMember(ctx context.Context, guildID, userID string) (IDomainMember, error)
	GuildMemberRoleAdd(ctx context.Context, guildID, userID, roleID string) error
	GuildRole(ctx context.Context, guildID, roleID string) (IDomainRole, error)
	GuildRoles(ctx context.Context, guildID string) ([]IDomainRole, error)
	GuildRoleCreate(ctx context.Context, guildID string, name string, colour int) (roleID string, err error)
	GuildRoleEdit(ctx context.Context, guildID, roleID, name string, colour int) error
	GuildRoleReorder(ctx context.Context, guildID string, roles []IDomainRole) error
}

// IDomainGuild describes information that IColoursDomain derives from discordgo.Guild.
type IDomainGuild interface {
	ID() string
}

// IDomainMember describes information that IColoursDomain derives from discordgo.Member.
// IDomainMember implements ICacheable.
type IDomainMember interface {
	Guild() IDomainGuild
	UserID() string
	Username() string
	Nick() string
	Roles() []IDomainRole

	CacheKey() string
}

// IDomainRole describes information that IColoursDomain derives from discordgo.Role.
// IDomainRole implements ICacheable.
type IDomainRole interface {
	ID() string
	Name() string
	Colour() *Colour

	CacheKey() string
}

// IDomainRepository describes methods that IColoursDomain uses to fetch/store data.
type IDomainRepository interface {
	FetchUserState(context.Context, IDomainMember, Reason) (time.Time, error)
	FetchUserHistory(context.Context, IDomainMember, time.Time) ([]*ColoursLogRecord, error)
	UpdateMutate(context.Context, IDomainMember, *Colour) error
	UpdateReroll(context.Context, IDomainMember, *Colour) error
	UpdateRerollPenalty(context.Context, IDomainMember, time.Time) error
	UpdateFreeze(context.Context, IDomainMember) error
	UpdateUnfreeze(context.Context, IDomainMember) error
}
