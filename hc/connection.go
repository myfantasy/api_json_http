package hc

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/myfantasy/mft"

	_ "github.com/myfantasy/api_json_http"
)

// Connection - to host
type Connection struct {
	Server              string        `json:"server"`
	IgnoreSSLValidation bool          `json:"ignore_ssql_validation"`
	QueryWait           time.Duration `json:"query_wait"`

	MaxConnsPerHost     int           `json:"max_conn"`
	MaxIdleConnDuration time.Duration `json:"max_idle_duration"`

	client *http.Client `json:"-"`
}

func CreateConnection(server string,
	ignoreSSLValidation bool,
	queryWait time.Duration,
	maxConnsPerHost int,
	maxIdleConnDuration time.Duration) (c *Connection) {

	c = &Connection{
		Server:              server,
		IgnoreSSLValidation: ignoreSSLValidation,
		QueryWait:           queryWait,
		MaxConnsPerHost:     maxConnsPerHost,
		MaxIdleConnDuration: maxIdleConnDuration,
	}

	return c
}

// Init connection
func (c *Connection) Init() *Connection {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.IgnoreSSLValidation,
		},
		MaxConnsPerHost:    c.MaxConnsPerHost,
		MaxIdleConns:       c.MaxConnsPerHost,
		IdleConnTimeout:    c.MaxIdleConnDuration,
		DisableCompression: true,
	}
	c.client = &http.Client{Transport: tr}

	return c
}

func ConnectionFromJson(body []byte) (c *Connection, err *mft.Error) {
	c = &Connection{}
	er0 := json.Unmarshal(body, c)
	if er0 != nil {
		return nil, mft.GenerateErrorE(10190000, er0)
	}
	return c, nil
}

// DoRawQuery do some query
func (c *Connection) DoRawQuery(queryWait time.Duration, path string, headersIn map[string]string, query []byte,
) (body []byte, headersOut map[string]string, statusCode int, err *mft.Error) {
	if queryWait == 0 {
		queryWait = c.QueryWait
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryWait)
	defer cancel()

	req, er0 := http.NewRequestWithContext(ctx, http.MethodPost, c.Server+path, bytes.NewBuffer(query))
	if er0 != nil {
		return body, nil, 0, mft.GenerateErrorE(20500100, er0, c.Server)
	}

	for k, v := range headersIn {
		req.Header.Add(k, v)
	}

	resp, er0 := c.client.Do(req)
	if er0 != nil {
		return body, nil, 0, mft.GenerateErrorE(20500101, er0, c.Server)
	}

	headersOut = map[string]string{}

	for k, v := range resp.Header {
		if len(v) == 0 {
			headersOut[k] = ""
		} else {
			headersOut[k] = v[0]
		}
	}

	body, er0 = io.ReadAll(resp.Body)
	if er0 != nil {
		return body, nil, 0, mft.GenerateErrorE(20500102, er0, c.Server)
	}

	return body, headersOut, resp.StatusCode, nil
}
