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
	"strconv"

	"github.com/piot/scrawl-go/src/token"
)

func (t *Tokenizer) parseNumber() (token.Token, error) {
	var a string
	startPosition := t.position
	start := t.nextRune()
	negative := false
	if start == '-' {
		negative = true
	} else {
		t.unreadRune()
	}
	for true {
		ch := t.nextRune()
		if !isDigit(ch) {
			t.unreadRune()
			break
		}
		a += string(ch)
	}
	v, vErr := strconv.Atoi(a)
	if negative {
		v = -v
	}
	if vErr != nil {
		return nil, vErr
	}
	return token.NewNumberToken(float64(v), startPosition), nil
}
