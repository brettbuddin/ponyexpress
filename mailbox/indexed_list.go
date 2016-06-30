package mailbox

import "container/list"

type indexedList struct {
	*list.List
	index map[string]*list.Element
}

func newIndexedList() *indexedList {
	return &indexedList{list.New(), map[string]*list.Element{}}
}

type keyer interface {
	Key() string
}

func (i *indexedList) PushBack(v keyer) *list.Element {
	e := i.List.PushBack(v)
	i.index[v.Key()] = e
	return e
}

func (i *indexedList) Remove(e *list.Element) keyer {
	v := i.List.Remove(e).(keyer)
	delete(i.index, v.Key())
	return v
}

func (i *indexedList) GetKey(key string) (*list.Element, bool) {
	e, ok := i.index[key]
	return e, ok
}
