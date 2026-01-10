package collector

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/kr-ilya/mtproxy-exporter/internal/parser"
	"github.com/kr-ilya/mtproxy-exporter/pkg/client"
	"github.com/prometheus/client_golang/prometheus"
)

// MTProxyCollector collects metrics from MTProxy
type MTProxyCollector struct {
	client *client.Client

	// General metrics
	up        *prometheus.Desc
	uptime    *prometheus.Desc
	startTime *prometheus.Desc

	// Exporter metrics
	scrapeDuration *prometheus.Desc
	scrapeErrors   *prometheus.Desc
	lastScrapeTime *prometheus.Desc

	// Connection metrics
	activeConnections          *prometheus.Desc
	inboundConnections         *prometheus.Desc
	outboundConnections        *prometheus.Desc
	connectionsCreated         *prometheus.Desc
	connectFailures            *prometheus.Desc
	activeInboundConnections   *prometheus.Desc
	activeOutboundConnections  *prometheus.Desc
	readyOutboundConnections   *prometheus.Desc
	inboundConnectionsAccepted *prometheus.Desc
	allocatedConnections       *prometheus.Desc

	// Network metrics
	tcpReadBytes   *prometheus.Desc
	tcpWriteBytes  *prometheus.Desc
	tcpReadvCalls  *prometheus.Desc
	tcpWritevCalls *prometheus.Desc

	// Performance metrics
	idlePercent         *prometheus.Desc
	recentIdlePercent   *prometheus.Desc
	epollCalls          *prometheus.Desc
	activeNetworkEvents *prometheus.Desc

	// Message metrics
	messagesTotal      *prometheus.Desc
	queriesForwarded   *prometheus.Desc
	responsesForwarded *prometheus.Desc
	queriesDropped     *prometheus.Desc
	responsesDropped   *prometheus.Desc

	// Buffer metrics
	bufferUsedBytes      *prometheus.Desc
	bufferAllocatedBytes *prometheus.Desc
	buffersUsed          *prometheus.Desc

	// Load metrics
	loadAverageUser  *prometheus.Desc
	loadAverageSys   *prometheus.Desc
	loadAverageTotal *prometheus.Desc
	loadRecentUser   *prometheus.Desc
	loadRecentSys    *prometheus.Desc
	loadRecentTotal  *prometheus.Desc
	threadLoadUser   *prometheus.Desc
	threadLoadSys    *prometheus.Desc

	// Memory metrics
	memoryVmsize *prometheus.Desc
	memoryVmrss  *prometheus.Desc
	memoryVmdata *prometheus.Desc

	// Jobs metrics
	jobsActive  *prometheus.Desc
	jobsCreated *prometheus.Desc

	// HTTP metrics
	httpConnections *prometheus.Desc
	httpQueries     *prometheus.Desc
	httpQPS         *prometheus.Desc

	// Error metrics
	proxyErrors            *prometheus.Desc
	connectionsFailedLru   *prometheus.Desc
	connectionsFailedFlood *prometheus.Desc

	// Target metrics
	readyTargets     *prometheus.Desc
	allocatedTargets *prometheus.Desc
	activeTargets    *prometheus.Desc
	inactiveTargets  *prometheus.Desc
}

