---
sidebar_position: 2
---

# Architecture
 <img alt="Architecture" src="/docs/docs/Screenshot 2025-04-23 131459.png" class="card w-1200"/>

## Codexray-node-agent

Codexray-node-agent is an open-source observability agent powered by eBPF. It collects metrics, logs, traces, and profiles from all containers running on a node.

The agent supports both pull and push modes for metrics: it exposes metrics in Prometheus format and can also send metrics directly to codexray using the Prometheus Remote Write protocol. 
Logs and traces are sent to codexray via the OpenTelemetry protocol, while profiles are transmitted using a custom HTTP-based protocol.

To ensure full coverage, the agent needs to be installed on every node in your cluster. 
If you’re using Kubernetes, it’s deployed as a DaemonSet, so it will automatically be added to each node.

## Codexray-cluster-agent

Codexray-cluster-agent is a dedicated tool for collecting cluster-wide telemetry data:

* It gathers database metrics by discovering databases through codexray's Service Map. Using the credentials provided by codexray, the agent connects to the identified databases such as Postgres, MySQL, Redis, Memcached, and MongoDB, collects database-specific metrics, and sends them to codexray using the Prometheus Remote Write protocol.
* In addition to the eBPF-based continuous profiler embedded in codexray-node-agent, codexray also supports application-level profiling. The Cluster Agent can discover Go applications annotated with codexray.com/profile-scrape and codexray.com/profile-port, and gather CPU and memory profiles from their instances.
* The agent can be integrated with AWS to discover RDS and ElastiCache clusters and collect their telemetry data.

## OpenTelemetry

Codexray supports the OpenTelemetry protocol (OTLP over HTTP) for logs and traces. If your applications are instrumented with OpenTelemetry SDKs, 
you can configure them to send data directly to codexray or route it through the OpenTelemetry Collector.

## Prometheus

Codexray uses Prometheus for storing metrics and is compatible with any Prometheus-compatible time series databases such as VictoriaMetrics, Thanos, or Grafana Mimir.

For faster access, codexray maintains its own on-disk metric cache, continuously retrieving metrics from Prometheus. 
As a result, codexray treats the time series database as a source for updating its cache. 
This allows you to configure Prometheus with a shorter retention period, such as a few hours.

## ClickHouse

codexray uses ClickHouse for storing logs, traces, and profiles. 
Thanks to the efficient data compression implemented by ClickHouse, you can expect a compression ratio of 10x or more for this telemetry data.
