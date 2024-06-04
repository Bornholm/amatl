package toc

type Branch[T any] struct {
	value    T
	parent   *Branch[T]
	branches []*Branch[T]
}

func (b *Branch[T]) Add(branch *Branch[T]) *Branch[T] {
	branch.parent = b
	b.branches = append(b.branches, branch)
	return branch
}

func (b *Branch[T]) Remove(branch *Branch[T]) (*Branch[T], bool) {
	for _, b := range b.branches {
		if b != branch {
			continue
		}

		b.parent = nil

		return b, true
	}

	return nil, false
}

func (b *Branch[T]) ParentAt(level int) *Branch[T] {
	if level >= b.Level() {
		return nil
	}

	parent := b.parent
	if parent == nil {
		return nil
	}

	if parent.Level() == level {
		return parent
	}

	return parent.ParentAt(level)
}

func (b *Branch[T]) Parent() *Branch[T] {
	return b.parent
}

func (b *Branch[T]) Branches() []*Branch[T] {
	return b.branches
}

func (b *Branch[T]) Set(v T) {
	b.value = v
}

func (b *Branch[T]) Get() T {
	return b.value
}

func (b *Branch[T]) Level() int {
	if b.parent == nil {
		return 0
	}

	level := 0
	branch := b

	for {
		parent := branch.parent
		if parent != nil {
			level += 1
			branch = parent
			continue
		}

		return level
	}
}

func (b *Branch[T]) Walk(walk func(branch *Branch[T])) {
	walk(b)
	for _, b := range b.branches {
		b.Walk(walk)
	}
}

func NewBranch[T any](branches ...*Branch[T]) *Branch[T] {
	return &Branch[T]{
		branches: branches,
	}
}
