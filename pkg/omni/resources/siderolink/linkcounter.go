// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package siderolink

import (
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/resource/protobuf"
	"github.com/cosi-project/runtime/pkg/resource/typed"

	"github.com/siderolabs/omni-client/api/omni/specs"
)

// NewLinkCounter creates new LinkCounter state.
func NewLinkCounter(ns, id string) *LinkCounter {
	return typed.NewResource[LinkCounterSpec, LinkCounterExtension](
		resource.NewMetadata(ns, LinkCounterType, id, resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.SiderolinkCounterSpec{}),
	)
}

// LinkCounterType is the type of LinkCounter resource.
//
// tsgen:SiderolinkCounterResourceType
const LinkCounterType = resource.Type("LinkCounters.omni.sidero.dev")

// LinkCounter resource keeps connected nodes state.
//
// LinkCounter resource ID is a machine UUID.
type LinkCounter = typed.Resource[LinkCounterSpec, LinkCounterExtension]

// LinkCounterSpec wraps specs.SiderolinkSpec.
type LinkCounterSpec = protobuf.ResourceSpec[specs.SiderolinkCounterSpec, *specs.SiderolinkCounterSpec]

// LinkCounterExtension providers auxiliary methods for LinkCounter resource.
type LinkCounterExtension struct{}

// ResourceDefinition implements [typed.Extension] interface.
func (LinkCounterExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{
		Type:             LinkCounterType,
		Aliases:          []resource.Type{},
		DefaultNamespace: CounterNamespace,
		PrintColumns: []meta.PrintColumn{
			{
				Name:     "RX",
				JSONPath: "{.bytesreceived}",
			},
			{
				Name:     "TX",
				JSONPath: "{.bytessent}",
			},
		},
	}
}
