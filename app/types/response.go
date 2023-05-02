package types

import (
	"encoding/json"
	"fmt"
	"io"

	dgo "github.com/bwmarrin/discordgo"
)

type ICommandResponse interface {
	Data() *dgo.InteractionResponse
	String() string
}

type Response struct {
	data *dgo.InteractionResponse
}

func NewResponse() *Response {
	d := &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{},
	}
	return &Response{data: d}
}

func (r *Response) TTS() *Response {
	r.data.Data.TTS = true
	return r
}

func (r *Response) Content(msg string) *Response {
	r.data.Data.Content = msg
	return r
}

func (r *Response) Contentf(msg string, v ...interface{}) *Response {
	r.data.Data.Content = fmt.Sprintf(msg, v...)
	return r
}

func (r *Response) Embeds(embeds ...IEmbed) *Response {
	for _, e := range embeds {
		r.data.Data.Embeds = append(r.data.Data.Embeds, e.Data())
	}
	return r
}

func (r *Response) File(name, contentType string, reader io.Reader) *Response {
	r.data.Data.Files = append(r.data.Data.Files, &dgo.File{
		Name:        name,
		ContentType: contentType,
		Reader:      reader,
	})
	return r
}

func (r *Response) Data() *dgo.InteractionResponse {
	return r.data
}

func (r *Response) String() string {
	data, _ := json.MarshalIndent(r.data, "| ", "  ")
	return string(data)
}
