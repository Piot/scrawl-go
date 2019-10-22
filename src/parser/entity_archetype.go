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

func (p *Parser) parseEntityArchetype(entityIndex definition.EntityIndex,
	validComponentTypes []string) (*definition.EntityArchetype, error) {

	name, nameErr := p.parseArchetypeNameAndStartScope()
	if nameErr != nil {
		return nil, nameErr
	}

	lods := make(map[int]*definition.EntityArchetypeLOD)
	expectedLevel := 0

	for {
		t, tErr := p.readNext()
		if tErr != nil {
			return nil, tErr
		}

		_, isEndOfScope := t.(token.EndScopeToken)
		if isEndOfScope {
			break
		}

		symbol, isSymbol := t.(token.SymbolToken)
		if !isSymbol {
			return nil, fmt.Errorf("expected end of scope or 'lod' %v", t)
		}

		if symbol.Symbol != "lod" {
			return nil, fmt.Errorf("expected 'lod' %v", symbol)
		}

		lod, err := p.parseLod(p.validComponentTypes)
		if err != nil {
			return nil, err
		}

		if lod.Level() != expectedLevel {
			return nil, fmt.Errorf("wrong lod:%v", lod)
		}

		lods[lod.Level()] = lod
		expectedLevel++
	}

	entity := definition.NewEntityArchetype(name, entityIndex, lods)
	return entity, nil
}
