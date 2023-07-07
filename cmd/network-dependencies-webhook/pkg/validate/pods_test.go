package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFromSuffixLabelKey(t *testing.T) {
	assert := assert.New(t)
	labels := map[string]string{
		"apps.openyurt.io/controller-revision-hash":                  "acceleration-network-demo-colibri-colibri-ffc5d5dc7",
		"apps.openyurt.io/pool-name":                                 "sparrow",
		"name.network.edgefarm.io/acceleration-network-demo-colibri": "",
		"network.edgefarm.io/type":                                   "leaf",
		"pod-template-hash":                                          "dbf46c65",
		"subnetwork.network.edgefarm.io/colibri":                     "",
	}
	l, err := GetFromSuffixLabelKey(labels, "name.network.edgefarm.io/")
	assert.Nil(err)
	assert.Equal("acceleration-network-demo-colibri", l)
	l, err = GetFromSuffixLabelKey(labels, "subnetwork.network.edgefarm.io/")
	assert.Nil(err)
	assert.Equal("colibri", l)
	_, err = GetFromSuffixLabelKey(labels, "missing/")
	assert.NotNil(err)
}
