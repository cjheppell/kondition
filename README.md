# Kondition

Kondition is a Go-based HTTP service that runs within a Kubernetes cluster, intended to determine the status of running deployments.

The Kondition server then makes it possible to inspect the status of an internal deployment at a given external URL path, without needing to expose your deployment services.

Think of it like a HTTP API that exposes the state of your internal deployment readiness status.

## How does it work

Create a `services.yaml` file which defines the deployment you want to track and the API path it will be accessible at in Kondition.

After that, hitting `http://<kondition-server>/<service>` will return one of two possibilities:

- HTTP 503 if the deployment is not marked as available
- HTTP 200 if the deployment is marked as available
