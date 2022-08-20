package types

import (
	"errors"
	"testing"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func Test_Options_FindOptions(t *testing.T) {
	cmd := NewCommand("test-command")
	cmd.Options(
		NewOption("test-opt").String(),
	)

	opt, ok := cmd.FindOption("test-opt")
	assert.True(t, ok)
	assert.NotNil(t, opt)
	assert.Equal(t, "test-opt", opt.Name())
}

func Test_Data(t *testing.T) {
	testCases := []struct {
		desc   string
		cmd    ICommand
		expect *dgo.ApplicationCommand
	}{
		{
			desc: "Desc() sets command description in data",
			cmd:  NewCommand("Name").Desc("Description"),
			expect: &dgo.ApplicationCommand{
				Name:        "Name",
				Description: "Description",
			},
		},
		{
			desc: "ForChat() sets command type as ChatApplicationCommand",
			cmd:  NewCommand("Name").ForChat(),
			expect: &dgo.ApplicationCommand{
				Name: "Name",
				Type: dgo.ChatApplicationCommand,
			},
		},
		{
			desc: "ForUser() sets command type as UserApplicationCommand",
			cmd:  NewCommand("Name").ForUser(),
			expect: &dgo.ApplicationCommand{
				Name: "Name",
				Type: dgo.UserApplicationCommand,
			},
		},
		{
			desc: "ForMessage() sets command type as MessageApplicationCommand",
			cmd:  NewCommand("Name").ForMessage(),
			expect: &dgo.ApplicationCommand{
				Name: "Name",
				Type: dgo.MessageApplicationCommand,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := tc.cmd.Data()
			assert.Equal(t, tc.expect.Name, tc.cmd.Name())
			assert.Equal(t, tc.expect, actual)
		})
	}
}

func Test_Handler(t *testing.T) {
	expectErr := errors.New(time.Now().Format(time.RFC3339Nano))

	fn := func(evt ICommandEvent) error { return expectErr }
	hdlr := Handler(fn)
	cmd := NewCommand("test").Handler(hdlr)

	actualErr := cmd.HandlerFunc()(nil)
	assert.Error(t, actualErr)
	assert.Equal(t, actualErr.Error(), expectErr.Error())
}
