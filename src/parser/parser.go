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
	"strings"

	"github.com/piot/scrawl-go/src/definition"
	"github.com/piot/scrawl-go/src/tokenize"
)

func setupTokenizer(text string) *tokenize.Tokenizer {
	ioReader := strings.NewReader(text)
	runeReader := tokenize.NewRuneReader(ioReader)
	tokenizer := tokenize.NewTokenizer(runeReader)
	return tokenizer
}

type ParserError struct {
	err      error
	position tokenize.Position
}

func (f ParserError) Error() string {
	return fmt.Sprintf("%v at %v", f.err, f.position)
}

type Parser struct {
	tokenizer *tokenize.Tokenizer
	root      *definition.Root
	lastToken tokenize.Token
}

func (p *Parser) readNext() (tokenize.Token, error) {
	token, err := p.tokenizer.ReadNext()
	if err != nil {
		return nil, err
	}
	p.lastToken = token
	return token, nil
}

func (p *Parser) parseSymbol() (string, error) {
	token, tokenErr := p.readNext()
	if tokenErr != nil {
		return "", tokenErr
	}
	symbolToken, wasSymbol := token.(tokenize.SymbolToken)
	if !wasSymbol {
		return "", fmt.Errorf("Wasn't a symbol")
	}

	return symbolToken.Symbol, nil
}

func (p *Parser) parseString() (string, error) {
	token, tokenErr := p.readNext()
	if tokenErr != nil {
		return "", tokenErr
	}
	stringToken, wasString := token.(tokenize.StringToken)
	if !wasString {
		return "", fmt.Errorf("Wasn't a string")
	}

	return stringToken.Text(), nil
}

func (p *Parser) parseMetaData() (definition.MetaData, error) {
	metaData := definition.MetaData{Values: make(map[string]string)}
	for true {
		token, tokenErr := p.readNext()
		if tokenErr != nil {
			return definition.MetaData{}, tokenErr
		}
		symbolToken, wasSymbol := token.(tokenize.SymbolToken)
		if wasSymbol {
			metaName := symbolToken.Symbol
			metaValue, metaValueErr := p.parseString()
			if metaValueErr != nil {
				return definition.MetaData{}, fmt.Errorf("Expected a meta value (%v)", metaValueErr)
			}
			metaData.Values[metaName] = metaValue
		} else {
			_, wasEndOfMetaData := token.(tokenize.EndMetaDataToken)
			if !wasEndOfMetaData {
				return definition.MetaData{}, fmt.Errorf("Expected a meta value or end of meta (%v)", token)
			}
			break
		}
	}

	return metaData, nil
}

func (p *Parser) parseField(index int, name string) (*definition.Field, error) {
	fieldType, fieldTypeErr := p.parseSymbol()
	if fieldTypeErr != nil {
		return &definition.Field{}, fmt.Errorf("Expected a field symbol (%v)", fieldTypeErr)
	}

	hopefullyLineDelimiter, hopefullyLineDelimiterErr := p.readNext()
	if hopefullyLineDelimiterErr != nil {
		return nil, hopefullyLineDelimiterErr
	}

	_, isStartMeta := hopefullyLineDelimiter.(tokenize.StartMetaDataToken)
	var metaData definition.MetaData
	if isStartMeta {
		var metaDataErr error
		metaData, metaDataErr = p.parseMetaData()
		if metaDataErr != nil {
			return nil, metaDataErr
		}
		hopefullyLineDelimiter, hopefullyLineDelimiterErr = p.readNext()
		if hopefullyLineDelimiterErr != nil {
			return nil, hopefullyLineDelimiterErr
		}
	}
	_, wasEndOfLine := hopefullyLineDelimiter.(tokenize.LineDelimiterToken)
	if !wasEndOfLine {
		return nil, fmt.Errorf("Must end lines after field type")
	}

	field := definition.NewField(index, name, fieldType, metaData)
	return field, nil
}

func (p *Parser) parseFieldsUntilEndScope() ([]*definition.Field, error) {
	var fields []*definition.Field

	for true {
		token, tokenErr := p.readNext()
		if tokenErr != nil {
			return nil, tokenErr
		}
		symbolToken, wasSymbol := token.(tokenize.SymbolToken)
		if !wasSymbol {
			_, wasEndScope := token.(tokenize.EndScopeToken)
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

func (p *Parser) parseStartScope() error {
	maybeStartScope, maybeStartScopeErr := p.readNext()
	if maybeStartScopeErr != nil {
		return maybeStartScopeErr
	}

	_, wasStartScope := maybeStartScope.(tokenize.StartScopeToken)
	if !wasStartScope {
		return fmt.Errorf("Missing start scope")
	}

	return nil
}

func (p *Parser) parseNameAndStartScope() (string, error) {
	name, symbolErr := p.parseSymbol()
	if symbolErr != nil {
		return "", symbolErr
	}
	startScopeErr := p.parseStartScope()
	if startScopeErr != nil {
		return "", startScopeErr
	}

	return name, nil
}

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

func (p *Parser) parseComponent() (*definition.Component, error) {
	name, fields, err := p.parseNameAndFields()
	if err != nil {
		return nil, err
	}
	component := definition.NewComponent(name, fields)
	return component, nil
}

func (p *Parser) parseUserType() (*definition.UserType, error) {
	name, fields, err := p.parseNameAndFields()
	if err != nil {
		return nil, err
	}
	userType := definition.NewUserType(name, fields)
	return userType, nil
}

func (p *Parser) parseEntity() (*definition.Entity, error) {
	name, fields, err := p.parseNameAndFields()
	if err != nil {
		return nil, err
	}
	var componentFields []*definition.ComponentField
	for _, fieldComponent := range fields {
		componentType := p.root.FindComponent(fieldComponent.FieldType())
		componentField := definition.NewComponentField(len(componentFields), fieldComponent.Name(), componentType)
		componentFields = append(componentFields, componentField)
	}
	entity := definition.NewEntity(name, componentFields)
	return entity, nil
}

func (p *Parser) next() (bool, error) {
	token, tokenErr := p.readNext()
	if token == nil {
		return true, nil
	}
	if tokenErr != nil {
		return false, tokenErr
	}
	symbolToken, wasSymbol := token.(tokenize.SymbolToken)
	if wasSymbol {
		switch symbolToken.Symbol {
		case "component":
			component, err := p.parseComponent()
			if err != nil {
				return false, err
			}
			p.root.AddComponent(component)
		case "type":
			userType, err := p.parseUserType()
			if err != nil {
				return false, err
			}
			p.root.AddUserType(userType)
		case "entity":
			entity, err := p.parseEntity()
			if err != nil {
				return false, err
			}
			p.root.AddEntity(entity)
		default:
			return false, fmt.Errorf("Unknown keyword %v", symbolToken)
		}
	} else {
		return false, fmt.Errorf("Unexpected token: %v", token)
	}
	return false, nil
}

func (p *Parser) Root() *definition.Root {
	return p.root
}

func NewParser(text string) (*Parser, error) {
	tokenizer := setupTokenizer(text)
	parser := &Parser{tokenizer: tokenizer, root: &definition.Root{}}
	done := false
	var err error
	err = nil
	for !done && err == nil {
		done, err = parser.next()
	}
	var parserErr error
	if err != nil {
		_, wasTokenizerError := err.(tokenize.TokenizerError)
		if !wasTokenizerError {
			parserErr = ParserError{err: err, position: parser.lastToken.Position()}
		} else {
			parserErr = err
		}
	}
	return parser, parserErr
}
