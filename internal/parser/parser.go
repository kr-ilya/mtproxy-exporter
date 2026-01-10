package parser

import (
	"bufio"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// Stats represents the parsed MTProxy statistics
type Stats struct {
	// General
	PID         int64
	StartTime   int64
	CurrentTime int64
	Uptime      int64

	// Idle and performance
	TotalIdleTime       float64
	AverageIdlePercent  float64
	RecentIdlePercent   float64
	ActiveNetworkEvents int64
	TimeAfterEpoll      float64
	EpollCalls          int64
	EpollIntr           int64

	// Connections
	ActiveConnections            int64
	ActiveDhConnections          int64
	OutboundConnections          int64
	ReadyOutboundConnections     int64
	ActiveOutboundConnections    int64
	OutboundConnectionsCreated   int64
	TotalConnectFailures         int64
	InboundConnections           int64
	ActiveInboundConnections     int64
	InboundConnectionsAccepted   int64
	ListeningConnections         int64
	UnusedConnectionsClosed      int64
	ReadyTargets                 int64
	AllocatedTargets             int64
	ActiveTargets                int64
	InactiveTargets              int64
	FreeTargets                  int64
	MaxConnections               int64
	ActiveSpecialConnections     int64
	MaxSpecialConnections        int64
	MaxAcceptRate                int64
	CurAcceptRateRemaining       float64
	MaxConnection                int64
	ConnGeneration               int64
	AllocatedConnections         int64
	AllocatedOutboundConnections int64
	AllocatedInboundConnections  int64
	AllocatedSocketConnections   int64
	AcceptCallsFailed            int64
	AcceptNonblockSetFailed      int64
	AcceptConnectionLimitFailed  int64
	AcceptRateLimitFailed        int64
	AcceptInitAcceptedFailed     int64

	// TCP stats
	TCPReadvCalls  int64
	TCPReadvIntr   int64
	TCPReadvBytes  int64
	TCPWritevCalls int64
	TCPWritevIntr  int64
	TCPWritevBytes int64
	FreeLaterSize  int64
	FreeLaterTotal int64

	// Raw messages
	RwmTotalMsgs     int64
	RwmTotalMsgParts int64

	// Buffers
	TotalUsedBuffersSize             int64
	TotalUsedBuffers                 int64
	AllocatedBufferBytes             int64
	BufferChunkAllocOps              int64
	AllocatedBufferChunks            int64
	MaxAllocatedBufferChunks         int64
	MaxBufferChunks                  int64
	MaxAllocatedBufferBytes          int64
	TotalNetworkBuffersUsedSize      int64
	TotalNetworkBuffersUsed          int64
	TotalNetworkBuffersAllocBytes    int64
	TotalNetworkBufferChunksAlloc    int64
	TotalNetworkBufferChunksAllocMax int64

	// RPC
	RpcQueriesReceived int64
	RpcAnswersError    int64
	RpcAnswersReceived int64
	RpcSentErrors      int64
	RpcSentAnswers     int64
	RpcSentQueries     int64
	TlInAllocated      int64
	TlOutAllocated     int64
	RpcQPS             float64
	DefaultRpcFlags    int64

	// Crypto
	AllocatedAesCrypto     int64
	AllocatedAesCryptoTemp int64
	AesPwdHash             string

	// Jobs
	JobTimersAllocated  int64
	JobsCreated         int64
	JobsActive          int64
	JobsAllocatedMemory int64
	TimerOps            int64
	TimerOpsScheduler   int64
	TotalTimers         int64

	// Load
	LoadAverageUser  float64
	LoadAverageSys   float64
	LoadAverageTotal float64
	LoadRecentUser   float64
	LoadRecentSys    float64
	LoadRecentTotal  float64
	MaxAverageUser   float64
	MaxAverageSys    float64
	MaxAverageTotal  float64
	MaxRecentUser    float64
	MaxRecentSys     float64
	MaxRecentTotal   float64

	// Thread loads (arrays)
	ThreadLoadAverageUser []float64
	ThreadLoadAverageSys  []float64
	ThreadLoadAverage     []float64
	ThreadLoadRecentUser  []float64
	ThreadLoadRecentSys   []float64
	ThreadLoadRecent      []float64

	// MP Queue
	MpqBlocksAllocated         int64
	MpqBlocksAllocatedMax      int64
	MpqBlocksAllocations       int64
	MpqBlocksTrueAllocations   int64
	MpqBlocksWasted            int64
	MpqBlocksPrepared          int64
	MpqSmallBlocksAllocated    int64
	MpqSmallBlocksAllocatedMax int64
	MpqActive                  int64
	MpqAllocated               int64

	// RPC Targets
	TotalRpcTargets              int64
	TotalConnectionsInRpcTargets int64

	// Memory
	VmsizeBytes int64
	VmrssBytes  int64
	VmdataBytes int64

	// Config
	ConfigFilename     string
	ConfigLoadedAt     int64
	ConfigSize         int64
	ConfigMD5          string
	ConfigAuthClusters int64

	// Queries
	Workers                 int64
	QueriesGet              int64
	QPSGet                  float64
	TotForwardedQueries     int64
	ExpiredForwardedQueries int64
	DroppedQueries          int64
	TotForwardedResponses   int64
	DroppedResponses        int64
	TotForwardedSimpleAcks  int64
	DroppedSimpleAcks       int64
	ActiveRpcsCreated       int64
	ActiveRpcs              int64
	RpcDroppedAnswers       int64
	RpcDroppedRunning       int64

	// Total metrics
	WindowClamp                       int64
	TotalReadyTargets                 int64
	TotalAllocatedTargets             int64
	TotalDeclaredTargets              int64
	TotalInactiveTargets              int64
	TotalConnections                  int64
	TotalEncryptedConnections         int64
	TotalAllocatedConnections         int64
	TotalAllocatedOutboundConnections int64
	TotalAllocatedInboundConnections  int64
	TotalAllocatedSocketConnections   int64
	TotalDhConnections                int64
	TotalSpecialConnections           int64
	TotalMaxSpecialConnections        int64
	TotalActiveNetworkEvents          int64
	ExtConnections                    int64
	ExtConnectionsCreated             int64

	// Errors
	MtprotoProxyErrors     int64
	ConnectionsFailedLru   int64
	ConnectionsFailedFlood int64

	// HTTP
	HttpConnections    int64
	PendingHttpQueries int64
	HttpQueries        int64
	HttpBadHeaders     int64
	HttpQPS            float64

	// Proxy
	ProxyMode   int64
	ProxyTagSet int64

	// Version
	Version string
}

// Parse parses the MTProxy stats output into a Stats struct
func Parse(data string) (*Stats, error) {
	stats := &Stats{}
	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Skip section markers
		if strings.HasPrefix(line, ">>>>>>") || strings.HasPrefix(line, "<<<<<<") {
			continue
		}

		// Parse key-value pairs
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		values := parts[1:]

		if err := parseField(stats, key, values); err != nil {
			// Log parsing errors for debugging
			slog.Debug("Failed to parse field", "key", key, "error", err)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning stats: %w", err)
	}

	return stats, nil
}

// parseField parses a single field and sets it in the Stats struct
func parseField(stats *Stats, key string, values []string) error {
	if len(values) == 0 {
		return fmt.Errorf("no values for key %s", key)
	}

	// Helper functions
	parseInt64 := func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 64)
	}

	parseFloat64 := func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	}

	parseFloat64Array := func(values []string) ([]float64, error) {
		result := make([]float64, 0, len(values))
		for _, v := range values {
			f, err := parseFloat64(v)
			if err != nil {
				return nil, err
			}
			result = append(result, f)
		}
		return result, nil
	}

	// Parse based on key
	switch key {
	case "pid":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.PID = v

	case "start_time":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.StartTime = v

	case "current_time":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.CurrentTime = v

	case "uptime":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.Uptime = v

	case "tot_idle_time":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.TotalIdleTime = v

	case "average_idle_percent":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.AverageIdlePercent = v

	case "recent_idle_percent":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.RecentIdlePercent = v

	case "active_network_events":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ActiveNetworkEvents = v

	case "time_after_epoll":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.TimeAfterEpoll = v

	case "epoll_calls":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.EpollCalls = v

	case "epoll_intr":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.EpollIntr = v

	case "active_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ActiveConnections = v

	case "active_dh_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ActiveDhConnections = v

	case "outbound_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.OutboundConnections = v

	case "ready_outbound_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ReadyOutboundConnections = v

	case "active_outbound_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ActiveOutboundConnections = v

	case "outbound_connections_created":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.OutboundConnectionsCreated = v

	case "total_connect_failures":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalConnectFailures = v

	case "inbound_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.InboundConnections = v

	case "active_inbound_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ActiveInboundConnections = v

	case "inbound_connections_accepted":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.InboundConnectionsAccepted = v

	case "listening_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ListeningConnections = v

	case "unused_connections_closed":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.UnusedConnectionsClosed = v

	case "ready_targets":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ReadyTargets = v

	case "allocated_targets":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AllocatedTargets = v

	case "active_targets":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ActiveTargets = v

	case "inactive_targets":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.InactiveTargets = v

	case "free_targets":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.FreeTargets = v

	case "max_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MaxConnections = v

	case "active_special_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ActiveSpecialConnections = v

	case "max_special_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MaxSpecialConnections = v

	case "max_accept_rate":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MaxAcceptRate = v

	case "cur_accept_rate_remaining":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.CurAcceptRateRemaining = v

	case "max_connection":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MaxConnection = v

	case "conn_generation":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ConnGeneration = v

	case "allocated_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AllocatedConnections = v

	case "allocated_outbound_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AllocatedOutboundConnections = v

	case "allocated_inbound_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AllocatedInboundConnections = v

	case "allocated_socket_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AllocatedSocketConnections = v

	case "tcp_readv_calls":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TCPReadvCalls = v

	case "tcp_readv_intr":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TCPReadvIntr = v

	case "tcp_readv_bytes":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TCPReadvBytes = v

	case "tcp_writev_calls":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TCPWritevCalls = v

	case "tcp_writev_intr":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TCPWritevIntr = v

	case "tcp_writev_bytes":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TCPWritevBytes = v

	case "free_later_size":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.FreeLaterSize = v

	case "free_later_total":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.FreeLaterTotal = v

	case "accept_calls_failed":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AcceptCallsFailed = v

	case "accept_nonblock_set_failed":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AcceptNonblockSetFailed = v

	case "accept_connection_limit_failed":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AcceptConnectionLimitFailed = v

	case "accept_rate_limit_failed":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AcceptRateLimitFailed = v

	case "accept_init_accepted_failed":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AcceptInitAcceptedFailed = v

	case "rwm_total_msgs":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.RwmTotalMsgs = v

	case "rwm_total_msg_parts":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.RwmTotalMsgParts = v

	case "total_used_buffers_size":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalUsedBuffersSize = v

	case "total_used_buffers":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalUsedBuffers = v

	case "allocated_buffer_bytes":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AllocatedBufferBytes = v

	case "buffer_chunk_alloc_ops":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.BufferChunkAllocOps = v

	case "allocated_buffer_chunks":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AllocatedBufferChunks = v

	case "max_allocated_buffer_chunks":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MaxAllocatedBufferChunks = v

	case "max_buffer_chunks":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MaxBufferChunks = v

	case "max_allocated_buffer_bytes":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MaxAllocatedBufferBytes = v

	case "rpc_queries_received":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.RpcQueriesReceived = v

	case "rpc_answers_error":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.RpcAnswersError = v

	case "rpc_answers_received":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.RpcAnswersReceived = v

	case "rpc_sent_errors":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.RpcSentErrors = v

	case "rpc_sent_answers":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.RpcSentAnswers = v

	case "rpc_sent_queries":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.RpcSentQueries = v

	case "tl_in_allocated":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TlInAllocated = v

	case "tl_out_allocated":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TlOutAllocated = v

	case "rpc_qps":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.RpcQPS = v

	case "default_rpc_flags":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.DefaultRpcFlags = v

	case "allocated_aes_crypto":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AllocatedAesCrypto = v

	case "allocated_aes_crypto_temp":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.AllocatedAesCryptoTemp = v

	case "aes_pwd_hash":
		stats.AesPwdHash = values[0]

	case "thread_average_idle_percent":
		arr, err := parseFloat64Array(values)
		if err != nil {
			return err
		}
		// Store first non-zero value as general metric
		for _, v := range arr {
			if v > 0 {
				// Already captured in average_idle_percent
				break
			}
		}

	case "thread_recent_idle_percent":
		arr, err := parseFloat64Array(values)
		if err != nil {
			return err
		}
		// Store first non-zero value as general metric
		for _, v := range arr {
			if v > 0 {
				// Already captured in recent_idle_percent
				break
			}
		}

	case "thread_load_average_user":
		arr, err := parseFloat64Array(values)
		if err != nil {
			return err
		}
		stats.ThreadLoadAverageUser = arr

	case "thread_load_average_sys":
		arr, err := parseFloat64Array(values)
		if err != nil {
			return err
		}
		stats.ThreadLoadAverageSys = arr

	case "thread_load_average":
		arr, err := parseFloat64Array(values)
		if err != nil {
			return err
		}
		stats.ThreadLoadAverage = arr

	case "thread_load_recent_user":
		arr, err := parseFloat64Array(values)
		if err != nil {
			return err
		}
		stats.ThreadLoadRecentUser = arr

	case "thread_load_recent_sys":
		arr, err := parseFloat64Array(values)
		if err != nil {
			return err
		}
		stats.ThreadLoadRecentSys = arr

	case "thread_load_recent":
		arr, err := parseFloat64Array(values)
		if err != nil {
			return err
		}
		stats.ThreadLoadRecent = arr

	case "load_average_user":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.LoadAverageUser = v

	case "load_average_sys":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.LoadAverageSys = v

	case "load_average_total":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.LoadAverageTotal = v

	case "load_recent_user":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.LoadRecentUser = v

	case "load_recent_sys":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.LoadRecentSys = v

	case "load_recent_total":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.LoadRecentTotal = v

	case "max_average_user":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.MaxAverageUser = v

	case "max_average_sys":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.MaxAverageSys = v

	case "max_average_total":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.MaxAverageTotal = v

	case "max_recent_user":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.MaxRecentUser = v

	case "max_recent_sys":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.MaxRecentSys = v

	case "max_recent_total":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.MaxRecentTotal = v

	case "job_timers_allocated":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.JobTimersAllocated = v

	case "jobs_created":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.JobsCreated = v

	case "jobs_active":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.JobsActive = v

	case "jobs_allocated_memory":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.JobsAllocatedMemory = v

	case "timer_ops":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TimerOps = v

	case "timer_ops_scheduler":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TimerOpsScheduler = v

	case "mpq_blocks_allocated":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MpqBlocksAllocated = v

	case "mpq_blocks_allocated_max":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MpqBlocksAllocatedMax = v

	case "mpq_blocks_allocations":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MpqBlocksAllocations = v

	case "mpq_blocks_true_allocations":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MpqBlocksTrueAllocations = v

	case "mpq_blocks_wasted":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MpqBlocksWasted = v

	case "mpq_blocks_prepared":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MpqBlocksPrepared = v

	case "mpq_small_blocks_allocated":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MpqSmallBlocksAllocated = v

	case "mpq_small_blocks_allocated_max":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MpqSmallBlocksAllocatedMax = v

	case "mpq_active":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MpqActive = v

	case "mpq_allocated":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MpqAllocated = v

	case "total_timers":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalTimers = v

	case "total_rpc_targets":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalRpcTargets = v

	case "total_connections_in_rpc_targets":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalConnectionsInRpcTargets = v

	case "vmsize_bytes":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.VmsizeBytes = v

	case "vmrss_bytes":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.VmrssBytes = v

	case "vmdata_bytes":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.VmdataBytes = v

	case "config_filename":
		stats.ConfigFilename = values[0]

	case "config_loaded_at":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ConfigLoadedAt = v

	case "config_size":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ConfigSize = v

	case "config_md5":
		stats.ConfigMD5 = values[0]

	case "config_auth_clusters":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ConfigAuthClusters = v

	case "workers":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.Workers = v

	case "queries_get":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.QueriesGet = v

	case "qps_get":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.QPSGet = v

	case "tot_forwarded_queries":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotForwardedQueries = v

	case "expired_forwarded_queries":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ExpiredForwardedQueries = v

	case "dropped_queries":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.DroppedQueries = v

	case "tot_forwarded_responses":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotForwardedResponses = v

	case "dropped_responses":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.DroppedResponses = v

	case "tot_forwarded_simple_acks":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotForwardedSimpleAcks = v

	case "dropped_simple_acks":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.DroppedSimpleAcks = v

	case "active_rpcs_created":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ActiveRpcsCreated = v

	case "active_rpcs":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ActiveRpcs = v

	case "rpc_dropped_answers":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.RpcDroppedAnswers = v

	case "rpc_dropped_running":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.RpcDroppedRunning = v

	case "window_clamp":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.WindowClamp = v

	case "total_ready_targets":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalReadyTargets = v

	case "total_allocated_targets":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalAllocatedTargets = v

	case "total_declared_targets":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalDeclaredTargets = v

	case "total_inactive_targets":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalInactiveTargets = v

	case "total_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalConnections = v

	case "total_encrypted_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalEncryptedConnections = v

	case "total_allocated_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalAllocatedConnections = v

	case "total_allocated_outbound_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalAllocatedOutboundConnections = v

	case "total_allocated_inbound_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalAllocatedInboundConnections = v

	case "total_allocated_socket_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalAllocatedSocketConnections = v

	case "total_dh_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalDhConnections = v

	case "total_special_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalSpecialConnections = v

	case "total_max_special_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalMaxSpecialConnections = v

	case "total_active_network_events":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalActiveNetworkEvents = v

	case "total_network_buffers_used_size":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalNetworkBuffersUsedSize = v

	case "total_network_buffers_allocated_bytes":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalNetworkBuffersAllocBytes = v

	case "total_network_buffers_used":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalNetworkBuffersUsed = v

	case "total_network_buffer_chunks_allocated":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalNetworkBufferChunksAlloc = v

	case "total_network_buffer_chunks_allocated_max":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.TotalNetworkBufferChunksAllocMax = v

	case "mtproto_proxy_errors":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.MtprotoProxyErrors = v

	case "connections_failed_lru":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ConnectionsFailedLru = v

	case "connections_failed_flood":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ConnectionsFailedFlood = v

	case "http_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.HttpConnections = v

	case "pending_http_queries":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.PendingHttpQueries = v

	case "http_queries":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.HttpQueries = v

	case "http_bad_headers":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.HttpBadHeaders = v

	case "http_qps":
		v, err := parseFloat64(values[0])
		if err != nil {
			return err
		}
		stats.HttpQPS = v

	case "proxy_mode":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ProxyMode = v

	case "proxy_tag_set":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ProxyTagSet = v

	case "ext_connections":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ExtConnections = v

	case "ext_connections_created":
		v, err := parseInt64(values[0])
		if err != nil {
			return err
		}
		stats.ExtConnectionsCreated = v

	case "version":
		stats.Version = strings.Join(values, " ")

	default:
		// Ignore unknown keys
	}

	return nil
}
