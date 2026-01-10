# Metrics Reference

This document describes all metrics exported by MTProxy Exporter.

## Exporter Metrics

These metrics provide information about the exporter itself:

| Metric | Type | Description |
|--------|------|-------------|
| `mtproxy_exporter_scrape_duration_seconds` | Gauge | Duration of the scrape operation in seconds |
| `mtproxy_exporter_scrape_errors_total` | Counter | Total number of scrape errors encountered |
| `mtproxy_exporter_last_scrape_timestamp` | Gauge | Unix timestamp of the last successful scrape |

## General Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `mtproxy_up` | Gauge | Whether the MTProxy is up (1) or down (0) |
| `mtproxy_uptime_seconds` | Gauge | Uptime of the MTProxy in seconds |
| `mtproxy_start_time` | Gauge | Unix timestamp when the MTProxy started |

## Connection Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `mtproxy_active_connections` | Gauge | Number of active connections |
| `mtproxy_inbound_connections` | Gauge | Number of inbound connections |
| `mtproxy_outbound_connections` | Gauge | Number of outbound connections |
| `mtproxy_active_inbound_connections` | Gauge | Number of active inbound connections |
| `mtproxy_active_outbound_connections` | Gauge | Number of active outbound connections |
| `mtproxy_ready_outbound_connections` | Gauge | Number of ready outbound connections |
| `mtproxy_allocated_connections` | Gauge | Number of allocated connections |
| `mtproxy_connections_created_total` | Counter | Total number of outbound connections created |
| `mtproxy_connect_failures_total` | Counter | Total number of connection failures |
| `mtproxy_inbound_connections_accepted_total` | Counter | Total number of accepted inbound connections |

## Network Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `mtproxy_tcp_read_bytes_total` | Counter | Total bytes read via TCP |
| `mtproxy_tcp_write_bytes_total` | Counter | Total bytes written via TCP |
| `mtproxy_tcp_readv_calls_total` | Counter | Total readv system calls |
| `mtproxy_tcp_writev_calls_total` | Counter | Total writev system calls |

## Performance Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `mtproxy_idle_percent` | Gauge | Average idle percentage |
| `mtproxy_recent_idle_percent` | Gauge | Recent idle percentage |
| `mtproxy_epoll_calls_total` | Counter | Total epoll system calls |
| `mtproxy_active_network_events` | Gauge | Number of active network events |

## Load Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `mtproxy_load_average_user` | Gauge | Average user CPU load |
| `mtproxy_load_average_sys` | Gauge | Average system CPU load |
| `mtproxy_load_average_total` | Gauge | Average total CPU load |
| `mtproxy_load_recent_user` | Gauge | Recent user CPU load |
| `mtproxy_load_recent_sys` | Gauge | Recent system CPU load |
| `mtproxy_load_recent_total` | Gauge | Recent total CPU load |
| `mtproxy_thread_load_user` | Gauge | User CPU load by thread (with `thread` label) |
| `mtproxy_thread_load_sys` | Gauge | System CPU load by thread (with `thread` label) |

## Message Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `mtproxy_messages_total` | Counter | Total number of messages processed |
| `mtproxy_queries_forwarded_total` | Counter | Total number of forwarded queries |
| `mtproxy_responses_forwarded_total` | Counter | Total number of forwarded responses |
| `mtproxy_queries_dropped_total` | Counter | Total number of dropped queries |
| `mtproxy_responses_dropped_total` | Counter | Total number of dropped responses |

## Buffer Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `mtproxy_buffer_used_bytes` | Gauge | Currently used buffer bytes |
| `mtproxy_buffer_allocated_bytes` | Gauge | Total allocated buffer bytes |
| `mtproxy_buffers_used` | Gauge | Number of buffers currently in use |

## Memory Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `mtproxy_memory_vmsize_bytes` | Gauge | Virtual memory size in bytes |
| `mtproxy_memory_vmrss_bytes` | Gauge | Resident set size in bytes |
| `mtproxy_memory_vmdata_bytes` | Gauge | Data segment size in bytes |

## Job Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `mtproxy_jobs_active` | Gauge | Number of active jobs |
| `mtproxy_jobs_created_total` | Counter | Total number of jobs created |

## HTTP Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `mtproxy_http_connections` | Gauge | Number of HTTP connections |
| `mtproxy_http_queries_total` | Counter | Total number of HTTP queries |
| `mtproxy_http_qps` | Gauge | HTTP queries per second |

## Error Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `mtproxy_errors_total` | Counter | Total number of MTProxy errors |
| `mtproxy_connections_failed_lru_total` | Counter | Total connections failed due to LRU |
| `mtproxy_connections_failed_flood_total` | Counter | Total connections failed due to flood |

## Target Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `mtproxy_ready_targets` | Gauge | Number of ready targets |
| `mtproxy_allocated_targets` | Gauge | Number of allocated targets |
| `mtproxy_active_targets` | Gauge | Number of active targets |
| `mtproxy_inactive_targets` | Gauge | Number of inactive targets |

## Example PromQL Queries

### Exporter Health
**Check if exporter is scraping successfully:**
```promql
mtproxy_up
```

**Scrape duration:**
```promql
mtproxy_exporter_scrape_duration_seconds
```

**Scrape error rate:**
```promql
rate(mtproxy_exporter_scrape_errors_total[5m])
```

### Connection Rate
```promql
rate(mtproxy_connections_created_total[5m])
```

### Error Rate
```promql
rate(mtproxy_connect_failures_total[5m])
```

### Network Throughput (Bytes/sec)
```promql
rate(mtproxy_tcp_read_bytes_total[5m]) + rate(mtproxy_tcp_write_bytes_total[5m])
```

### CPU Usage
```promql
100 - mtproxy_idle_percent
```

### Memory Usage (MB)
```promql
mtproxy_memory_vmrss_bytes / 1024 / 1024
```

### Messages Per Second
```promql
rate(mtproxy_messages_total[5m])
```

### Connection Success Rate
```promql
rate(mtproxy_connections_created_total[5m]) / (rate(mtproxy_connections_created_total[5m]) + rate(mtproxy_connect_failures_total[5m]))
```

### Thread Load Per Core
**User CPU load for thread 0:**
```promql
mtproxy_thread_load_user{thread="0"}
```

**Total CPU load across all threads:**
```promql
sum(mtproxy_thread_load_user) + sum(mtproxy_thread_load_sys)
```

## Labels

Some metrics include labels for additional dimensionality:

- **`thread`**: Thread identifier (0, 1, 2, etc.) used in `mtproxy_thread_load_user` and `mtproxy_thread_load_sys` metrics
