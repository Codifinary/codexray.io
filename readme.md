# Codexray

Codexray is an observability platform inspired by Coroot, designed for collecting and analyzing logs, traces, performance data, and more from microservices environments. It extends these capabilities to browser and mobile realms through BRUM (Browser Real User Monitoring) and MRUM (Mobile Real User Monitoring) agents, enabling comprehensive client-side telemetry alongside server-side insights.

## Prerequisites

- [Go](https://golang.org/) installed
- [Docker](https://www.docker.com/) installed
- Access to a Kubernetes cluster with `kubectl` configured
- Properly configured ClickHouse and Prometheus integrations

## Running Codexray via `go run`

Clone the repository and navigate to the project directory:

```bash
git clone https://github.com/Codifinary/codexray.io.git
cd codexray
```

Install dependencies and run the application:

```bash
go mod download
go run main.go
```

Codexray will start on the default address (e.g., `0.0.0.0:8080`).

## Running Codexray via Docker

Clone the repository and navigate to the project directory:

```bash
git clone https://github.com/Codifinary/codexray.io.git
cd codexray
```

Build the Docker image and run the container:

```bash
docker build -t codexray:latest . && docker run -p 8080:8080  codexray:latest
```

## Running Codexray on Kubernetes

Codexray includes a Kubernetes deployment configuration file located at `deploy/k8s/codexray.yaml`.

Ensure your Docker image is built and pushed to a registry:

```bash
docker build -t your-registry/codexray:latest . && docker push your-registry/codexray:latest
```

Update the image reference in `deploy/k8s/codexray.yaml` to point to your pushed image, then apply the Kubernetes manifests:

```bash
kubectl apply -f deploy/k8s/codexray.yaml
```

Verify the deployment:

```bash
kubectl get pods
kubectl get svc
```

Ensure that the Codexray pods are running and the service is exposed as expected.
