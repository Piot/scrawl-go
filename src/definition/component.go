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

type Component struct {
	name              string
	index             uint8
	fields            []*Field
	eventReferences   []*EventReference
	commandReferences []*CommandReference
}

func NewComponent(name string, index uint8, fields []*Field, eventReferences []*EventReference, commandReferences []*CommandReference) *Component {
	return &Component{name: name, index: index, fields: fields, eventReferences: eventReferences, commandReferences: commandReferences}
}

func (c *Component) Name() string {
	return c.name
}

func (c *Component) Index() uint8 {
	return c.index
}

func (c *Component) Fields() []*Field {
	return c.fields
}

func (c *Component) EventReferences() []*EventReference {
	return c.eventReferences
}

func (c *Component) CommandReferences() []*CommandReference {
	return c.commandReferences
}

func (c *Component) String() string {
	var s string
	s += fmt.Sprintf("[component '%v' fields:%d]\n", c.name, len(c.fields))
	for _, field := range c.fields {
		s += "    " + field.String() + "\n"
	}

	return s
}
