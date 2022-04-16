package types

import (
	dgo "github.com/bwmarrin/discordgo"
)

type IArgs interface {
	Int(key string) (int, bool)
	Float(key string) (float64, bool)
	String(key string) (string, bool)
	Bool(key string) (bool, bool)
	Channel(key string) (*dgo.Channel, bool)
	Role(key string) (*dgo.Role, bool)
	User(key string) (*dgo.User, bool)
}

// args implements IArgs
type args struct {
	mapping map[string]*dgo.ApplicationCommandInteractionDataOption
}

func NewArgs(m map[string]*dgo.ApplicationCommandInteractionDataOption) IArgs {
	return &args{mapping: m}
}
func (a *args) fetch(key string) interface{} {
	if v, ok := a.mapping[key]; ok {
		return v.Value
	} else {
		return v
	}
}
func (a *args) Int(key string) (int, bool)       { v, ok := a.fetch(key).(int); return v, ok }
func (a *args) Float(key string) (float64, bool) { v, ok := a.fetch(key).(float64); return v, ok }
func (a *args) String(key string) (string, bool) { v, ok := a.fetch(key).(string); return v, ok }
func (a *args) Bool(key string) (bool, bool)     { v, ok := a.fetch(key).(bool); return v, ok }

func (a *args) Channel(key string) (*dgo.Channel, bool) {
	v, ok := a.fetch(key).(*dgo.Channel)
	return v, ok
}
func (a *args) Role(key string) (*dgo.Role, bool) {
	v, ok := a.fetch(key).(*dgo.Role)
	return v, ok
}
func (a *args) User(key string) (*dgo.User, bool) {
	v, ok := a.fetch(key).(*dgo.User)
	return v, ok
}
