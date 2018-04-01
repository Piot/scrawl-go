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

func (p *Parser) parseField(index int, name string) (*definition.Field, error) {
	fieldType, fieldTypeErr := p.parseSymbol()
	if fieldTypeErr != nil {
		return &definition.Field{}, fmt.Errorf("Expected a field symbol (%v)", fieldTypeErr)
	}

	hopefullyLineDelimiter, hopefullyLineDelimiterErr := p.readNext()
	if hopefullyLineDelimiterErr != nil {
		return nil, hopefullyLineDelimiterErr
	}

	_, isStartMeta := hopefullyLineDelimiter.(token.StartMetaDataToken)
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
	_, wasEndOfLine := hopefullyLineDelimiter.(token.LineDelimiterToken)
	if !wasEndOfLine {
		return nil, fmt.Errorf("Must end lines after field type")
	}

	field := definition.NewField(index, name, fieldType, metaData)
	return field, nil
}
