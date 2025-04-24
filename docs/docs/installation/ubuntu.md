---
sidebar_position: 6
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Ubuntu & Debian

<Tabs queryString="edition">
  <TabItem value="ce" label="Community Edition" default>

**Step #1: Installing ClickHouse**

```bash
sudo apt install -y apt-transport-https ca-certificates curl gnupg
curl -fsSL 'https://packages.clickhouse.com/rpm/lts/repodata/repomd.xml.key' | sudo gpg --dearmor -o /usr/share/keyrings/clickhouse-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/clickhouse-keyring.gpg] https://packages.clickhouse.com/deb stable main" | sudo tee /etc/apt/sources.list.d/clickhouse.list
sudo apt update
sudo DEBIAN_FRONTEND=noninteractive apt install -y clickhouse-server clickhouse-client
sudo service clickhouse-server start
```

**Step #2: Installing Prometheus**

codexray requires Prometheus with support for Remote Write Receiver, which has been available since v2.25.0.

```bash
sudo apt install -y prometheus
sudo service prometheus start
```

Enable Remote Write Receiver in Prometheus by adding the `--enable-feature=remote-write-receiver` argument to the `/etc/default/prometheus` file:

```bash
# Set the command-line arguments to pass to the server.
# Due to shell escaping, to pass backslashes for regexes, you need to double
# them (\\d for \d). If running under systemd, you need to double them again
# (\\\\d to mean \d), and escape newlines too.
ARGS="--enable-feature=remote-write-receiver"
```

Restart Prometheus:

```bash
sudo service prometheus restart
```

**Step #3: Installing codexray**

```bash
curl -sfL https://raw.githubusercontent.com/codexray/codexray/main/deploy/install.sh | \
  BOOTSTRAP_PROMETHEUS_URL="http://127.0.0.1:9090" \
  BOOTSTRAP_REFRESH_INTERVAL=15s \
  BOOTSTRAP_CLICKHOUSE_ADDRESS=127.0.0.1:9000 \
  sh -
```

**Step #4: Installing codexray-node-agent**

```bash
curl -sfL https://raw.githubusercontent.com/codexray/codexray-node-agent/main/install.sh | \
  COLLECTOR_ENDPOINT=http://127.0.0.1:8080 \
  SCRAPE_INTERVAL=15s \
  sh -
```

**Step #5: Accessing codexray**

Access codexray at: http://NODE_IP:8080.

**Uninstall codexray**

To uninstall codexray run the following command:

```bash
/usr/bin/codexray-uninstall.sh
```

Uninstall codexray-node-agent:

```bash
/usr/bin/codexray-node-agent-uninstall.sh
```
  </TabItem>

  <TabItem value="ee" label="Enterprise Edition">

:::info
codexray Enterprise Edition is a paid subscription (from $1 per CPU core/month) that offers extra features and priority support.
To install the Enterprise Edition, you'll need a license. [Start](https://codexray.com/account) your free trial today.
:::

**Step #1: Installing ClickHouse**

```bash
sudo apt install -y apt-transport-https ca-certificates curl gnupg
curl -fsSL 'https://packages.clickhouse.com/rpm/lts/repodata/repomd.xml.key' | sudo gpg --dearmor -o /usr/share/keyrings/clickhouse-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/clickhouse-keyring.gpg] https://packages.clickhouse.com/deb stable main" | sudo tee /etc/apt/sources.list.d/clickhouse.list
sudo apt update
sudo DEBIAN_FRONTEND=noninteractive apt install -y clickhouse-server clickhouse-client
sudo service clickhouse-server start
```

**Step #2: Installing Prometheus**

codexray requires Prometheus with support for Remote Write Receiver, which has been available since v2.25.0.

```bash
sudo apt install -y prometheus
sudo service prometheus start
```

Enable Remote Write Receiver in Prometheus by adding the `--enable-feature=remote-write-receiver` argument to the `/etc/default/prometheus` file:

```bash
# Set the command-line arguments to pass to the server.
# Due to shell escaping, to pass backslashes for regexes, you need to double
# them (\\d for \d). If running under systemd, you need to double them again
# (\\\\d to mean \d), and escape newlines too.
ARGS="--enable-feature=remote-write-receiver"
```

Restart Prometheus:

```bash
sudo service prometheus restart
```

**Step #3: Installing Codexray**

```bash
curl -sfL https://raw.githubusercontent.com/codexray/codexray-ee/main/deploy/install.sh | \
  LICENSE_KEY="codexray-LICENSE-KEY-HERE" \
  BOOTSTRAP_PROMETHEUS_URL="http://127.0.0.1:9090" \
  BOOTSTRAP_REFRESH_INTERVAL=15s \
  BOOTSTRAP_CLICKHOUSE_ADDRESS=127.0.0.1:9000 \
  sh -
```

**Step #4: Installing codexray-node-agent**

```bash
curl -sfL https://raw.githubusercontent.com/codexray/codexray-node-agent/main/install.sh | \
  COLLECTOR_ENDPOINT=http://127.0.0.1:8080 \
  SCRAPE_INTERVAL=15s \
  sh -
```

**Step #5: Accessing Codexray**

Access codexray at: http://NODE_IP:8080.

**Uninstall codexray**

To uninstall codexray run the following command:

```bash
/usr/bin/codexray-ee-uninstall.sh
```

Uninstall codexray-node-agent:

```bash
/usr/bin/codexray-node-agent-uninstall.sh
```
</TabItem>
</Tabs>
