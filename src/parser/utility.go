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
	"fmt"

	"github.com/piot/scrawl-go/src/definition"
	"github.com/piot/scrawl-go/src/token"
)

func (p *Parser) parseNameAndFields() (string, []*definition.Field, error) {
	name, nameErr := p.parseNameAndStartScope()
	if nameErr != nil {
		return "", nil, nameErr
	}
	fields, fieldsErr := p.parseFieldsUntilEndScope()
	if fieldsErr != nil {
		return "", nil, fieldsErr
	}
	return name, fields, nil
}

func (p *Parser) parseArchetypeNameAndItems() (string, []*definition.EntityArchetypeItem, error) {
	name, nameErr := p.parseNameAndStartScope()
	if nameErr != nil {
		return "", nil, nameErr
	}
	fields, fieldsErr := p.parseEntityArchetypeItemsUntilEndScope()
	if fieldsErr != nil {
		return "", nil, fieldsErr
	}
	return name, fields, nil
}

func (p *Parser) parseIntegerAndFields() (int, []*definition.Field, error) {
	name, nameErr := p.parseIntegerAndStartScope()
	if nameErr != nil {
		return 0, nil, nameErr
	}
	fields, fieldsErr := p.parseFieldsUntilEndScope()
	if fieldsErr != nil {
		return 0, nil, fieldsErr
	}
	return name, fields, nil
}

func (p *Parser) parseIntegerAndStartScope() (int, error) {
	number, symbolErr := p.parseInteger()
	if symbolErr != nil {
		return 0, symbolErr
	}
	startScopeErr := p.parseStartScope()
	if startScopeErr != nil {
		return 0, startScopeErr
	}

	return number, nil
}

func (p *Parser) parseNameAndStartScope() (string, error) {
	name, symbolErr := p.parseSymbol()
	if symbolErr != nil {
		return "", symbolErr
	}
	fmt.Printf("name:%v\n", name)
	startScopeErr := p.parseStartScope()
	if startScopeErr != nil {
		return "", startScopeErr
	}

	return name, nil
}

func (p *Parser) parseFieldsUntilEndScope() ([]*definition.Field, error) {
	var fields []*definition.Field

	for true {
		t, tokenErr := p.readNext()
		if tokenErr != nil {
			return nil, tokenErr
		}
		symbolToken, wasSymbol := t.(token.SymbolToken)
		if !wasSymbol {
			_, wasEndScope := t.(token.EndScopeToken)
			if wasEndScope {
				return fields, nil
			}
			return nil, fmt.Errorf("Expected fieldname or end of scope")
		}

		parsedField, parseFieldErr := p.parseField(len(fields), symbolToken.Symbol)
		if parseFieldErr != nil {
			return nil, parseFieldErr
		}
		fields = append(fields, parsedField)
	}
	return nil, nil
}

func convertEntityArchetypeItem(root *definition.Root, validComponentTypes []string, typeName string) (*definition.EntityArchetypeItem, error) {
	componentDataTypeReference := root.FindComponentDataType(typeName)
	var archetypeItem *definition.EntityArchetypeItem
	if componentDataTypeReference == nil {
		if !Contains(validComponentTypes, typeName) {
			return nil, fmt.Errorf("unknown component type:%v", typeName)
		}
		archetypeItem = definition.NewEntityArchetypeItemUsingFieldType(typeName)
	} else {
		archetypeItem = definition.NewEntityArchetypeItemUsingComponentDataTypeReference(componentDataTypeReference)
	}
	return archetypeItem, nil
}

func (p *Parser) parseEntityArchetypeItemsUntilEndScope() ([]*definition.EntityArchetypeItem, error) {
	var items []*definition.EntityArchetypeItem

	for true {
		t, tokenErr := p.readNext()
		if tokenErr != nil {
			return nil, tokenErr
		}
		symbolToken, wasSymbol := t.(token.SymbolToken)
		if !wasSymbol {
			_, wasEndScope := t.(token.EndScopeToken)
			if wasEndScope {
				return items, nil
			}
			return nil, fmt.Errorf("Expected component data type field reference or end of scope")
		}
		fmt.Printf("checking symbol:%v\n", symbolToken)

		parsedField, parseFieldErr := p.parseEntityArchetypeItem(len(items), symbolToken.Symbol)
		if parseFieldErr != nil {
			return nil, parseFieldErr
		}
		items = append(items, parsedField)
	}
	return nil, nil
}

func (p *Parser) parseEnumConstantsUntilEndScope() ([]*definition.EnumConstant, error) {
	var fields []*definition.EnumConstant

	for true {
		t, tokenErr := p.readNext()
		if tokenErr != nil {
			return nil, tokenErr
		}
		symbolToken, wasSymbol := t.(token.SymbolToken)
		if !wasSymbol {
			_, wasEndScope := t.(token.EndScopeToken)
			if wasEndScope {
				return fields, nil
			}
			return nil, fmt.Errorf("Expected enum name or end of scope %v", t)
		}

		t, tokenErr = p.readNext()
		if tokenErr != nil {
			return nil, tokenErr
		}
		numberToken, wasNumber := t.(token.NumberToken)
		if !wasNumber {
			return nil, fmt.Errorf("Expected enum constant value")
		}
		hopefullyLineDelimiterErr := p.expectLineDelimiter()
		if hopefullyLineDelimiterErr != nil {
			return nil, fmt.Errorf("enum constants:%v", hopefullyLineDelimiterErr)
		}

		index := len(fields)
		enumConstant := definition.NewEnumConstant(index, symbolToken.Symbol, numberToken.Integer(), nil)
		fields = append(fields, enumConstant)
	}
	return nil, nil
}
