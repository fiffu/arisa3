package api

import (
	"fmt"
	"net/url"
)

func GetPostURL(post *Post) string {
	return fmt.Sprint(apiHostHTTP, "/posts/", post.ID)
}
func GetSearchURL(queryStr string) string {
	return fmt.Sprint(apiHostHTTP, "/posts?utf8=%E2%9C%93&tags=", url.PathEscape(queryStr))
}
