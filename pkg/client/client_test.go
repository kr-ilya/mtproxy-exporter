package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		wantErr bool
	}{
		{
			name:    "valid http URL",
			baseURL: "http://localhost:8888",
			wantErr: false,
		},
		{
			name:    "valid https URL",
			baseURL: "https://localhost:8888",
			wantErr: false,
		},
		{
			name:    "invalid URL",
			baseURL: "://invalid",
			wantErr: true,
		},
		{
			name:    "invalid scheme",
			baseURL: "ftp://localhost:8888",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.baseURL, 10*time.Second)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, client)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, client)
				assert.Equal(t, tt.baseURL, client.baseURL)
				assert.NotNil(t, client.httpClient)
			}
		})
	}
}

func TestGetStats(t *testing.T) {
	tests := []struct {
		name       string
		serverFunc func(w http.ResponseWriter, r *http.Request)
		wantErr    bool
		wantData   string
	}{
		{
			name: "successful request",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/stats", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("pid\t12345\nuptime\t1000"))
			},
			wantErr:  false,
			wantData: "pid\t12345\nuptime\t1000",
		},
		{
			name: "server returns 404",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantErr: true,
		},
		{
			name: "server returns 500",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "empty response",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantErr:  false,
			wantData: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverFunc))
			defer server.Close()

			client, err := New(server.URL, 5*time.Second)
			require.NoError(t, err)
			ctx := context.Background()

			data, err := client.GetStats(ctx)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantData, data)
		})
	}
}

func TestGetStatsWithContext(t *testing.T) {
	t.Run("context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client, err := New(server.URL, 5*time.Second)
		require.NoError(t, err)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err = client.GetStats(ctx)
		require.Error(t, err)
	})
}

func TestGetStatsTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := New(server.URL, 50*time.Millisecond)
	require.NoError(t, err)
	ctx := context.Background()

	_, err = client.GetStats(ctx)
	require.Error(t, err)
}

func BenchmarkGetStats(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pid\t12345\nuptime\t1000"))
	}))
	defer server.Close()

	client, err := New(server.URL, 5*time.Second)
	if err != nil {
		b.Fatal(err)
	}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.GetStats(ctx)
	}
}
