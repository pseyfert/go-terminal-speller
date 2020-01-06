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
		{"the ∑ xᵢ i² ≥ 11"},
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
		{"🖖🏿", ":raised_hand_with_part_between_middle_and_ring_fingers::emoji_modifier_fitzpatrick_type-6:"},
		{"😅", ":smiling_face_with_open_mouth_and_cold_sweat:"},
		{"👍", ":thumbs_up_sign:"},
		{"🙈🥑", ":see-no-evil_monkey::avocado:"},
		{"you know what we should do? 🚵‍♀️ ← that's my pick!", "you know what we should do? :mountain_bicyclist:‍:female_sign:️ ← that's my pick!"},
		{"where do I find a 🏦", "where do I find a :bank:"},
		{"💖, really cool!", ":sparkling_heart:, really cool!"},
		{"🇪🇸 flags are tricky", ":regional_indicator_symbol_letter_e::regional_indicator_symbol_letter_s: flags are tricky"},
		{"∫ ☺ ⌢ ∼ …", "∫ :white_smiling_face: :frown: ∼ …"},
	}
	for _, testcase := range testcases {
		var buffer bytes.Buffer
		writer := NewTranslator(&buffer)
		writer.Write([]byte(testcase.Msg)) // CHECKME
		if strings.Compare(buffer.String(), testcase.Expected) != 0 {
			t.Errorf("Emoji in '%s' didn't get spelled out as expected: got '%s' (expected '%s')", testcase.Msg, buffer.String(), testcase.Expected)
		}
	}
}
