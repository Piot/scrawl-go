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

type Position struct {
	line   int
	column int
}

func (p Position) String() string {
	return fmt.Sprintf("[%d:%d]", p.line, p.column)
}

// Tokenizer :
type Tokenizer struct {
	r                     *RuneReader
	indentation           int
	lastParsedIndentation int
	targetIndentation     int
	position              Position
	oldPosition           Position
	lastTokenWasDelimiter bool
}

// Token :
type Token interface {
	Position() Position
}

type TokenizerError struct {
	err      error
	position Position
}

func (f TokenizerError) Error() string {
	return fmt.Sprintf("%v at %v", f.err, f.position)
}

// SymbolToken :
type SymbolToken struct {
	Symbol   string
	position Position
}

func (s SymbolToken) Position() Position {
	return s.position
}

func (s SymbolToken) String() string {
	return fmt.Sprintf("Symbol:%s", s.Symbol)
}

// NumberToken :
type NumberToken struct {
	number   float64
	position Position
}

func (s NumberToken) Position() Position {
	return s.position
}

func (s NumberToken) String() string {
	return fmt.Sprintf("Number:%f", s.number)
}

// StartScopeToken :
type StartScopeToken struct {
	position Position
}

func (s StartScopeToken) Position() Position {
	return s.position
}
func (StartScopeToken) String() string {
	return "StartScope"
}

// EndScopeToken :
type EndScopeToken struct {
	position Position
}

func (s EndScopeToken) Position() Position {
	return s.position
}
func (EndScopeToken) String() string {
	return "EndScope"
}

// LineDelimiterToken :
type LineDelimiterToken struct {
	position Position
}

func (s LineDelimiterToken) Position() Position {
	return s.position
}
func (LineDelimiterToken) String() string {
	return "LineDelimiter"
}

// StartMetaDataToken :
type StartMetaDataToken struct {
	position Position
}

func (s StartMetaDataToken) Position() Position {
	return s.position
}
func (StartMetaDataToken) String() string {
	return "StartMetaDataToken"
}

// EndMetaDataToken :
type EndMetaDataToken struct {
	position Position
}

func (s EndMetaDataToken) Position() Position {
	return s.position
}
func (EndMetaDataToken) String() string {
	return "EndMetaDataToken"
}

// StringToken :
type StringToken struct {
	text     string
	position Position
}

func (s StringToken) Position() Position {
	return s.position
}

func (s StringToken) String() string {
	return fmt.Sprintf("string:%s", s.text)
}

func (s StringToken) Text() string {
	return s.text
}

func isIndentation(ch rune) bool {
	return ch == ' '
}

func isNewLine(ch rune) bool {
	return ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isLegalIndentationWhiteSpace(ch rune) bool {
	return !(ch == ' ')
}

func isWhitespaceExceptNewLine(ch rune) bool {
	return (ch == ' ' || ch == '\t')
}

func isIllegalDuringIndentation(ch rune) bool {
	return ch == '\t'
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isSymbol(ch rune) bool {
	return isLetter(ch) || isDigit(ch) || ch == '-'
}

func isStartMetaData(ch rune) bool {
	return ch == '['
}

func isEndMetaData(ch rune) bool {
	return ch == ']'
}

func isStartString(ch rune) bool {
	return ch == '\'' || ch == '"'
}

// NewTokenizer :
func NewTokenizer(r *RuneReader) *Tokenizer {
	return &Tokenizer{r: r, position: Position{line: 1, column: 1}, lastTokenWasDelimiter: true}
}

func (t *Tokenizer) nextRune() rune {
	t.oldPosition = t.position
	ch := t.r.read()
	if ch == '\n' {
		t.position.line++
		t.position.column = 1
	} else {
		t.position.column++
	}
	return ch
}

func (t *Tokenizer) unreadRune() {
	t.r.unread()
	t.position = t.oldPosition
}

func (t *Tokenizer) parseIndentation() (int, error) {
	indentation := 0
	for true {
		ch := t.nextRune()
		if !isIndentation(ch) {
			if isIllegalDuringIndentation(ch) {
				return 0, fmt.Errorf("Illegal indentation character: %c (%v)", ch, ch)
			}
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
	startPosition := t.position

	for true {
		ch := t.nextRune()
		if !isSymbol(ch) {
			t.unreadRune()
			break
		}
		a += string(ch)
	}
	return SymbolToken{Symbol: a, position: startPosition}, nil
}

func (t *Tokenizer) parseNumber() (Token, error) {
	var a string
	startPosition := t.position
	for true {
		ch := t.nextRune()
		if !isDigit(ch) {
			t.unreadRune()
			break
		}
		a += string(ch)
	}
	return NumberToken{number: 0.0, position: startPosition}, nil
}

func (t *Tokenizer) parseString(startStringRune rune) (Token, error) {
	var a string
	startPosition := t.position
	for true {
		ch := t.nextRune()
		if ch == startStringRune {
			break
		}
		if ch == 0 {
			return nil, fmt.Errorf("Unexpected end while finding end of string")
		}
		a += string(ch)
	}
	return StringToken{text: a, position: startPosition}, nil
}

func (t *Tokenizer) parseNewLine() (Token, error) {
	indentation, indentationErr := t.parseIndentation()
	if indentationErr != nil {
		return nil, indentationErr
	}
	if indentation == -1 {
		return t.internalReadNext()
	}
	if indentation > t.indentation+1 {
		return nil, fmt.Errorf("Too much indentation")
	}
	t.targetIndentation = indentation
	if t.indentation >= t.targetIndentation && !t.lastTokenWasDelimiter {
		t.lastTokenWasDelimiter = true
		return LineDelimiterToken{position: t.position}, nil
	}
	return t.internalReadNext()
}

func (t *Tokenizer) internalReadNext() (Token, error) {
	if t.indentation != t.targetIndentation {
		diff := t.targetIndentation - t.indentation
		if diff > 0 {
			t.indentation++
			t.lastTokenWasDelimiter = true
			return StartScopeToken{position: t.position}, nil
		}
		t.indentation--
		return EndScopeToken{position: t.position}, nil
	}

	r := t.nextRune()
	if isNewLine(r) {
		return t.parseNewLine()
	} else {
		t.lastTokenWasDelimiter = false
		if isLetter(r) {
			t.unreadRune()
			return t.parseSymbol()
		} else if isDigit(r) {
			t.unreadRune()
			return t.parseNumber()
		} else if isStartString(r) {
			return t.parseString(r)
		} else if isStartMetaData(r) {
			return StartMetaDataToken{position: t.position}, nil
		} else if isEndMetaData(r) {
			return EndMetaDataToken{position: t.position}, nil
		} else if isWhitespaceExceptNewLine(r) {
			return t.internalReadNext()
		} else if r == ',' {
			return t.internalReadNext()
		} else if r == 0 {
			return nil, nil
		}
	}
	return nil, fmt.Errorf("Unknown rune '%c' %v", r, r)
}

func (t *Tokenizer) ReadNext() (Token, error) {
	token, err := t.internalReadNext()
	if err != nil {
		return nil, TokenizerError{err: err, position: t.position}
	}
	// fmt.Printf("return: %v\n", token)
	return token, nil
}
