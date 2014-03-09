// Copyright 2014 Martin Hilton. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package solr

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// A Conn is a connection to a SOLR server.
type Conn struct {
	// HTTPClient can be used to specify a different http.Client to use to
	// query the SOLR server.
	//
	// If HTTPClient is nil, the http.DefaultClient will be used.
	HTTPClient *http.Client

	// URL is the base URL of the SOLR server to which you wish to connect.
	URL *url.URL
}

// NewConn creates a new Conn pointing to the specified url. An error is
// returned if url cannot be parsed.
func NewConn(base string) (Conn, error) {
	u, err := url.Parse(base)
	return Conn{URL: u}, err
}

// Query performs the specified query against the SOLR server, using the
// default select handler.
func (c Conn) Query(q Query) (*QueryResponse, error) {
	return c.QueryHandler(q, "select")
}

// QueryHandler performs the specified query againt the SOLR server,
// using the specified handler.
func (c Conn) QueryHandler(q Query, h string) (*QueryResponse, error) {
	cl := c.HTTPClient
	if cl == nil {
		cl = http.DefaultClient
	}

	u, err := c.URL.Parse(h)
	if err != nil {
		return nil, err
	}

	// copy the query into a url.Values because it's rude to modify someone elses map
	p := url.Values{}
	for k, v := range q {
		p[k] = v
	}

	p.Set("wt", "json")
	p.Set("json.nl", "arrarr")

	u.RawQuery = p.Encode()

	resp, err := cl.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer flush(resp.Body)

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, ErrHTTPStatus(resp.Status)
	}

	qr := new(QueryResponse)
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(qr)
	if err != nil {
		return nil, err
	}

	return qr, err
}

// flush reads any remaining content and then closes the ReadCloser.
// SOLR may send additional whitespace and the Body needs to be read
// completely before the connection can be reused.
func flush(rc io.ReadCloser) error {
	ioutil.ReadAll(rc)
	return rc.Close()
}

// ErrHTTPStatus is returned if an HTTP request responds with a non ok response.
type ErrHTTPStatus string

func (err ErrHTTPStatus) Error() string {
	return "bad HTTP response: " + string(err)
}
