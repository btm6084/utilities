package http

import (
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MockTransport struct{}

func (*MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := http.Response{
		StatusCode: http.StatusOK,
		Header:     req.Header,
		Body:       req.Body,
	}

	return &r, nil
}

func TestRequest(t *testing.T) {
	transport := &MockTransport{}
	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	r := NewRequestor(&client)
	resp, err := r.DoRequest(context.Background(), "GET", "http://localhost:8080", RequestOptions{})

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "text/plain", resp.ContentType)
}

func TestJSONRequest(t *testing.T) {
	transport := &MockTransport{}
	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	r := NewRequestor(&client)
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	b := []byte(`{"testKey":"testString"}`)
	resp, err := r.DoRequest(context.Background(), "GET", "http://localhost:8080", RequestOptions{Headers: h, Body: b})

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "application/json", resp.ContentType)
	require.Equal(t, b, resp.Body)
}

func TestGZIPRequest(t *testing.T) {
	transport := &MockTransport{}
	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	r := NewRequestor(&client)
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	h.Set("Content-Encoding", "gzip")
	b := []byte(`{"testKey0":0,"testKey1":"testString","testKey2":2,"testKey3":false,"testKey4":true,"testKey5":[0,1,2,3,4,5,6,7,8,9],"testKey6":{"zero":0,"one":1}}`)

	var buffer bytes.Buffer
	w := gzip.NewWriter(&buffer)
	w.Write(b)
	w.Close()

	resp, err := r.DoRequest(context.Background(), "GET", "http://localhost:8080", RequestOptions{Headers: h, Body: buffer.Bytes()})

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "application/json", resp.ContentType)
	require.Equal(t, b, resp.Body)
}

func TestCookieRequest(t *testing.T) {
	transport := &MockTransport{}
	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	r := NewRequestor(&client)
	cookie := http.Cookie{}
	cookies := []*http.Cookie{&cookie}

	resp, err := r.DoRequest(context.Background(), "GET", "http://localhost:8080", RequestOptions{Cookies: cookies})

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "text/plain", resp.ContentType)
}

func TestGZIPFail(t *testing.T) {
	transport := &MockTransport{}
	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	r := NewRequestor(&client)
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	h.Set("Content-Encoding", "gzip")
	b := []byte(`{"testKey0":0,"testKey1":"testString","testKey2":2,"testKey3":false,"testKey4":true,"testKey5":[0,1,2,3,4,5,6,7,8,9],"testKey6":{"zero":0,"one":1}}`)

	resp, err := r.DoRequest(context.Background(), "GET", "http://localhost:8080", RequestOptions{Headers: h, Body: b})

	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "application/json", resp.ContentType)
	require.NotEqual(t, b, resp.Body)
}
