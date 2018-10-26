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

func (p *Parser) parseCommandNameParameterAndReturn() (string, string, string, error) {
	name, symbolErr := p.parseSymbol()
	if symbolErr != nil {
		return "", "", "", symbolErr
	}
	parameter, parameterErr := p.parseSymbol()
	if parameterErr != nil {
		return "", "", "", parameterErr
	}
	returnType, returnTypeErr := p.parseSymbol()
	if returnTypeErr != nil {
		return "", "", "", returnTypeErr
	}

	hopefullyLineDelimiter, hopefullyLineDelimiterErr := p.readNext()
	if hopefullyLineDelimiterErr != nil {
		return "", "", "", hopefullyLineDelimiterErr
	}

	_, wasEndOfLine := hopefullyLineDelimiter.(token.LineDelimiterToken)
	if !wasEndOfLine {
		return "", "", "", fmt.Errorf("Must end lines after field type")
	}

	return name, parameter, returnType, nil
}

func (p *Parser) parseCommand(index definition.CommandTypeIndex) (*definition.Command, error) {
	name, parameter, returnType, err := p.parseCommandNameParameterAndReturn()
	if err != nil {
		return nil, err
	}
	method := definition.NewCommand(index, name, parameter, returnType)
	return method, nil
}