// NewMTProxyCollector creates a new MTProxy collector
func NewMTProxyCollector(client *client.Client) *MTProxyCollector {
	return &MTProxyCollector{
		client: client,

		up: prometheus.NewDesc(
			"mtproxy_up",
			"Whether the MTProxy is up (1) or down (0)",
			nil, nil,
		),
		uptime: prometheus.NewDesc(
			"mtproxy_uptime_seconds",
			"Uptime of the MTProxy in seconds",
			nil, nil,
		),
		startTime: prometheus.NewDesc(
			"mtproxy_start_time",
			"Unix timestamp when the MTProxy started",
			nil, nil,
		),

		// Exporter metrics
		scrapeDuration: prometheus.NewDesc(
			"mtproxy_exporter_scrape_duration_seconds",
			"Duration of the scrape in seconds",
			nil, nil,
		),
		scrapeErrors: prometheus.NewDesc(
			"mtproxy_exporter_scrape_errors_total",
			"Total number of scrape errors",
			nil, nil,
		),
		lastScrapeTime: prometheus.NewDesc(
			"mtproxy_exporter_last_scrape_timestamp",
			"Unix timestamp of the last successful scrape",
			nil, nil,
		),

		// Connections
		activeConnections: prometheus.NewDesc(
			"mtproxy_active_connections",
			"Number of active connections",
			nil, nil,
		),
		inboundConnections: prometheus.NewDesc(
			"mtproxy_inbound_connections",
			"Number of inbound connections",
			nil, nil,
		),
		outboundConnections: prometheus.NewDesc(
			"mtproxy_outbound_connections",
			"Number of outbound connections",
			nil, nil,
		),
		connectionsCreated: prometheus.NewDesc(
			"mtproxy_connections_created_total",
			"Total number of outbound connections created",
			nil, nil,
		),
		connectFailures: prometheus.NewDesc(
			"mtproxy_connect_failures_total",
			"Total number of connection failures",
			nil, nil,
		),
		activeInboundConnections: prometheus.NewDesc(
			"mtproxy_active_inbound_connections",
			"Number of active inbound connections",
			nil, nil,
		),
		activeOutboundConnections: prometheus.NewDesc(
			"mtproxy_active_outbound_connections",
			"Number of active outbound connections",
			nil, nil,
		),
		readyOutboundConnections: prometheus.NewDesc(
			"mtproxy_ready_outbound_connections",
			"Number of ready outbound connections",
			nil, nil,
		),
		inboundConnectionsAccepted: prometheus.NewDesc(
			"mtproxy_inbound_connections_accepted_total",
			"Total number of accepted inbound connections",
			nil, nil,
		),
		allocatedConnections: prometheus.NewDesc(
			"mtproxy_allocated_connections",
			"Number of allocated connections",
			nil, nil,
		),

		// Network
		tcpReadBytes: prometheus.NewDesc(
			"mtproxy_tcp_read_bytes_total",
			"Total bytes read via TCP",
			nil, nil,
		),
		tcpWriteBytes: prometheus.NewDesc(
			"mtproxy_tcp_write_bytes_total",
			"Total bytes written via TCP",
			nil, nil,
		),
		tcpReadvCalls: prometheus.NewDesc(
			"mtproxy_tcp_readv_calls_total",
			"Total readv system calls",
			nil, nil,
		),
		tcpWritevCalls: prometheus.NewDesc(
			"mtproxy_tcp_writev_calls_total",
			"Total writev system calls",
			nil, nil,
		),

		// Performance
		idlePercent: prometheus.NewDesc(
			"mtproxy_idle_percent",
			"Average idle percentage",
			nil, nil,
		),
		recentIdlePercent: prometheus.NewDesc(
			"mtproxy_recent_idle_percent",
			"Recent idle percentage",
			nil, nil,
		),
		epollCalls: prometheus.NewDesc(
			"mtproxy_epoll_calls_total",
			"Total epoll system calls",
			nil, nil,
		),
		activeNetworkEvents: prometheus.NewDesc(
			"mtproxy_active_network_events",
			"Number of active network events",
			nil, nil,
		),

		// Messages
		messagesTotal: prometheus.NewDesc(
			"mtproxy_messages_total",
			"Total number of messages processed",
			nil, nil,
		),
		queriesForwarded: prometheus.NewDesc(
			"mtproxy_queries_forwarded_total",
			"Total number of forwarded queries",
			nil, nil,
		),
		responsesForwarded: prometheus.NewDesc(
			"mtproxy_responses_forwarded_total",
			"Total number of forwarded responses",
			nil, nil,
		),
		queriesDropped: prometheus.NewDesc(
			"mtproxy_queries_dropped_total",
			"Total number of dropped queries",
			nil, nil,
		),
		responsesDropped: prometheus.NewDesc(
			"mtproxy_responses_dropped_total",
			"Total number of dropped responses",
			nil, nil,
		),

		// Buffers
		bufferUsedBytes: prometheus.NewDesc(
			"mtproxy_buffer_used_bytes",
			"Currently used buffer bytes",
			nil, nil,
		),
		bufferAllocatedBytes: prometheus.NewDesc(
			"mtproxy_buffer_allocated_bytes",
			"Total allocated buffer bytes",
			nil, nil,
		),
		buffersUsed: prometheus.NewDesc(
			"mtproxy_buffers_used",
			"Number of buffers currently in use",
			nil, nil,
		),

		// Load
		loadAverageUser: prometheus.NewDesc(
			"mtproxy_load_average_user",
			"Average user CPU load",
			nil, nil,
		),
		loadAverageSys: prometheus.NewDesc(
			"mtproxy_load_average_sys",
			"Average system CPU load",
			nil, nil,
		),
		loadAverageTotal: prometheus.NewDesc(
			"mtproxy_load_average_total",
			"Average total CPU load",
			nil, nil,
		),
		loadRecentUser: prometheus.NewDesc(
			"mtproxy_load_recent_user",
			"Recent user CPU load",
			nil, nil,
		),
		loadRecentSys: prometheus.NewDesc(
			"mtproxy_load_recent_sys",
			"Recent system CPU load",
			nil, nil,
		),
		loadRecentTotal: prometheus.NewDesc(
			"mtproxy_load_recent_total",
			"Recent total CPU load",
			nil, nil,
		),
		threadLoadUser: prometheus.NewDesc(
			"mtproxy_thread_load_user",
			"User CPU load by thread",
			[]string{"thread"}, nil,
		),
		threadLoadSys: prometheus.NewDesc(
			"mtproxy_thread_load_sys",
			"System CPU load by thread",
			[]string{"thread"}, nil,
		),

		// Memory
		memoryVmsize: prometheus.NewDesc(
			"mtproxy_memory_vmsize_bytes",
			"Virtual memory size in bytes",
			nil, nil,
		),
		memoryVmrss: prometheus.NewDesc(
			"mtproxy_memory_vmrss_bytes",
			"Resident set size in bytes",
			nil, nil,
		),
		memoryVmdata: prometheus.NewDesc(
			"mtproxy_memory_vmdata_bytes",
			"Data segment size in bytes",
			nil, nil,
		),

		// Jobs
		jobsActive: prometheus.NewDesc(
			"mtproxy_jobs_active",
			"Number of active jobs",
			nil, nil,
		),
		jobsCreated: prometheus.NewDesc(
			"mtproxy_jobs_created_total",
			"Total number of jobs created",
			nil, nil,
		),

		// HTTP
		httpConnections: prometheus.NewDesc(
			"mtproxy_http_connections",
			"Number of HTTP connections",
			nil, nil,
		),
		httpQueries: prometheus.NewDesc(
			"mtproxy_http_queries_total",
			"Total number of HTTP queries",
			nil, nil,
		),
		httpQPS: prometheus.NewDesc(
			"mtproxy_http_qps",
			"HTTP queries per second",
			nil, nil,
		),

		// Errors
		proxyErrors: prometheus.NewDesc(
			"mtproxy_errors_total",
			"Total number of MTProxy errors",
			nil, nil,
		),
		connectionsFailedLru: prometheus.NewDesc(
			"mtproxy_connections_failed_lru_total",
			"Total connections failed due to LRU",
			nil, nil,
		),
		connectionsFailedFlood: prometheus.NewDesc(
			"mtproxy_connections_failed_flood_total",
			"Total connections failed due to flood",
			nil, nil,
		),

		// Targets
		readyTargets: prometheus.NewDesc(
			"mtproxy_ready_targets",
			"Number of ready targets",
			nil, nil,
		),
		allocatedTargets: prometheus.NewDesc(
			"mtproxy_allocated_targets",
			"Number of allocated targets",
			nil, nil,
		),
		activeTargets: prometheus.NewDesc(
			"mtproxy_active_targets",
			"Number of active targets",
			nil, nil,
		),
		inactiveTargets: prometheus.NewDesc(
			"mtproxy_inactive_targets",
			"Number of inactive targets",
			nil, nil,
		),
	}
}

