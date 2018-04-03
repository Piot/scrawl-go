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

package writer

import (
	"fmt"

	"github.com/piot/scrawl-go/src/definition"
)

func WriteCSharp(root *definition.Root) {
	for _, component := range root.Components() {
		fmt.Printf("public class %s \n{\n", component.Name())
		for _, field := range component.Fields() {
			fmt.Printf(" public %s %s;\n", field.FieldType(), field.Name())
		}

		fmt.Printf("}\n")
	}

	for _, entity := range root.Entities() {
		fmt.Printf("public class %s\n{\n", entity.Name())
		for _, field := range entity.Components() {
			fmt.Printf(" public %s %s;\n", field.Component().Name(), field.Name())
		}
		fmt.Printf("}\n")
	}
}
