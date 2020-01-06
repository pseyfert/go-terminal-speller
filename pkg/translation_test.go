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
	"bytes"
	"strings"
	"testing"
)

func TestNoChange(t *testing.T) {
	testcases := []struct {
		Msg string
	}{
		{"τ→µµµ"},
		{"Kein Mülleimer → Leute können Müll nicht richtig entsorgen → überall liegt Müll rum"},
	}
	for _, testcase := range testcases {
		var buffer bytes.Buffer
		writer := NewTranslator(&buffer)
		writer.Write([]byte(testcase.Msg)) // CHECKME
		if bytes.Compare(buffer.Bytes(), []byte(testcase.Msg)) != 0 {
			t.Errorf("Message that wasn't to be altered got altered: got '%s' (expected '%s')", buffer.Bytes(), testcase.Msg)
		}
	}
}

func TestEmojis(t *testing.T) {
	testcases := []struct {
		Msg      string
		Expected string
	}{
		{"🤷", ":shrug:"},
	}
	for _, testcase := range testcases {
		var buffer bytes.Buffer
		writer := NewTranslator(&buffer)
		writer.Write([]byte(testcase.Msg)) // CHECKME
		if strings.Compare(buffer.String(), testcase.Expected) != 0 {
			t.Errorf("Emoji didn't get spelled out as expected: got '%s' (expected '%s')", buffer.String(), testcase.Expected)
		}
	}
}
