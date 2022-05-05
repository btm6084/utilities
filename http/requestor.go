package http

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/btm6084/gojson"
	"github.com/btm6084/utilities/logging"
	log "github.com/sirupsen/logrus"
)

var (
	DefaultTimeout = 10 * time.Second
)

// RequestOptions is the data to configure an http request sent by a Requestor
type RequestOptions struct {
	Body    []byte
	Headers http.Header
	Cookies []*http.Cookie
}

// RequestResponse is the data returned by a Requestor
type RequestResponse struct {
	Body        []byte
	ContentType string
	StatusCode  int
	Headers     http.Header
}

// Requestor is an interface for making http requests
type Requestor interface {
	DoRequest(ctx context.Context, method, url string, options RequestOptions) (RequestResponse, error)
}

// HttpRequestor implements the Requestor interface
type HttpRequestor struct {
	c *http.Client
}

// NewRequestor returns a Requestor with a default client if one is not passed
func NewRequestor(c *http.Client) Requestor {
	if c == nil {
		t := http.DefaultTransport.(*http.Transport).Clone()
		t.MaxIdleConns = 100
		t.MaxConnsPerHost = 100
		t.MaxIdleConnsPerHost = 100

		c = &http.Client{
			Timeout:   10 * time.Second,
			Transport: t,
		}
	}
	return &HttpRequestor{c: c}
}

// DoRequest executes and parsed the response of an http request
// The context passed in will be used for request timeouts, UNLESS there is no timeout set, at which point the default will be used.
func (r *HttpRequestor) DoRequest(ctx context.Context, method, url string, options RequestOptions) (RequestResponse, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, DefaultTimeout)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(options.Body))
	if err != nil {
		log.WithFields(logging.TxnFields(ctx)).WithFields(log.Fields{"passthrough_url": url}).Println(err)
		return RequestResponse{}, err
	}

	for n, headers := range options.Headers {
		for _, h := range headers {
			req.Header.Add(n, h)
		}
	}

	req.Header.Set("X-Transaction-ID", logging.TransactionFromContext(ctx))
	req.Header.Set("Accept-Encoding", "gzip")

	for _, v := range options.Cookies {
		req.AddCookie(v)
	}

	res, err := r.c.Do(req)
	if err != nil {
		log.WithFields(logging.TxnFields(ctx)).WithFields(log.Fields{"passthrough_url": url}).Println(err)
		return RequestResponse{}, err
	}

	defer res.Body.Close()

	// Decompress gzip content
	var resBody io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		resBody, err = gzip.NewReader(res.Body)
		if err != nil {
			resBody = res.Body
		}
	default:
		resBody = res.Body
	}

	b, err := io.ReadAll(resBody)
	if err != nil {
		log.WithFields(logging.TxnFields(ctx)).WithFields(log.Fields{"passthrough_url": url}).Println(err)
		return RequestResponse{Body: b, StatusCode: res.StatusCode}, err
	}

	ct := res.Header.Get("Content-Type")
	if ct == "" {
		if gojson.IsJSON(b) {
			ct = "application/json"
		} else {
			ct = "text/plain"
		}
	}

	return RequestResponse{
		Body:        b,
		ContentType: ct,
		StatusCode:  res.StatusCode,
		Headers:     res.Header,
	}, nil
}
