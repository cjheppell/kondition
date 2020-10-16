# Kondition

Kondition is a Go-based HTTP service that runs within a Kubernetes cluster, intended to determine the status of running containers.

The Kondition server then makes it possible to inspect the status of an internal pod at a given URL path.

Think of it like an publicly accessible service that displays the state of your internal pod readiness probes.