// Describe implements prometheus.Collector
func (c *MTProxyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.uptime
	ch <- c.startTime
	ch <- c.scrapeDuration
	ch <- c.scrapeErrors
	ch <- c.lastScrapeTime
	ch <- c.activeConnections
	ch <- c.inboundConnections
	ch <- c.outboundConnections
	ch <- c.connectionsCreated
	ch <- c.connectFailures
	ch <- c.activeInboundConnections
	ch <- c.activeOutboundConnections
	ch <- c.readyOutboundConnections
	ch <- c.inboundConnectionsAccepted
	ch <- c.allocatedConnections
	ch <- c.tcpReadBytes
	ch <- c.tcpWriteBytes
	ch <- c.tcpReadvCalls
	ch <- c.tcpWritevCalls
	ch <- c.idlePercent
	ch <- c.recentIdlePercent
	ch <- c.epollCalls
	ch <- c.activeNetworkEvents
	ch <- c.messagesTotal
	ch <- c.queriesForwarded
	ch <- c.responsesForwarded
	ch <- c.queriesDropped
	ch <- c.responsesDropped
	ch <- c.bufferUsedBytes
	ch <- c.bufferAllocatedBytes
	ch <- c.buffersUsed
	ch <- c.loadAverageUser
	ch <- c.loadAverageSys
	ch <- c.loadAverageTotal
	ch <- c.loadRecentUser
	ch <- c.loadRecentSys
	ch <- c.loadRecentTotal
	ch <- c.threadLoadUser
	ch <- c.threadLoadSys
	ch <- c.memoryVmsize
	ch <- c.memoryVmrss
	ch <- c.memoryVmdata
	ch <- c.jobsActive
	ch <- c.jobsCreated
	ch <- c.httpConnections
	ch <- c.httpQueries
	ch <- c.httpQPS
	ch <- c.proxyErrors
	ch <- c.connectionsFailedLru
	ch <- c.connectionsFailedFlood
	ch <- c.readyTargets
	ch <- c.allocatedTargets
	ch <- c.activeTargets
	ch <- c.inactiveTargets
}

