// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package topo

import (
	"encoding/json"
	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LoadKind(t *testing.T) {
	jsonData := `{"somekind": {"type": "kind", "name": "SomeKind", "labels": {"env": "test", "what": "kindly"}, "onos.topo.Location": {"lat": 3.14, "lng": 6.28}}}`

	var data interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	assert.NoError(t, err)

	jsonObjects := data.(map[string]interface{})
	for k, v := range jsonObjects {
		object, err := parseObject(topoapi.ID(k), v)
		assert.NoError(t, err)
		assert.NotNil(t, object)
		assert.Equal(t, topoapi.Object_KIND, object.Type)
		assert.Equal(t, topoapi.ID("somekind"), object.ID)
		assert.Equal(t, "SomeKind", object.GetKind().Name)
		assert.Equal(t, 2, len(object.Labels))
		assert.Equal(t, "test", object.Labels["env"])
		assert.Equal(t, "kindly", object.Labels["what"])
		assert.Equal(t, 1, len(object.Aspects))
		assert.NotNil(t, object.Aspects["onos.topo.Location"])
	}
}

func Test_LoadEntity(t *testing.T) {
	jsonData := `{"foo": {"type": "entity", "kind": "somekind", "labels": {"env": "test", "what": "something"}, "onos.topo.Location": {"lat": 3.14, "lng": 6.28}}}`

	var data interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	assert.NoError(t, err)

	jsonObjects := data.(map[string]interface{})
	for k, v := range jsonObjects {
		object, err := parseObject(topoapi.ID(k), v)
		assert.NoError(t, err)
		assert.NotNil(t, object)
		assert.Equal(t, topoapi.Object_ENTITY, object.Type)
		assert.Equal(t, topoapi.ID("foo"), object.ID)
		assert.Equal(t, topoapi.ID("somekind"), object.GetEntity().KindID)
		assert.Equal(t, 2, len(object.Labels))
		assert.Equal(t, "test", object.Labels["env"])
		assert.Equal(t, "something", object.Labels["what"])
		assert.Equal(t, 1, len(object.Aspects))
		assert.NotNil(t, object.Aspects["onos.topo.Location"])
	}
}

func Test_LoadRelation(t *testing.T) {
	jsonData := `{"rel": {"type": "relation", "kind": "somekind", "source": "foo", "target": "bar", "labels": {"env": "test", "what": "relative"}, "onos.topo.Location": {"lat": 3.14, "lng": 6.28}}}`

	var data interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	assert.NoError(t, err)

	jsonObjects := data.(map[string]interface{})
	for k, v := range jsonObjects {
		object, err := parseObject(topoapi.ID(k), v)
		assert.NoError(t, err)
		assert.NotNil(t, object)
		assert.Equal(t, topoapi.Object_RELATION, object.Type)
		assert.Equal(t, topoapi.ID("rel"), object.ID)
		assert.Equal(t, topoapi.ID("somekind"), object.GetRelation().KindID)
		assert.Equal(t, topoapi.ID("foo"), object.GetRelation().SrcEntityID)
		assert.Equal(t, topoapi.ID("bar"), object.GetRelation().TgtEntityID)
		assert.Equal(t, 2, len(object.Labels))
		assert.Equal(t, "test", object.Labels["env"])
		assert.Equal(t, "relative", object.Labels["what"])
		assert.Equal(t, 1, len(object.Aspects))
		assert.NotNil(t, object.Aspects["onos.topo.Location"])
	}
}
