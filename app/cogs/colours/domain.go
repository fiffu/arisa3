package colours

import (
	"errors"
	"time"

	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"
	"github.com/rs/zerolog/log"
)

// domain.go implements IColoursDomain defined by interfaces.go.

var (
	ErrCooldownPending   = errors.New("cooldown is still in progress")
	ErrInvalidRoleHeight = errors.New("invalid target role height, it should be >=0")
)

type domain struct {
	cog               types.ICog
	repo              IDomainRepository
	maxHeightRoleName string
	maxRoleHeight     int

	mutateCooldownMins time.Duration
	rerollCooldownMins time.Duration
	rerollPenaltyMins  time.Duration
}

// NewColoursDomain implements IColoursDomain
func NewColoursDomain(c types.ICog, repo IDomainRepository, cfg *Config) IColoursDomain {
	return &domain{
		cog:               c,
		repo:              repo,
		maxHeightRoleName: cfg.MaxRoleHeightName,
		maxRoleHeight:     -1,
	}
}

func (d *domain) GetLastFrozen(mem IDomainMember) (time.Time, error) {
	return d.repo.FetchUserState(mem, Freeze)
}

func (d *domain) GetLastMutate(mem IDomainMember) (time.Time, bool, error) {
	last, err := d.repo.FetchUserState(mem, Mutate)
	if err != nil {
		return last, false, err
	}
	cooldownPeriod := d.mutateCooldownMins * time.Minute
	return last, d.hasCooldownFinished(last, cooldownPeriod), nil
}

func (d *domain) GetLastReroll(mem IDomainMember) (time.Time, bool, error) {
	last, err := d.repo.FetchUserState(mem, Reroll)
	if err != nil {
		return last, false, err
	}
	cooldownPeriod := d.rerollCooldownMins * time.Minute
	return last, d.hasCooldownFinished(last, cooldownPeriod), nil
}

func (d *domain) hasCooldownFinished(cooldownStartTime time.Time, cooldownPeriod time.Duration) bool {
	if cooldownStartTime == Never {
		return true
	}
	cooldownEndTime := cooldownStartTime.Add(cooldownPeriod)
	finished := time.Now().After(cooldownEndTime)
	return finished
}

func (d *domain) Mutate(s IDomainSession, mem IDomainMember) (*Colour, error) {
	// No role, no mutate
	role := d.GetColourRole(mem)
	if role == nil {
		return nil, nil
	}

	// Check cooldown
	_, ok, err := d.GetLastMutate(mem)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrCooldownPending
	}

	// Generate new colour, apply cooldown
	newColour := role.Colour().Nudge()
	if err := d.repo.UpdateMutate(mem, newColour); err != nil {
		return newColour, err
	}

	// API call
	if err := s.GuildRoleEdit(
		mem.Guild().ID(),
		role.ID(),
		role.Name(),
		newColour.ToDecimal(),
	); err != nil {
		return newColour, err
	}
	return newColour, nil
}

func (d *domain) Reroll(s IDomainSession, mem IDomainMember) (*Colour, error) {
	// Check cooldown
	last, ok, err := d.GetLastReroll(mem)
	if err != nil {
		return nil, err
	}

	// Apply penalty if reroll not
	if !ok {
		// Skip DB call if no penalty configured
		if d.rerollPenaltyMins > 0*time.Minute {
			addedPenalty := last.Add(d.rerollPenaltyMins)
			if err := d.repo.UpdateRerollPenalty(mem, addedPenalty); err != nil {
				return nil, err
			}
		}
		return nil, ErrCooldownPending
	}

	// Generate new colour, apply cooldown
	newColour := (&Colour{}).Random()
	if err := d.repo.UpdateReroll(mem, newColour); err != nil {
		return newColour, err
	}

	// Edit existing role or assign a new role
	if d.HasColourRole(mem) {
		role := d.GetColourRole(mem)
		err = s.GuildRoleEdit(
			mem.Guild().ID(),
			role.ID(),
			role.Name(),
			newColour.ToDecimal(),
		)
		return newColour, err

	} else {
		role, err := d.CreateColourRole(s, mem, newColour)
		if err != nil {
			return newColour, err
		}
		err = d.AssignColourRole(s, mem, role)
		if err != nil {
			return newColour, err
		}
		return newColour, nil
	}
}

