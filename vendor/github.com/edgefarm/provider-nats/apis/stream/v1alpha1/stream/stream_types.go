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

package stream

// +kubebuilder:object:generate=true
// StreamConfig will determine the properties for a stream.
// There are sensible defaults for most.
type StreamConfig struct {
	// Description is a human readable description of the stream.
	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`

	// Subjects is a list of subjects to consume, supports wildcards.
	// +kubebuilder:validation:Optional
	Subjects []string `json:"subjects,omitempty"`

	// Retention defines the retention policy for the stream.
	// +kubebuilder:validation:Enum=Limits;Interest;WorkQueue
	// +kubebuilder:validation:Required
	// +kubebuilder:default=Limits
	Retention string `json:"retention"`

	// MaxConsumers defines how many Consumers can be defined for a given Stream.
	// Define -1 for unlimited.
	// +kubebuilder:default=-1
	// +kubebuilder:validation:Required
	MaxConsumers int `json:"maxConsumers"`

	// MaxMsgs defines how many messages may be in a Stream.
	// Adheres to Discard Policy, removing oldest or refusing new messages if the Stream exceeds this number of messages.
	// +kubebuilder:default=-1
	// +kubebuilder:validation:Required
	MaxMsgs int64 `json:"maxMsgs"`

	// MaxBytes defines how many bytes the Stream may contain.
	// Adheres to Discard Policy, removing oldest or refusing new messages if the Stream exceeds this size.
	// +kubebuilder:default=-1
	// +kubebuilder:validation:Required
	MaxBytes int64 `json:"maxBytes"`

	// Discard defines the behavior of discarding messages when any streams' limits have been reached.
	// Old (default): This policy will delete the oldest messages in order to maintain the limit. For example, if MaxAge is set to one minute, the server will automatically delete messages older than one minute with this policy.
	// New: This policy will reject new messages from being appended to the stream if it would exceed one of the limits. An extension to this policy is DiscardNewPerSubject which will apply this policy on a per-subject basis within the stream.
	// +kubebuilder:validation:Enum=Old;New
	// +kubebuilder:default=Old
	// +kubebuilder:validation:Required
	Discard string `json:"discard"`

	// DiscardOldPerSubject will discard old messages per subject.
	// +kubebuilder:default=false
	// +kubebuilder:validation:Optional
	DiscardNewPerSubject bool `json:"discardNewPerSubject"`

	// MaxAge is the maximum age of a message in the stream.
	// Format is a string duration, e.g. 1h, 1m, 1s, 1h30m or 2h3m4s.
	// +kubebuilder:validation:Pattern="([0-9]+h)?([0-9]+m)?([0-9]+s)?"
	// +kubebuilder:default="0s"
	// +kubebuilder:validation:Optional
	MaxAge string `json:"maxAge"`

	// MaxMsgsPerSubject defines the limits how many messages in the stream to retain per subject.
	// +kubebuilder:default=-1
	// +kubebuilder:validation:Minimum=-1
	// +kubebuilder:validation:Optional
	MaxMsgsPerSubject int64 `json:"maxMsgsPerSubject"`

	// MaxBytesPerSubject defines the largest message that will be accepted by the Stream.
	// +kubebuilder:default=-1
	// +kubebuilder:validation:Minimum=-1
	// +kubebuilder:validation:Optional
	MaxMsgSize int32 `json:"maxMsgSize"`

	// Storage defines the storage type for stream data..
	// +kubebuilder:validation:Enum=File;Memory
	// +kubebuilder:default=File
	Storage string `json:"storage"`

	// Replicas defines how many replicas to keep for each message in a clustered JetStream.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=5
	// +kubebuilder:validation:Optional
	Replicas int `json:"replicas"`

	// NoAck is a flag to disable acknowledging messages that are received by the Stream.
	// +kubebuilder:default=false
	// +kubebuilder:validation:Optional
	NoAck bool `json:"noAck"`

	// Template is the owner of the template associated with this stream.
	// +kubebuilder:validation:Optional
	TemplateOwner string `json:"template,omitempty"`

	// Duplicates defines the time window within which to track duplicate messages.
	// +kubebuilder:default="2m0s"
	// +kubebuilder:validation:Pattern="^(([0-9]+[smh]){1,3})$"
	// +kubebuilder:validation:Optional
	Duplicates string `json:"duplicates,omitempty"`

	// Placement is the placement policy for the stream.
	// +kubebuilder:validation:Optional
	Placement *Placement `json:"placement,omitempty"`

	// Mirror is the mirror configuration for the stream.
	// +kubebuilder:validation:Optional
	Mirror *StreamSource `json:"mirror,omitempty"`

	// Sources is the list of one or more sources configurations for the stream.
	// +kubebuilder:validation:Optional
	Sources []*StreamSource `json:"sources,omitempty"`

	// Sealed is a flag to prevent message deletion from  the stream  via limits or API.
	// +kubebuilder:validation:Optional
	Sealed bool `json:"sealed,omitempty"`

	// DenyDelete is a flag to restrict the ability to delete messages from a stream via the API.
	// +kubebuilder:validation:Optional
	DenyDelete bool `json:"denyDelete,omitempty"`

	// DenyPurge is a flag to restrict the ability to purge messages from a stream via the API.
	// +kubebuilder:validation:Optional
	DenyPurge bool `json:"denyPurge,omitempty"`

	// AllowRollup is a flag to allow the use of the Nats-Rollup header to replace all contents of a stream, or subject in a stream, with a single new message.
	// +kubebuilder:validation:Optional
	AllowRollup bool `json:"allowRollup,omitempty"`

	// Allow republish of the message after being sequenced and stored.
	// +kubebuilder:validation:Optional
	RePublish *RePublish `json:"rePublish,omitempty"`

	// AllowDirect is a flag that if true and the stream has more than one replica, each replica will respond to direct get requests for individual messages, not only the leader.
	// +kubebuilder:validation:Optional
	AllowDirect bool `json:"allowDirect,omitempty"`

	// MirrorDirect is a flag that if true, and the stream is a mirror, the mirror will participate in a serving direct get requests for individual messages from origin stream.
	// +kubebuilder:validation:Optional
	MirrorDirect bool `json:"mirrorDirect,omitempty"`
}

// RePublish is for republishing messages once committed to a stream. The original subject cis remapped from the subject pattern to the destination pattern.
// For information on RePublish see https://docs.nats.io/nats-concepts/jetstream/streams#republish
type RePublish struct {
	// Source is an optional subject pattern which is a subset of the subjects bound to the stream. It defaults to all messages in the stream, e.g. >.
	// +kubebuilder:default=">"
	Source string `json:"source"`

	// Destination is the destination subject messages will be re-published to. The source and destination must be a valid subject mapping.
	// For information on subject mapping see https://docs.nats.io/jetstream/concepts/subjects#subject-mapping
	Destination string `json:"destination"`

	// HeadersOnly defines if true, that the message data will not be included in the re-published message, only an additional header Nats-Msg-Size indicating the size of the message in bytes.
	// +kubebuilder:validation:Optional
	HeadersOnly bool `json:"headersOnly,omitempty"`
}

// Placement is used to guide placement of streams in clustered JetStream.
// For information on Placement see https://docs.nats.io/nats-concepts/jetstream/streams#placement
type Placement struct {
	// Cluster is the name of the Jetstream cluster.
	Cluster string `json:"cluster"`

	// Tags defines a list of server tags.
	// +kubebuilder:validation:Optional
	Tags []string `json:"tags,omitempty"`
}

// StreamSource dictates how streams can source from other streams.
type StreamSource struct {
	// Name of the origin stream to source messages from.
	Name string `json:"name"`

	// StartSeq is an optional start sequence the of the origin stream to start mirroring from.
	// +kubebuilder:validation:Optional
	StartSeq uint64 `json:"startSeq,omitempty"`

	// StartTime is an optional message start time to start mirroring from. Any messages that are equal to or greater than the start time will be included.
	// The time format is RFC 3339, e.g. 2023-01-09T14:48:32Z
	// +kubebuilder:validation:Pattern="^((?:(\\d{4}-\\d{2}-\\d{2})T(\\d{2}:\\d{2}:\\d{2}(?:\\.\\d+)?))(Z|[\\+-]\\d{2}:\\d{2})?)$"
	// +kubebuilder:validation:Optional
	StartTime string `json:"startTime,omitempty"`

	// FilterSubject is an optional filter subject which will include only messages that match the subject, typically including a wildcard.
	// +kubebuilder:validation:Optional
	FilterSubject string `json:"filterSubject,omitempty"`

	// Domain is the JetStream domain of where the origin stream exists. This is commonly used between a cluster/supercluster and a leaf node/cluster.
	Domain string `json:"domain,omitempty"`

	// External is the external stream configuration.
	// +kubebuilder:validation:Optional
	External *ExternalStream `json:"external,omitempty"`
}

// ExternalStream allows you to qualify access to a stream source in another
// account.
type ExternalStream struct {
	// APIPrefix is the prefix for the API of the external stream.
	APIPrefix string `json:"apiPrefix"`

	// DeliverPrefix is the prefix for the deliver subject of the external stream.
	// +kubebuilder:validation:Optional
	DeliverPrefix string `json:"deliverPrefix,omitempty"`
}

type StreamObservationState struct {
	// Mesasges is the number of messages in the stream.
	Messages uint64 `json:"messages"`
	// Bytes is the number of bytes in the stream.
	Bytes string `json:"bytes"`
	// FirstSequence is the first sequence number in the stream.
	FirstSequence uint64 `json:"firstSequence"`
	// FirstTimestamp is the first timestamp in the stream.
	FirstTimestamp string `json:"firstTimestamp,omitempty"`
	// LastSequence is the last sequence number in the stream.
	LastSequence uint64 `json:"lastSequence"`
	// LastTimestamp is the last timestamp in the stream.
	LastTimestamp string `json:"lastTimestamp,omitempty"`
	// ConsumerCount is the number of consumers in the stream.
	ConsumerCount int `json:"consumerCount"`
	// Deleted TBD
	Deleted []uint64 `json:"deleted,omitempty"`
	// NumDeleted TBD
	NumDeleted int `json:"numDeleted,omitempty"`
	// NumSubjects is the number of subjects in the stream.
	NumSubjects uint64 `json:"numSubjects,omitempty"`
	// Subjects is a map of subjects to their number of messages.
	Subjects map[string]uint64 `json:"subjects,omitempty"`
}

// StreamObservationClusterInfo shows information about the underlying set of servers
// that make up the stream or consumer.
type StreamObservationClusterInfo struct {
	// Name is the name of the cluster.
	Name string `json:"name,omitempty"`
	// Leader is the leader of the cluster.
	Leader string `json:"leader,omitempty"`
	// Replicas are the replicas of the cluster.
	Replicas []*PeerInfo `json:"replicas,omitempty"`
}

type StreamObservationConnection struct {
	// Address is the address of the connection.
	Address string `json:"address"`
	// AccountPublicKey is the public key of the used account.
	AccountPublicKey string `json:"accountPublicKey"`
	// UserPublicKey is the public key of the used user.
	UserPublicKey string `json:"userPublicKey"`
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
