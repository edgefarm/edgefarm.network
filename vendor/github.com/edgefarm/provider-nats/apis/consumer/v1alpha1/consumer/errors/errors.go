package errors

var (
	// ErrJetStreamNotEnabled is an error returned when both push and pull consumer are defined in the same consumer configuration.
	PushAndPullConsumerError ConsumerError = &consumerError{message: "invalid configuration: cannot define both 'push' and 'pull' consumer"}
)

// ConsumerError is an error result that happens when using consumer.
type ConsumerError interface {
	error
}

type consumerError struct {
	message string
}

func (err *consumerError) Error() string {
	return err.message
}
