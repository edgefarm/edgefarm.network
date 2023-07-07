# network-dependencies-webhook

This helm chart is part of edgefarm.network and contains network-dependencies-webhook, a kubernetes webhook that prevents edgenetwork pods from deltetion when there are still NATS stream resources available.
Note, that for this to work you need to have Crossplane (>=v1.11.3) installed and configured.

## Prerequisites

    Kubernetes 1.22+
    Helm 3.2.0+
    Crossplane 1.11.3+