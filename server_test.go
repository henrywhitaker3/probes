package probes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServerRespondsToHealthProbes(t *testing.T) {
	p := New(ProbeOpts{}).WithHealthies(Demo)
	srv := server(p)
	defer srv.Close()

	resp := get(t, srv, "/healthz")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Nil(t, p.Unhealthy(Demo))

	resp = get(t, srv, "/healthz")
	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	require.Nil(t, p.Healthy(Demo))

	resp = get(t, srv, "/healthz")
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServerRepsondsToStartupProbes(t *testing.T) {
	p := New(ProbeOpts{}).WithStartups(Demo)
	srv := server(p)
	defer srv.Close()

	resp := get(t, srv, "/startedz")
	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	require.Nil(t, p.Started(Demo))

	resp = get(t, srv, "/startedz")
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServerRespondsToReadinessProbes(t *testing.T) {
	p := New(ProbeOpts{}).WithReadies(Demo)
	srv := server(p)
	defer srv.Close()

	resp := get(t, srv, "/readyz")
	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	require.Nil(t, p.Ready(Demo))

	resp = get(t, srv, "/readyz")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	require.Nil(t, p.NotReady(Demo))

	resp = get(t, srv, "/readyz")
	require.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
}

func get(t *testing.T, srv *httptest.Server, url string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", srv.URL, url), nil)
	require.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	return resp
}

func server(p *Probes) *httptest.Server {
	srv := NewServer(ServerOpts{
		Probes: p,
	})
	testsrv := httptest.NewServer(srv.server.Handler)
	return testsrv
}
