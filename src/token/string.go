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

package token

import "fmt"

// StringToken :
type StringToken struct {
	text     string
	position Position
}

func NewStringToken(text string, position Position) StringToken {
	return StringToken{text: text, position: position}
}

func (s StringToken) IsEqual(other Token) bool {
	otherNumber, isNumber := other.(StringToken)
	if !isNumber {
		return false
	}

	return otherNumber.text == s.text
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
