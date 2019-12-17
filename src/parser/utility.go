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

func (p *Parser) parseNameOptionalMetaAndFields() (string, definition.MetaData, []*definition.Field, error) {
	name, meta, nameErr := p.parseNameOptionalMetaAndStartScope()
	if nameErr != nil {
		return "", definition.MetaData{}, nil, nameErr
	}
	fields, fieldsErr := p.parseFieldsUntilEndScope()
	if fieldsErr != nil {
		return "", definition.MetaData{}, nil, fieldsErr
	}
	return name, meta, fields, nil
}

func (p *Parser) parseArchetypeNameAndStartScope() (string, definition.MetaData, error) {
	name, meta, nameErr := p.parseNameOptionalMetaAndStartScope()
	if nameErr != nil {
		return "", definition.MetaData{}, nameErr
	}
	return name, meta, nil
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

func (p *Parser) parseNameOptionalMetaAndStartScope() (string, definition.MetaData, error) {
	name, symbolErr := p.parseSymbol()
	if symbolErr != nil {
		return "", definition.MetaData{}, symbolErr
	}

	maybeMetaOrStartScope, tokErr := p.readNext()
	if tokErr != nil {
		return "", definition.MetaData{}, tokErr
	}

	_, isStartMeta := maybeMetaOrStartScope.(token.StartMetaDataToken)

	var metaData definition.MetaData

	if isStartMeta {
		var metaDataErr error

		metaData, metaDataErr = p.parseMetaData()
		if metaDataErr != nil {
			return "", definition.MetaData{}, metaDataErr
		}
		maybeMetaOrStartScope, metaDataErr = p.readNext()
		if metaDataErr != nil {
			return "", definition.MetaData{}, metaDataErr
		}
	}
	_, wasStartScope := maybeMetaOrStartScope.(token.StartScopeToken)
	if !wasStartScope {
		return "", definition.MetaData{}, fmt.Errorf("needs to have start scope after archetype name")
	}

	return name, metaData, nil
}

func (p *Parser) parseFieldsUntilEndScope() ([]*definition.Field, error) {
	var fields []*definition.Field

	for {
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
}

func convertEntityArchetypeItem(root *definition.Root, validComponentTypes []string, typeName string, itemIndex int, meta definition.MetaData) (*definition.EntityArchetypeItem, error) {
	componentDataTypeReference := root.FindComponentDataType(typeName)
	var archetypeItem *definition.EntityArchetypeItem
	if componentDataTypeReference == nil {
		if !Contains(validComponentTypes, typeName) {
			return nil, fmt.Errorf("unknown component type:%v", typeName)
		}
		archetypeItem = definition.NewEntityArchetypeItemUsingFieldType(itemIndex, typeName, meta)
	} else {
		archetypeItem = definition.NewEntityArchetypeItemUsingComponentDataTypeReference(componentDataTypeReference, meta)
	}
	return archetypeItem, nil
}

func (p *Parser) readMetaOrNewline() (definition.MetaData, bool, error) {
	hopefullyLineDelimiter, hopefullyLineDelimiterErr := p.readNext()
	if hopefullyLineDelimiterErr != nil {
		return definition.MetaData{}, false, hopefullyLineDelimiterErr
	}

	_, isStartMeta := hopefullyLineDelimiter.(token.StartMetaDataToken)

	var metaData definition.MetaData

	if isStartMeta {
		var metaDataErr error

		metaData, metaDataErr = p.parseMetaData()
		if metaDataErr != nil {
			return definition.MetaData{}, false, metaDataErr
		}
		hopefullyLineDelimiter, hopefullyLineDelimiterErr = p.readNext()
		if hopefullyLineDelimiterErr != nil {
			return definition.MetaData{}, false, hopefullyLineDelimiterErr
		}
	}
	token, wasEndOfLine := hopefullyLineDelimiter.(token.LineDelimiterToken)
	if !wasEndOfLine {
		return definition.MetaData{}, false, fmt.Errorf("must end lines or have meta information:%v", token)
	}

	return metaData, !isStartMeta, nil
}

func (p *Parser) symbolOrEndOfScope() (token.SymbolToken, bool, error) {
	t, tokenErr := p.readNext()
	if tokenErr != nil {
		return token.SymbolToken{}, false, tokenErr
	}

	switch ct := t.(type) {
	case token.EndScopeToken:
		return token.SymbolToken{}, true, nil
	case token.SymbolToken:
		return ct, false, nil
	}

	return token.SymbolToken{}, false, fmt.Errorf("expected symbol or end of scope %v %T", t, t)
}

func (p *Parser) parseEntityArchetypeItemsUntilEndScope() ([]*definition.EntityArchetypeItem, error) {
	var items []*definition.EntityArchetypeItem

	for {
		symbolToken, wasEndScope, symbolErr := p.symbolOrEndOfScope()
		if symbolErr != nil {
			return nil, symbolErr
		}

		if wasEndScope {
			return items, nil
		}

		parsedField, parseFieldErr := p.parseEntityArchetypeItem(len(items), symbolToken.Symbol)

		if parseFieldErr != nil {
			return nil, parseFieldErr
		}

		items = append(items, parsedField)
	}
}

func (p *Parser) parseEnumConstantsUntilEndScope() ([]*definition.EnumConstant, error) {
	var fields []*definition.EnumConstant

	for {
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
}
