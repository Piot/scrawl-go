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

type Type struct {
	name string
}

type Field struct {
	name      string
	fieldType string
}

func NewField(name string, fieldType string) *Field {
	return &Field{name: name, fieldType: fieldType}
}

func (c *Field) Name() string {
	return c.name
}

func (c *Field) FieldType() string {
	return c.fieldType
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

type Entity struct {
	name   string
	fields []*Field
}

func NewEntity(name string, fields []*Field) *Entity {
	return &Entity{name: name, fields: fields}
}

type UserType struct {
	name   string
	fields []*Field
}

func NewUserType(name string, fields []*Field) *UserType {
	return &UserType{name: name, fields: fields}
}

type Root struct {
	components []*Component
	userTypes  []*UserType
	entities   []*Entity
}

func (r *Root) Components() []*Component {
	return r.components
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
