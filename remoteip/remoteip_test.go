package remoteip

import (
	"net/http"
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	testCases := []struct {
		Raw      string
		Expected string
	}{
		{"192.168.1.1", "192.168.1.1"},
		{"192.168.1.1:4080", "192.168.1.1"},
		{"001:db8:3333:4444:5555:7777:8888:9999", "1:db8:3333:4444:5555:7777:8888:9999"},
		{"[002:db8:3333:4444:5555:7777:8888:9999]:8080", "2:db8:3333:4444:5555:7777:8888:9999"},
		{"[003:db8:3333:4444:5555:7777:8888:9999]:", "3:db8:3333:4444:5555:7777:8888:9999"},
		{"[004:db8:3333:4444:5555:7777:8888:9999]", "-"},
		{"005:db8:3333:4444:5555:7777:8888:9999:8080", "-"},
		{"", "-"},
		{"192.162.1111", "-"},
	}

	for n, tc := range testCases {
		t.Run("X-Client-IP-"+cast.ToString(n), func(t *testing.T) {
			r := &http.Request{Header: http.Header{}}

			r.Header.Set("X-Client-IP", tc.Raw)
			assert.Equal(t, tc.Expected, Get(r))
		})
	}

	for n, tc := range testCases {
		t.Run("X-Forwarded-For-"+cast.ToString(n), func(t *testing.T) {
			r := &http.Request{Header: http.Header{}}

			r.Header.Set("X-Forwarded-For", tc.Raw)
			assert.Equal(t, tc.Expected, Get(r))
		})
	}

	for n, tc := range testCases {
		t.Run("X-Real-IP-"+cast.ToString(n), func(t *testing.T) {
			r := &http.Request{Header: http.Header{}}

			r.Header.Set("X-Real-IP", tc.Raw)
			assert.Equal(t, tc.Expected, Get(r))
		})
	}

	for n, tc := range testCases {
		t.Run("RemoteAddr-"+cast.ToString(n), func(t *testing.T) {
			r := &http.Request{RemoteAddr: tc.Raw}
			assert.Equal(t, tc.Expected, Get(r))
		})
	}
}