// Collect implements prometheus.Collector
func (c *MTProxyCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	var scrapeErrors float64

	// Create context with timeout for the scrape
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := c.client.GetStats(ctx)
	scrapeDuration := time.Since(start).Seconds()

	if err != nil {
		slog.Error("Failed to fetch stats", "error", err)
		scrapeErrors = 1
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		ch <- prometheus.MustNewConstMetric(c.scrapeDuration, prometheus.GaugeValue, scrapeDuration)
		ch <- prometheus.MustNewConstMetric(c.scrapeErrors, prometheus.CounterValue, scrapeErrors)
		return
	}

	stats, err := parser.Parse(data)
	if err != nil {
		slog.Error("Failed to parse stats", "error", err)
		scrapeErrors = 1
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		ch <- prometheus.MustNewConstMetric(c.scrapeDuration, prometheus.GaugeValue, scrapeDuration)
		ch <- prometheus.MustNewConstMetric(c.scrapeErrors, prometheus.CounterValue, scrapeErrors)
		return
	}

	// Report that we're up
	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(c.scrapeDuration, prometheus.GaugeValue, scrapeDuration)
	ch <- prometheus.MustNewConstMetric(c.scrapeErrors, prometheus.CounterValue, 0)
	ch <- prometheus.MustNewConstMetric(c.lastScrapeTime, prometheus.GaugeValue, float64(time.Now().Unix()))

	// General metrics
	ch <- prometheus.MustNewConstMetric(c.uptime, prometheus.GaugeValue, float64(stats.Uptime))
	ch <- prometheus.MustNewConstMetric(c.startTime, prometheus.GaugeValue, float64(stats.StartTime))

	// Connection metrics
	ch <- prometheus.MustNewConstMetric(c.activeConnections, prometheus.GaugeValue, float64(stats.ActiveConnections))
	ch <- prometheus.MustNewConstMetric(c.inboundConnections, prometheus.GaugeValue, float64(stats.InboundConnections))
	ch <- prometheus.MustNewConstMetric(c.outboundConnections, prometheus.GaugeValue, float64(stats.OutboundConnections))
	ch <- prometheus.MustNewConstMetric(c.connectionsCreated, prometheus.CounterValue, float64(stats.OutboundConnectionsCreated))
	ch <- prometheus.MustNewConstMetric(c.connectFailures, prometheus.CounterValue, float64(stats.TotalConnectFailures))
	ch <- prometheus.MustNewConstMetric(c.activeInboundConnections, prometheus.GaugeValue, float64(stats.ActiveInboundConnections))
	ch <- prometheus.MustNewConstMetric(c.activeOutboundConnections, prometheus.GaugeValue, float64(stats.ActiveOutboundConnections))
	ch <- prometheus.MustNewConstMetric(c.readyOutboundConnections, prometheus.GaugeValue, float64(stats.ReadyOutboundConnections))
	ch <- prometheus.MustNewConstMetric(c.inboundConnectionsAccepted, prometheus.CounterValue, float64(stats.InboundConnectionsAccepted))
	ch <- prometheus.MustNewConstMetric(c.allocatedConnections, prometheus.GaugeValue, float64(stats.AllocatedConnections))

	// Network metrics
	ch <- prometheus.MustNewConstMetric(c.tcpReadBytes, prometheus.CounterValue, float64(stats.TCPReadvBytes))
	ch <- prometheus.MustNewConstMetric(c.tcpWriteBytes, prometheus.CounterValue, float64(stats.TCPWritevBytes))
	ch <- prometheus.MustNewConstMetric(c.tcpReadvCalls, prometheus.CounterValue, float64(stats.TCPReadvCalls))
	ch <- prometheus.MustNewConstMetric(c.tcpWritevCalls, prometheus.CounterValue, float64(stats.TCPWritevCalls))

	// Performance metrics
	ch <- prometheus.MustNewConstMetric(c.idlePercent, prometheus.GaugeValue, stats.AverageIdlePercent)
	ch <- prometheus.MustNewConstMetric(c.recentIdlePercent, prometheus.GaugeValue, stats.RecentIdlePercent)
	ch <- prometheus.MustNewConstMetric(c.epollCalls, prometheus.CounterValue, float64(stats.EpollCalls))
	ch <- prometheus.MustNewConstMetric(c.activeNetworkEvents, prometheus.GaugeValue, float64(stats.ActiveNetworkEvents))

	// Message metrics
	ch <- prometheus.MustNewConstMetric(c.messagesTotal, prometheus.CounterValue, float64(stats.RwmTotalMsgs))
	ch <- prometheus.MustNewConstMetric(c.queriesForwarded, prometheus.CounterValue, float64(stats.TotForwardedQueries))
	ch <- prometheus.MustNewConstMetric(c.responsesForwarded, prometheus.CounterValue, float64(stats.TotForwardedResponses))
	ch <- prometheus.MustNewConstMetric(c.queriesDropped, prometheus.CounterValue, float64(stats.DroppedQueries))
	ch <- prometheus.MustNewConstMetric(c.responsesDropped, prometheus.CounterValue, float64(stats.DroppedResponses))

	// Buffer metrics
	ch <- prometheus.MustNewConstMetric(c.bufferUsedBytes, prometheus.GaugeValue, float64(stats.TotalUsedBuffersSize))
	ch <- prometheus.MustNewConstMetric(c.bufferAllocatedBytes, prometheus.GaugeValue, float64(stats.AllocatedBufferBytes))
	ch <- prometheus.MustNewConstMetric(c.buffersUsed, prometheus.GaugeValue, float64(stats.TotalUsedBuffers))

	// Load metrics
	ch <- prometheus.MustNewConstMetric(c.loadAverageUser, prometheus.GaugeValue, stats.LoadAverageUser)
	ch <- prometheus.MustNewConstMetric(c.loadAverageSys, prometheus.GaugeValue, stats.LoadAverageSys)
	ch <- prometheus.MustNewConstMetric(c.loadAverageTotal, prometheus.GaugeValue, stats.LoadAverageTotal)
	ch <- prometheus.MustNewConstMetric(c.loadRecentUser, prometheus.GaugeValue, stats.LoadRecentUser)
	ch <- prometheus.MustNewConstMetric(c.loadRecentSys, prometheus.GaugeValue, stats.LoadRecentSys)
	ch <- prometheus.MustNewConstMetric(c.loadRecentTotal, prometheus.GaugeValue, stats.LoadRecentTotal)

	// Memory metrics
	ch <- prometheus.MustNewConstMetric(c.memoryVmsize, prometheus.GaugeValue, float64(stats.VmsizeBytes))
	ch <- prometheus.MustNewConstMetric(c.memoryVmrss, prometheus.GaugeValue, float64(stats.VmrssBytes))
	ch <- prometheus.MustNewConstMetric(c.memoryVmdata, prometheus.GaugeValue, float64(stats.VmdataBytes))

	// Jobs metrics
	ch <- prometheus.MustNewConstMetric(c.jobsActive, prometheus.GaugeValue, float64(stats.JobsActive))
	ch <- prometheus.MustNewConstMetric(c.jobsCreated, prometheus.CounterValue, float64(stats.JobsCreated))

	// HTTP metrics
	ch <- prometheus.MustNewConstMetric(c.httpConnections, prometheus.GaugeValue, float64(stats.HttpConnections))
	ch <- prometheus.MustNewConstMetric(c.httpQueries, prometheus.CounterValue, float64(stats.HttpQueries))
	ch <- prometheus.MustNewConstMetric(c.httpQPS, prometheus.GaugeValue, stats.HttpQPS)

	// Error metrics
	ch <- prometheus.MustNewConstMetric(c.proxyErrors, prometheus.CounterValue, float64(stats.MtprotoProxyErrors))
	ch <- prometheus.MustNewConstMetric(c.connectionsFailedLru, prometheus.CounterValue, float64(stats.ConnectionsFailedLru))
	ch <- prometheus.MustNewConstMetric(c.connectionsFailedFlood, prometheus.CounterValue, float64(stats.ConnectionsFailedFlood))

	// Target metrics
	ch <- prometheus.MustNewConstMetric(c.readyTargets, prometheus.GaugeValue, float64(stats.ReadyTargets))
	ch <- prometheus.MustNewConstMetric(c.allocatedTargets, prometheus.GaugeValue, float64(stats.AllocatedTargets))
	ch <- prometheus.MustNewConstMetric(c.activeTargets, prometheus.GaugeValue, float64(stats.ActiveTargets))
	ch <- prometheus.MustNewConstMetric(c.inactiveTargets, prometheus.GaugeValue, float64(stats.InactiveTargets))

	// Thread load metrics with labels
	for i, load := range stats.ThreadLoadRecentUser {
		threadID := fmt.Sprintf("%d", i)
		ch <- prometheus.MustNewConstMetric(c.threadLoadUser, prometheus.GaugeValue, load, threadID)
	}
	for i, load := range stats.ThreadLoadRecentSys {
		threadID := fmt.Sprintf("%d", i)
		ch <- prometheus.MustNewConstMetric(c.threadLoadSys, prometheus.GaugeValue, load, threadID)
	}
}
