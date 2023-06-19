# provider-nats

This helm chart is part of edgefarm.network and contains provider-nats, a Crossplane provider for NATS to manage NATS  jetstream (stream, consumer) resources.
Note, that for this to work you need to have Crossplane (>=v1.11.3) installed and configured.

## Prerequisites

    Kubernetes 1.22+
    Helm 3.2.0+
    Crossplane 1.11.3+
    NATS 2.9.14+