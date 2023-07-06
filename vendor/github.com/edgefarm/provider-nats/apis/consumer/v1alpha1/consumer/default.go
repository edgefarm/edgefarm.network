package consumer

const (
	// DefaultMaxWait is the default number of waiting pull requests.
	DefaultMaxWait = 512
)

type ConsumerType int

const (
	ConsumerTypePull ConsumerType = iota
	ConsumerTypePush
	ConsumerTypeNone
)

func (c *ConsumerConfig) SetDefaults(t ConsumerType) {
	c.DeliverPolicy = "All"
	c.AckPolicy = "Explicit"
	c.AckWait = "30s"
	c.Replicas = 0
	c.ReplayPolicy = "Instant"
	c.FilterSubject = ""
	c.MaxDeliver = -1
	c.MaxAckPending = 1000
	c.BackOff = ""

	if t == ConsumerTypePush {
		c.PushConsumer.HeadersOnly = false
	} else if t == ConsumerTypePull {
		c.PullConsumer.MaxWaiting = func() *int { i := DefaultMaxWait; return &i }()
		c.PullConsumer.MaxRequestExpires = ""
		c.PullConsumer.MaxRequestBatch = 0
		c.PullConsumer.MaxRequestMaxBytes = 0
	}

}
