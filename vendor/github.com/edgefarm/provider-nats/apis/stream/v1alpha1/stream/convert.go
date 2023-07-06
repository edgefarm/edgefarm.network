package stream

import (
	"time"

	"github.com/nats-io/nats.go"

	convert "github.com/edgefarm/provider-nats/internal/convert"
)

func convertRetentionPolicy(retention string) nats.RetentionPolicy {
	switch retention {
	case "Limits":
		return nats.LimitsPolicy
	case "Interest":
		return nats.InterestPolicy
	case "WorkQueue":
		return nats.WorkQueuePolicy
	default:
		return nats.LimitsPolicy
	}
}

func convertDiscardPolicy(discard string) nats.DiscardPolicy {
	switch discard {
	case "Old":
		return nats.DiscardOld
	case "New":
		return nats.DiscardNew
	default:
		return nats.DiscardOld
	}
}

func convertStorage(storage string) nats.StorageType {
	switch storage {
	case "File":
		return nats.FileStorage
	case "Memory":
		return nats.MemoryStorage
	default:
		return nats.FileStorage
	}
}

func convertBase(name string, config *StreamConfig) *nats.StreamConfig {
	return &nats.StreamConfig{
		Name:                 name,
		Description:          config.Description,
		Subjects:             config.Subjects,
		Retention:            convertRetentionPolicy(config.Retention),
		MaxConsumers:         config.MaxConsumers,
		MaxMsgs:              config.MaxMsgs,
		MaxBytes:             config.MaxBytes,
		Discard:              convertDiscardPolicy(config.Discard),
		DiscardNewPerSubject: config.DiscardNewPerSubject,
		MaxMsgsPerSubject:    config.MaxMsgsPerSubject,
		MaxMsgSize:           config.MaxMsgSize,
		Storage:              convertStorage(config.Storage),
		Replicas:             config.Replicas,
		NoAck:                config.NoAck,
		Template:             config.TemplateOwner,
		Sealed:               config.Sealed,
		DenyDelete:           config.DenyDelete,
		DenyPurge:            config.DenyPurge,
		AllowRollup:          config.AllowRollup,
		AllowDirect:          config.AllowDirect,
		MirrorDirect:         config.MirrorDirect,
	}
}

func convertDurations(in *StreamConfig, out *nats.StreamConfig) error {
	if in.MaxAge != "" {
		maxAge, err := time.ParseDuration(in.MaxAge)
		if err != nil {
			return err
		}
		out.MaxAge = maxAge
	}

	if in.Duplicates != "" {
		duplicates, err := time.ParseDuration(in.Duplicates)
		if err != nil {
			return err
		}
		out.Duplicates = duplicates
	}
	return nil
}

func convertMirror(in *StreamConfig, out *nats.StreamConfig) error {
	if in.Mirror != nil {
		var optStartTime *time.Time
		if in.Mirror.StartTime != "" {
			var err error
			optStartTime, err = convert.RFC3339ToTime(in.Mirror.StartTime)
			if err != nil {
				return err
			}
		}
		mirrorConfig := &nats.StreamSource{
			Name:          in.Mirror.Name,
			OptStartSeq:   in.Mirror.StartSeq,
			OptStartTime:  optStartTime,
			FilterSubject: in.Mirror.FilterSubject,
			Domain:        in.Mirror.Domain,
		}
		if in.Mirror.External != nil {
			mirrorConfig.External = &nats.ExternalStream{
				APIPrefix:     in.Mirror.External.APIPrefix,
				DeliverPrefix: in.Mirror.External.DeliverPrefix,
			}
		}
		out.Mirror = mirrorConfig
	}
	return nil
}

func convertSources(in *StreamConfig, out *nats.StreamConfig) error {
	if in.Sources != nil {
		sources := []*nats.StreamSource{}
		for _, source := range in.Sources {
			var optStartTime *time.Time
			if source.StartTime != "" {
				var err error
				optStartTime, err = convert.RFC3339ToTime(source.StartTime)
				if err != nil {
					return err
				}
			}
			streamSource := &nats.StreamSource{
				Name:          source.Name,
				OptStartSeq:   source.StartSeq,
				OptStartTime:  optStartTime,
				FilterSubject: source.FilterSubject,
				Domain:        source.Domain,
			}
			if source.External != nil {
				streamSource.External = &nats.ExternalStream{
					APIPrefix:     source.External.APIPrefix,
					DeliverPrefix: source.External.DeliverPrefix,
				}
			}
			sources = append(sources, streamSource)
		}
		out.Sources = sources
	}
	return nil
}

func ConfigV1Alpha1ToNats(name string, config *StreamConfig) (*nats.StreamConfig, error) {
	natsConfig := convertBase(name, config)
	err := convertDurations(config, natsConfig)
	if err != nil {
		return &nats.StreamConfig{}, err
	}

	if config.Placement != nil {
		placement := &nats.Placement{
			Cluster: config.Placement.Cluster,
			Tags:    config.Placement.Tags,
		}
		natsConfig.Placement = placement
	}

	if config.RePublish != nil {
		republish := &nats.RePublish{
			Source:      config.RePublish.Source,
			Destination: config.RePublish.Destination,
			HeadersOnly: config.RePublish.HeadersOnly,
		}
		natsConfig.RePublish = republish
	}

	err = convertMirror(config, natsConfig)
	if err != nil {
		return &nats.StreamConfig{}, err
	}

	err = convertSources(config, natsConfig)
	if err != nil {
		return &nats.StreamConfig{}, err
	}

	return natsConfig, nil
}

func ConvertPeerInfo(peer *nats.PeerInfo) *PeerInfo {
	return &PeerInfo{
		Name:    peer.Name,
		Current: peer.Current,
		Offline: peer.Offline,
		Active:  peer.Active.String(),
		Lag:     peer.Lag,
	}
}
