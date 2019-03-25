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

type ComponentFieldType uint8

const (
	ComponentReferenceType ComponentFieldType = iota
	RawReferenceType
)

type ComponentField struct {
	index     int
	name      string
	component *Component
	rawType   string
	fieldType ComponentFieldType
}

func (c *ComponentField) Index() int {
	return c.index
}
func (c *ComponentField) Name() string {
	return c.name
}

func (c *ComponentField) Component() *Component {
	if c.fieldType != ComponentReferenceType {
		panic("wrong component field type")
	}

	return c.component
}

func (c *ComponentField) HasRawReference() bool {
	return c.fieldType == RawReferenceType
}

func (c *ComponentField) RawType() string {
	if c.fieldType != RawReferenceType {
		panic("wrong component field type")
	}
	return c.rawType
}

func NewComponentField(index int, name string, component *Component) *ComponentField {
	return &ComponentField{index: index, fieldType: ComponentReferenceType, name: name, component: component}
}

func NewComponentFieldRawType(index int, name string, rawType string) *ComponentField {
	return &ComponentField{index: index, fieldType: RawReferenceType, name: name, rawType: rawType}
}

func (c *ComponentField) String() string {
	var s string
	s += fmt.Sprintf("[componentfield '%v' %s   ]", c.name, c.Component())
	return s
}
