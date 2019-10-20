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

type ComponentTypeVariant uint8

const (
	ComponentTypeComponentReference ComponentTypeVariant = iota
	ComponentTypeField
)

type ComponentDataType struct {
	component *ComponentDataType
	field     *Field
	variant   ComponentTypeVariant
}

func NewComponentTypeUsingField(field *Field) *ComponentDataType {
	return &ComponentDataType{variant: ComponentTypeField, field: field}
}

func NewComponentTypeUsingComponent(component *ComponentDataType) *ComponentDataType {
	return &ComponentDataType{variant: ComponentTypeComponentReference, component: component}
}

func (c *ComponentDataType) HasFieldReference() bool {
	return c.variant == ComponentTypeField
}

func (c *ComponentDataType) FieldReference() *Field {
	if c.variant != ComponentTypeField {
		panic("wrong component field type")
	}
	return c.field
}

func (c *ComponentDataType) Name() string {
	switch c.variant {
	case ComponentTypeComponentReference:
		return c.component.Name()
	case ComponentTypeField:
		return c.field.FieldType()
	}
	panic("unknown component type")
}

func (c *ComponentDataType) HasComponentReference() bool {
	return c.variant == ComponentTypeComponentReference
}

func (c *ComponentDataType) ComponentDataType() *ComponentDataType {
	return c.component
}

func (c *ComponentDataType) String() string {
	var s string
	switch c.variant {
	case ComponentTypeComponentReference:
		s += fmt.Sprintf("[componenttype '%v']", c.component)
	case ComponentTypeField:
		s += fmt.Sprintf("[componenttype '%v']", c.field)
	}
	return s
}
