/*
 * Copyright (c) 2014 Martin Hilton <martin.hilton@zepler.net>
 *
 * Permission to use, copy, modify, and distribute this software for any
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

package solr

import (
	"testing"
)

func TestQuerySetAndGet(t *testing.T) {
	q := Query{}

	q.Set("q", "test query")

	if q.Get("q") != "test query" {
		t.Error("Incorrect value returned for set value")
	}

	if q.Get("fl") != "" {
		t.Error("Incorrect value returned for not-set value")
	}
}

func TestNewQuery(t *testing.T) {
	q := NewQuery("test query")
	if q.Get("q") != "test query" {
		t.Error("Incorrect value for query: " + q.Get("q"))
	}
}

func TestSetFields(t *testing.T) {
	q := Query{}

	q.SetFields("f1", "f2", "f3")
	if q.Get("fl") != "f1 f2 f3" {
		t.Error("Incorrect value for fl: " + q.Get("fl"))
	}
}

func TestPhraseQueryEscape(t *testing.T) {
	tests := []struct{ in, out string }{
		{"test 1", "test 1"},
		{`"`, `\"`},
		{`\`, `\\`},
		{`phrase containing"a quote`, `phrase containing\"a quote`},
	}

	for _, test := range tests {
		r := PhraseQueryEscape(test.in)
		if r != test.out {
			t.Errorf("Unexpected escape of %q: expected %q got %q", test.in, test.out, r)
		}
	}
}
