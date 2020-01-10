/*
 * Copyright (c) 2019 Paul Seyfert <pseyfert.mathphys@gmail.com>
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package terminalspeller

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/pseyfert/go-http-redirect-resolve/resolve"
	"golang.org/x/text/unicode/runenames"
)

type translator struct {
	w            io.Writer
	Didsomething bool
}

func NewTranslator(writer io.Writer) translator {
	retval := translator{
		w:            writer,
		Didsomething: false,
	}
	return retval
}

func EmojipediaUrl(emoji string) (string, error) {
	resolved, err := resolve.Resolve(fmt.Sprintf("http://📙.la/%s", emoji))
	return resolved, err
}

func RuneReplace(r rune, sb io.StringWriter) bool {
	if unicode.IsOneOf([]*unicode.RangeTable{unicode.So, unicode.Sk}, r) {
		sb.WriteString(":")
		sb.WriteString(strings.ReplaceAll(strings.ToLower(runenames.Name(r)), " ", "_"))
		sb.WriteString(":")
		return true
	} else {
		sb.WriteString(string(r))
		return false
	}
}

func StringForceTranslate(p_string string) (string, error) {
	var sb strings.Builder
	var combiner strings.Builder
	didsomething := false
	for _, r := range p_string {
		if combiner.Len() != 0 {
			// check if a combination continues (or if a new emoji starts)
			if unicode.Is(unicode.Sk, r) { // modifier
				combiner.WriteRune(r)
			} else if str := []rune(combiner.String()); str[len(str)-1] == []rune("\u200d")[0] { // zero width joiner
				combiner.WriteRune(r)
			} else if r == '\ufe0f' { // print previous character as emoji https://emojipedia.org/emoji/%EF%B8%8F/
				combiner.WriteRune(r)
			} else if r == '\u200d' { // zero width joiner https://emojipedia.org/zero-width-joiner/
				combiner.WriteRune(r)
			} else {
				url, err := EmojipediaUrl(combiner.String())
				if err == nil {
					sb.WriteString(fmt.Sprintf("[%s]", url))
				}
				combiner.Reset()
			}
		} else if unicode.Is(unicode.So, r) {
			combiner.WriteRune(r)
		}

		didsomething = RuneReplace(r, &sb) || didsomething
	}
	if didsomething {
		return sb.String(), nil
	}
	if sb.Len() != 0 {
		sb.WriteString(fmt.Sprintf("[http://📙.la/%s]", combiner.String()))
	}
	return p_string, ErrNoReplacment
}

func (myself *translator) WriteString(p_string string) (int, error) {
	translation, err := StringForceTranslate(p_string)
	if err != ErrNoReplacment {
		myself.Didsomething = true
	}
	return io.WriteString(myself.w, translation)
}

func (myself *translator) Write(p []byte) (int, error) {
	// FIXME: how to handle unterminated unicode code points?
	translation, err := StringForceTranslate(string(p))
	if err != ErrNoReplacment {
		myself.Didsomething = true
	}
	return io.WriteString(myself.w, translation)
}

var ErrNoReplacment = errors.New("No replacement has been done despite being requested")
