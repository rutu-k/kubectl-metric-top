# kubectl-metric-top


Provides a handy plugin to examine performance of application running on Kubernetes Pods.

### Prerequisites

- Prometheus Operator/Deployment
- Prometheus Adapter

### Build the binary
```
go build ./cmd/cli
```

### Transfer binary to local/bin to create kubectl plugin
```
cp cli /usr/local/bin/kubectl-metric-top
```

### To enable new custom metric, add the required metric in the main.go and build the binary




