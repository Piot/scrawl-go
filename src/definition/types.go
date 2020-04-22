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

type Hash uint32

type Root struct {
	componentDataTypes []*ComponentDataType
	userTypes          []*UserType
	archetypes         []*EntityArchetype
	commands           []*Command
	events             []*Event
	enums              []*Enum
	hash               Hash
	namespace          string
}

func (r *Root) SetHash(hash Hash) {
	r.hash = hash
}

func (r *Root) SetNamespace(namespace string) {
	r.namespace = namespace
}

func (r *Root) Hash() Hash {
	return r.hash
}

func (r *Root) FindComponentDataType(name string) *ComponentDataType {
	for _, component := range r.componentDataTypes {
		if component.Name() == name {
			return component
		}
	}
	return nil
}

func (r *Root) FindEntity(name string) *EntityArchetype {
	for _, entity := range r.archetypes {
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

	s += fmt.Sprintf("ComponentDataType-types: %d\n", len(r.componentDataTypes))
	s += fmt.Sprintf("Archetypes: %d\n", len(r.archetypes))
	s += fmt.Sprintf("User-types: %d\n", len(r.userTypes))

	for _, entity := range r.archetypes {
		s += entity.String() + "\n"
	}

	return s
}

func (r *Root) Namespace() string {
	return r.namespace
}

func (r *Root) ComponentDataTypes() []*ComponentDataType {
	return r.componentDataTypes
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

func (r *Root) Archetypes() []*EntityArchetype {
	return r.archetypes
}

func (r *Root) UserTypes() []*UserType {
	return r.userTypes
}

func (r *Root) AddComponentDataType(c *ComponentDataType) {
	r.componentDataTypes = append(r.componentDataTypes, c)
}

func (r *Root) AddUserType(c *UserType) {
	r.userTypes = append(r.userTypes, c)
}

func (r *Root) AddArchetype(c *EntityArchetype) {
	r.archetypes = append(r.archetypes, c)
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
