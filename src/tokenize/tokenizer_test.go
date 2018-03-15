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

package tokenize

import (
	"strings"
	"testing"
)

func setup(x string) *Tokenizer {
	ioReader := strings.NewReader(x)
	runeReader := NewRuneReader(ioReader)
	tokenizer := NewTokenizer(runeReader)
	return tokenizer
}

func TestSimpleSymbol(t *testing.T) {
	tokenizer := setup("          someSymbol")
	expectSymbol(t, tokenizer, "someSymbol")
	expectEnd(t, tokenizer)
}

func TestSimpleString(t *testing.T) {
	tokenizer := setup("          'this is a test' ")
	expectString(t, tokenizer, "this is a test")
	expectEnd(t, tokenizer)
}

func expectSymbol(t *testing.T, tokenizer *Tokenizer, expectedString string) {
	hopefullySymbolToken, hopefullySymbolTokenErr := tokenizer.ReadNext()
	if hopefullySymbolTokenErr != nil {
		t.Error(hopefullySymbolTokenErr)
	}
	_, ok := hopefullySymbolToken.(SymbolToken)
	if !ok {
		t.Errorf("Wrong type. Expected Symbol but was %v", hopefullySymbolToken)
	}
}

func expectString(t *testing.T, tokenizer *Tokenizer, expectedString string) {
	hopefullyStringToken, hopefullyStringTokenErr := tokenizer.ReadNext()
	if hopefullyStringTokenErr != nil {
		t.Error(hopefullyStringTokenErr)
	}
	_, ok := hopefullyStringToken.(StringToken)
	if !ok {
		t.Errorf("Wrong type. Expected String but was %v", hopefullyStringToken)
	}
}

func expectStartScope(t *testing.T, tokenizer *Tokenizer) {
	maybeStartScope, maybeStartScopeErr := tokenizer.ReadNext()
	if maybeStartScopeErr != nil {
		t.Error(maybeStartScopeErr)
	}

	_, ok := maybeStartScope.(StartScopeToken)
	if !ok {
		t.Errorf("Couldn't cast to start scope token")
	}
}

func expectEndScope(t *testing.T, tokenizer *Tokenizer) {
	maybeEndScope, maybeEndScopeErr := tokenizer.ReadNext()
	if maybeEndScopeErr != nil {
		t.Error(maybeEndScopeErr)
	}

	_, ok := maybeEndScope.(EndScopeToken)
	if !ok {
		t.Errorf("Couldn't cast to end scope token")
	}
}

func expectLineDelimiter(t *testing.T, tokenizer *Tokenizer) {
	maybeStartScope, maybeStartScopeErr := tokenizer.ReadNext()
	if maybeStartScopeErr != nil {
		t.Error(maybeStartScopeErr)
	}

	_, ok := maybeStartScope.(LineDelimiterToken)
	if !ok {
		t.Errorf("Expected LineDelimiter but got %v", maybeStartScope)
	}
}

func expectEnd(t *testing.T, tokenizer *Tokenizer) {
	maybeEnd, maybeEndErr := tokenizer.ReadNext()
	if maybeEndErr != nil {
		t.Error(maybeEndErr)
	}
	if maybeEnd != nil {
		t.Errorf("Expected stream end")
	}
}

func TestIndentationSymbol(t *testing.T) {
	tokenizer := setup(
		`
namespace TheSpace
  component Something

    hello int32

    type Another
`)

	expectSymbol(t, tokenizer, "namespace")
	expectSymbol(t, tokenizer, "Thespace")
	expectStartScope(t, tokenizer)
	expectSymbol(t, tokenizer, "component")
	expectSymbol(t, tokenizer, "Something")
	expectStartScope(t, tokenizer)
	expectSymbol(t, tokenizer, "hello")
	expectSymbol(t, tokenizer, "int32")
	expectLineDelimiter(t, tokenizer)
	expectSymbol(t, tokenizer, "type")
	expectSymbol(t, tokenizer, "Another")
	expectLineDelimiter(t, tokenizer)
	expectEndScope(t, tokenizer)
	expectEndScope(t, tokenizer)
	expectEnd(t, tokenizer)
}
