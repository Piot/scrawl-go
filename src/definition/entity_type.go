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
	"hash/fnv"
)

// EntityIndex :
type EntityIndex struct {
	id uint8
}

func NewEntityIndex(id uint8) EntityIndex {
	return EntityIndex{id: id}
}

func (e EntityIndex) Value() uint8 {
	return e.id
}

func (e EntityIndex) String() string {
	return fmt.Sprintf("[entityindex %v]", e.id)
}

// EntityTypeID :
type EntityTypeID struct {
	id uint16
}

func (e EntityTypeID) Value() uint16 {
	return e.id
}

func (e EntityTypeID) String() string {
	return fmt.Sprintf("[entitytypeid %v]", e.id)
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func typeHash(name string) uint16 {
	v := hash(name)
	w := uint16((v >> 16) ^ (v & 0xffff))
	return w
}

func NewEntityTypeIDFromString(name string) EntityTypeID {
	return EntityTypeID{id: typeHash(name)}
}
