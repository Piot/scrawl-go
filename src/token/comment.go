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

// CommentToken :
type CommentToken struct {
	text     string
	position Position
}

func NewCommentToken(text string, position Position) CommentToken {
	return CommentToken{text: text, position: position}
}

func (s CommentToken) IsEqual(other Token) bool {
	otherNumber, isNumber := other.(CommentToken)
	if !isNumber {
		return false
	}

	return otherNumber.text == s.text
}

func (s CommentToken) Position() Position {
	return s.position
}

func (s CommentToken) String() string {
	return fmt.Sprintf("comment:%s", s.text)
}

func (s CommentToken) Text() string {
	return s.text
}
