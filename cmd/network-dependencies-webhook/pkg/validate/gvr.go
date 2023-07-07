package validate

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	consumerResource = schema.GroupVersionResource{
		Group:    "streams.network.edgefarm.io",
		Version:  "v1alpha1",
		Resource: "xconsumers",
	}

	streamResource = schema.GroupVersionResource{
		Group:    "streams.network.edgefarm.io",
		Version:  "v1alpha1",
		Resource: "xstreams",
	}

	userResource = schema.GroupVersionResource{
		Group:    "issue.natssecrets.crossplane.io",
		Version:  "v1alpha1",
		Resource: "users",
	}

	accountResource = schema.GroupVersionResource{
		Group:    "issue.natssecrets.crossplane.io",
		Version:  "v1alpha1",
		Resource: "accounts",
	}

	providerConfigResource = schema.GroupVersionResource{
		Group:    "nats.crossplane.io",
		Version:  "v1alpha1",
		Resource: "providerconfigs",
	}
)
