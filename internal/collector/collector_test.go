package collector

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kr-ilya/mtproxy-exporter/pkg/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewMTProxyCollector(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pid\t12345"))
	}))
	defer server.Close()

	mtproxyClient, err := client.New(server.URL, 5*time.Second)
	assert.NoError(t, err)
	collector := NewMTProxyCollector(mtproxyClient)

	assert.NotNil(t, collector)
	assert.NotNil(t, collector.client)
}

func TestCollector_Describe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mtproxyClient, err := client.New(server.URL, 5*time.Second)
	assert.NoError(t, err)
	collector := NewMTProxyCollector(mtproxyClient)

	ch := make(chan *prometheus.Desc, 100)
	go func() {
		collector.Describe(ch)
		close(ch)
	}()

	count := 0
	for range ch {
		count++
	}

	assert.Greater(t, count, 0, "Should have registered metrics")
}

func TestCollector_Collect_Success(t *testing.T) {
	statsData := `pid	35
start_time	1768054068
current_time	1768055903
uptime	1835
active_connections	65
outbound_connections	64
inbound_connections	1
tcp_readv_bytes	244761
tcp_writev_bytes	436498`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(statsData))
	}))
	defer server.Close()

	mtproxyClient, err := client.New(server.URL, 5*time.Second)
	assert.NoError(t, err)
	collector := NewMTProxyCollector(mtproxyClient)

	// Create a registry and register the collector
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// Verify that metrics are collected
	count := testutil.CollectAndCount(collector)
	assert.Greater(t, count, 0, "Should have collected metrics")

	// Collect metrics and verify
	ch := make(chan prometheus.Metric, 100)
	go func() {
		collector.Collect(ch)
		close(ch)
	}()

	hasUpMetric := false
	for range ch {
		hasUpMetric = true
		break
	}
	assert.True(t, hasUpMetric, "Should have at least one metric")
}

func TestCollector_Collect_ServerDown(t *testing.T) {
	// Server that always returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	mtproxyClient, err := client.New(server.URL, 5*time.Second)
	assert.NoError(t, err)
	collector := NewMTProxyCollector(mtproxyClient)

	// Register the collector
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	// The collector should still work but report mtproxy_up as 0
	ch := make(chan prometheus.Metric, 100)
	go func() {
		collector.Collect(ch)
		close(ch)
	}()

	// Collect all metrics
	metrics := []prometheus.Metric{}
	for m := range ch {
		metrics = append(metrics, m)
	}

	assert.Greater(t, len(metrics), 0, "Should have at least the 'up' metric")
}

func TestCollector_Collect_ParseError(t *testing.T) {
	// Server returns invalid data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid data that cannot be parsed"))
	}))
	defer server.Close()

	mtproxyClient, err := client.New(server.URL, 5*time.Second)
	assert.NoError(t, err)
	collector := NewMTProxyCollector(mtproxyClient)

	ch := make(chan prometheus.Metric, 100)
	go func() {
		collector.Collect(ch)
		close(ch)
	}()

	// Should not panic and should return at least the 'up' metric
	count := 0
	for range ch {
		count++
	}
	assert.Greater(t, count, 0)
}

func BenchmarkCollector_Collect(b *testing.B) {
	statsData := `pid	35
uptime	1835
active_connections	65
tcp_readv_bytes	244761`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(statsData))
	}))
	defer server.Close()

	mtproxyClient, err := client.New(server.URL, 5*time.Second)
	if err != nil {
		b.Fatal(err)
	}
	collector := NewMTProxyCollector(mtproxyClient)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch := make(chan prometheus.Metric, 100)
		go func() {
			collector.Collect(ch)
			close(ch)
		}()
		for range ch {
		}
	}
}
