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

	"github.com/edgefarm/provider-nats/apis/consumer/v1alpha1/consumer"
)

// ConsumerParameters are the configurable fields of a consumer.
type ConsumerParameters struct {
	// Stream is the name of the Jetstream stream the consumer is created for.
	// +kubebuilder:validation:Required
	Stream string `json:"stream"`

	// Domain is the domain of the Jetstream stream the consumer is created for.
	// +kubebuilder:validation:Optional
	Domain string `json:"domain,omitempty"`

	// Config is the consumer configuration.
	// +kubebuilder:validation:Required
	Config consumer.ConsumerConfig `json:"config"`
}

// ConsumerObservation are the observable fields of a consumer.
type ConsumerObservation struct {
	// State is the current state of the consumer
	State consumer.ConsumerObservationState `json:"state,omitempty"`
}

// A ConsumerSpec defines the desired state of a consumer.
type ConsumerSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       ConsumerParameters `json:"forProvider"`
}

// A ConsumerStatus represents the observed state of a consumer.
type ConsumerStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ConsumerObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:object:generate=true
// +genclient
// +genclient:nonNamespaced

// A Consumer is an example API type.
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="DOMAIN",type="string",JSONPath=".spec.forProvider.domain"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="PUSH CONSUMER",type="string",priority=1,JSONPath=".status.atProvider.state.pushBound"
// +kubebuilder:printcolumn:name="STREAM",type="string",priority=1,JSONPath=".status.atProvider.state.streamName"
// +kubebuilder:printcolumn:name="UNPROCESSED",type="string",priority=1,JSONPath=".status.atProvider.state.numPending"
// +kubebuilder:printcolumn:name="REDELIVERERD",type="string",priority=1,JSONPath=".status.atProvider.state.numRedelivered"
// +kubebuilder:printcolumn:name="ACK PENDING",type="string",priority=1,JSONPath=".status.atProvider.state.numAckPending"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,nats}
type Consumer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConsumerSpec   `json:"spec"`
	Status ConsumerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ConsumerList contains a list of consumer
type ConsumerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Consumer `json:"items"`
}

// Consumer type metadata.
var (
	ConsumerKind             = reflect.TypeOf(Consumer{}).Name()
	ConsumerGroupKind        = schema.GroupKind{Group: Group, Kind: ConsumerKind}.String()
	ConsumerKindAPIVersion   = ConsumerKind + "." + SchemeGroupVersion.String()
	ConsumerGroupVersionKind = SchemeGroupVersion.WithKind(ConsumerKind)
)

func init() {
	SchemeBuilder.Register(&Consumer{}, &ConsumerList{})
}
