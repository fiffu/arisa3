package colours

// interfaces.go describes the interfaces of the Colour Roles domain.
// The domain logic should operate on these primitives.

import (
	"time"
)

// IColoursDomain describes the colour roles domain
type IColoursDomain interface {
	// Get a member's last mutate time.
	GetLastMutate(IDomainMember) (time.Time, error)
	// Get a member's last reroll time.
	GetLastReroll(IDomainMember) (time.Time, error)
	// Get a member's last freeze time. If Never is returned, member's colour is not frozen.
	GetLastFrozen(IDomainMember) (time.Time, error)

	// Check whether Mutate cooldown expired for the member.
	CanMutate(IDomainMember) (bool, error)
	// Apply a mutation on member's colour role. If frozen, mutation is not allowed.
	Mutate(IDomainMember) (*Colour, error)

	// Check whether Reroll cooldown expired for the member.
	CanReroll(IDomainMember) (bool, error)
	// Reroll the colour for a member's colour role.
	Reroll(IDomainMember) (*Colour, error)

	// Freeze a member's colour role, i.e. disable mutations.
	Freeze(IDomainMember) error
	// Unfreeze a member's colour role, i.e. enable mutations.
	Unfreeze(IDomainMember) error

	// Search member's roles for one that looks like a Colour Role, based on name.
	GetColourRole(IDomainMember) IDomainRole
	// Generate a colour role name based on the member's nickname or username.
	CreateColourRole(IDomainSession, IDomainMember, *Colour) (IDomainRole, error)
	// Create and assign a colour role to a member.
	AssignColourRole(IDomainSession, IDomainMember, IDomainRole) error
	// What height should
	GetColourRoleHeight(s IDomainSession, guildID string) (int, error)
}

// IDomainSession wraps methods of discordgo.Session that IColoursDomain will use.
type IDomainSession interface {
	GuildMember(guildID, userID string) (IDomainMember, error)
	GuildMemberRoleAdd(guildID, userID, roleID string) error
	GuildRole(guildID, roleID string) (IDomainRole, error)
	GuildRoles(guildID string) ([]IDomainRole, error)
	GuildRoleCreate(guildID string) (roleID string, err error)
	GuildRoleEdit(guildID, roleID, name string, color int) error
	GuildRoleReorder(guildID, roleID string, height int) error
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
	CacheDuration() time.Duration
}

// IDomainRole describes information that IColoursDomain derives from discordgo.Role.
// IDomainRole implements ICacheable.
type IDomainRole interface {
	ID() string
	Name() string
	Colour() *Colour

	CacheKey() string
	CacheDuration() time.Duration
}

// IDomainRepository describes methods that IColoursDomain uses to fetch/store data.
type IDomainRepository interface {
	FetchUserState(IDomainMember, Reason) (time.Time, error)
	UpdateMutate(IDomainMember, *Colour) error
	UpdateReroll(IDomainMember, *Colour) error
	UpdateFreeze(IDomainMember) error
	UpdateUnfreeze(IDomainMember) error
}
