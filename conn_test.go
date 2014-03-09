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
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestSimpleQuery(t *testing.T) {
	rt := TestRoundTripper{}
	rt.RequestCheck = func(req *http.Request) {
		v := req.URL.Query()
		if v.Get("q") != "test query" {
			t.Error("Incorrect query: expected \"test query\", got \"" + v.Get("q") + "\"")
		}

		if v.Get("wt") != "json" {
			t.Error("Incorrect wt: expected \"json\", got \"" + v.Get("wt") + "\"")
		}

		if v.Get("json.nl") != "arrarr" {
			t.Error("Incorrect json.nl: expected \"arr\", got \"" + v.Get("json.nl") + "\"")
		}
	}

	rt.Response = makeResponse(200, "200 OK", `{
  "response_header":{
    "status":20,
    "QTime":10,
    "params":{
      "wt":"json",
      "q":"test query",
      "json.nl":"arrarr"}},
  "response":{
    "numFound":1,
    "start":2,
    "docs":[{"id":"testdoc"}]}} `)

	cl := &http.Client{Transport: rt}
	c, err := NewConn("http://example.com/solr/")
	if err != nil {
		t.Fatal("Unexpected error creating connection: " + err.Error())
	}
	c.HTTPClient = cl

	q := NewQuery("test query")
	qr, err := c.Query(q)
	if err != nil {
		t.Fatal("Unexpected error performing query: " + err.Error())
	}

	if !rt.Response.Body.(*TestReadCloser).EOF {
		t.Error("The response body was not read fully")
	}

	if !rt.Response.Body.(*TestReadCloser).Closed {
		t.Error("The response body was not closed")
	}

	if qr.ResponseHeader.Status != 20 {
		t.Errorf("Unexpected ResponseHeader.Status: expected 20 got %v", qr.ResponseHeader.Status)
	}

	if qr.ResponseHeader.QTime != 10 {
		t.Errorf("Unexpected ResponseHeader.QTime: expected 10 got %v", qr.ResponseHeader.QTime)
	}

	if len(qr.ResponseHeader.Params) != 3 {
		t.Errorf("Unexpeced number of ResponseHeader.Params: expected 3 items, got %v", qr.ResponseHeader.Params)
	}

	if qr.Response.NumFound != 1 {
		t.Errorf("Unexpected Response.NumFound: expected 1, got %v", qr.Response.NumFound)
	}

	if qr.Response.Start != 2 {
		t.Errorf("Unexpected Response.Start: expected 2, got %v", qr.Response.Start)
	}

	if len(qr.Response.Docs) != 1 {
		t.Fatalf("Unexpected number of Response.Docs: expected 1, got %v", qr.Response.Docs)
	}

	doc := qr.Response.Docs[0]
	if doc["id"] != "testdoc" {
		t.Errorf("Unexpected doc[\"id\"]: expected \"testdoc\" got %q", doc["id"])
	}
}

