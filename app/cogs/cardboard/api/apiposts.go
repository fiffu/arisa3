package api

import (
	"fmt"
	"strings"
)

const (
	MaxPosts = 100
)

type Post struct {
	ID            int    `json:"id"`
	MD5           string `json:"md5"`
	FileExt       string `json:"file_ext"`
	AllTags       string `json:"tag_string"`
	GeneralTags   string `json:"tag_string_general"`
	CharacterTags string `json:"tag_string_character"`
	CopyrightTags string `json:"tag_string_copyright"`
	ArtistTags    string `json:"tag_string_artist"`
	File          string `json:"file_url"`
	FileLarge     string `json:"large_file_url"`
	FilePreview   string `json:"preview_file_url"`

	url string
}

func (p *Post) TagsList() []string {
	return strings.Split(p.AllTags, " ")
}

func (p *Post) GetFileURL() string {
	if p.url != "" {
		return p.url
	}
	switch {
	case p.File != "":
		p.setURL(p.File)
	case p.FileLarge != "":
		p.setURL(p.FileLarge)
	case p.FilePreview != "":
		p.setURL(p.FilePreview)
	default:
		return ""
	}
	return p.GetFileURL()
}

func (p *Post) setURL(url string) {
	p.url = p.repairURL(url)
}

func (p *Post) repairURL(url string) string {
	// Fix relative file reference
	if len(url) < 4 {
		return ""
	}
	if url[:4] != "http" {
		url = apiHostHTTPS + url
	}

	// Sometimes url returns with repeated forward slash
	url = strings.Replace(url, "//data/", "/data/", 1)

	return url
}

// GetPosts lists posts matching the given query.
func (c *client) GetPosts(tags []string) ([]*Post, error) {
	ctx, cancel := c.httpContext()
	defer cancel()

	var result []*Post
	builder := c.requestPosts().
		Param("tags", commaJoin(tags)).
		Param("limit", fmt.Sprint(MaxPosts)).
		ToJSON(&result)

	err := c.fetch(ctx, builder)
	return result, err
}
