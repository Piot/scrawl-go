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

package parser

import (
	"testing"

	"github.com/piot/scrawl-go/src/writer"
)

func setup(x string) (*Parser, error) {
	return NewParser(x)
}

func TestIndentationSymbol(t *testing.T) {
	parser, err := setup(
		`

component Something
  hello	 int32

  type   Another

`)

	if err != nil {
		t.Error(err)
	}

	def := parser.Root()

	if len(def.Components()) != 1 {
		t.Errorf("Wrong number of components:%v", len(def.Components()))
	}

	firstComponent := def.Components()[0]
	firstComponentName := firstComponent.Name()

	if firstComponentName != "Something" {
		t.Errorf("Wrong component-name:%v", firstComponentName)
	}

	fields := def.Components()[0].Fields()
	fieldCount := len(fields)
	if fieldCount != 2 {
		t.Errorf("Wrong number of fields:%v", fieldCount)
	}

	if fields[0].Name() != "hello" {
		t.Errorf("Wrong field name:%v", fields[0].Name())
	}
	if fields[0].FieldType() != "int32" {
		t.Errorf("Wrong field name:%v", fields[0].FieldType())
	}
	if fields[1].Name() != "type" {
		t.Errorf("Wrong field name:%v", fields[1].Name())
	}
	if fields[1].FieldType() != "Another" {
		t.Errorf("Wrong field name:%v", fields[1].FieldType())
	}
}

func TestEverything(t *testing.T) {
	parser, err := setup(
		`
type SomeType
  hello	 int32
  type   Another

component AnotherComponent
  usingTheType int
  Name string

entity ThisISTheEntity2
  shouldTypeSomething Here
`)

	if err != nil {
		t.Error(err)
	}

	def := parser.Root()
	if len(def.Components()) != 1 {
		t.Errorf("Wrong number of components:%v", len(def.Components()))
	}

	firstComponent := def.Components()[0]
	firstComponentName := firstComponent.Name()
	if firstComponentName != "AnotherComponent" {
		t.Errorf("Wrong component-name:%v", firstComponentName)
	}

	fields := def.Components()[0].Fields()
	fieldCount := len(fields)
	if fieldCount != 2 {
		t.Errorf("Wrong number of fields:%v", fieldCount)
	}

	if fields[0].Name() != "usingTheType" {
		t.Errorf("Wrong field name:%v", fields[0].Name())
	}
	if fields[0].FieldType() != "int" {
		t.Errorf("Wrong field name:%v", fields[0].FieldType())
	}
	if fields[1].Name() != "Name" {
		t.Errorf("Wrong field name:%v", fields[1].Name())
	}
	if fields[1].FieldType() != "string" {
		t.Errorf("Wrong field name:%v", fields[1].FieldType())
	}

	writer.WriteCSharp(def)
}
