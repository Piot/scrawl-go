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

type Event struct {
	id     EventTypeIndex
	name   string
	meta   MetaData
	fields []*Field
}

func NewEvent(id EventTypeIndex, name string, meta MetaData, fields []*Field) *Event {
	return &Event{id: id, name: name, meta: meta, fields: fields}
}

func (e *Event) TypeIndex() EventTypeIndex {
	return e.id
}

func (e *Event) Name() string {
	return e.name
}

func (e *Event) Meta() MetaData {
	return e.meta
}

func (e *Event) Fields() []*Field {
	return e.fields
}

func (e *Event) String() string {
	var s string

	s += fmt.Sprintf("[event %v]", e.fields)

	return s
}
