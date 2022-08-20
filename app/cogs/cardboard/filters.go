package cardboard

import (
	"math/rand"

	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/fiffu/arisa3/lib/functional"
)

type opsHelper struct {
	tags2operation OperationsMap
}

func (h *opsHelper) hasOperation(post *api.Post, hasOper TagOperation) bool {
	for _, tag := range post.TagsList() {
		if oper, ok := h.tags2operation[tag]; ok {
			if oper == hasOper {
				return true
			}
		}
	}
	return false
}

func (h *opsHelper) gatherByOperation(posts []*api.Post, oper TagOperation) (gathered, remaining []*api.Post) {
	gathered = make([]*api.Post, 0)
	for i, post := range posts {
		if h.hasOperation(post, oper) {
			gathered = append(gathered, post)
			if i+1 > len(posts) {
				posts = posts[:i]
			} else {
				posts = append(posts[:i], posts[i+1:]...)
			}
		}
	}
	remaining = posts
	return /* gathered, remaining */
}

func postsFilter(predicate func(*api.Post) bool) Filter {
	return func(posts []*api.Post) []*api.Post {
		allowed := make([]*api.Post, 0)
		for _, post := range posts {
			if predicate(post) {
				allowed = append(allowed, post)
			}
		}
		return allowed
	}
}

func PromoteFilter(h *opsHelper) Filter {
	return func(posts []*api.Post) []*api.Post {
		gathered, remaining := h.gatherByOperation(posts, Promote)
		return append(gathered, remaining...)
	}
}

func DemoteFilter(h *opsHelper) Filter {
	return func(posts []*api.Post) []*api.Post {
		gathered, remaining := h.gatherByOperation(posts, Demote)
		return append(gathered, remaining...)
	}
}

func OmitFilter(h *opsHelper) Filter {
	return func(posts []*api.Post) []*api.Post {
		_, remaining := h.gatherByOperation(posts, Omit)
		return remaining
	}
}

func Shuffle() Filter {
	return func(posts []*api.Post) []*api.Post {
		rand.Shuffle(
			len(posts),
			func(i, j int) {
				posts[i], posts[j] = posts[j], posts[i]
			},
		)
		return posts
	}
}

func HasMediaFile() Filter {
	condition := func(post *api.Post) bool {
		hasMediaExt := functional.Contains(api.MediaFileExts, post.FileExt)
		hasURL := post.GetFileURL() != ""
		return hasMediaExt && hasURL
	}
	return postsFilter(condition)
}

func HasURL() Filter {
	condition := func(post *api.Post) bool {
		return post.GetFileURL() != ""
	}
	return postsFilter(condition)
}
