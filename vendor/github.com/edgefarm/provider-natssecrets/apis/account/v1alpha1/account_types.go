/*
Copyright 2022 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/edgefarm/vault-plugin-secrets-nats/pkg/claims/account/v1alpha1"
)

// AccountParameters are the configurable fields of a Account.
type AccountParameters struct {
	Operator string `json:"operator"`
	// +kubebuilder:validation:Optional
	UseSigningKey string `json:"useSigningKey,omitempty"`
	// +kubebuilder:validation:Optional
	Claims v1alpha1.AccountClaims `json:"claims,omitempty"`
}

// AccountObservation are the observable fields of a Account.
type AccountObservation struct {
	Operator   string `json:"operator,omitempty"`
	Account    string `json:"account,omitempty"`
	Issue      string `json:"issue,omitempty"`
	NKey       string `json:"nkey,omitempty"`
	JWT        string `json:"jwt,omitempty"`
	NKeyPath   string `json:"nkeyPath,omitempty"`
	JWTPath    string `json:"jwtPath,omitempty"`
	Pushed     string `json:"pushed,omitempty"`
	LastPushed string `json:"lastPushed,omitempty"`
}

// A AccountSpec defines the desired state of a Account.
type AccountSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       AccountParameters `json:"forProvider"`
}

// A AccountStatus represents the observed state of a Account.
type AccountStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          AccountObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Account is an example API type.
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="OPERATOR",type="string",JSONPath=".status.atProvider.operator"
// +kubebuilder:printcolumn:name="NKEY",type="string",priority=1,JSONPath=".status.atProvider.nkey"
// +kubebuilder:printcolumn:name="JWT",type="string",priority=1,JSONPath=".status.atProvider.jwt"
// +kubebuilder:printcolumn:name="PUSHED",type="string",priority=1,JSONPath=".status.atProvider.pushed"
// +kubebuilder:printcolumn:name="LAST PUSHED",type="string",priority=1,JSONPath=".status.atProvider.lastPushed"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,natssecrets}
type Account struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AccountSpec   `json:"spec"`
	Status AccountStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AccountList contains a list of Account
type AccountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Account `json:"items"`
}

// Account type metadata.
var (
	AccountKind             = reflect.TypeOf(Account{}).Name()
	AccountGroupKind        = schema.GroupKind{Group: Group, Kind: AccountKind}.String()
	AccountKindAPIVersion   = AccountKind + "." + SchemeGroupVersion.String()
	AccountGroupVersionKind = SchemeGroupVersion.WithKind(AccountKind)
)

func init() {
	SchemeBuilder.Register(&Account{}, &AccountList{})
}
