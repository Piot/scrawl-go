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

func TestIgnoringCarriageReturn(t *testing.T) {
	setup, err := setup(
		"component SomeType\r\n" +
			"  hello	 int32\r\n" +
			"  type \r  Another\r\n\r" +
			"")
	if err != nil {
		t.Error(err)
	}
	components := setup.Root().Components()
	if len(components) != 1 {
		t.Errorf("Should have been 1")
	}

	if components[0].Name() != "SomeType" {
		t.Errorf("Wrong name")
	}
}

func TestEverything(t *testing.T) {
	parser, err := setup(
		`
type SomeType
  hello	 int32
  type   Another

component AnotherComponent
  usingTheType int [range "34-45", debug "hello world"]
  Name string
  event Jump

component SomeOtherComponent
  event Jump

entity ThisISTheEntity2
  shouldTypeSomething AnotherComponent

event Jump
  where Position
  energy int [range "10-34"]

type ReturnType
  hello	 int32
  type   Another

command CommandName SomeType ReturnType
`)

	if err != nil {
		t.Error(err)
	}

	def := parser.Root()
	if len(def.Components()) != 2 {
		t.Errorf("Wrong number of components:%v", len(def.Components()))
	}

	event := def.Events()[0]
	if event.Name() != "Jump" {
		t.Errorf("Wrong event name:%v", event)
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
	metaDataValues := fields[0].MetaData().Values
	if metaDataValues["range"] != "34-45" {
		t.Errorf("Wrong meta data name:%v", metaDataValues["range"])
	}
	if metaDataValues["debug"] != "hello world" {
		t.Errorf("Wrong meta data name:%v", metaDataValues["debug"])
	}

	if fields[1].Name() != "Name" {
		t.Errorf("Wrong field name:%v", fields[1].Name())
	}
	if fields[1].FieldType() != "string" {
		t.Errorf("Wrong field name:%v", fields[1].FieldType())
	}

	eventReferences := def.Components()[0].EventReferences()
	if len(eventReferences) != 1 {
		t.Errorf("Must have one event reference")
	}

	firstRef := eventReferences[0]
	if firstRef.ReferencedType() != "Jump" {
		t.Errorf("should be Jump")
	}

	if firstRef.ReferencedType() != "Jump" {
		t.Errorf("wrong referenced type")
	}

	cmd := def.Commands()[0]
	if cmd.Name() != "CommandName" {
		t.Errorf("wrong command name %v", cmd)
	}

	if cmd.ParameterTypeName() != "SomeType" {
		t.Errorf("Wrong argument type %v", cmd)
	}

	if cmd.ReturnTypeName() != "ReturnType" {
		t.Errorf("wrong command return %v", cmd)
	}
}
