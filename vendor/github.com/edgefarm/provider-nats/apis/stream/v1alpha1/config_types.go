package v1alpha1

/*
Copyright 2023 The Crossplane Authors.

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

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"

	"github.com/edgefarm/provider-nats/apis/stream/v1alpha1/stream"
)

// StreamParameters are the configurable fields of a Stream.
type StreamParameters struct {
	// Domain is the Jetstream domain in which the stream is created.
	// +kubebuilder:validation:Optional
	Domain string `json:"domain,omitempty"`

	// Config is the stream configuration.
	Config stream.StreamConfig `json:"config"`
}

// StreamObservation are the observable fields of a Stream.
type StreamObservation struct {
	// Domain is the Jetstream domain in which the stream is created.
	Domain string `json:"domain,omitempty"`

	// State is the current state of the stream
	State stream.StreamObservationState `json:"state,omitempty"`

	// ClusterInfo shows information about the underlying set of servers that make up the stream.
	ClusterInfo stream.StreamObservationClusterInfo `json:"clusterInfo,omitempty"`

	// Connection shows information about the connection to the stream.
	Connection stream.StreamObservationConnection `json:"connection,omitempty"`
}

// A StreamSpec defines the desired state of a Stream.
type StreamSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       StreamParameters `json:"forProvider"`
}

// A StreamStatus represents the observed state of a Stream.
type StreamStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          StreamObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:object:generate=true
// +genclient
// +genclient:nonNamespaced

// A Stream is an example API type.
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DOMAIN",type="string",JSONPath=".spec.forProvider.domain"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="ADDRESS",type="string",priority=1,JSONPath=".status.atProvider.connection.address"
// +kubebuilder:printcolumn:name="ACCOUNT PUB KEY",type="string",priority=1,JSONPath=".status.atProvider.connection.accountPublicKey"
// +kubebuilder:printcolumn:name="MESSAGES",type="string",priority=1,JSONPath=".status.atProvider.state.messages"
// +kubebuilder:printcolumn:name="BYTES",type="string",priority=1,JSONPath=".status.atProvider.state.bytes"
// +kubebuilder:printcolumn:name="CONSUMERS",type="string",priority=1,JSONPath=".status.atProvider.state.consumerCount"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,nats}
type Stream struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StreamSpec   `json:"spec"`
	Status StreamStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StreamList contains a list of Stream
type StreamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Stream `json:"items"`
}

// Stream type metadata.
var (
	StreamKind             = reflect.TypeOf(Stream{}).Name()
	StreamGroupKind        = schema.GroupKind{Group: Group, Kind: StreamKind}.String()
	StreamKindAPIVersion   = StreamKind + "." + SchemeGroupVersion.String()
	StreamGroupVersionKind = SchemeGroupVersion.WithKind(StreamKind)
)

func init() {
	SchemeBuilder.Register(&Stream{}, &StreamList{})
}
