package consumer

import (
	"strings"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/edgefarm/provider-nats/internal/convert"

	errors "github.com/edgefarm/provider-nats/apis/consumer/v1alpha1/consumer/errors"
)

func convertDeliverPolicy(delivery string) nats.DeliverPolicy {
	switch delivery {
	case "All":
		return nats.DeliverAllPolicy
	case "Last":
		return nats.DeliverLastPolicy
	case "New":
		return nats.DeliverNewPolicy
	case "ByStartSequence":
		return nats.DeliverByStartSequencePolicy
	case "ByStartTime":
		return nats.DeliverByStartTimePolicy
	case "LastPerSubject":
		return nats.DeliverLastPerSubjectPolicy
	default:
		return nats.DeliverAllPolicy
	}
}

func convertAckPolicy(policy string) nats.AckPolicy {
	switch policy {
	case "None":
		return nats.AckNonePolicy
	case "All":
		return nats.AckAllPolicy
	case "Explicit":
		return nats.AckExplicitPolicy
	default:
		return nats.AckExplicitPolicy
	}
}

func convertReplayPolicy(replay string) nats.ReplayPolicy {
	switch replay {
	case "Instant":
		return nats.ReplayInstantPolicy
	case "Original":
		return nats.ReplayOriginalPolicy
	default:
		return nats.ReplayInstantPolicy
	}
}

func convertBase(name string, config *ConsumerConfig) *nats.ConsumerConfig {
	return &nats.ConsumerConfig{
		Durable:         name,
		Name:            name,
		Description:     config.Description,
		DeliverPolicy:   convertDeliverPolicy(config.DeliverPolicy),
		OptStartSeq:     config.OptStartSeq,
		AckPolicy:       convertAckPolicy(config.AckPolicy),
		MaxDeliver:      config.MaxDeliver,
		FilterSubject:   config.FilterSubject,
		ReplayPolicy:    convertReplayPolicy(config.ReplayPolicy),
		SampleFrequency: config.SampleFrequency,
		MaxAckPending:   config.MaxAckPending,
		Replicas:        config.Replicas,
		MemoryStorage:   config.MemoryStorage,
	}
}

func convertTimes(in *ConsumerConfig, out *nats.ConsumerConfig) error {
	if in.OptStartTime != "" {
		optStartTime, err := convert.RFC3339ToTime(in.OptStartTime)
		if err != nil {
			return err
		}
		out.OptStartTime = optStartTime
	}

	return nil
}

func convertPushConsumer(in *ConsumerConfig, out *nats.ConsumerConfig) error {
	if in.PushConsumer != nil {
		out.DeliverSubject = in.PushConsumer.DeliverSubject
		out.DeliverGroup = in.PushConsumer.DeliverGroup
		out.FlowControl = in.PushConsumer.FlowControl
		if in.PushConsumer.IdleHeartbeat != "" {
			dur, err := time.ParseDuration(in.PushConsumer.IdleHeartbeat)
			if err != nil {
				return err
			}
			out.Heartbeat = dur
		}
		out.RateLimit = in.PushConsumer.RateLimit
		out.HeadersOnly = in.PushConsumer.HeadersOnly
	}
	return nil
}

func convertPullConsumer(in *ConsumerConfig, out *nats.ConsumerConfig) error {
	if in.PullConsumer != nil {
		out.MaxRequestBatch = in.PullConsumer.MaxRequestBatch
		if in.PullConsumer.MaxRequestExpires != "" {
			dur, err := time.ParseDuration(in.PullConsumer.MaxRequestExpires)
			if err != nil {
				return err
			}
			out.MaxRequestExpires = dur
		}
		out.MaxRequestMaxBytes = in.PullConsumer.MaxRequestMaxBytes
		if in.PullConsumer.MaxWaiting != nil {
			out.MaxWaiting = *in.PullConsumer.MaxWaiting
		} else {
			out.MaxWaiting = DefaultMaxWait
		}
		return nil
	}

	// Ensure that even if no PullConsumer is defined, we set a default value
	out.MaxWaiting = DefaultMaxWait
	return nil
}

func convertDurations(in *ConsumerConfig, out *nats.ConsumerConfig) error {
	if in.AckWait != "" {
		dur, err := time.ParseDuration(in.AckWait)
		if err != nil {
			return err
		}
		out.AckWait = dur
	}

	if in.InactiveThreshold != "" {
		dur, err := time.ParseDuration(in.InactiveThreshold)
		if err != nil {
			return err
		}
		out.InactiveThreshold = dur
	}

	if in.BackOff != "" {
		var backOff []time.Duration
		split := strings.Split(in.BackOff, ",")
		for _, b := range split {

			dur, err := time.ParseDuration(b)
			if err != nil {
				return err
			}
			backOff = append(backOff, dur)
		}
		out.BackOff = backOff
	}

	return nil
}

func ConfigV1Alpha1ToNats(name string, config *ConsumerConfig) (*nats.ConsumerConfig, error) {
	natsConfig := convertBase(name, config)
	err := convertDurations(config, natsConfig)
	if err != nil {
		return &nats.ConsumerConfig{}, err
	}

	err = convertTimes(config, natsConfig)
	if err != nil {
		return &nats.ConsumerConfig{}, err
	}

	if config.PushConsumer != nil && config.PullConsumer != nil {
		return nil, errors.PushAndPullConsumerError
	}

	err = convertPushConsumer(config, natsConfig)
	if err != nil {
		return &nats.ConsumerConfig{}, err
	}

	if config.PushConsumer == nil {
		err = convertPullConsumer(config, natsConfig)
		if err != nil {
			return &nats.ConsumerConfig{}, err
		}
	}
	return natsConfig, nil
}
