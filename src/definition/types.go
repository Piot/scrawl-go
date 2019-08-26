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
)

type Type struct {
	name string
}
type Hash uint32

type Root struct {
	components []*Component
	userTypes  []*UserType
	entities   []*Entity
	commands   []*Command
	events     []*Event
	enums      []*Enum
	hash       Hash
}

func (r *Root) SetHash(hash Hash) {
	r.hash = hash
}

func (r *Root) Hash() Hash {
	return r.hash
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

func (r *Root) FindUserType(name string) *UserType {
	for _, userType := range r.userTypes {
		if userType.name == name {
			return userType
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

func (r *Root) Enums() []*Enum {
	return r.enums
}

func (r *Root) Events() []*Event {
	return r.events
}

func (r *Root) Commands() []*Command {
	return r.commands
}

func (r *Root) Entities() []*Entity {
	return r.entities
}

func (r *Root) UserTypes() []*UserType {
	return r.userTypes
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

func (r *Root) AddEvent(c *Event) {
	r.events = append(r.events, c)
}

func (r *Root) AddMethod(c *Command) {
	r.commands = append(r.commands, c)
}

func (r *Root) AddEnum(c *Enum) {
	r.enums = append(r.enums, c)
}
