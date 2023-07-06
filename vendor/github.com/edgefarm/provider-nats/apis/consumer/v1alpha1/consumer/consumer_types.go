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

package consumer

// +kubebuilder:object:generate=true
// ConsumerConfig will determine the properties for a JetStream consumer.
// For more information see https://docs.nats.io/jetstream/concepts/consumers
type ConsumerConfig struct {
	// Description is a human readable description of the consumer.
	// This can be particularly useful for ephemeral consumers to indicate their purpose since the durable name cannot be provided.
	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`

	// DeliverPolicy defines the point in the stream to receive messages from, either All, Last, New, ByStartSequence, ByStartTime, or LastPerSubject.
	// Fore more information see https://docs.nats.io/jetstream/concepts/consumers#deliverpolicy
	// +kubebuilder:validation:Enum=All;Last;New;ByStartSequence;ByStartTime;LastPerSubject
	// +kubebuilder:default=All
	// +kubebuilder:validation:Required
	DeliverPolicy string `json:"deliverPolicy"`

	// OptStartSeq is an optional start sequence number and is used with the DeliverByStartSequence deliver policy.
	// +kubebuilder:validation:Optional
	OptStartSeq uint64 `json:"optStartSeq,omitempty"`

	// OptStartTime is an optional start time and is used with the DeliverByStartTime deliver policy.
	// The time format is RFC 3339, e.g. 2023-01-09T14:48:32Z
	// +kubebuilder:validation:Pattern="^((?:(\\d{4}-\\d{2}-\\d{2})T(\\d{2}:\\d{2}:\\d{2}(?:\\.\\d+)?))(Z|[\\+-]\\d{2}:\\d{2})?)$"
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	OptStartTime string `json:"optStartTime,omitempty"`

	// AckPolicy describes the requirement of client acknowledgements, either Explicit, None, or All.
	// For more information see https://docs.nats.io/nats-concepts/jetstream/consumers#ackpolicy
	// +kubebuilder:validation:Enum=Explicit;None;All
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:default=Explicit
	AckPolicy string `json:"ackPolicy"`

	// AckWait is the duration that the server will wait for an ack for any individual message once it has been delivered to a consumer.
	// If an ack is not received in time, the message will be redelivered.
	// Format is a string duration, e.g. 1h, 1m, 1s, 1h30m or 2h3m4s.
	// +kubebuilder:validation:Pattern="([0-9]+h)?([0-9]+m)?([0-9]+s)?"
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:default="30s"
	AckWait string `json:"ackWait"`

	// MaxDeliver is the maximum number of times a specific message delivery will be attempted.
	// Applies to any message that is re-sent due to ack policy (i.e. due to a negative ack, or no ack sent by the client).
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=-1
	MaxDeliver int `json:"maxDeliver,omitempty"`

	// Backoff is a list of time durations that represent the time to delay based on delivery count.
	// Format of the durations is a string duration, e.g. 1h, 1m, 1s, 1h30m or 2h3m4s where multiple durations are separated by commas.
	// Example: `1s,2s,3s,4s,5s`.
	// +kubebuilder:validation:Pattern="^(([0-9]+h)?([0-9]+m)?([0-9]+s)?)(?:,\\s*(([0-9]+h)?([0-9]+m)?([0-9]+s)?))*$"
	// +kubebuilder:validation:Optional
	BackOff string `json:"backoff,omitempty"`

	// FilterSubject defines an overlapping subject with the subjects bound to the stream which will filter the set of messages received by the consumer.
	FilterSubject string `json:"filterSubject,omitempty"`

	// ReplayPolicy is used to define the mode of message replay.
	// If the policy is Instant, the messages will be pushed to the client as fast as possible while adhering to the Ack Policy, Max Ack Pending and the client's ability to consume those messages.
	// If the policy is Original, the messages in the stream will be pushed to the client at the same rate that they were originally received, simulating the original timing of messages.
	// +kubebuilder:validation:Enum=Instant;Original
	// +kubebuilder:validation:Required
	// +kubebuilder:default=Instant
	ReplayPolicy string `json:"replayPolicy"`

	// SampleFrequency sets the percentage of acknowledgements that should be sampled for observability.
	// +kubebuilder:validation:Pattern="^([1-9][0-9]?|100)%?$"
	// +kubebuilder:validation:Optional
	SampleFrequency string `json:"sampleFreq,omitempty"`

	// MaxAckPending sets the number of outstanding acks that are allowed before message delivery is halted.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=1000
	MaxAckPending int `json:"maxAckPending,omitempty"`

	// InactiveThreshold defines the duration that instructs the server to cleanup consumers that are inactive for that long.
	// Format is a string duration, e.g. 1h, 1m, 1s, 1h30m or 2h3m4s.
	// +kubebuilder:validation:Pattern="([0-9]+h)?([0-9]+m)?([0-9]+s)?"
	// +kubebuilder:validation:Optional
	InactiveThreshold string `json:"inactiveThreshold,omitempty"`

	// Replicas sets the number of replicas for the consumer's state.
	// By default, when the value is set to zero, consumers inherit the number of replicas from the stream.
	// +kubebuilder:validation:Required
	// +kubebuilder:default=0
	Replicas int `json:"numReplicas"`

	// MemoryStorage if set, forces the consumer state to be kept in memory rather than inherit the storage type of the stream (file in this case).
	// +kubebuilder:validation:Optional
	MemoryStorage bool `json:"memStorage,omitempty"`

	// PullConsumer defines the pull-based consumer configuration.
	// +kubebuilder:validation:Optional
	PullConsumer *PullConsumerSpec `json:"pull,omitempty"`

	// PushConsumer defines the push-based consumer configuration.
	// +kubebuilder:validation:Optional
	PushConsumer *PushConsumerSpec `json:"push,omitempty"`
}

// PushConsumerSpec defines the pull-based consumer configuration.
// For more information, see https://docs.nats.io/nats-concepts/jetstream/consumers#push-specific
type PushConsumerSpec struct {
	// RateLimit is used to throttle the delivery of messages to the consumer, in bits per second.
	// +kubebuilder:validation:Optional
	RateLimit uint64 `json:"rateLimitBps,omitempty"`

	// HeadersOnly delivers, if set, only the headers of messages in the stream and not the bodies.
	// Additionally adds Nats-Msg-Size header to indicate the size of the removed payload.
	// +kubebuilder:validation:Optional
	HeadersOnly bool `json:"headersOnly,omitempty"`

	// DeliverSubject defines the subject to deliver messages to.
	// Note, setting this field implicitly decides whether the consumer is push or pull-based.
	// With a deliver subject, the server will push messages to client subscribed to this subject.
	// This is a push consumer specific setting.
	// +kubebuilder:validation:Required
	DeliverSubject string `json:"deliverSubject,omitempty"`

	// DeliverGroup defines the queue group name which, if specified, is then used to distribute the messages between the subscribers to the consumer.
	// This is analogous to a queue group in core NATS. See https://docs.nats.io/nats-concepts/core-nats/queue for more information on queue groups.
	// This is a push consumer specific setting.
	// +kubebuilder:validation:Optional
	DeliverGroup string `json:"deliverGroup,omitempty"`

	// FlowControl enables per-subscription flow control using a sliding-window protocol.
	// This protocol relies on the server and client exchanging messages to regulate when and how many messages are pushed to the client.
	// This one-to-one flow control mechanism works in tandem with the one-to-many flow control imposed by MaxAckPending across all subscriptions bound to a consumer.
	// This is a push consumer specific setting.
	// +kubebuilder:validation:Optional
	FlowControl bool `json:"flowControl,omitempty"`

	// IdleHeartbeat defines, if set, that the server will regularly send a status message to the client (i.e. when the period has elapsed) while there are no new messages to send.
	// This lets the client know that the JetStream service is still up and running, even when there is no activity on the stream.
	// The message status header will have a code of 100. Unlike FlowControl, it will have no reply to address.
	// It may have a description such "Idle Heartbeat".
	// Note that this heartbeat mechanism is all handled transparently by supported clients and does not need to be handled by the application.
	// Format is a string duration, e.g. 1h, 1m, 1s, 1h30m or 2h3m4s.
	// This is a push consumer specific setting.
	// +kubebuilder:validation:Pattern="([0-9]+h)?([0-9]+m)?([0-9]+s)?"
	// +kubebuilder:validation:Optional
	IdleHeartbeat string `json:"idleHeartbeat,omitempty"`
}

// PullConsumerSpec defines the pull-based consumer configuration.
// For more information, see https://docs.nats.io/nats-concepts/jetstream/consumers#pull-specific
type PullConsumerSpec struct {
	// MaxWaiting defines the maximum number of waiting pull requests.
	// This is a pull consumer specific setting.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=512
	MaxWaiting *int `json:"maxWaiting,omitempty"`

	// MaxRequestExpires defines the maximum duration a single pull request will wait for messages to be available to pull.
	// This is a pull consumer specific setting.
	// +kubebuilder:validation:Pattern="([0-9]+h)?([0-9]+m)?([0-9]+s)?"
	// +kubebuilder:validation:Optional
	MaxRequestExpires string `json:"maxExpires,omitempty"`

	// MaxRequestBatch defines th maximum batch size a single pull request can make.
	// When set with MaxRequestMaxBytes, the batch size will be constrained by whichever limit is hit first.
	// This is a pull consumer specific setting.
	// +kubebuilder:validation:Optional
	MaxRequestBatch int `json:"maxBatch,omitempty"`

	// MaxRequestMaxBytes defines the  maximum total bytes that can be requested in a given batch.
	// When set with MaxRequestBatch, the batch size will be constrained by whichever limit is hit first.
	// This is a pull consumer specific setting.
	// +kubebuilder:validation:Optional
	MaxRequestMaxBytes int `json:"maxBytes,omitempty"`
}

// ConsumerInfo is the info from a JetStream consumer.
type ConsumerObservationState struct {
	// Domain is the domain of the consumer.
	Domain string `json:"domain"`
	// Stream is the stream name.
	Stream string `json:"streamName"`
	// Name is the consumer name.
	Name string `json:"name"`
	// Durable is the durable name.
	Durable string `json:"durableName"`
	// Created is the time the consumer was created.
	// needs to be converted to time.Time
	Created string `json:"created"`
	// Delivered is the consumer sequence and last activity.
	Delivered SequenceInfo `json:"delivered"`
	// AckFloor TBD
	AckFloor SequenceInfo `json:"ackFloor"`
	// NumAckPending is the number of messages pending acknowledgement.
	NumAckPending int `json:"numAckPending"`
	// NumRedelivered is the number of redelivered messages.
	NumRedelivered int `json:"numRedelivered"`
	// NumWaiting is the number of messages waiting to be delivered.
	NumWaiting int `json:"numWaiting"`
	// NumPending is the number of messages pending.
	NumPending uint64 `json:"numPending"`
	// Cluster is the cluster information.
	Cluster *ClusterInfo `json:"cluster,omitempty"`
	// PushBound is whether the consumer is push bound.
	PushBound string `json:"pushBound,omitempty"`
}

// SequenceInfo has both the consumer and the stream sequence and last activity.
type SequenceInfo struct {
	// Consumer is the consumer name
	Consumer uint64 `json:"consumerSeq"`
	// Stream is the name of the stream
	Stream uint64 `json:"streamSeq"`
	// Last is the last time the consumer was active
	// needs to be converted to time.Time
	Last string `json:"lastActive,omitempty"`
}

// StreamObservationClusterInfo shows information about the underlying set of servers
// that make up the stream or consumer.
type ClusterInfo struct {
	// Name is the name of the cluster.
	Name string `json:"name,omitempty"`
	// Leader is the leader of the cluster.
	Leader string `json:"leader,omitempty"`
	// Replicas are the replicas of the cluster.
	Replicas []*PeerInfo `json:"replicas,omitempty"`
}

// PeerInfo shows information about all the peers in the cluster that
// are supporting the stream or consumer.
type PeerInfo struct {
	Name    string `json:"name"`
	Current bool   `json:"current"`
	Offline bool   `json:"offline,omitempty"`
	Active  string `json:"active"`
	Lag     uint64 `json:"lag,omitempty"`
}
