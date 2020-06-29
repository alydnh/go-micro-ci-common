package tests

import (
	"github.com/alydnh/go-micro-ci-common/yaml"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestOpenCI(t *testing.T) {
	path := filepath.Join("files", "micro-ci.yaml")
	ci, err := yaml.OpenCI(path, false)
	if !assert.Nil(t, err) {
		t.Fatal("Unexpected Error", err)
	}
	assert.Equal(t, "micro-ci-test", ci.Name())
	assert.Equal(t, "v1-value", ci.Variables["v1"])
	assert.Equal(t, "v2-value", ci.Variables["v2"])
	assert.Equal(t, "cev1-value", ci.CommonEnvs["cev1"])
	assert.Equal(t, "cev2-value", ci.CommonEnvs["cev2"])
	assert.Equal(t, "v2-value-consul", ci.Registry.Address)
}
