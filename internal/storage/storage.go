package storage

import (
	"slices"
	"sync"
	"time"
)

type Document[V any] struct {
	Value     V
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewDocument[V any](value V) *Document[V] {
	return &Document[V]{
		Value:     value,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

type Storage[V any] struct {
	mu     sync.RWMutex
	values []*Document[V]
}

func New[V any](capacity int) *Storage[V] {
	return &Storage[V]{
		values: make([]*Document[V], 0, capacity),
	}
}

func (s *Storage[V]) All() []*Document[V] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return slices.Clone(s.values)
}

func (s *Storage[V]) GetDocument(f func(V) bool) (*Document[V], bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, v := range s.values {
		if f(v.Value) {
			return v, true
		}
	}

	return nil, false
}

func (s *Storage[V]) GetDocuments(f func(V) bool) []*Document[V] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*Document[V], 0, len(s.values))

	for _, v := range s.values {
		if f(v.Value) {
			result = append(result, v)
		}
	}

	return result
}

func (s *Storage[V]) Get(f func(V) bool) (V, bool) {
	document, ok := s.GetDocument(f)
	if ok {
		return document.Value, true
	}

	var v V

	return v, false
}

func (s *Storage[V]) GetMany(f func(V) bool) []V {
	documents := s.GetDocuments(f)

	result := make([]V, len(documents))
	for i := range documents {
		result[i] = documents[i].Value
	}

	return result
}

func (s *Storage[V]) Insert(v V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.values = append(s.values, NewDocument(v))
}

func (s *Storage[V]) Replace(v V, f func(V) bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	index := slices.IndexFunc(s.values, func(d *Document[V]) bool {
		return f(d.Value)
	})
	if index == -1 {
		return false
	}

	document := s.values[index]
	document.Value = v
	document.UpdatedAt = time.Now().UTC()

	s.values[index] = document

	return true
}
