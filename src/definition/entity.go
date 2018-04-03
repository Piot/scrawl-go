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

import "fmt"

type Entity struct {
	name            string
	entityTypeID    EntityTypeID
	componentFields []*ComponentField
}

func NewEntity(name string, componentFields []*ComponentField) *Entity {
	return &Entity{name: name, entityTypeID: NewEntityTypeIDFromString(name), componentFields: componentFields}
}

func (c *Entity) String() string {
	var s string
	s += fmt.Sprintf("[entity '%v' components:%d\n", c.name, len(c.componentFields))
	for _, field := range c.componentFields {
		s += "  " + field.String() + "\n"
	}
	s += "]"
	return s
}

func (c *Entity) Name() string {
	return c.name
}

func (c *Entity) ID() EntityTypeID {
	return c.entityTypeID
}

func (c *Entity) Component(index int) *ComponentField {
	return c.componentFields[index]
}

func (c *Entity) Components() []*ComponentField {
	return c.componentFields
}
