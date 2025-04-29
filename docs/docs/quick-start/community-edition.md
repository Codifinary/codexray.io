---
sidebar_position: 1
slug: /
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Community Edition

This guide provides a quick overview of launching CodeXray Community Edition with default options. For more details and customization options check out the Installation section.

<Tabs queryString="env">
  <TabItem value="kubernetes" label="Kubernetes" default>

Add the CodeXray helm chart repo:

```bash
helm repo add codexray https://codexray.github.io/helm-charts
helm repo update codexray
```

Next, install the codexray Operator:

```bash
helm install -n codexray --create-namespace codexray-operator codexray/codexray-operator
```

Install the CodeXray Community Edition. This chart creates a minimal [codexray Custom Resource](/installation/k8s-operator):

```bash
helm install -n codexray codexray codexray/codexray-ce \
  --set "clickhouse.shards=2,clickhouse.replicas=2"
```

Forward the CodeXray port to your machine:

```bash
kubectl port-forward -n codexray service/codexray-codexray 8080:8080
```

Then, you can access CodeXray at http://localhost:8080

  </TabItem>

  <TabItem value="docker" label="Docker">

To deploy CodeXray using Docker Compose, run the following command. Before applying it, you can review the configuration file in CodeXray's GitHub repository: docker-compose.yaml

```bash
curl -fsS https://raw.githubusercontent.com/codexray/codexray/main/deploy/docker-compose.yaml | \
  docker compose -f - up -d
```

If you installed CodeXray on your desktop machine, you can access it at http://localhost:8080/. If codexray is deployed on a remote node, replace `NODE_IP_ADDRESS` with the IP address of the node in the following URL: http://NODE_IP_ADDRESS:8080/.

  </TabItem>

  <TabItem value="docker-swarm" label="Docker Swarm">

Deploy the CodeXray stack to your cluster by running the following command on the manager node. Before applying, you can review the configuration file in CodeXray's GitHub repository: docker-swarm-stack.yaml

```bash
curl -fsS https://raw.githubusercontent.com/codexray/codexray/main/deploy/docker-swarm-stack.yaml | \
  docker stack deploy -c - codexray
```

Since Docker Swarm doesn't support privileged containers, you'll have to manually deploy codexray-node-agent on each cluster node. Just replace `NODE_IP` with any node's IP address in the Docker Swarm cluster.

```bash
docker run --detach --name codexray-node-agent \
  --pull=always \
  --privileged --pid host \
  -v /sys/kernel/tracing:/sys/kernel/tracing:rw \
  -v /sys/kernel/debug:/sys/kernel/debug:rw \
  -v /sys/fs/cgroup:/host/sys/fs/cgroup:ro \
  ghcr.io/codexray/codexray-node-agent \
  --cgroupfs-root=/host/sys/fs/cgroup \
  --collector-endpoint=http://NODE_IP:8080
```
Access CodeXray through any node in your Docker Swarm cluster using its published port: http://NODE_IP:8080.
  </TabItem>

  <TabItem value="ubuntu" label="Ubuntu & Debian">

CodeXray requires a Prometheus server with the Remote Write Receiver enabled, along with a ClickHouse server. 
For detailed steps on installing all the necessary components on an Ubuntu/Debian node, refer to the [full instructions](/installation/ubuntu).

To install codexray, run the following command:

```bash
curl -sfL https://raw.githubusercontent.com/codexray/codexray/main/deploy/install.sh | \
  BOOTSTRAP_PROMETHEUS_URL="http://PROMETHEUS_IP:9090" \
  BOOTSTRAP_REFRESH_INTERVAL=15s \
  BOOTSTRAP_CLICKHOUSE_ADDRESS=CLICKHOUSE_IP:9000 \
  sh -
```

Install codexray-node-agent to every node within your infrastructure:

```bash
curl -sfL https://raw.githubusercontent.com/codexray/codexray-node-agent/main/install.sh | \
  COLLECTOR_ENDPOINT=http://codexray_NODE_IP:8080 \
  SCRAPE_INTERVAL=15s \
  sh -
```

Access CodeXray at: http://codexray_NODE_IP:8080.
</TabItem>

<TabItem value="rhel" label="RHEL & CentOS">

CodeXray requires a Prometheus server with the Remote Write Receiver enabled, along with a ClickHouse server. 
For detailed steps on installing all the necessary components on an Ubuntu/Debian node, refer to the [full instructions](/installation/rhel).

To install CodeXray, run the following command:

```bash
curl -sfL https://raw.githubusercontent.com/codexray/codexray/main/deploy/install.sh | \
  BOOTSTRAP_PROMETHEUS_URL="http://PROMETHEUS_IP:9090" \
  BOOTSTRAP_REFRESH_INTERVAL=15s \
  BOOTSTRAP_CLICKHOUSE_ADDRESS=CLICKHOUSE_IP:9000 \
  sh -
```

Install codexray-node-agent to every node within your infrastructure:

```bash
curl -sfL https://raw.githubusercontent.com/codexray/codexray-node-agent/main/install.sh | \
  COLLECTOR_ENDPOINT=http://codexray_NODE_IP:8080 \
  SCRAPE_INTERVAL=15s \
  sh -
```
Access CodeXray at: http://codexray_NODE_IP:8080.
</TabItem>
  


</Tabs>
