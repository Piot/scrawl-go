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
	"fmt"
	"strings"

	"github.com/piot/scrawl-go/src/runestream"
	"github.com/piot/scrawl-go/src/token"
)

// Tokenizer :
type Tokenizer struct {
	r                     *runestream.RuneReader
	indentation           int
	lastParsedIndentation int
	targetIndentation     int
	position              token.Position
	oldPosition           token.Position
	lastTokenWasDelimiter bool
}

type TokenizerError struct {
	err      error
	position token.Position
}

func (f TokenizerError) Error() string {
	return fmt.Sprintf("%v at %v", f.err, f.position)
}

func SetupTokenizer(text string) *Tokenizer {
	ioReader := strings.NewReader(text)
	runeReader := runestream.NewRuneReader(ioReader)
	tokenizer := NewTokenizer(runeReader)
	return tokenizer
}

func FetchAllTokens(x string) ([]token.Token, error) {
	t := SetupTokenizer(x)
	tokens, readErr := t.ReadAll()

	return tokens, readErr
}

// NewTokenizer :
func NewTokenizer(r *runestream.RuneReader) *Tokenizer {
	return &Tokenizer{r: r, position: token.NewPositionTopLeft(), lastTokenWasDelimiter: true}
}

func (t *Tokenizer) nextRune() rune {
	t.oldPosition = t.position
	ch := t.r.Read()
	if ch == '\n' {
		t.position = t.position.NextLine()
		t.position = t.position.FirstColumn()
	} else {
		t.position = t.position.NextColumn()
	}
	return ch
}

func (t *Tokenizer) unreadRune() {
	t.r.Unread()
	t.position = t.oldPosition
}

func (t *Tokenizer) parseComment() (token.Token, error) {
	var a string
	startPosition := t.position
	for true {
		ch := t.nextRune()
		if isNewLine(ch) {
			t.unreadRune()
			t.lastTokenWasDelimiter = false
			break
		}
		a += string(ch)
	}
	a = strings.TrimSpace(a)
	return token.NewCommentToken(a, startPosition), nil
}

func (t *Tokenizer) parseString(startStringRune rune) (token.Token, error) {
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
	return token.NewStringToken(a, startPosition), nil
}

func (t *Tokenizer) parseNewLine() (token.Token, error) {
	indentation, indentationErr := t.parseIndentation()
	if indentationErr != nil {
		return nil, indentationErr
	}
	//fmt.Printf("indentation:%v target:%v\n", indentation, t.targetIndentation)
	if indentation == -1 {
		return t.internalReadNext()
	}
	if indentation > t.indentation+1 {
		return nil, fmt.Errorf("Too much indentation")
	}
	t.targetIndentation = indentation
	//fmt.Printf("indentation after:%v target:%v\n", t.indentation, t.targetIndentation)
	if t.indentation >= t.targetIndentation && !t.lastTokenWasDelimiter {
		t.lastTokenWasDelimiter = true
		return token.NewLineDelimiter(t.position), nil
	}
	return t.internalReadNext()
}

func (t *Tokenizer) internalReadNext() (token.Token, error) {
	if t.indentation != t.targetIndentation {
		diff := t.targetIndentation - t.indentation
		if diff > 0 {
			t.indentation++
			t.lastTokenWasDelimiter = true
			return token.NewStartScopeToken(t.position), nil
		}
		t.indentation--
		return token.NewEndScopeToken(t.position), nil
	}

	r := t.nextRune()
	if isNewLine(r) {
		return t.parseNewLine()
	} else {
		if isWhitespaceExceptNewLine(r) {
			return t.internalReadNext()
		}
		t.lastTokenWasDelimiter = false
		if isLetter(r) {
			t.unreadRune()
			return t.parseSymbol()
		} else if isDigit(r) || r == '-' {
			t.unreadRune()
			return t.parseNumber()
		} else if isStartString(r) {
			return t.parseString(r)
		} else if isStartMetaData(r) {
			return token.NewStartMetaDataToken(t.position), nil
		} else if isEndMetaData(r) {
			return token.NewEndMetaDataToken(t.position), nil
		} else if r == ',' {
			return t.internalReadNext()
		} else if r == '.' {
			return token.NewOperatorToken(r, t.position), nil
		} else if r == '#' {
			return t.parseComment()
		} else if r == 0 {
			return nil, nil
		}
	}
	return nil, fmt.Errorf("Unknown rune '%c' %v", r, r)
}

func (t *Tokenizer) ReadNext() (token.Token, error) {
	token, err := t.internalReadNext()
	if err != nil {
		return nil, TokenizerError{err: err, position: t.position}
	}
	//fmt.Printf("return: %v\n", token)
	return token, nil
}

func (t *Tokenizer) ReadAll() ([]token.Token, error) {
	var tokens []token.Token
	for {
		tok, tokErr := t.ReadNext()
		if tokErr != nil {
			return nil, tokErr
		}
		if tok == nil {
			break
		}
		tokens = append(tokens, tok)
	}

	return tokens, nil
}
