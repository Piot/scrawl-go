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

type EntityArchetypeItemVariant uint8

const (
	ComponentTypeComponentReference EntityArchetypeItemVariant = iota
	ComponentTypeField
)

type EntityArchetypeItem struct {
	componentDataType *ComponentDataType
	fieldType         string
	variant           EntityArchetypeItemVariant
}

func NewEntityArchetypeItemUsingFieldType(fieldType string) *EntityArchetypeItem {
	return &EntityArchetypeItem{variant: ComponentTypeField, fieldType: fieldType}
}

func NewEntityArchetypeItemUsingComponentDataTypeReference(componentDataType *ComponentDataType) *EntityArchetypeItem {
	return &EntityArchetypeItem{variant: ComponentTypeComponentReference, componentDataType: componentDataType}
}

func (c *EntityArchetypeItem) HasFieldReference() bool {
	return c.variant == ComponentTypeField
}

func (c *EntityArchetypeItem) FieldReference() string {
	if c.variant != ComponentTypeField {
		panic("wrong component field type")
	}
	return c.fieldType
}

func (c *EntityArchetypeItem) Name() string {
	switch c.variant {
	case ComponentTypeComponentReference:
		return c.componentDataType.Name()
	case ComponentTypeField:
		return c.fieldType
	}
	panic("unknown component type")
}

func (c *EntityArchetypeItem) HasComponentReference() bool {
	return c.variant == ComponentTypeComponentReference
}

func (c *EntityArchetypeItem) ComponentDataType() *ComponentDataType {
	if c.variant != ComponentTypeComponentReference {
		panic("wrong component field type")
	}
	return c.componentDataType
}

func (c *EntityArchetypeItem) String() string {
	var s string
	switch c.variant {
	case ComponentTypeComponentReference:
		s += fmt.Sprintf("[componenttype '%v']", c.componentDataType)
	case ComponentTypeField:
		s += fmt.Sprintf("[componenttype '%v']", c.fieldType)
	}
	return s
}
