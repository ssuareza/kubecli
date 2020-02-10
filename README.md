# Kubernetes simplified client

This is just a simplified Kubernetes client made in Go.

To build:
```sh
cd cmd/kubecli
# Linux
GOARCH=amd64 GOOS=linux go build -o kubecli

# Mac
GOARCH=amd64 GOOS=darwin go build -o kubecli
```