package stream

func (s *StreamConfig) SetDefaults() {
	s.Retention = "Limits"
	s.Discard = "DiscardOld"
	s.Storage = "File"
	s.Replicas = 1
	s.MaxConsumers = -1
	s.MaxMsgs = -1
	s.MaxBytes = -1
	s.MaxMsgSize = -1
	s.MaxAge = "0s"
	s.MaxMsgsPerSubject = -1
	s.Duplicates = "2m0s"
}
