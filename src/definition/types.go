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
	"strconv"
)

type Type struct {
	name string
}

type Field struct {
	index     int
	name      string
	fieldType string
	metaData  MetaData
}

func NewField(index int, name string, fieldType string, metaData MetaData) *Field {
	return &Field{index: index, name: name, fieldType: fieldType, metaData: metaData}
}
func (c *Field) Index() int {
	return c.index
}

func (c *Field) Name() string {
	return c.name
}

func (c *Field) FieldType() string {
	return c.fieldType
}

func (c *Field) MetaData() MetaData {
	return c.metaData
}

func (c *Field) String() string {
	var s string
	s += fmt.Sprintf("[field '%v' %v]", c.name, c.fieldType)

	return s
}

type Component struct {
	name   string
	fields []*Field
}

func NewComponent(name string, fields []*Field) *Component {
	return &Component{name: name, fields: fields}
}

func (c *Component) Name() string {
	return c.name
}

func (c *Component) Fields() []*Field {
	return c.fields
}

func (c *Component) String() string {
	var s string
	s += fmt.Sprintf("[component '%v' fields:%d]\n", c.name, len(c.fields))
	for _, field := range c.fields {
		s += "    " + field.String() + "\n"
	}

	return s
}

type ComponentField struct {
	index     int
	name      string
	component *Component
}

func (c *ComponentField) Index() int {
	return c.index
}
func (c *ComponentField) Name() string {
	return c.name
}

func (c *ComponentField) Component() *Component {
	return c.component
}

func NewComponentField(index int, name string, component *Component) *ComponentField {
	return &ComponentField{index: index, name: name, component: component}
}

func (c *ComponentField) String() string {
	var s string
	s += fmt.Sprintf("[componentfield '%v' %s   ]", c.name, c.Component())
	return s
}

type Entity struct {
	name            string
	componentFields []*ComponentField
}

func NewEntity(name string, componentFields []*ComponentField) *Entity {
	return &Entity{name: name, componentFields: componentFields}
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

func (c *Entity) Component(index int) *ComponentField {
	return c.componentFields[index]
}

func (c *Entity) Components() []*ComponentField {
	return c.componentFields
}

type UserType struct {
	name   string
	fields []*Field
}

func NewUserType(name string, fields []*Field) *UserType {
	return &UserType{name: name, fields: fields}
}

type MetaData struct {
	Values map[string]string
}

func (m *MetaData) Int(name string) (int, error) {
	s := m.Values[name]
	return strconv.Atoi(s)
}

type Root struct {
	components []*Component
	userTypes  []*UserType
	entities   []*Entity
}

func (r *Root) FindComponent(name string) *Component {
	for _, component := range r.components {
		if component.Name() == name {
			return component
		}
	}
	return nil
}

func (r *Root) FindEntity(name string) *Entity {
	for _, entity := range r.entities {
		if entity.Name() == name {
			return entity
		}
	}
	return nil
}

func (r *Root) String() string {
	var s string

	s += fmt.Sprintf("Components: %d\n", len(r.components))
	s += fmt.Sprintf("Entities: %d\n", len(r.entities))
	s += fmt.Sprintf("UserTypes: %d\n", len(r.userTypes))

	for _, entity := range r.entities {
		s += entity.String() + "\n"
	}

	return s
}

func (r *Root) Components() []*Component {
	return r.components
}

func (r *Root) Entities() []*Entity {
	return r.entities
}

func (r *Root) AddComponent(c *Component) {
	r.components = append(r.components, c)
}

func (r *Root) AddUserType(c *UserType) {
	r.userTypes = append(r.userTypes, c)
}

func (r *Root) AddEntity(c *Entity) {
	r.entities = append(r.entities, c)
}