func TestFacetQuery(t *testing.T) {
	rt := TestRoundTripper{}
	rt.RequestCheck = func(req *http.Request) {
		v := req.URL.Query()
		if v.Get("q") != "test query" {
			t.Error("Incorrect query: expected \"test query\", got \"" + v.Get("q") + "\"")
		}

		if v.Get("facet") != "true" {
			t.Error("Unexpected facet: expected true, got " + v.Get("facet"))
		}

		if v.Get("facet.field") != "facet_field" {
			t.Error("Unexpected facet.field: expected \"fascet_field\", got \"" + v.Get("facet.mincount") + "\"")
		}

		if v.Get("facet.mincount") != "1" {
			t.Error("Unexpected facet.mincount: expected 1, got " + v.Get("facet.mincount"))
		}

		if v.Get("wt") != "json" {
			t.Error("Incorrect wt: expected \"json\", got \"" + v.Get("wt") + "\"")
		}

		if v.Get("json.nl") != "arrarr" {
			t.Error("Incorrect json.nl: expected \"arr\", got \"" + v.Get("json.nl") + "\"")
		}
	}

	rt.Response = makeResponse(200, "200 OK", `{
  "response_header":{
    "status":0,
    "QTime":10,
    "params":{
      "wt":"json",
      "facet":true,
      "facet.mincount":1,
      "facet.field":["facet_field"],
      "q":"test query",
      "json.nl":"arrarr"}},
  "response":{
    "numFound":20,
    "start":0,
    "docs":[]},
  "facet_counts":{
    "facet_fields":{
      "facet_field":[
        ["value1", 10],
        ["value2", 8],
        ["value4", 2]]}}} `)

	cl := &http.Client{Transport: rt}
	c, err := NewConn("http://example.com/solr/")
	if err != nil {
		t.Fatal("Unexpected error creating connection: " + err.Error())
	}
	c.HTTPClient = cl

	q := NewQuery("test query")
	q.SetFacet(true)
	q.AddFacetField("facet_field")
	q.Set("facet.mincount", "1")
	qr, err := c.Query(q)
	if err != nil {
		t.Fatal("Unexpected error performing query: " + err.Error())
	}

	if !rt.Response.Body.(*TestReadCloser).EOF {
		t.Error("The response body was not read fully")
	}

	if !rt.Response.Body.(*TestReadCloser).Closed {
		t.Error("The response body was not closed")
	}

	if qr.ResponseHeader.Status != 0 {
		t.Errorf("Unexpected ResponseHeader.Status: expected 0 got %v", qr.ResponseHeader.Status)
	}

	if qr.ResponseHeader.QTime != 10 {
		t.Errorf("Unexpected ResponseHeader.QTime: expected 10 got %v", qr.ResponseHeader.QTime)
	}

	if len(qr.ResponseHeader.Params) != 6 {
		t.Errorf("Unexpeced number of ResponseHeader.Params: expected 6 items, got %v", qr.ResponseHeader.Params)
	}

	if qr.Response.NumFound != 20 {
		t.Errorf("Unexpected Response.NumFound: expected 20, got %v", qr.Response.NumFound)
	}

	if qr.Response.Start != 0 {
		t.Errorf("Unexpected Response.Start: expected 0, got %v", qr.Response.Start)
	}

	if len(qr.Response.Docs) != 0 {
		t.Errorf("Unexpected number of Response.Docs: expected 0, got %v", qr.Response.Docs)
	}

	if len(qr.FacetCounts.FacetFields) != 1 {
		t.Fatalf("Unexpected number of FacetFields: expected 1, got %v", qr.FacetCounts.FacetFields)
	}

	field := qr.FacetCounts.FacetFields["facet_field"]
	if len(field) != 3 {
		t.Fatalf("Unexpected number of FacetCounts in field: expected 3 got %v", field)
	}

	if field[0].Value != "value1" {
		t.Errorf("Unexpected value for FacetField[0]: expected \"value1\", got \"%s\"", field[0].Value)
	}

	if field[0].Count != 10 {
		t.Errorf("Unexpected count for FacetField[0]: expected 10, got %d", field[0].Count)
	}

	if field[1].Value != "value2" {
		t.Errorf("Unexpected value for FacetField[1]: expected \"value2\", got \"%s\"", field[1].Value)
	}

	if field[1].Count != 8 {
		t.Errorf("Unexpected count for FacetField[1]: expected 8, got %d", field[1].Count)
	}

	if field[2].Value != "value4" {
		t.Errorf("Unexpected value for FacetField[2]: expected \"value4\", got \"%s\"", field[2].Value)
	}

	if field[2].Count != 2 {
		t.Errorf("Unexpected count for FacetField[2]: expected 2, got %d", field[2].Count)
	}
}

func makeResponse(statusCode int, status string, content string) *http.Response {
	return &http.Response{
		Status:     fmt.Sprintf("%d %s", statusCode, status),
		StatusCode: statusCode,
		Proto:      "HTTP",
		Body:       &TestReadCloser{Data: content},
	}
}

type TestReadCloser struct {
	Data   string
	pos    int
	EOF    bool
	Closed bool
}

func (t *TestReadCloser) Read(p []byte) (n int, err error) {
	if t.pos >= len(t.Data) {
		t.EOF = true
		return 0, io.EOF
	}

	n = len(p)
	if n > len(t.Data)-t.pos {
		n = len(t.Data) - t.pos
	}

	for i := 0; i < n; i++ {
		p[i] = t.Data[t.pos+i]
	}
	t.pos += n

	return
}

func (t *TestReadCloser) Close() error {
	t.Closed = true
	return nil
}

type TestRoundTripper struct {
	RequestCheck func(*http.Request)
	Response     *http.Response
	Err          error
}

func (t TestRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.RequestCheck != nil {
		t.RequestCheck(req)
	}

	return t.Response, t.Err
}
