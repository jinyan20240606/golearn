package course

type Course struct {
	Name string
}

func (c Course) GetName() string {
	return c.Name
}
