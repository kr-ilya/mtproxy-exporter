package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Stats
		wantErr bool
	}{
		{
			name: "basic stats parsing",
			input: `pid	12345
start_time	1768054068
current_time	1768055903
uptime	1835
tot_idle_time	1827.153
average_idle_percent	99.572
recent_idle_percent	98.688
active_network_events	0
time_after_epoll	0.000040
epoll_calls	1528389
active_connections	65
outbound_connections	64
inbound_connections	1
tcp_readv_bytes	244761
tcp_writev_bytes	436498
version mtproxy-0.02 compiled at Jan 10 2026`,
			want: &Stats{
				PID:                 12345,
				StartTime:           1768054068,
				CurrentTime:         1768055903,
				Uptime:              1835,
				TotalIdleTime:       1827.153,
				AverageIdlePercent:  99.572,
				RecentIdlePercent:   98.688,
				ActiveNetworkEvents: 0,
				TimeAfterEpoll:      0.000040,
				EpollCalls:          1528389,
				ActiveConnections:   65,
				OutboundConnections: 64,
				InboundConnections:  1,
				TCPReadvBytes:       244761,
				TCPWritevBytes:      436498,
				Version:             "mtproxy-0.02 compiled at Jan 10 2026",
			},
			wantErr: false,
		},
		{
			name: "with section markers",
			input: `pid	100
>>>>>>connections>>>>>> start
active_connections	50
<<<<<<connections<<<<<< end
>>>>>>jobs>>>>>> start
jobs_active	10
<<<<<<jobs<<<<<<	end`,
			want: &Stats{
				PID:               100,
				ActiveConnections: 50,
				JobsActive:        10,
			},
			wantErr: false,
		},
		{
			name: "with float values",
			input: `average_idle_percent	99.572
load_average_user	1.034
load_average_sys	1.319
load_average_total	2.354`,
			want: &Stats{
				AverageIdlePercent: 99.572,
				LoadAverageUser:    1.034,
				LoadAverageSys:     1.319,
				LoadAverageTotal:   2.354,
			},
			wantErr: false,
		},
		{
			name: "with thread arrays",
			input: `thread_load_average_user	0.000 1.000 2.000 3.000
thread_load_average_sys	0.100 1.100 2.100 3.100`,
			want: &Stats{
				ThreadLoadAverageUser: []float64{0.000, 1.000, 2.000, 3.000},
				ThreadLoadAverageSys:  []float64{0.100, 1.100, 2.100, 3.100},
			},
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			want:    &Stats{},
			wantErr: false,
		},
		{
			name: "full example from MTProxy",
			input: `pid	35
start_time	1768054068
current_time	1768055903
uptime	1835
tot_idle_time	1827.153
average_idle_percent	99.572
recent_idle_percent	98.688
active_network_events	0
>>>>>>connections>>>>>> start
active_connections	65
active_dh_connections	0
outbound_connections	64
outbound_connections_created	3213
total_connect_failures	568
inbound_connections	1
<<<<<<connections<<<<<< end
tcp_readv_calls	13256
tcp_readv_bytes	244761
tcp_writev_calls	5318
tcp_writev_bytes	436498
>>>>>>raw_msg>>>>>>	start
rwm_total_msgs	359
rwm_total_msg_parts	289
<<<<<<raw_msg<<<<<<	end
>>>>>>jobs>>>>>>	start
jobs_created	6515
jobs_active	153
<<<<<<jobs<<<<<<	end
vmsize_bytes	803565568
vmrss_bytes	9830400
workers	1
http_queries	32
version mtproxy-0.02 compiled at Jan 10 2026`,
			want: &Stats{
				PID:                        35,
				StartTime:                  1768054068,
				CurrentTime:                1768055903,
				Uptime:                     1835,
				TotalIdleTime:              1827.153,
				AverageIdlePercent:         99.572,
				RecentIdlePercent:          98.688,
				ActiveNetworkEvents:        0,
				ActiveConnections:          65,
				ActiveDhConnections:        0,
				OutboundConnections:        64,
				OutboundConnectionsCreated: 3213,
				TotalConnectFailures:       568,
				InboundConnections:         1,
				TCPReadvCalls:              13256,
				TCPReadvBytes:              244761,
				TCPWritevCalls:             5318,
				TCPWritevBytes:             436498,
				RwmTotalMsgs:               359,
				RwmTotalMsgParts:           289,
				JobsCreated:                6515,
				JobsActive:                 153,
				VmsizeBytes:                803565568,
				VmrssBytes:                 9830400,
				Workers:                    1,
				HttpQueries:                32,
				Version:                    "mtproxy-0.02 compiled at Jan 10 2026",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseField(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		values  []string
		check   func(*testing.T, *Stats)
		wantErr bool
	}{
		{
			name:   "parse integer",
			key:    "pid",
			values: []string{"12345"},
			check: func(t *testing.T, s *Stats) {
				assert.Equal(t, int64(12345), s.PID)
			},
			wantErr: false,
		},
		{
			name:   "parse float",
			key:    "average_idle_percent",
			values: []string{"99.572"},
			check: func(t *testing.T, s *Stats) {
				assert.Equal(t, 99.572, s.AverageIdlePercent)
			},
			wantErr: false,
		},
		{
			name:   "parse string",
			key:    "config_filename",
			values: []string{"/data/proxy-multi.conf"},
			check: func(t *testing.T, s *Stats) {
				assert.Equal(t, "/data/proxy-multi.conf", s.ConfigFilename)
			},
			wantErr: false,
		},
		{
			name:   "parse multi-word string",
			key:    "version",
			values: []string{"mtproxy-0.02", "compiled", "at", "Jan", "10", "2026"},
			check: func(t *testing.T, s *Stats) {
				assert.Equal(t, "mtproxy-0.02 compiled at Jan 10 2026", s.Version)
			},
			wantErr: false,
		},
		{
			name:    "parse invalid integer",
			key:     "pid",
			values:  []string{"invalid"},
			wantErr: true,
		},
		{
			name:    "parse invalid float",
			key:     "average_idle_percent",
			values:  []string{"invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := &Stats{}
			err := parseField(stats, tt.key, tt.values)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.check != nil {
				tt.check(t, stats)
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	input := `pid	35
start_time	1768054068
current_time	1768055903
uptime	1835
tot_idle_time	1827.153
average_idle_percent	99.572
recent_idle_percent	98.688
active_connections	65
outbound_connections	64
inbound_connections	1
tcp_readv_bytes	244761
tcp_writev_bytes	436498
jobs_created	6515
jobs_active	153
vmsize_bytes	803565568
version mtproxy-0.02 compiled`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Parse(input)
	}
}
