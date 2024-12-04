package colours

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/lib/functional"
)

// domain.go implements IColoursDomain defined by interfaces.go.

var (
	ErrMutateFrozen          = errors.New("colour is frozen")
	ErrMutateCooldownPending = errors.New("mutate cooldown is still in progress")
	ErrRerollCooldownPending = errors.New("reroll cooldown is still in progress")
	ErrInvalidRoleHeight     = errors.New("invalid target role height, it should be >=0")

	rolePattern = regexp.MustCompile(`\w+#(0|\d{4})`)
)

type domain struct {
	now               func() time.Time
	cog               types.ICog
	repo              IDomainRepository
	maxHeightRoleName string
	maxRoleHeight     int

	mutateCooldownMins int
	rerollCooldownMins int
	rerollPenaltyMins  int
}

// NewColoursDomain implements IColoursDomain
func NewColoursDomain(c types.ICog, repo IDomainRepository, cfg *Config) IColoursDomain {
	return &domain{
		now:               time.Now,
		cog:               c,
		repo:              repo,
		maxHeightRoleName: cfg.MaxRoleHeightName,
		maxRoleHeight:     -1,

		mutateCooldownMins: cfg.MutateCooldownMins,
		rerollCooldownMins: cfg.RerollCooldownMins,
		rerollPenaltyMins:  cfg.RerollPenaltyMins,
	}
}

func (d *domain) GetLastFrozen(ctx context.Context, mem IDomainMember) (time.Time, error) {
	return d.repo.FetchUserState(ctx, mem, Freeze)
}

func (d *domain) GetLastMutate(ctx context.Context, mem IDomainMember) (time.Time, bool, error) {
	last, err := d.repo.FetchUserState(ctx, mem, Mutate)
	if err != nil {
		return last, false, err
	}
	cooldownPeriod := time.Duration(d.mutateCooldownMins) * time.Minute
	return last, d.hasCooldownFinished(last, cooldownPeriod), nil
}

func (d *domain) GetLastReroll(ctx context.Context, mem IDomainMember) (time.Time, bool, error) {
	last, err := d.repo.FetchUserState(ctx, mem, Reroll)
	if err != nil {
		return last, false, err
	}
	cooldownPeriod := time.Duration(d.rerollCooldownMins) * time.Minute
	return last, d.hasCooldownFinished(last, cooldownPeriod), nil
}

func (d *domain) hasCooldownFinished(cooldownStartTime time.Time, cooldownPeriod time.Duration) bool {
	if cooldownStartTime == Never {
		return true
	}
	if cooldownPeriod == 0 {
		return true
	}
	cooldownEndTime := d.offsetTime(cooldownStartTime, cooldownPeriod)
	finished := time.Now().After(cooldownEndTime)
	return finished
}

func (d *domain) offsetTime(startTime time.Time, cooldownPeriod time.Duration) time.Time {
	return startTime.Add(cooldownPeriod)
}

func (d *domain) GetRerollCooldownEndTime(ctx context.Context, mem IDomainMember) (time.Time, error) {
	last, _, err := d.GetLastReroll(ctx, mem)
	if err != nil {
		return time.Time{}, err
	}
	cooldownPeriod := time.Duration(d.rerollCooldownMins) * time.Minute
	endTime := d.offsetTime(last, cooldownPeriod)
	return endTime, nil
}

func (d *domain) GetHistory(ctx context.Context, mem IDomainMember) (*History, error) {
	endTime := d.now()
	startTime := endTime.Add(-14 * 24 * time.Hour)

	logs, err := d.repo.FetchUserHistory(ctx, mem, startTime)
	if err != nil {
		return nil, err
	}

	cols, err := functional.Slice[*ColoursLogRecord, *Colour](logs).
		Filter(func(c *ColoursLogRecord) bool {
			return isPartOfHistory(Reason(c.Reason))
		}), nil
	return &History{cols, startTime, endTime}, err
}

func (d *domain) Mutate(ctx context.Context, s IDomainSession, mem IDomainMember) (*Colour, error) {
	// No role, no mutate
	role := d.GetColourRole(ctx, mem)
	if role == nil {
		return nil, nil
	}

	// Check frozen
	if lastFrozen, err := d.GetLastFrozen(ctx, mem); err != nil {
		return nil, err
	} else if lastFrozen != Never {
		return nil, ErrMutateFrozen
	}

	// Check cooldown
	if _, isCooldownDone, err := d.GetLastMutate(ctx, mem); err != nil {
		return nil, err
	} else if !isCooldownDone {
		return nil, ErrMutateCooldownPending
	}

	// Generate new colour, apply cooldown
	newColour := role.Colour().Nudge()
	if err := d.repo.UpdateMutate(ctx, mem, newColour); err != nil {
		return newColour, err
	}

	// API call
	if err := s.GuildRoleEdit(
		ctx,
		mem.Guild().ID(),
		role.ID(),
		role.Name(),
		newColour.ToDecimal(),
	); err != nil {
		return newColour, err
	}
	return newColour, nil
}

