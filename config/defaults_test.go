package config_test

import (
	"testing"

	"github.com/EscherAuth/escher/config"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaults_SetAllTheDefaultValuesThatAreEmptyString(t *testing.T) {

	c := config.Config{}
	config.SetDefaults(&c)

	assert.Equal(t, "ESR", c.AlgoPrefix)
	assert.Equal(t, "SHA256", c.HashAlgo)
	assert.Equal(t, "Escher", c.VendorKey)
	assert.Equal(t, "X-Escher-Auth", c.AuthHeaderName)
	assert.Equal(t, "X-Escher-Date", c.DateHeaderName)

}

func TestSetDefaults_DoesNotSetAlgoPrefixIfPresent(t *testing.T) {

	c := config.Config{AlgoPrefix: "Cat"}
	config.SetDefaults(&c)

	assert.Equal(t, "Cat", c.AlgoPrefix)
	assert.Equal(t, "SHA256", c.HashAlgo)
	assert.Equal(t, "Escher", c.VendorKey)
	assert.Equal(t, "X-Escher-Auth", c.AuthHeaderName)
	assert.Equal(t, "X-Escher-Date", c.DateHeaderName)

}

func TestSetDefaults_DoesNotSetHashAlgoIfPresent(t *testing.T) {

	c := config.Config{HashAlgo: "Cat"}
	config.SetDefaults(&c)

	assert.Equal(t, "ESR", c.AlgoPrefix)
	assert.Equal(t, "Cat", c.HashAlgo)
	assert.Equal(t, "Escher", c.VendorKey)
	assert.Equal(t, "X-Escher-Auth", c.AuthHeaderName)
	assert.Equal(t, "X-Escher-Date", c.DateHeaderName)

}

func TestSetDefaults_DoesNotSetVendorKeyIfPresent(t *testing.T) {

	c := config.Config{VendorKey: "Cat"}
	config.SetDefaults(&c)

	assert.Equal(t, "ESR", c.AlgoPrefix)
	assert.Equal(t, "SHA256", c.HashAlgo)
	assert.Equal(t, "Cat", c.VendorKey)
	assert.Equal(t, "X-Escher-Auth", c.AuthHeaderName)
	assert.Equal(t, "X-Escher-Date", c.DateHeaderName)

}

func TestSetDefaults_DoesNotSetAuthHeaderIfPresent(t *testing.T) {

	c := config.Config{AuthHeaderName: "X-Cat-Auth"}
	config.SetDefaults(&c)

	assert.Equal(t, "ESR", c.AlgoPrefix)
	assert.Equal(t, "SHA256", c.HashAlgo)
	assert.Equal(t, "Escher", c.VendorKey)
	assert.Equal(t, "X-Cat-Auth", c.AuthHeaderName)
	assert.Equal(t, "X-Escher-Date", c.DateHeaderName)

}

func TestSetDefaults_DoesNotSetDateHeaderIfPresent(t *testing.T) {

	c := config.Config{DateHeaderName: "X-Cat-Date"}
	config.SetDefaults(&c)

	assert.Equal(t, "ESR", c.AlgoPrefix)
	assert.Equal(t, "SHA256", c.HashAlgo)
	assert.Equal(t, "Escher", c.VendorKey)
	assert.Equal(t, "X-Escher-Auth", c.AuthHeaderName)
	assert.Equal(t, "X-Cat-Date", c.DateHeaderName)

}

func TestSetDefaults_ThereISNoDefaultValueForCredentialScope(t *testing.T) {

	c := config.Config{}
	config.SetDefaults(&c)

	assert.Equal(t, "ESR", c.AlgoPrefix)
	assert.Equal(t, "SHA256", c.HashAlgo)
	assert.Equal(t, "Escher", c.VendorKey)
	assert.Equal(t, "X-Escher-Auth", c.AuthHeaderName)
	assert.Equal(t, "X-Escher-Date", c.DateHeaderName)
	assert.Equal(t, "", c.CredentialScope)

}
