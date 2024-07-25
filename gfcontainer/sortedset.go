package gfcontainer

import "slices"

type SortedSet[T comparable] struct {
	mp  map[T]struct{}
	seq []*T
}

func NewSortedSet[T comparable]() *SortedSet[T] {
	return NewSortedSetWithCapacity[T](0)
}

func NewSortedSetWithCapacity[T comparable](capacity int) *SortedSet[T] {
	return &SortedSet[T]{
		mp:  make(map[T]struct{}, capacity),
		seq: make([]*T, 0, capacity),
	}
}

func (s *SortedSet[T]) List() []T {
	ret := make([]T, 0, len(s.seq))
	for _, v := range s.seq {
		if v == nil || !s.Has(*v) {
			continue
		}
		ret = append(ret, *v)
	}
	return ret
}

func (s *SortedSet[T]) Has(v T) bool {
	_, ok := s.mp[v]
	return ok
}

func (s *SortedSet[T]) Add(v T) {
	if s.Has(v) {
		return
	}
	s.mp[v] = struct{}{}
	s.seq = append(s.seq, &v)
}

func (s *SortedSet[T]) Del(v T) {
	delete(s.mp, v)
	idx := slices.Index(s.seq, &v)
	s.seq = slices.Delete(s.seq, idx, idx+1)
}
