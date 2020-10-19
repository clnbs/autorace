package engine

type Menu struct {
	Name string
	Items []*MenuItem
	SelectedItem int
}

type MenuItem struct {
	Name string
	Action func()
}
