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
	"io"
	"strings"
	"unicode"

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

func StringForceTranslate(p_string string) (string, error) {
	var sb strings.Builder
	didsomething := false
	for _, r := range p_string {
		if unicode.IsOneOf([]*unicode.RangeTable{unicode.So, unicode.Sk}, r) {
			sb.WriteRune(':')
			sb.WriteString(strings.ReplaceAll(strings.ToLower(runenames.Name(r)), " ", "_"))
			didsomething = true
			sb.WriteRune(':')
		} else {
			sb.WriteRune(r)
		}
	}
	if didsomething {
		return sb.String(), nil
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
