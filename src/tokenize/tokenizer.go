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

import "fmt"

// Tokenizer :
type Tokenizer struct {
	r                     *RuneReader
	indentation           int
	lastParsedIndentation int
	targetIndentation     int
}

// Token :
type Token interface {
}

// SymbolToken :
type SymbolToken struct {
	Symbol string
}

// StartScopeToken :
type StartScopeToken struct {
}

// EndScopeToken :
type EndScopeToken struct {
}

func isIndentation(ch rune) bool {
	return ch == ' ' || ch == '\t'
}

func isNewLine(ch rune) bool {
	return ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isSymbol(ch rune) bool {
	return isLetter(ch) || isDigit(ch) || ch == '-'
}

// NewTokenizer :
func NewTokenizer(r *RuneReader) *Tokenizer {
	return &Tokenizer{r: r}
}

func (t *Tokenizer) nextRune() rune {
	return t.r.read()
}

func (t *Tokenizer) unreadRune() {
	t.r.unread()
}

func (t *Tokenizer) parseIndentation() (int, error) {
	indentation := 0
	for true {
		ch := t.nextRune()
		if !isIndentation(ch) {
			t.unreadRune()
			if isNewLine(ch) {
				return -1, nil
			}
			if ch == 0 {
				return 0, nil
			}
			break
		}
		indentation++

	}
	if indentation%2 != 0 {
		return 0, fmt.Errorf("Must have double spaces as indentation")
	}
	return indentation / 2, nil
}

func (t *Tokenizer) parseSymbol() (Token, error) {
	var a string
	for true {
		ch := t.nextRune()
		if !isSymbol(ch) {
			t.unreadRune()
			break
		}
		a += string(ch)
	}
	return SymbolToken{Symbol: a}, nil
}

func (t *Tokenizer) readNext() (Token, error) {
	if t.indentation != t.targetIndentation {
		diff := t.targetIndentation - t.indentation
		if diff > 0 {
			t.indentation++
			return StartScopeToken{}, nil
		}
		t.indentation--
		return EndScopeToken{}, nil
	}

	r := t.nextRune()
	if isNewLine(r) {
		indentation, indentationErr := t.parseIndentation()
		if indentationErr != nil {
			return nil, indentationErr
		}
		if indentation == -1 {
			return t.readNext()
		}
		if indentation > t.indentation+1 {
			return nil, fmt.Errorf("Too much indentation")
		}
		t.targetIndentation = indentation
		return t.readNext()
	} else if isLetter(r) {
		t.unreadRune()
		return t.parseSymbol()
	} else if r == 0 {
		return nil, nil
	}
	return t.readNext()
}
