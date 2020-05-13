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
	return NewParser(x, []string{"WorldPosition"}, []string{"WorldPositionComponent"})
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

	if len(def.ComponentDataTypes()) != 1 {
		t.Errorf("Wrong number of components:%v", len(def.ComponentDataTypes()))
	}

	firstComponent := def.ComponentDataTypes()[0]
	firstComponentName := firstComponent.Name()

	if firstComponentName != "Something" {
		t.Errorf("Wrong component-name:%v", firstComponentName)
	}

	fields := def.ComponentDataTypes()[0].Fields()
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
	components := setup.Root().ComponentDataTypes()
	if len(components) != 1 {
		t.Errorf("Should have been 1")
	}

	if components[0].Name() != "SomeType" {
		t.Errorf("Wrong name")
	}
}

func TestEmptyComponent(t *testing.T) {
	setup, err := setup(
		`
component EmptyComponent 
  # Intentionally empty 
`)
	if err != nil {
		t.Error(err)
	}
	components := setup.Root().ComponentDataTypes()
	if len(components) != 1 {
		t.Errorf("Should have been 1")
	}

	if components[0].Name() != "EmptyComponent" {
		t.Errorf("Wrong name")
	}
}

func TestAnotherEmptyComponent(t *testing.T) {
	setup, err := setup(
		`
component EmptyComponent 
  # Intentionally empty

component SomethingElse
  # Empty
`)
	if err != nil {
		t.Error(err)
	}
	components := setup.Root().ComponentDataTypes()
	if len(components) != 2 {
		t.Errorf("Should have been 2")
	}

	if components[0].Name() != "EmptyComponent" {
		t.Errorf("Wrong name")
	}
}

func TestCommentField(t *testing.T) {
	setup, err := setup(
		`
component EmptyComponent 
  speed Float # this is a comment that should be ignored

`)
	if err != nil {
		t.Error(err)
	}
	components := setup.Root().ComponentDataTypes()
	if len(components) != 1 {
		t.Errorf("Should have been 1")
	}

	if components[0].Name() != "EmptyComponent" {
		t.Errorf("Wrong name")
	}
}

