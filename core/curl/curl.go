package curl

import (
	"fmt"

	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	MethodGet     = "GET"
	MethodHead    = "HEAD"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH" // RFC 5789
	MethodDelete  = "DELETE"
	MethodConnect = "CONNECT"
	MethodOptions = "OPTIONS"
	MethodTrace   = "TRACE"
)

const (
	TypeJSON       = "application/json"
	TypeUrlEncoded = "application/x-www-form-urlencoded"
)

type Client struct {
	client http.Client
}

type Option func(c *Client)

func WithTimeout(timeout int) Option {
	return func(c *Client) {
		c.client.Timeout = time.Duration(timeout) * time.Second
	}
}

// Next add CircuitBreaker, SingleFlight here

func NewClient(options ...Option) *Client {
	defaultClient := &Client{
		client: http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, function := range options {
		function(defaultClient)
	}

	fmt.Printf("%+v\n", defaultClient.client)
	return defaultClient
}

func (c *Client) Call(ctx context.Context, req *Request) ([]byte, int, error) {
	u, err := url.Parse(req.url)
	if err != nil {
		return nil, 0, err
	}
	u.RawQuery = req.params.Encode()

	fmt.Printf("%+v\n", req.body)
	request, err := http.NewRequest(req.method, u.String(), req.body)
	if err != nil {
		return nil, 0, err
	}
	for key, value := range req.headers {
		request.Header.Set(key, value)
	}

	response, err := c.client.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	}

	return contents, response.StatusCode, err
}

type Request struct {
	url     string
	method  string
	headers map[string]string
	params  url.Values
	body    io.Reader
}

func NewRequest(method string, url string) *Request {
	return &Request{url: url, method: method, headers: map[string]string{}}
}

func (r *Request) AddHeaders(maps map[string]string) {
	r.headers = maps
}

func (r *Request) AddContentType(contentType string) {
	r.headers["Content-Type"] = contentType
}

func (r *Request) WithUrlParameter(maps map[string]string) {
	r.params = mapToUrlValue(maps)
}

func (r *Request) WithBodyParameter(maps map[string]string) {
	params := mapToUrlValue(maps)
	r.body = strings.NewReader(params.Encode())
}

func (r *Request) WithBodyByte(body []byte) {
	r.body = bytes.NewReader(body)
}

func (r *Request) WithBodyJSON(data interface{}) {
	body, _ := json.Marshal(data)
	r.body = bytes.NewReader(body)
}

func mapToUrlValue(maps map[string]string) url.Values {
	param := url.Values{}
	for key, value := range maps {
		param.Add(key, value)
	}
	return param
}
