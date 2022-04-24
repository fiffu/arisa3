package colours

// interfaces.go describes the interfaces of the Colour Roles domain.
// The domain logic should operate on these primitives.

import (
	"time"
)

// IColoursDomain describes the colour roles domain
type IColoursDomain interface {
	// Get a member's last freeze time. If Never is returned, member's colour is not frozen.
	GetLastFrozen(IDomainMember) (time.Time, error)
	// Get a member's last mutate time.
	GetLastMutate(IDomainMember) (time.Time, bool, error)
	// Get a member's last reroll time.
	GetLastReroll(IDomainMember) (time.Time, bool, error)

	// Apply a mutation on member's colour role. If frozen, mutation is not allowed.
	Mutate(IDomainSession, IDomainMember) (*Colour, error)
	// Reroll the colour for a member's colour role.
	Reroll(IDomainSession, IDomainMember) (*Colour, error)

	// Freeze a member's colour role, i.e. disable mutations.
	Freeze(IDomainMember) error
	// Unfreeze a member's colour role, i.e. enable mutations.
	Unfreeze(IDomainMember) error

	// Returns whether member has colour role.
	HasColourRole(IDomainMember) bool
	// Search member's roles for one that looks like a colour role, based on name.
	GetColourRole(IDomainMember) IDomainRole
	// Generate a colour role name based on the member's nickname or username.
	CreateColourRole(IDomainSession, IDomainMember, *Colour) (IDomainRole, error)
	// Get height that colour roles should be at, based on position of role with maxHeightRoleName
	GetColourRoleHeight(IDomainSession, IDomainGuild) (int, error)
	// Set the height of a role
	SetRoleHeight(IDomainSession, IDomainGuild, string, int) error
	// Create and assign a colour role to a member.
	AssignColourRole(IDomainSession, IDomainMember, IDomainRole) error
}

// IDomainSession wraps methods of discordgo.Session that IColoursDomain will use.
type IDomainSession interface {
	GuildMember(guildID, userID string) (IDomainMember, error)
	GuildMemberRoleAdd(guildID, userID, roleID string) error
	GuildRole(guildID, roleID string) (IDomainRole, error)
	GuildRoles(guildID string) ([]IDomainRole, error)
	GuildRoleCreate(guildID string) (roleID string, err error)
	GuildRoleEdit(guildID, roleID, name string, color int) error
	GuildRoleReorder(guildID string, roles []IDomainRole) error
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
	CacheData() interface{}
	CacheDuration() time.Duration
}

// IDomainRole describes information that IColoursDomain derives from discordgo.Role.
// IDomainRole implements ICacheable.
type IDomainRole interface {
	ID() string
	Name() string
	Colour() *Colour

	CacheKey() string
	CacheData() interface{}
	CacheDuration() time.Duration
}

// IDomainRepository describes methods that IColoursDomain uses to fetch/store data.
type IDomainRepository interface {
	FetchUserState(IDomainMember, Reason) (time.Time, error)
	UpdateMutate(IDomainMember, *Colour) error
	UpdateReroll(IDomainMember, *Colour) error
	UpdateRerollPenalty(IDomainMember, time.Time) error
	UpdateFreeze(IDomainMember) error
	UpdateUnfreeze(IDomainMember) error
}
