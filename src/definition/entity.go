/*

MIT License

Copyright (c) 2017 Peter Bjorklund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package definition

import (
	"fmt"
	"sort"
)

type EntityLod struct {
	lodLevel        int
	componentFields []*ComponentField
}

func (c *EntityLod) Component(index int) *ComponentField {
	return c.componentFields[index]
}

func (c *EntityLod) Components() []*ComponentField {
	return c.componentFields
}

func (c *EntityLod) String() string {
	var s string
	s += fmt.Sprintf("[lod%d components:%d\n", c.lodLevel, len(c.componentFields))
	for _, field := range c.componentFields {
		s += "  " + field.String() + "\n"
	}
	s += "]"
	return s
}

func NewEntityLod(lodLevel int, componentFields []*ComponentField) *EntityLod {
	return &EntityLod{lodLevel: lodLevel, componentFields: componentFields}
}

type Entity struct {
	name         string
	entityTypeID EntityTypeID
	index        EntityIndex
	lods         map[int]*EntityLod
}

func NewEntity(name string, index EntityIndex, componentFields []*ComponentField) *Entity {
	lods := make(map[int]*EntityLod, 1)
	lods[0] = NewEntityLod(0, componentFields)
	return &Entity{name: name, index: index, entityTypeID: NewEntityTypeIDFromString(name), lods: lods}
}

func (c *Entity) String() string {
	var s string
	s += fmt.Sprintf("[entity '%v' lodLevels:%d\n", c.name, len(c.lods))

	keys := make([]int, 0, len(c.lods))
	for k := range c.lods {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, key := range keys {
		lod := c.lods[key]
		s += lod.String()
	}

	return s
}

func (c *Entity) Name() string {
	return c.name
}

func (c *Entity) NewLod(lodLevel int, componentFields []*ComponentField) (*EntityLod, error) {
	_, doesExist := c.lods[lodLevel]
	if doesExist {
		return nil, fmt.Errorf("lod level %d already exists", lodLevel)
	}
	lod := NewEntityLod(lodLevel, componentFields)
	c.lods[lodLevel] = lod
	return lod, nil
}

func (c *Entity) ID() EntityTypeID {
	return c.entityTypeID
}

func (c *Entity) Index() EntityIndex {
	return c.index
}

func (c *Entity) HighestLevelOfDetail() *EntityLod {
	return c.lods[0]
}

func (c *Entity) Lod(lodLevel int) *EntityLod {
	return c.lods[lodLevel]
}
