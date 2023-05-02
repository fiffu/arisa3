package types

import (
	"context"
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
			cmd:  NewCommand("name").Desc("Description"),
			expect: &dgo.ApplicationCommand{
				Name:        "name",
				Description: "Description",
			},
		},
		{
			desc: "ForChat() sets command type as ChatApplicationCommand",
			cmd:  NewCommand("name").ForChat(),
			expect: &dgo.ApplicationCommand{
				Name: "name",
				Type: dgo.ChatApplicationCommand,
			},
		},
		{
			desc: "ForUser() sets command type as UserApplicationCommand",
			cmd:  NewCommand("name").ForUser(),
			expect: &dgo.ApplicationCommand{
				Name: "name",
				Type: dgo.UserApplicationCommand,
			},
		},
		{
			desc: "ForMessage() sets command type as MessageApplicationCommand",
			cmd:  NewCommand("name").ForMessage(),
			expect: &dgo.ApplicationCommand{
				Name: "name",
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

	fn := func(context.Context, ICommandEvent) error { return expectErr }
	hdlr := Handler(fn)
	cmd := NewCommand("test").Handler(hdlr)

	actualErr := cmd.HandlerFunc()(nil, nil)
	assert.Error(t, actualErr)
	assert.Equal(t, actualErr.Error(), expectErr.Error())
}

func Test_mustValidate(t *testing.T) {
	testCases := []struct {
		desc        string
		commandName string
		expectOK    bool
	}{
		{
			desc:        "lowercase and compliant",
			commandName: "asdf",
			expectOK:    true,
		},
		{
			desc:        "must utilize lowercase variant of any letters used",
			commandName: "ASDF",
			expectOK:    false,
		},
		{
			desc:        "empty name should not work",
			commandName: "",
			expectOK:    false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			if tc.expectOK {
				NewCommand(tc.commandName)
			} else {
				assert.Panics(t, func() {
					NewCommand(tc.commandName)
				})
			}
		})
	}
}
