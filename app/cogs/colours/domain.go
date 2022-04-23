package colours

import (
	"time"
)

type domain struct {
	repo IDomainRepository
}

// NewColoursDomain implements IColoursDomain
func NewColoursDomain(repo IDomainRepository) IColoursDomain {
	return &domain{
		repo: repo,
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

func (d *domain) Mutate(m IDomainMember) (*Colour, error) {
	// if role := m.ColourRole(); role != nil {
	// 	current := role.Colour()
	// 	newCol := current.Nudge()
	// }
	return nil, nil
}

func (d *domain) Reroll(m IDomainMember) (*Colour, error) {
	return nil, nil
}

func (d *domain) MakeRoleName(mem IDomainMember) string {
	if nick := mem.Nick(); nick != "" {
		return nick
	}
	return mem.UserName()
}
