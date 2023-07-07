package convert

import (
	"time"
)

func TimeToRFC3339(t *time.Time) (string, error) {
	str, err := t.MarshalText()
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func RFC3339ToTime(t string) (*time.Time, error) {
	r := &time.Time{}
	err := r.UnmarshalText([]byte(t))
	if err != nil {
		return r, err
	}
	return r, nil
}
