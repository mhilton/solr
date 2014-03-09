// Copyright 2014 Martin Hilton. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package solr

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

// A Query stores the parameters used to perform the query.
type Query map[string][]string

// NewQuery creates a new query using q as the query string
func NewQuery(q string) Query {
	return Query{"q": []string{q}}
}

// Add adds the param to value. It appends to any existing values
func (q Query) Add(param, value string) {
	q[param] = append(q[param], value)
}

// Set sets the param to value. It replaces any existing values.
func (q Query) Set(param, value string) {
	q[param] = []string{value}
}

// Get gets the first value associated with the givn param. If there are
// no values associated with the param Get returns the empty string.
func (q Query) Get(param string) string {
	if len(q[param]) == 0 {
		return ""
	}

	return q[param][0]
}

// SetQuery sets the q parameter to query.
func (q Query) SetQuery(query string) {
	q["q"] = []string{query}
}

// SetFields sets the fl param to fields
func (q Query) SetFields(fields ...string) {
	q["fl"] = []string{strings.Join(fields, " ")}
}

// SetRows sets the number of rows to return from the query.
func (q Query) SetRows(rows int) {
	q["rows"] = []string{strconv.Itoa(rows)}
}

//SetStart set the first row to return from the query.
func (q Query) SetStart(start int) {
	q["start"] = []string{strconv.Itoa(start)}
}

// SetFacet sets the facet field to on.
func (q Query) SetFacet(on bool) {
	if on {
		q["facet"] = []string{"true"}
	} else {
		q["facet"] = []string{"false"}
	}
}

// AddFacetField adds a new field to the list of facet fields
func (q Query) AddFacetField(field string) {
	q["facet.field"] = append(q["facet.field"], field)
}

type QueryResponse struct {
	ResponseHeader ResponseHeader `json:"response_header"`
	Response       Response
	FacetCounts    FacetCounts `json:"facet_counts"`
}

type ResponseHeader struct {
	Status int
	QTime  int
	Params map[string]interface{}
}

type Response struct {
	NumFound int
	Start    int
	Docs     []Document
}

type FacetCounts struct {
	FacetFields map[string][]FacetCount `json:"facet_fields"`
}

type FacetCount struct {
	Value string
	Count int
}

func (fc *FacetCount) UnmarshalJSON(b []byte) error {
	var d []interface{}
	json.Unmarshal(b, &d)

	if len(d) < 2 {
		return errors.New("not enough elements in FacetCount")
	}

	if v, ok := d[0].(string); ok {
		fc.Value = v
	} else {
		return errors.New("unexpected type for facet value")
	}

	if c, ok := d[1].(float64); ok {
		fc.Count = int(c)
	} else {
		return errors.New("unexpected type for facet count")
	}

	return nil
}

// PhraseQueryEscape formats a string for use in a SOLR phrase query
func PhraseQueryEscape(s string) string {
	var escapes int

	for _, c := range s {
		if c == '\\' || c == '"' {
			escapes++
		}
	}

	if escapes == 0 {
		return s
	}

	b := make([]byte, len(s)+escapes)
	j := 0
	for _, c := range []byte(s) {
		if c == '\\' || c == '"' {
			b[j] = '\\'
			j++
		}
		b[j] = c
		j++
	}

	return string(b)
}
