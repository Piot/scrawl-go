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

type EntityArchetype struct {
	name         string
	entityTypeID EntityArchetypeID
	index        EntityIndex
	lods         []*EntityArchetypeLOD
	meta         MetaData
}

func NewEntityArchetype(name string, index EntityIndex, lods []*EntityArchetypeLOD, meta MetaData) *EntityArchetype {
	return &EntityArchetype{name: name, index: index, entityTypeID: NewEntityArchetypeIDFromString(name), lods: lods, meta: meta}
}

func (c *EntityArchetype) String() string {
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

func (c *EntityArchetype) Name() string {
	return c.name
}

func (c *EntityArchetype) NewLod(lodLevel int, items []*EntityArchetypeItem) (*EntityArchetypeLOD, error) {
	lastIndex := len(c.lods)
	if lastIndex+1 != lodLevel {
		return nil, fmt.Errorf("must add lods in order %v (%v)", lodLevel, lastIndex)
	}

	lod := NewEntityArchetypeLOD(lodLevel, items)
	c.lods = append(c.lods, lod)

	return lod, nil
}

func (c *EntityArchetype) ID() EntityArchetypeID {
	return c.entityTypeID
}

func (c *EntityArchetype) Index() EntityIndex {
	return c.index
}

func (c *EntityArchetype) HighestLevelOfDetail() *EntityArchetypeLOD {
	return c.lods[0]
}

func (c *EntityArchetype) Lod(lodLevel int) *EntityArchetypeLOD {
	return c.lods[lodLevel]
}

func (c *EntityArchetype) Lods() []*EntityArchetypeLOD {
	return c.lods
}

func (c *EntityArchetype) Meta() MetaData {
	return c.meta
}
