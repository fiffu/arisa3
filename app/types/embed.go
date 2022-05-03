package types

import (
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

type IEmbed interface {
	Data() *dgo.MessageEmbed
	URL(string) IEmbed
	Title(string) IEmbed
	Description(string) IEmbed
	Timestamp(time.Time) IEmbed
	Colour(int) IEmbed
	Image(url string) IEmbed
	Video(url string) IEmbed
	Thumbnail(url string) IEmbed
	Author(url, name, iconURL string) IEmbed
	Field(key, value string, inline bool) IEmbed // put a field
	Footer(text, iconURL string) IEmbed
}

type embed struct {
	d *dgo.MessageEmbed
}

func NewEmbed() IEmbed {
	d := &dgo.MessageEmbed{
		Fields: make([]*dgo.MessageEmbedField, 0),
	}
	return &embed{d}
}

// Getters

func (e *embed) Data() *dgo.MessageEmbed { return e.d }

// Simple setters

func (e *embed) URL(s string) IEmbed          { e.d.URL = s; return e }
func (e *embed) Title(s string) IEmbed        { e.d.Title = s; return e }
func (e *embed) Description(s string) IEmbed  { e.d.Description = s; return e }
func (e *embed) Timestamp(t time.Time) IEmbed { e.d.Timestamp = t.Format(time.RFC3339); return e }
func (e *embed) Colour(i int) IEmbed          { e.d.Color = i; return e }

// Complex setters

func (e *embed) Image(url string) IEmbed {
	e.d.Image = &dgo.MessageEmbedImage{URL: url}
	return e
}
func (e *embed) Video(url string) IEmbed {
	e.d.Video = &dgo.MessageEmbedVideo{URL: url}
	return e
}
func (e *embed) Thumbnail(url string) IEmbed {
	e.d.Thumbnail = &dgo.MessageEmbedThumbnail{URL: url}
	return e
}
func (e *embed) Author(url, name, iconURL string) IEmbed {
	e.d.Author = &dgo.MessageEmbedAuthor{
		URL:     url,
		Name:    name,
		IconURL: iconURL,
	}
	return e
}
func (e *embed) Field(k, v string, inline bool) IEmbed {
	fld := &dgo.MessageEmbedField{
		Name:   k,
		Value:  v,
		Inline: inline,
	}
	e.d.Fields = append(e.d.Fields, fld)
	return e
}
func (e *embed) Footer(text, iconURL string) IEmbed {
	e.d.Footer = &dgo.MessageEmbedFooter{
		Text:    text,
		IconURL: iconURL,
	}
	return e
}
