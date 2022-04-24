package types

import (
	"encoding/json"
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
)

type ICommandResponse interface {
	Data() *dgo.InteractionResponse
	TTS() ICommandResponse
	Content(string) ICommandResponse
	Embeds(...IEmbed) ICommandResponse
	// Components() ICommandResponse
	// Files() ICommandResponse
	String() string // stringify
}

type Response struct {
	data *dgo.InteractionResponse
}

func NewResponse() ICommandResponse {
	d := &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{},
	}
	return &Response{data: d}
}

func (r *Response) Data() *dgo.InteractionResponse      { return r.data }
func (r *Response) TTS() ICommandResponse               { r.data.Data.TTS = true; return r }
func (r *Response) Content(msg string) ICommandResponse { r.data.Data.Content = msg; return r }
func (r *Response) Contentf(msg string, v ...interface{}) ICommandResponse {
	r.data.Data.Content = fmt.Sprintf(msg, v...)
	return r
}
func (r *Response) Embeds(embeds ...IEmbed) ICommandResponse {
	for _, e := range embeds {
		r.data.Data.Embeds = append(r.data.Data.Embeds, e.Data())
	}
	return r
}
func (r *Response) String() string {
	embeds := make([]dgo.MessageEmbed, 0)
	if r.data.Data.Embeds != nil {
		for _, e := range r.data.Data.Embeds {
			embeds = append(embeds, *e)
		}
	}
	emb, _ := json.Marshal(embeds)
	return fmt.Sprintf(
		"content='%s' embeds=%s",
		r.data.Data.Content, string(emb),
	)
}
