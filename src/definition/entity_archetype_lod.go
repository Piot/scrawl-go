package definition

import "fmt"

type EntityArchetypeLOD struct {
	lodLevel int
	items    []*EntityArchetypeItem
}

func (c *EntityArchetypeLOD) Level() int {
	return c.lodLevel
}

func (c *EntityArchetypeLOD) EntityArchetypeItem(index int) *EntityArchetypeItem {
	return c.items[index]
}

func (c *EntityArchetypeLOD) Items() []*EntityArchetypeItem {
	return c.items
}

func (c *EntityArchetypeLOD) String() string {
	var s string
	s += fmt.Sprintf("[lod%d items:%d\n", c.lodLevel, len(c.items))

	for _, field := range c.items {
		s += "  " + field.String() + "\n"
	}

	s += "]"

	return s
}

func NewEntityArchetypeLOD(lodLevel int, items []*EntityArchetypeItem) *EntityArchetypeLOD {
	return &EntityArchetypeLOD{lodLevel: lodLevel, items: items}
}
