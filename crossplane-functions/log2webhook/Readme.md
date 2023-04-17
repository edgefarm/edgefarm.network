Use this in a Crossplane Composition as a function to dump the contents of a resource to a webhook.

This works for normal webhook services like `https://webhook.site` but could also work with others (untested).

```yaml
piVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: <your composition name>
  labels:
    crossplane.io/xrd: <your xrd name>
spec:
#   ...
  functions:
    - name: log2wehook
      type: Container
      config:
        apiVersion: log2webhook.xfn.edgefarm.io/v1alpha1
        kind: Config
        spec:
          webhookURL: https://webhook.site/<your webhook id>
      container:
        image: ghcr.io/edgefarm/edgefarm.network/xfn-log2webhook:latest
        network:
          policy: Runner
```