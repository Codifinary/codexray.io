---
sidebar_position: 6
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Docker Swarm

<Tabs queryString="edition">
  <TabItem value="ce" label="Community Edition" default>

**Step #1: Initialize Docker Swarm**

If you haven't already initialized Docker Swarm on your manager node, run the following command on the manager node:

```bash
docker swarm init
```

This initializes a new Docker Swarm and joins the current node as a manager.

**Step #2: Deploy the Codexray Stack**

Deploy the codexray stack to your cluster by running the following command on the manager node. 
Before applying, you can review the configuration file in codexray's GitHub repository: docker-swarm-stack.yaml

```bash
curl -fsS https://raw.githubusercontent.com/codexray/codexray/main/deploy/docker-swarm-stack.yaml | \
  docker stack deploy -c - codexray
```

**Step #3: Validate the deployment**

After deploying the stack, you can use docker stack ls to list the deployed stacks in your Docker Swarm cluster. 
Here's an example of how the output might look:

```bash
NAME      SERVICES
codexray    3
```

**Step #4: Installing codexray-node-agent**

Since Docker Swarm doesn't support privileged containers, you'll have to manually deploy codexray-node-agent on each cluster node. 
Just replace `NODE_IP` with any node's IP address in the Docker Swarm cluster.

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

**Step #5: Accessing Codexray**

Access codexray through any node in your Docker Swarm cluster using its published port: http://NODE_IP:8080.

**Uninstall Codexray**

To uninstall codexray run the following command:

```bash
docker stack rm codexray
```
  </TabItem>

  <TabItem value="ee" label="Enterprise Edition">

:::info
Codexray Enterprise Edition is a paid subscription that offers extra features and priority support.
To install the Enterprise Edition, you'll need a license. [Start](https://codexray.com/account) your free trial today.
:::

**Step #1: Initialize Docker Swarm**

If you haven't already initialized Docker Swarm on your manager node, run the following command on the manager node:

```
docker swarm init
```

This initializes a new Docker Swarm and joins the current node as a manager.

**Step #2: Deploy the Codexray Stack**

Deploy the codexray stack to your cluster by running the following command on the manager node. Before applying, you can review the configuration file in Codexray's GitHub repository: docker-swarm-stack.yaml

```
curl -fsS https://raw.githubusercontent.com/codexray/codexray-ee/main/deploy/docker-swarm-stack.yaml | \
  LICENSE_KEY="codexray-LICENSE-KEY-HERE" docker stack deploy -c - codexray-ee
```

**Step #3: Validate the deployment**

After deploying the stack, you can use docker stack ls to list the deployed stacks in your Docker Swarm cluster. 
Here's an example of how the output might look:

```
NAME        SERVICES
codexray-ee   4
```

**Step #4: Installing codexray-node-agent**

Since Docker Swarm doesn't support privileged containers, you'll have to manually deploy codexray-node-agent on each cluster node. 
Just replace `NODE_IP` with any node's IP address in the Docker Swarm cluster.

```
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

**Step #5: Accessing Codexray**

Access codexray through any node in your Docker Swarm cluster using its published port: http://NODE_IP:8080.

**Uninstall Codexray**

To uninstall codexray run the following command:

```
docker stack rm codexray-ee
```
</TabItem>
</Tabs>
