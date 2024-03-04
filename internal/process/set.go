package process

type Set map[string]struct{}

func (set Set) Add(item string) {
	set[item] = struct{}{}
}

func (set Set) Remove(item string) {
	delete(set, item)
}

func (set Set) Contains(item string) bool {
	_, ok := set[item]
	return ok
}
