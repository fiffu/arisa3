package colours

import (
	"time"
)

type domain struct {
	repo              IDomainRepository
	maxHeightRoleName string
	maxRoleHeight     int
}

// NewColoursDomain implements IColoursDomain
func NewColoursDomain(repo IDomainRepository, cfg Config) IColoursDomain {
	return &domain{
		repo:              repo,
		maxHeightRoleName: cfg.MaxRoleHeightName,
		maxRoleHeight:     -1,
	}
}

func (d *domain) GetLastMutate(mem IDomainMember) (time.Time, error) {
	return d.repo.FetchUserState(mem, Mutate)
}

func (d *domain) GetLastReroll(mem IDomainMember) (time.Time, error) {
	return d.repo.FetchUserState(mem, Reroll)
}

func (d *domain) GetLastFrozen(mem IDomainMember) (time.Time, error) {
	return d.repo.FetchUserState(mem, Freeze)
}

func (d *domain) CanMutate(m IDomainMember) (bool, error) {
	// TODO
	return false, nil
}

func (d *domain) Mutate(m IDomainMember) (*Colour, error) {
	// TODO
	// if role := m.ColourRole(); role != nil {
	// 	current := role.Colour()
	// 	newCol := current.Nudge()
	// }
	return nil, nil
}

func (d *domain) CanReroll(m IDomainMember) (bool, error) {
	// TODO
	return false, nil
}

func (d *domain) Reroll(m IDomainMember) (*Colour, error) {
	// TODO
	return nil, nil
}

func (d *domain) Freeze(m IDomainMember) error {
	// TODO
	return nil
}

func (d *domain) Unfreeze(m IDomainMember) error {
	// TODO
	return nil
}

func (d *domain) GetColourRole(mem IDomainMember) IDomainRole {
	expectName := d.getColourRoleName(mem)
	for _, role := range mem.Roles() {
		if role.Name() == expectName {
			return role
		}
	}
	return nil
}

func (d *domain) getColourRoleName(mem IDomainMember) string {
	roleName := mem.Nick()
	if roleName == "" {
		roleName = mem.Username()
	}
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
	height, err := d.GetColourRoleHeight(s, guildID)
	if err != nil {
		return nil, err
	}
	err = s.GuildRoleReorder(guildID, id, height)
	if err != nil {
		return nil, err
	}
	return NewDomainRole(id, roleName, col), nil
}

func (d *domain) AssignColourRole(s IDomainSession, mem IDomainMember, role IDomainRole) error {
	return s.GuildMemberRoleAdd(mem.Guild().ID(), mem.UserID(), role.ID())
}

func (d *domain) GetColourRoleHeight(s IDomainSession, guildID string) (int, error) {
	if d.maxRoleHeight > -1 {
		return d.maxRoleHeight, nil
	}
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return -1, err
	}
	for i, role := range roles {
		if role.Name() == d.maxHeightRoleName {
			d.maxRoleHeight = i
			return i, nil
		}
	}
	return -1, nil
}
