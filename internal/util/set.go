package util

type Set[K comparable] struct {
	m map[K]struct{}
}

// NewSet returns an empty set.
func NewSet[K comparable]() Set[K] {
	return Set[K]{
		m: make(map[K]struct{}),
	}
}

// SetOf returns a new set initialized with the given elements
func SetOf[K comparable](elements ...K) Set[K] {
	s := NewSet[K]()
	for _, element := range elements {
		s.Put(element)
	}

	return s
}

// Put adds the element to the set.
func (s Set[K]) Put(element K) {
	s.m[element] = struct{}{}
}

// Has returns true only if the element is in the set.
func (s Set[K]) Has(element K) bool {
	_, ok := s.m[element]
	return ok
}

// Remove removes the element from the set.
func (s Set[K]) Remove(element K) {
	delete(s.m, element)
}

// Clear removes all elements from the set.
func (s Set[K]) Clear() {
	for k := range s.m {
		delete(s.m, k)
	}
}

// Size returns the number of elements in the set.
func (s Set[K]) Size() int {
	return len(s.m)
}

// Each calls 'fn' on every element in the set in no particular order.
func (s Set[K]) Each(fn func(key K)) {
	for k := range s.m {
		fn(k)
	}
}

func (s Set[K]) ToSlice() []K {
	keys := make([]K, 0, len(s.m))
	for elem := range s.m {
		keys = append(keys, elem)
	}

	return keys
}
