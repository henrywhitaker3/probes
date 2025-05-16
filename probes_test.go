package probes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	Demo Subject = "demo"
)

func TestItErrorsWithUnknownSubjects(t *testing.T) {
	p := New(ProbeOpts{})
	require.NotNil(t, p.Healthy(Demo))
	require.NotNil(t, p.Unhealthy(Demo))
	require.NotNil(t, p.Ready(Demo))
	require.NotNil(t, p.NotReady(Demo))
	require.NotNil(t, p.Started(Demo))
}

func TestItIsStartedWithNoSubjects(t *testing.T) {
	p := New(ProbeOpts{})
	require.True(t, p.IsStarted())
}

func TestItIsReadyWithNoSubjects(t *testing.T) {
	p := New(ProbeOpts{})
	require.True(t, p.IsReady())
}

func TestItIsHealthyWithNoSubjects(t *testing.T) {
	p := New(ProbeOpts{})
	require.True(t, p.IsHealthy())
}

func TestItIsNotStartedWithASubject(t *testing.T) {
	p := New(ProbeOpts{}).WithStartups(Demo)
	require.False(t, p.IsStarted())
	require.Nil(t, p.Started(Demo))
	require.True(t, p.IsStarted())
}

func TestItIsNotReadyWithASubject(t *testing.T) {
	p := New(ProbeOpts{}).WithReadies(Demo)
	require.False(t, p.IsReady())
	require.Nil(t, p.Ready(Demo))
	require.True(t, p.IsReady())
	require.Nil(t, p.NotReady(Demo))
	require.False(t, p.IsReady())
}

func TestItIsHealthyWithSubject(t *testing.T) {
	p := New(ProbeOpts{}).WithHealthies(Demo)
	require.True(t, p.IsHealthy())
	require.Nil(t, p.Healthy(Demo))
	require.True(t, p.IsHealthy())
	require.Nil(t, p.Unhealthy(Demo))
	require.False(t, p.IsHealthy())
}
