package types

import (
	dgo "github.com/bwmarrin/discordgo"
)

type IOption interface {
	Data() *dgo.ApplicationCommandOption
	Desc(s string) IOption
	Min(n float64) IOption
	Max(n float64) IOption
	Required(n float64) IOption
	typed(t dgo.ApplicationCommandOptionType) IOption
	Int() IOption
	String() IOption
	Bool() IOption
	User() IOption
	Role() IOption
	Mention() IOption
	Number() IOption
	Attachment() IOption
	Channel() IOption
	ChannelType(n []dgo.ChannelType) IOption
}

type Option struct {
	data *dgo.ApplicationCommandOption
}

func NewOption(name string) IOption {
	d := &dgo.ApplicationCommandOption{
		Name:        name,
		Description: "(no description)", // this is mandatory
	}
	return &Option{data: d}
}

func (co *Option) Data() *dgo.ApplicationCommandOption              { return co.data }
func (co *Option) Desc(s string) IOption                            { co.data.Description = s; return co }
func (co *Option) Min(n float64) IOption                            { co.data.MinValue = &n; return co }
func (co *Option) Max(n float64) IOption                            { co.data.MaxValue = n; return co }
func (co *Option) Required(n float64) IOption                       { co.data.Required = true; return co }
func (co *Option) typed(t dgo.ApplicationCommandOptionType) IOption { co.data.Type = t; return co }
func (co *Option) Int() IOption                                     { return co.typed(dgo.ApplicationCommandOptionInteger) }
func (co *Option) String() IOption                                  { return co.typed(dgo.ApplicationCommandOptionString) }
func (co *Option) Bool() IOption                                    { return co.typed(dgo.ApplicationCommandOptionBoolean) }
func (co *Option) User() IOption                                    { return co.typed(dgo.ApplicationCommandOptionUser) }
func (co *Option) Role() IOption                                    { return co.typed(dgo.ApplicationCommandOptionRole) }
func (co *Option) Mention() IOption                                 { return co.typed(dgo.ApplicationCommandOptionMentionable) }
func (co *Option) Number() IOption                                  { return co.typed(dgo.ApplicationCommandOptionNumber) }
func (co *Option) Attachment() IOption                              { return co.typed(dgo.ApplicationCommandOptionAttachment) }
func (co *Option) Channel() IOption                                 { return co.typed(dgo.ApplicationCommandOptionChannel) }
func (co *Option) ChannelType(n []dgo.ChannelType) IOption          { co.data.ChannelTypes = n; return co } // me lazy
func (co *Option) Choice(k string, v interface{}) IOption {
	choice := &dgo.ApplicationCommandOptionChoice{Name: k, Value: v}
	co.data.Choices = append(co.data.Choices, choice)
	return co
}
