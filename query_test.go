// Copyright 2014 Martin Hilton. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

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
