// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package provisioner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestZip(t *testing.T) {
	artifacts := map[string][]byte{
		"p4info": []byte("kjsdhfakljshfkjashfkjashfkjhsfkjashfdlkajhsdfkjahsfkajhsdf"),
		"p4bin":  []byte("asdasdasdasdadkjsdhasdadadasdadasdadadadasdasdfakljshfkjashfkjashfkjhsfkjashfdlkajhsdfkjahsfkajhsdf"),
	}
	err := writeArtifacts("/tmp/artifacts.tgz", artifacts)
	assert.NoError(t, err)

	artifacts, err = readArtifacts("/tmp/artifacts.tgz")
	assert.NoError(t, err)

	err = writeArtifacts("/tmp/copy.tgz", artifacts)
	assert.NoError(t, err)

	//af1, err := os.Stat("/tmp/artifacts.tgz")
	//assert.NoError(t, err)
	//
	//af2, err := os.Stat("/tmp/copy.tgz")
	//assert.NoError(t, err)
	//
	//assert.Equal(t, af1.Size(), af2.Size())
}
