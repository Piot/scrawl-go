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

	"github.com/piot/scrawl-go/src/beautify"

	"github.com/piot/scrawl-go/src/scrawlhash"
	"github.com/piot/scrawl-go/src/token"

	"github.com/piot/scrawl-go/src/definition"
	"github.com/piot/scrawl-go/src/tokenize"
)

type ParserError struct {
	err      error
	position token.Position
}

func (f ParserError) Error() string {
	return fmt.Sprintf("%v at %v", f.err, f.position)
}

type Parser struct {
	tokenizer            *tokenize.Tokenizer
	root                 *definition.Root
	lastToken            token.Token
	lastEntity           *definition.EntityArchetype
	validComponentTypes  []string
	validComponentFields []string
}

func (p *Parser) readNextEvenComments() (token.Token, error) {
	token, err := p.tokenizer.ReadNext()
	if err != nil {
		return nil, err
	}
	p.lastToken = token
	return token, nil
}

func (p *Parser) readNext() (token.Token, error) {
	for {
		foundToken, tokenErr := p.readNextEvenComments()
		if tokenErr != nil {
			return nil, tokenErr
		}
		_, isComment := foundToken.(token.CommentToken)
		if !isComment {
			return foundToken, tokenErr
		}
	}
}

func (p *Parser) next() (bool, error) {
	t, readErr := p.readNext()
	if readErr != nil {
		return false, readErr
	}
	if t == nil {
		return true, nil
	}
	symbolToken, wasSymbol := t.(token.SymbolToken)
	if wasSymbol {
		switch symbolToken.Symbol {
		case "component":
			index := uint8(len(p.root.ComponentDataTypes()))
			component, err := p.parseComponentDataType(index)
			if err != nil {
				return false, err
			}
			p.root.AddComponentDataType(component)
		case "type":
			userType, err := p.parseUserType()
			if err != nil {
				return false, err
			}
			p.root.AddUserType(userType)

		case "archetype":
			index := uint8(len(p.root.Archetypes()))
			entityIndex := definition.NewEntityIndex(index)
			entity, err := p.parseEntityArchetype(entityIndex)
			if err != nil {
				return false, err
			}
			p.root.AddArchetype(entity)
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
		case "enum":
			enum, err := p.parseEnum()
			if err != nil {
				return false, err
			}
			p.root.AddEnum(enum)

		default:
			return false, fmt.Errorf("Unknown keyword %v", symbolToken)
		}
	} else {
		_, isEndOfLine := t.(token.LineDelimiterToken)
		if isEndOfLine {
			return p.next()
		}
		return false, fmt.Errorf("Unexpected token: %v", t)
	}
	return false, nil
}

func (p *Parser) Root() *definition.Root {
	return p.root
}

func calculateHash(text string) (uint32, error) {
	hashTokens, fetchErr := tokenize.FetchAllTokens(text)
	if fetchErr != nil {
		return 0, fetchErr
	}
	beautified := beautify.Process(hashTokens, beautify.DiscardComments)
	//fmt.Printf("beautified:\n%s\n", beautified)
	hashValue := scrawlhash.CalculateHash([]byte(beautified))
	return hashValue, nil
}

func setHash(root *definition.Root, text string) error {
	hashValue, calculateErr := calculateHash(text)
	if calculateErr != nil {
		return calculateErr
	}
	definitionHash := definition.Hash(hashValue)
	root.SetHash(definitionHash)
	return nil
}

func NewParser(text string, allowedComponentFields []string, allowedComponentTypes []string) (*Parser, error) {
	tokenizer := tokenize.SetupTokenizer(text)
	parser := &Parser{tokenizer: tokenizer, root: &definition.Root{}, validComponentTypes: allowedComponentTypes,
		validComponentFields: allowedComponentFields}
	done := false
	var err error
	err = nil

	for !done && err == nil {
		done, err = parser.next()
	}
	if err != nil {
		var parserErr error
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

		return nil, parserErr
	}

	hashErr := setHash(parser.root, text)
	if hashErr != nil {
		return nil, hashErr
	}

	return parser, nil
}
