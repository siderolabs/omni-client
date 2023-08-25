// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package models

import (
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/siderolabs/gen/pair"

	"github.com/siderolabs/omni-client/api/omni/specs"
	"github.com/siderolabs/omni-client/pkg/constants"
	"github.com/siderolabs/omni-client/pkg/omni/resources"
	"github.com/siderolabs/omni-client/pkg/omni/resources/omni"
)

// MachineSet is a base model for controlplane and workers.
type MachineSet struct {
	Meta `yaml:",inline"`

	// Name is the name of the machine set. When empty, the default name will be used.
	Name string `yaml:"name"`

	// ControlPlane machines.
	Machines MachineIDList `yaml:"machines"`

	// Controlplane patches.
	Patches PatchList `yaml:"patches"`
}

// Translate the model.
func (machineset *MachineSet) translate(ctx TranslateContext, nameSuffix, roleLabel string) ([]resource.Resource, error) {
	id := omni.AdditionalWorkersResourceID(ctx.ClusterName, nameSuffix)

	machineSet := omni.NewMachineSet(resources.DefaultNamespace, id)
	machineSet.Metadata().Labels().Set(omni.LabelCluster, ctx.ClusterName)
	machineSet.Metadata().Labels().Set(roleLabel, "")

	machineSet.TypedSpec().Value.UpdateStrategy = specs.MachineSetSpec_Rolling

	resourceList := []resource.Resource{machineSet}

	for _, machineID := range machineset.Machines {
		machineSetNode := omni.NewMachineSetNode(resources.DefaultNamespace, string(machineID), machineSet)

		_, locked := ctx.LockedMachines[machineID]
		if locked {
			machineSetNode.Metadata().Annotations().Set(omni.MachineLocked, "")
		}

		resourceList = append(resourceList, machineSetNode)
	}

	patches, err := machineset.Patches.Translate(
		id,
		constants.PatchBaseWeightMachineSet,
		pair.MakePair(omni.LabelCluster, ctx.ClusterName),
		pair.MakePair(omni.LabelMachineSet, id),
	)
	if err != nil {
		return nil, err
	}

	return append(resourceList, patches...), nil
}
