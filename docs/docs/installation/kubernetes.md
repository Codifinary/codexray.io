---
sidebar_position: 3
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Kubernetes

<Tabs queryString="edition">
  <TabItem value="ce" label="Community Edition (operator)" default>

Add the codexray helm chart repo:

```bash
helm repo add codexray https://codexray.github.io/helm-charts
helm repo update codexray
```

Next, install the codexray Operator:

```bash
helm install -n codexray --create-namespace codexray-operator codexray/codexray-operator
```

Install the codexray Community Edition. This chart creates a minimal [codexray Custom Resource](/installation/k8s-operator):

```bash
helm install -n codexray codexray codexray/codexray-ce \
  --set "clickhouse.shards=2,clickhouse.replicas=2"
```

The helm chart 

Forward the codexray port to your machine:

```bash
kubectl port-forward -n codexray service/codexray-codexray 8080:8080
```

Then, you can access codexray at http://localhost:8080

**Upgrade**

The codexray Operator for Kubernetes automatically upgrades all components.

**Uninstall**

To uninstall codexray run the following command:

```bash
helm uninstall codexray -n codexray
helm uninstall codexray-operator -n codexray
```
  </TabItem>

  <TabItem value="ee" label="Enterprise Edition (operator)">

:::info
codexray Enterprise Edition is a paid subscription (from $1 per CPU core/month) that offers extra features and priority support.
To install the Enterprise Edition, you'll need a license. [Start](https://codexray.com/account) your free trial today.
:::

Add the codexray helm chart repo:

```bash
helm repo add codexray https://codexray.github.io/helm-charts
helm repo update codexray
```

Next, install the codexray Operator:

```bash
helm install -n codexray --create-namespace codexray-operator codexray/codexray-operator
```

Install the codexray Enterprise Edition.This chart creates a minimal [codexray Custom Resource](/installation/k8s-operator):

```
helm install -n codexray codexray codexray/codexray-ee \
  --set "licenseKey=codexray-LICENSE-KEY-HERE,clickhouse.shards=2,clickhouse.replicas=2"
```

Forward the codexray port to your machine:

```
kubectl port-forward -n codexray service/codexray-codexray 8080:8080
```

Then, you can access codexray at http://localhost:8080

**Upgrade**

The codexray Operator for Kubernetes automatically upgrades all components.

**Uninstall**

To uninstall codexray run the following command:

```
helm uninstall codexray -n codexray
helm uninstall codexray-operator -n codexray
```
  </TabItem>

<TabItem value="ce-helm" label="Community Edition (Helm, deprecated)">

:::warning
Installing codexray via the Helm chart is deprecated. Please use the codexray Operator instead.
:::

Add the codexray helm chart repo:

```
helm repo add codexray https://codexray.github.io/helm-charts
helm repo update codexray
```

Next, install the chart that includes:

```
helm install --namespace codexray --create-namespace codexray codexray/codexray
```

Forward the codexray port to your machine:

```
kubectl port-forward -n codexray service/codexray 8080:8080
```

Then, you can access codexray at http://localhost:8080

**Upgrade**

To upgrade codexray to the latest version:

```
helm repo update codexray
helm upgrade codexray --namespace codexray codexray/codexray
```

**Uninstall**

To uninstall codexray run the following command:

```
helm uninstall codexray -n codexray
```
  </TabItem>



</Tabs>
