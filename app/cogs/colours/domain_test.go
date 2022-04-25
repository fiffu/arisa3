package colours

import (
	"testing"
	"time"

	"github.com/fiffu/arisa3/app/types"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var Any = gomock.Any()

func newTestingDomain(t *testing.T, cfg *Config) (
	*gomock.Controller, *types.MockICog, *MockIDomainRepository, IColoursDomain) {

	t.Helper()
	ctrl := gomock.NewController(t)

	cog := types.NewMockICog(ctrl)
	cog.EXPECT().Name().AnyTimes().Return("test cog")

	repo := NewMockIDomainRepository(ctrl)

	return ctrl, cog, repo, NewColoursDomain(cog, repo, cfg)
}

func newTestingMember(ctrl *gomock.Controller, hasColourRole bool) *MockIDomainMember {
	mem := NewMockIDomainMember(ctrl)
	mem.EXPECT().UserID().AnyTimes().Return("123123123123")
	mem.EXPECT().Username().AnyTimes().Return("test#1234")
	mem.EXPECT().Guild().AnyTimes().Return(NewDomainGuild("87979878098098908"))

	roles := make([]IDomainRole, 0)
	if hasColourRole {
		roleName := (&domain{}).GetColourRoleName(mem)
		roleColour := (&Colour{}).Random()
		role := &colourRole{name: roleName, colour: roleColour}
		roles = append(roles, role)
	}
	mem.EXPECT().Roles().AnyTimes().Return(roles)
	return mem
}

func newTestingConfig() *Config {
	return &Config{
		MaxRoleHeightName:  "[Arisa] Max colour role height",
		MutateCooldownMins: 240,
		RerollCooldownMins: 720,
		RerollPenaltyMins:  30,
	}
}

func Test_Reroll(t *testing.T) {
	const (
		Disallow  = 0
		Reuse     = 1
		Provision = 2
	)
	type testCases struct {
		name              string
		cooldownStartTime time.Time
		hasColourRole     bool

		expectOutcome int
	}

	tests := []testCases{
		{
			name:              "cooldown unfinished, has colourRole: disallow",
			cooldownStartTime: time.Now().Add(-1 * time.Minute), // 1 min ago
			hasColourRole:     true,
			expectOutcome:     Disallow,
		},
		{
			name:              "cooldown unfinished, no colourRole: disallow",
			cooldownStartTime: time.Now().Add(-1 * time.Minute), // 1 min ago
			hasColourRole:     false,
			expectOutcome:     Disallow,
		},
		{
			name:              "no cooldown, has colourRole: reroll with reuse",
			cooldownStartTime: Never,
			hasColourRole:     true,
			expectOutcome:     Reuse,
		},
		{
			name:              "no cooldown, no colourRole: reroll with provision",
			cooldownStartTime: Never,
			hasColourRole:     false,
			expectOutcome:     Provision,
		},
		{
			name:              "cooldown finished, has colourRole: reroll with reuse",
			cooldownStartTime: time.Now().Add(-3000 * time.Minute), // 3000 min ago
			hasColourRole:     true,
			expectOutcome:     Reuse,
		},
		{
			name:              "cooldown finished, no colourRole: reroll with provision",
			cooldownStartTime: time.Now().Add(-3000 * time.Minute), // 3000 min ago
			hasColourRole:     false,
			expectOutcome:     Provision,
		},
	}

	setup := func(hasColourRoles bool) (
		*types.MockICog, *MockIDomainRepository, IColoursDomain,
		*MockIDomainSession, *MockIDomainMember) {

		ctrl, cog, repo, d := newTestingDomain(t, newTestingConfig())
		s := NewMockIDomainSession(ctrl)

		mem := newTestingMember(ctrl, hasColourRoles)
		return cog, repo, d, s, mem
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			_, repo, d, s, mem := setup(tc.hasColourRole)

			repo.EXPECT().FetchUserState(Any, Reroll).
				AnyTimes().Return(tc.cooldownStartTime, nil)

			var expectError error
			switch tc.expectOutcome {
			case Disallow:
				repo.EXPECT().UpdateRerollPenalty(Any, Any).Return(nil)
				expectError = ErrCooldownPending

			case Provision:
				repo.EXPECT().UpdateReroll(Any, Any).Return(nil)
				s.EXPECT().GuildRoleCreate(Any)
				s.EXPECT().GuildRoleEdit(Any, Any, Any, Any)
				s.EXPECT().GuildRoles(Any)
				// s.EXPECT().GuildRoleReorder(Any, Any)  // commented out; lazy to mock guild roles
				s.EXPECT().GuildMemberRoleAdd(Any, Any, Any)

			case Reuse:
				repo.EXPECT().UpdateReroll(Any, Any).Return(nil)
				s.EXPECT().GuildRoleEdit(Any, Any, Any, Any)
			}

			_, err := d.Reroll(s, mem)

			if expectError != nil {
				assert.Equal(
					t,
					expectError.Error(),
					err.Error(),
				)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func Test_Mutate(t *testing.T) {
	const (
		Noop     = 0
		Disallow = 1
		Allow    = 2
	)
	type testCases struct {
		name              string
		cooldownStartTime time.Time
		hasColourRole     bool

		expectOutcome int
	}

	tests := []testCases{
		{
			name:              "cooldown not finished should not allow mutate",
			cooldownStartTime: time.Now().Add(-1 * time.Minute), // 1 min ago
			hasColourRole:     true,
			expectOutcome:     Disallow,
		},
		{
			name:              "no cooldown, has colourRole: mutate",
			cooldownStartTime: Never,
			hasColourRole:     true,
			expectOutcome:     Allow,
		},
		{
			name:              "no cooldown, no colourRole: do nothing",
			cooldownStartTime: Never,
			hasColourRole:     false,
			expectOutcome:     Noop,
		},
		{
			name:              "cooldown finished, has colourRole: mutate",
			cooldownStartTime: time.Now().Add(-3000 * time.Minute), // 3000 min ago
			hasColourRole:     true,
			expectOutcome:     Allow,
		},
		{
			name:              "cooldown finished, no colourRole: reroll with provision",
			cooldownStartTime: time.Now().Add(-3000 * time.Minute), // 3000 min ago
			hasColourRole:     false,
			expectOutcome:     Noop,
		},
	}

	setup := func(hasColourRoles bool) (
		*types.MockICog, *MockIDomainRepository, IColoursDomain,
		*MockIDomainSession, *MockIDomainMember) {

		ctrl, cog, repo, d := newTestingDomain(t, newTestingConfig())
		s := NewMockIDomainSession(ctrl)

		cog.EXPECT().Name().AnyTimes().Return("test cog")

		mem := newTestingMember(ctrl, hasColourRoles)
		return cog, repo, d, s, mem
	}
	Any := gomock.Any()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			_, repo, d, s, mem := setup(tc.hasColourRole)

			repo.EXPECT().FetchUserState(Any, Mutate).
				AnyTimes().Return(tc.cooldownStartTime, nil)

			var expectError error
			switch tc.expectOutcome {
			case Noop:
				break

			case Disallow:
				expectError = ErrCooldownPending

			case Allow:
				repo.EXPECT().UpdateMutate(Any, Any).Return(nil)
				s.EXPECT().GuildRoleEdit(Any, Any, Any, Any)
			}

			_, err := d.Mutate(s, mem)

			if expectError != nil {
				assert.Error(t, err)
				assert.Equal(
					t,
					expectError.Error(),
					err.Error(),
				)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func Test_GetLastFrozen(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockIDomainRepository(ctrl)
	d := &domain{repo: repo}
	rsn := Freeze
	{
		// happy case
		var expect = time.Now()
		var expectErr error
		repo.EXPECT().FetchUserState(Any, rsn).Return(expect, expectErr)
		actual, err := d.GetLastFrozen(nil)
		assert.Equal(t, expect, actual)
		assert.Equal(t, expectErr, err)
	}
	{
		// error case
		var expect time.Time
		var expectErr = assert.AnError
		repo.EXPECT().FetchUserState(Any, rsn).Return(expect, expectErr)
		actual, err := d.GetLastFrozen(nil)
		assert.Equal(t, expect, actual)
		assert.Equal(t, expectErr, err)
	}
}

func Test_GetLastMutate(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockIDomainRepository(ctrl)
	d := &domain{repo: repo, rerollCooldownMins: 1}
	rsn := Mutate
	{
		// happy case
		var expect = time.Now().Add(-10 * time.Minute)
		var expectOK = true
		var expectErr error
		repo.EXPECT().FetchUserState(Any, rsn).Return(expect, expectErr)
		actual, ok, err := d.GetLastMutate(nil)
		assert.Equal(t, expect, actual)
		assert.Equal(t, expectOK, ok)
		assert.Equal(t, expectErr, err)
	}
	{
		// happy case
		var expect time.Time
		var expectOK = false
		var expectErr = assert.AnError
		repo.EXPECT().FetchUserState(Any, rsn).Return(expect, expectErr)
		actual, ok, err := d.GetLastMutate(nil)
		assert.Equal(t, expect, actual)
		assert.Equal(t, expectOK, ok)
		assert.Equal(t, expectErr, err)
	}
}

func Test_GetLastReroll(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockIDomainRepository(ctrl)
	d := &domain{repo: repo, rerollCooldownMins: 1}
	rsn := Reroll
	{
		// happy case
		var expect = time.Now().Add(-10 * time.Minute)
		var expectOK = true
		var expectErr error
		repo.EXPECT().FetchUserState(Any, rsn).Return(expect, expectErr)
		actual, ok, err := d.GetLastReroll(nil)
		assert.Equal(t, expect, actual)
		assert.Equal(t, expectOK, ok)
		assert.Equal(t, expectErr, err)
	}
	{
		// happy case
		var expect time.Time
		var expectOK = false
		var expectErr = assert.AnError
		repo.EXPECT().FetchUserState(Any, rsn).Return(expect, expectErr)
		actual, ok, err := d.GetLastReroll(nil)
		assert.Equal(t, expect, actual)
		assert.Equal(t, expectOK, ok)
		assert.Equal(t, expectErr, err)
	}
}

func Test_Freeze(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockIDomainRepository(ctrl)
	d := &domain{repo: repo}
	repo.EXPECT().UpdateFreeze(nil).Return(assert.AnError)
	actual := d.Freeze(nil)
	assert.Error(t, actual)
}

func Test_Unfreeze(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := NewMockIDomainRepository(ctrl)
	d := &domain{repo: repo}
	repo.EXPECT().UpdateUnfreeze(nil).Return(assert.AnError)
	actual := d.Unfreeze(nil)
	assert.Error(t, actual)
}

func Test_hasCooldownFinished(t *testing.T) {
	d := &domain{}
	now := time.Now()

	type testCase struct {
		cooldownStartTime time.Time
		cooldownDuration  time.Duration
		expectFinished    bool
	}

	tests := []testCase{
		{
			now.Add(-5 * time.Minute), // 5 mins ago
			300 * time.Minute,         // 300 mins cd
			false,                     // cooldown should not be finished
		},
		{
			now.Add(-300 * time.Minute), // 300 mins ago
			5 * time.Minute,             // 5 mins cd
			true,                        // cooldown should be finished
		},
		{
			Never, // cooldown never started
			300 * time.Minute,
			true,
		},
		{
			time.Now(),
			0, // 0 mins should always pass
			true,
		},
	}
	for _, tc := range tests {
		actual := d.hasCooldownFinished(tc.cooldownStartTime, tc.cooldownDuration)
		assert.Equal(t, tc.expectFinished, actual)
	}
}

func Test_HasColourRole(t *testing.T) {
	for _, hasColourRole := range []bool{true, false} {
		ctrl, _, _, d := newTestingDomain(t, newTestingConfig())
		mem := newTestingMember(ctrl, hasColourRole)
		actual := d.HasColourRole(mem)
		assert.Equal(t, hasColourRole, actual)
	}
}
