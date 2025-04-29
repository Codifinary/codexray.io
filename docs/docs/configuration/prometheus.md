---
sidebar_position: 4
---

# Prometheus

Codexray uses Prometheus to store metrics. To integrate Codexray with Prometheus, go to the **Project Settings**, 
click on **Prometheus**, and define the Prometheus address and credentials as shown in the following example:

<img alt="Prometheus Configuration" src="/docs/docs/Doc_Prometheus_Integration.png" class="card w-1200"/>

## Multi-tenancy mode

Codexray supports a multi-tenancy mode, allowing a single Prometheus server to store metrics for multiple projects (or clusters).

In this mode, all Codexray agents (both `codexray-node-agent` and `codexray-cluster-agent`) are configured to push metrics 
to Codexray using the Prometheus Remote Write Protocol. 
Codexray automatically adds the `codexray_project_id` label to each metric and uses `{codexray_project_id="XXXX"}` as an additional 
selector when querying metrics for a specific project.


## Metric cache
For faster access, codexray maintains its own on-disk metric cache, continuously retrieving metrics from Prometheus. 
As a result, Codexray treats the time series database as a source for updating its cache. 
This allows you to configure Prometheus with a shorter retention period, such as a few hours.

The retention of codexray's metric cache can be configured using the `--cache-ttl` CLI argument of the `CACHE_TTL` environment variable. 
Check the [CLI arguments](/configuration/cli-arguments) section for more details.

