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

	"github.com/piot/scrawl-go/src/token"

	"strings"

	"github.com/piot/scrawl-go/src/definition"
	"github.com/piot/scrawl-go/src/runestream"
	"github.com/piot/scrawl-go/src/tokenize"
)

func setupTokenizer(text string) *tokenize.Tokenizer {
	ioReader := strings.NewReader(text)
	runeReader := runestream.NewRuneReader(ioReader)
	tokenizer := tokenize.NewTokenizer(runeReader)
	return tokenizer
}

type ParserError struct {
	err      error
	position token.Position
}

func (f ParserError) Error() string {
	return fmt.Sprintf("%v at %v", f.err, f.position)
}

type Parser struct {
	tokenizer  *tokenize.Tokenizer
	root       *definition.Root
	lastToken  token.Token
	lastEntity *definition.Entity
}

func (p *Parser) readNext() (token.Token, error) {
	token, err := p.tokenizer.ReadNext()
	if err != nil {
		return nil, err
	}
	p.lastToken = token
	return token, nil
}

func (p *Parser) next() (bool, error) {
	t, tokenErr := p.readNext()
	if t == nil {
		return true, nil
	}
	if tokenErr != nil {
		return false, tokenErr
	}
	symbolToken, wasSymbol := t.(token.SymbolToken)
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
		case "lod":
			_, err := p.parseLod(p.lastEntity)
			if err != nil {
				return false, err
			}
		case "entity":
			index := uint8(len(p.root.Entities()))
			entityIndex := definition.NewEntityIndex(index)
			entity, err := p.parseEntity(entityIndex)
			if err != nil {
				return false, err
			}
			p.root.AddEntity(entity)
			p.lastEntity = entity
		case "event":
			event, err := p.parseEvent()
			if err != nil {
				return false, err

			}
			p.root.AddEvent(event)
		case "command":
			method, err := p.parseCommand(definition.NewCommandTypeIndex(len(p.root.Commands())))
			if err != nil {
				return false, err
			}
			p.root.AddMethod(method)
		default:
			return false, fmt.Errorf("Unknown keyword %v", symbolToken)
		}
	} else {
		return false, fmt.Errorf("Unexpected token: %v", t)
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
			if parser.lastToken == nil {
				parserErr = ParserError{err: err, position: token.Position{}}
			} else {
				parserErr = ParserError{err: err, position: parser.lastToken.Position()}
			}
		} else {
			parserErr = err
		}
	}
	return parser, parserErr
}
