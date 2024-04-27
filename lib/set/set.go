package set

type (
	Set[T comparable]  struct {
		hash map[T]nothing
	}

	nothing struct{}
)

func New[T comparable](initial ...T) *Set[T] {
	s := &Set[T]{make(map[T]nothing)}

	for _, v := range initial {
		s.Insert(v)
	}

	return s
}

func (this *Set[T]) Has(element T) bool {
	_, exists := this.hash[element]
	return exists
}

func (this *Set[T]) Insert(element T) {
	this.hash[element] = nothing{}
}

func (this *Set[T]) Remove(element T) {
	delete(this.hash, element)
}

func (this *Set[T]) ToSlice() []T {
	var arr = make([]T, len(this.hash))
	i:= 0
	for k := range this.hash {
    	arr[i] = k
    	i++
	}
	return arr
}

func (this *Set[T]) Do(f func(T)) {
	for k := range this.hash {
		f(k)
	}
}

func (this *Set[T]) Difference(set *Set[T]) *Set[T] {
	hash := make(map[T]nothing)

	for k := range this.hash {
		_, exists := set.hash[k]
		if  !exists {
			hash[k] = nothing{}
		}
	}

	return &Set[T]{ hash }
}

func (this *Set[T]) Intersection(set *Set[T]) *Set[T] {
	hash := make(map[T]nothing)

	for k := range this.hash {
		if _, exists := set.hash[k]; exists {
			hash[k] = nothing{}
		}
	}

	return &Set[T]{ hash }
}

func (this *Set[T]) Len() int {
	return len(this.hash)
}

func (this *Set[T]) SubsetOf(set *Set[T]) bool {
	if this.Len() > set.Len() {
		return false
	}
	for k := range this.hash {
		if _, exists := set.hash[k]; !exists {
			return false
		}
	}
	return true
}

func (this *Set[T]) ProperSubsetOf(set *Set[T]) bool {
	return this.SubsetOf(set) && this.Len() < set.Len()
}

func (this *Set[T]) Union(set *Set[T]) *Set[T] {
	hash := make(map[T]nothing)

	for k := range this.hash {
		hash[k] = nothing{}
	}
	for k := range set.hash {
		hash[k] = nothing{}
	}

	return &Set[T]{ hash }
}