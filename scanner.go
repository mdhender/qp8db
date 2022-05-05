// qp8dp - a pure Go data query parser
//
// Copyright (c) 2022 Michael D Henderson
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package qp8db

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"
)

func Scan(b []byte) [][]byte {
	for len(b) > 0 {
		token, buffer := next(b)
		if len(token) != 0 {
			fmt.Printf("token(%q)\n", string(token))
		}
		b = buffer
	}
	return nil
}

// next returns the next chunk and remainder of the buffer.
// a chunk is either end of input, end of line, text, or spaces.
// both the empty buffer and the semicolon are treated as end of input.
// spaces is loosely defined to include spaces and invalid UTF-8.
// text is everything but end of file, end of line, or spaces.
// returns nil, nil on end of input.
// otherwise returns the chunk and the remainder of the buffer.
func next(b []byte) (token, buffer []byte) {
	// both the empty buffer and the semicolon are treated as end of input.
	if len(b) == 0 || b[0] == ';' {
		return nil, nil
	}

	// delimiters are always returned separately.
	delimiters := []byte{'.', ',', '(', ')', ';', '=', '+', '-', '*', '/', '\n'}
	if bytes.IndexByte(delimiters, b[0]) != -1 {
		return b[:1], b[1:]
	}

	// look for a run of spaces or text.
	r, w := utf8.DecodeRune(b)
	length := w
	if r == utf8.RuneError || unicode.IsSpace(r) {
		// consume runes up to the next non-space rune, end of line, or end of input.
		for length < len(b) && b[length] != '\n' {
			if r, w = utf8.DecodeRune(b[length:]); r == utf8.RuneError || unicode.IsSpace(r) {
				length += w
				continue
			}
			break
		}
	} else if r == '\'' {
		// return a string of all bytes up to the closing quote, end of line, or end of input.
		for length < len(b) && b[length] != '\n' {
			if b[length] != '\'' {
				r, w = utf8.DecodeRune(b[length:])
				length += w
				continue
			}
			length++ // consume the quote mark
			if length < len(b) && b[length] == '\'' {
				// two quote marks in a row counts as an escaped quote mark and is part of the string
				length++
				continue
			}
			break // otherwise we must terminate the string
		}
	} else {
		// return a run of runes up to the next space, delimiter, end of line, or end of input.
		for length < len(b) && b[length] != '\n' {
			if bytes.IndexByte(delimiters, b[length]) != -1 {
				break
			} else if r, w = utf8.DecodeRune(b[length:]); r == utf8.RuneError || unicode.IsSpace(r) {
				break
			}
			length += w
		}
	}

	// as a convenience, return nil for the buffer if we've exhausted it.
	if len(b[length:]) == 0 {
		return b[:length], nil
	}

	return b[:length], b[length:]
}
