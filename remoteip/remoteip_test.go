package remoteip

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Run("X-Forwarded-For", func(t *testing.T) {
		r := &http.Request{Header: http.Header{}, RemoteAddr: "10.10.1.247"}
		r.Header.Set("X-Forwarded-For", "127.0.0.1, 54.146.177.1,192.168.0.1,54.146.177.2, 10.10.1.246,192.168.1.1")
		require.Equal(t, "54.146.177.2", Get(r))
	})

	t.Run("With Ports", func(t *testing.T) {
		r := &http.Request{Header: http.Header{}, RemoteAddr: "10.10.1.247"}
		r.Header.Set("X-Forwarded-For", "127.0.0.1, 54.146.177.1,192.168.0.1:20,54.146.177.2:8080, 10.10.1.246,192.168.1.1")
		require.Equal(t, "54.146.177.2", Get(r))
	})

	t.Run("Prefer RemoteAddr", func(t *testing.T) {
		r := &http.Request{Header: http.Header{}, RemoteAddr: "54.146.177.3"}
		r.Header.Set("X-Forwarded-For", "127.0.0.1, 54.146.177.1,192.168.0.1,54.146.177.2, 10.10.1.246,192.168.1.1")
		require.Equal(t, "54.146.177.3", Get(r))
	})

	t.Run("Client", func(t *testing.T) {
		r := &http.Request{Header: http.Header{}, RemoteAddr: "54.146.177.172"}
		require.Equal(t, "54.146.177.172", Get(r))
	})

	t.Run("Client", func(t *testing.T) {
		r := &http.Request{Header: http.Header{}, RemoteAddr: "127.0.0.1"}
		require.Equal(t, "127.0.0.1", Get(r))
	})
}
