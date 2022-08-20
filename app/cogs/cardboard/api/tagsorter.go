package api

import "sort"

type TagComparer func(*Tag, *Tag) (leftShouldBeBeforeRight bool)

var (
	ByAlphabeticalOrder TagComparer = func(x, y *Tag) bool { return x.Name < y.Name }
	ByTagLength         TagComparer = func(x, y *Tag) bool {
		i, j := len(x.Name), len(y.Name)
		if i == j {
			return ByAlphabeticalOrder(x, y)
		}
		return i < j
	}
)

// TagsSorter implements sort.Interface.
type TagsSorter struct {
	Data    []*Tag
	Compare TagComparer
}

func (s TagsSorter) Len() int           { return len(s.Data) }
func (s TagsSorter) Less(i, j int) bool { return s.Compare(s.Data[i], s.Data[j]) }
func (s TagsSorter) Swap(i, j int)      { s.Data[i], s.Data[j] = s.Data[j], s.Data[i] }
func (s TagsSorter) Sorted() []*Tag     { sort.Sort(s); return s.Data }
