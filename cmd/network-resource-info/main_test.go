package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPath(t *testing.T) {
	assert := assert.New(t)
	path := []string{"", "domains", "namespace", "default", "network", "default", "subnetwork", "default"}
	p, err := getPath(path)
	assert.Nil(err)
	assert.Equal("default", p.Namespace)
	assert.Equal("default", p.Network)
	assert.Equal("default", p.Subnetwork)
	assert.Nil(p.Stream)

	path = []string{"", "domains", "namespace", "mynamespace", "network", "mynetwork", "subnetwork", "mysubnetwork", "stream", "mystream"}
	p, err = getPath(path)
	assert.Nil(err)
	assert.Equal("mynamespace", p.Namespace)
	assert.Equal("mynetwork", p.Network)
	assert.Equal("mysubnetwork", p.Subnetwork)
	assert.Equal("mystream", *p.Stream)

}