func (d *domain) Reroll(ctx context.Context, s IDomainSession, mem IDomainMember) (*Colour, error) {
	// Check cooldown
	last, cooldownFinished, err := d.GetLastReroll(ctx, mem)
	log.Infof(ctx,
		"%s last roll was %s, %d mins cooldown finished? %v",
		mem.Username(), last.Format(time.RFC3339), d.rerollCooldownMins, cooldownFinished,
	)
	if err != nil {
		return nil, err
	}

	// Apply penalty if reroll cooldown not finished
	if !cooldownFinished {
		// Skip DB call if no penalty configured
		log.Infof(ctx, "Applying %v mins penalty on %s", d.rerollPenaltyMins, mem.Username())
		if d.rerollPenaltyMins > 0 {
			addedPenalty := last.Add(time.Duration(d.rerollPenaltyMins) * time.Minute)
			if err := d.repo.UpdateRerollPenalty(ctx, mem, addedPenalty); err != nil {
				return nil, err
			}
		}
		return nil, ErrRerollCooldownPending
	}

	// Generate new colour, apply cooldown
	newColour := (&Colour{}).Random()
	if err := d.repo.UpdateReroll(ctx, mem, newColour); err != nil {
		return newColour, err
	}

	// Edit existing role or assign a new role
	if d.HasColourRole(ctx, mem) {
		role := d.GetColourRole(ctx, mem)
		err = s.GuildRoleEdit(
			ctx,
			mem.Guild().ID(),
			role.ID(),
			role.Name(),
			newColour.ToDecimal(),
		)
		return newColour, err

	} else {
		role, err := d.CreateColourRole(ctx, s, mem, newColour)
		if err != nil {
			return newColour, err
		}
		err = d.AssignColourRole(ctx, s, mem, role)
		if err != nil {
			return newColour, err
		}
		return newColour, nil
	}
}

func (d *domain) Freeze(ctx context.Context, mem IDomainMember) error {
	return d.repo.UpdateFreeze(ctx, mem)
}

func (d *domain) Unfreeze(ctx context.Context, mem IDomainMember) error {
	return d.repo.UpdateUnfreeze(ctx, mem)
}

func (d *domain) HasColourRole(ctx context.Context, mem IDomainMember) bool {
	hasRole := d.GetColourRole(ctx, mem) != nil
	return hasRole
}

func (d *domain) GetColourRole(ctx context.Context, mem IDomainMember) IDomainRole {
	for _, role := range mem.Roles() {
		roleName := role.Name()
		if rolePattern.MatchString(roleName) {
			return role
		}
	}
	return nil
}

func (d *domain) GetColourRoleName(ctx context.Context, mem IDomainMember) string {
	roleName := mem.Username()
	return roleName
}

func (d *domain) CreateColourRole(ctx context.Context, s IDomainSession, mem IDomainMember, colour *Colour) (IDomainRole, error) {
	roleName := d.GetColourRoleName(ctx, mem)
	guildID := mem.Guild().ID()

	// Create role
	col := colour.ToDecimal()
	id, err := s.GuildRoleCreate(ctx, guildID, roleName, col)
	if err != nil {
		return nil, err
	}

	// Set height
	height, err := d.GetColourRoleHeight(ctx, s, mem.Guild())
	if err != nil {
		return nil, err
	}
	if height > -1 {
		err = d.SetRoleHeight(ctx, s, mem.Guild(), id, height)
		if err != nil {
			return nil, err
		}
	}
	return NewDomainRole(id, roleName, col), nil
}

func (d *domain) AssignColourRole(ctx context.Context, s IDomainSession, mem IDomainMember, role IDomainRole) error {
	return s.GuildMemberRoleAdd(ctx, mem.Guild().ID(), mem.UserID(), role.ID())
}

func (d *domain) GetColourRoleHeight(ctx context.Context, s IDomainSession, guild IDomainGuild) (int, error) {
	if d.maxRoleHeight > -1 {
		return d.maxRoleHeight, nil
	}
	roles, err := s.GuildRoles(ctx, guild.ID())
	if err != nil {
		return -1, err
	}

	log.Debugf(ctx, "Checking height of role: %s", d.maxHeightRoleName)
	for i, role := range roles {
		if role.Name() == d.maxHeightRoleName {
			d.maxRoleHeight = i
			log.Debugf(ctx, "Found height of role: %s (= %d)", d.maxHeightRoleName, i)
			return i, nil
		}
	}
	return -1, nil
}

func (d *domain) SetRoleHeight(ctx context.Context, s IDomainSession, g IDomainGuild, newRoleID string, height int) error {
	if height <= -1 {
		return ErrInvalidRoleHeight
	}
	guildID := g.ID()
	allRoles, err := s.GuildRoles(ctx, guildID)
	if err != nil {
		return err
	}

	numRoles := len(allRoles)
	if numRoles == 0 {
		return nil
	}
	if height < 0 {
		return nil
	}
	if height >= numRoles {
		height = numRoles - 1
	}

	var theRole IDomainRole
	found := false
	for idx, role := range allRoles {
		if role.ID() == newRoleID {
			theRole = role
			allRoles = append(allRoles[:idx], allRoles[idx+1:]...)
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	payload := make([]IDomainRole, 0)
	payload = append(payload, allRoles[:height]...)
	payload = append(payload, theRole)
	payload = append(payload, allRoles[height:]...)
	return s.GuildRoleReorder(ctx, guildID, payload)
}