func (d *domain) Freeze(mem IDomainMember) error {
	return d.repo.UpdateFreeze(mem)
}

func (d *domain) Unfreeze(mem IDomainMember) error {
	return d.repo.UpdateUnfreeze(mem)
}

func (d *domain) HasColourRole(mem IDomainMember) bool {
	hasRole := d.GetColourRole(mem) != nil
	return hasRole
}

func (d *domain) GetColourRole(mem IDomainMember) IDomainRole {
	expectName := d.getColourRoleName(mem)
	who := mem.Username()
	for _, role := range mem.Roles() {
		if role.Name() == expectName {
			engine.CogLog(d.cog, log.Info()).Msgf("%s has colour role (%s)", who, expectName)
			return role
		}
	}
	engine.CogLog(d.cog, log.Info()).Msgf("%s doesn't have colour role (expected: %s)", who, expectName)
	return nil
}

func (d *domain) getColourRoleName(mem IDomainMember) string {
	roleName := mem.Username()
	return roleName
}

func (d *domain) CreateColourRole(s IDomainSession, mem IDomainMember, colour *Colour) (IDomainRole, error) {
	roleName := d.getColourRoleName(mem)
	guildID := mem.Guild().ID()

	// Create role
	id, err := s.GuildRoleCreate(guildID)
	if err != nil {
		return nil, err
	}

	// Set colour
	col := colour.ToDecimal()
	err = s.GuildRoleEdit(guildID, id, roleName, col)
	if err != nil {
		return nil, err
	}

	// Set height
	height, err := d.GetColourRoleHeight(s, mem.Guild())
	if err != nil {
		return nil, err
	}
	if height > -1 {
		err = d.SetRoleHeight(s, mem.Guild(), id, height)
		if err != nil {
			return nil, err
		}
	}
	return NewDomainRole(id, roleName, col), nil
}

func (d *domain) AssignColourRole(s IDomainSession, mem IDomainMember, role IDomainRole) error {
	return s.GuildMemberRoleAdd(mem.Guild().ID(), mem.UserID(), role.ID())
}

func (d *domain) GetColourRoleHeight(s IDomainSession, guild IDomainGuild) (int, error) {
	if d.maxRoleHeight > -1 {
		return d.maxRoleHeight, nil
	}
	roles, err := s.GuildRoles(guild.ID())
	if err != nil {
		return -1, err
	}

	engine.CogLog(d.cog, log.Debug()).Msgf("Checking height of role: %s", d.maxHeightRoleName)
	for i, role := range roles {
		if role.Name() == d.maxHeightRoleName {
			d.maxRoleHeight = i
			engine.CogLog(d.cog, log.Info()).Msgf(
				"Found height of role: %s (= %d)", d.maxHeightRoleName, i,
			)
			return i, nil
		}
	}
	return -1, nil
}

func (d *domain) SetRoleHeight(s IDomainSession, g IDomainGuild, roleID string, height int) error {
	if height <= -1 {
		return ErrInvalidRoleHeight
	}
	guildID := g.ID()
	allRoles, err := s.GuildRoles(guildID)
	if err != nil {
		return err
	}

	var theRole IDomainRole
	found := false
	for idx, role := range allRoles {
		if role.ID() == roleID {
			allRoles = append(allRoles[:idx], allRoles[idx+1:]...)
			found = true
			break
		}
	}
	if !found {
		return nil
	}
	payload := append(allRoles[:height], theRole)
	payload = append(payload, allRoles[height:]...)
	return s.GuildRoleReorder(guildID, payload)
}
