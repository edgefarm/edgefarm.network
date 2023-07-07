package v1alpha1

import "github.com/edgefarm/vault-plugin-secrets-nats/pkg/claims/user/v1alpha1"

func FixEmptySlices(params *v1alpha1.UserClaims) {
	if params == nil {
		return
	}
	if params.Permissions.Pub.Allow == nil {
		params.Permissions.Pub.Allow = []string{}
	}
	if params.Permissions.Pub.Deny == nil {
		params.Permissions.Pub.Deny = []string{}
	}
	if params.Permissions.Sub.Allow == nil {
		params.Permissions.Sub.Allow = []string{}
	}
	if params.Permissions.Sub.Deny == nil {
		params.Permissions.Sub.Deny = []string{}
	}
}