func TestCommentLines(t *testing.T) {
	setup, err := setup(
		`
component EmptyComponent 
  speed Float # this is a comment that should be ignored

# Ignore this
# and this

`)
	if err != nil {
		t.Error(err)
	}
	components := setup.Root().ComponentDataTypes()
	if len(components) != 1 {
		t.Errorf("Should have been 1")
	}

	if components[0].Name() != "EmptyComponent" {
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
  usingTheType int [range '34-45', debug 'hello world']
  Name string

archetype ThisISTheEntity2
  lod 0
    AnotherComponent

event Jump
  hello int32

command Fire
  target String

buffer Tile
  index int32

type ReturnType
  hello	 int32
  type   Another
`)

	if err != nil {
		t.Error(err)
	}

	def := parser.Root()
	if len(def.ComponentDataTypes()) != 1 {
		t.Errorf("Wrong number of components:%v", len(def.ComponentDataTypes()))
	}

	event := def.Events()[0]
	if event.Name() != "Jump" {
		t.Errorf("Wrong event name:%v", event)
	}

	cmd := def.Commands()[0]
	if cmd.Name() != "Fire" {
		t.Errorf("Wrong command name:%v", event)
	}

	cmdField := cmd.Fields()[0]
	if cmdField.FieldType() != "String" {
		t.Errorf("Wrong command field:%v", cmdField)
	}

	buffer := def.Buffers()[0]
	if buffer.Name() != "Tile" {
		t.Errorf("wrong buffer %v", buffer)
	}

	bufferField := buffer.Fields()[0]
	if bufferField.FieldType() != "int32" {
		t.Errorf("wrong field type for buffer %v", bufferField.FieldType())
	}
	firstComponent := def.ComponentDataTypes()[0]
	firstComponentName := firstComponent.Name()
	if firstComponentName != "AnotherComponent" {
		t.Errorf("Wrong component-name:%v", firstComponentName)
	}

	fields := def.ComponentDataTypes()[0].Fields()
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
}

func TestEnum(t *testing.T) {
	parser, err := setup(
		`
enum MovementState 
  Attacking 1
  Idle 0

component AnotherComponent
  usingTheType int [range "34-45", debug "hello world"]
  Name string
  state MovementState
`)

	if err != nil {
		t.Error(err)
	}

	def := parser.Root()
	if len(def.Enums()) != 1 {
		t.Errorf("Wrong number of enums:%v", len(def.Enums()))
	}
	firstEnum := def.Enums()[0]
	if firstEnum.Name() != "MovementState" {
		t.Errorf("wrong enum name")
	}

	firstConstant := firstEnum.Constants()[0]
	if firstConstant.Name() != "Attacking" {
		t.Errorf("wrong constant name")
	}

	if firstConstant.Value() != 1 {
		t.Errorf("wrong constant value")
	}
	if firstConstant.Index() != 0 {
		t.Errorf("wrong constant value")
	}
	secondConstant := firstEnum.Constants()[1]
	if secondConstant.Name() != "Idle" {
		t.Errorf("wrong constant name")
	}

	if secondConstant.Value() != 0 {
		t.Errorf("wrong constant value")
	}
	if secondConstant.Index() != 1 {
		t.Errorf("wrong constant value")
	}
}

func TestNamespace(t *testing.T) {
	parser, err := setup(
		`
namespace Some.Namespace.Here

archetype ThisISTheEntity2
  lod 0
    WorldPosition  [range   "0-122", priority "low"]
  lod 1
    WorldPosition  [range   "0-122"]
`)

	if err != nil {
		t.Fatal(err)
	}

	def := parser.Root()
	if def.Namespace() != "Some.Namespace.Here" {
		t.Fatalf("wrong namespace")
	}
}

func TestName(t *testing.T) {
	parser, err := setup(
		`
name MyCoolName

namespace Some.Namespace.Here

archetype ThisISTheEntity2
  lod 0
    WorldPosition  [range   "0-122", priority "low"]
  lod 1
    WorldPosition  [range   "0-122"]
`)

	if err != nil {
		t.Fatal(err)
	}

	def := parser.Root()
	if def.Namespace() != "Some.Namespace.Here" {
		t.Fatalf("wrong namespace")
	}

	if def.Name() != "MyCoolName" {
		t.Fatalf("wrong name")
	}
}

func TestTypeInsteadOfComponent(t *testing.T) {
	parser, err := setup(
		`
component SomeOtherComponent [priority "high"]
  something Integer

archetype ThisISTheEntity2
  lod 0
    WorldPosition  [range   "0-122", priority "low"]
    SomeOtherComponent [priority "mid"]
  lod 1
    WorldPosition  [range   "0-122"]
`)

	if err != nil {
		t.Fatal(err)
	}

	def := parser.Root()

	firstComponent := def.ComponentDataTypes()[0]
	if firstComponent.Name() != "SomeOtherComponent" {
		t.Errorf("wrong component name:%v", firstComponent)
	}

	firstMeta := firstComponent.Meta()
	firstPriority := firstMeta.Field("priority")
	if firstPriority != "high" {
		t.Errorf("wrong priority %v", firstMeta)
	}

	archetypeToTest := def.Archetypes()[0]
	firstLod := archetypeToTest.HighestLevelOfDetail()
	secondComponentData := firstLod.EntityArchetypeItem(1)
	raw := secondComponentData.Name()

	secondItemInFirstLod := firstLod.Items()[1]
	secondItemMeta := secondItemInFirstLod.Meta()
	if secondItemMeta.Field("priority") != "mid" {
		t.Errorf("wrong priority :%v", secondItemMeta)
	}

	if raw != "SomeOtherComponent" {
		t.Errorf("wrong type %v", raw)
	}

	secondLod := archetypeToTest.Lod(1)
	secondLodItem := secondLod.EntityArchetypeItem(0)
	raw2 := secondLodItem.Name()
	if raw2 != "WorldPosition" {
		t.Errorf("wrong type %v", raw2)
	}
}
