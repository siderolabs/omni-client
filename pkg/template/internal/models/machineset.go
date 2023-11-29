// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package models

import (
	"fmt"
	"strconv"

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

	BootstrapSpec *BootstrapSpec `yaml:"bootstrapSpec,omitempty"`

	// MachineSet machines.
	Machines MachineIDList `yaml:"machines"`

	MachineClass *MachineClassConfig `yaml:"machineClass,omitempty"`

	// MachineSet patches.
	Patches PatchList `yaml:"patches"`
}

// MachineClassConfig defines the model for setting the machine class based machine selector in the machine set.
type MachineClassConfig struct {
	// Name defines used machine class name.
	Name string `yaml:"name"`

	// Size sets the number of machines to be pulled from the machine class.
	Size Size `yaml:"size"`
}

// BootstrapSpec defines the model for setting the bootstrap specification, i.e. restoring from a backup, in the machine set.
// Only valid for the control plane machine set.
type BootstrapSpec struct {
	// ClusterUUID defines the UUID of the cluster to restore from.
	ClusterUUID string `yaml:"clusterUUID"`

	// Snapshot defines the snapshot file name to restore from.
	Snapshot string `yaml:"snapshot"`
}

// Size extends protobuf generated allocation type enum to parse string constants.
type Size struct {
	value          uint32
	allocationType specs.MachineSetSpec_MachineClass_AllocationType
}

// UnmarshalYAML implements yaml.Unmarshaller.
func (c *Size) UnmarshalYAML(unmarshal func(any) error) error {
	var value string

	if err := unmarshal(&value); err != nil {
		return err
	}

	switch value {
	case "unlimited", "∞", "infinity":
		value = "Unlimited"
	}

	v, ok := specs.MachineSetSpec_MachineClass_AllocationType_value[value]

	if !ok {
		c.allocationType = specs.MachineSetSpec_MachineClass_Static

		count, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return fmt.Errorf("invalid machine count %s: %w", value, err)
		}

		c.value = uint32(count)
	}

	c.allocationType = specs.MachineSetSpec_MachineClass_AllocationType(v)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (c Size) MarshalYAML() (any, error) {
	if c.allocationType != specs.MachineSetSpec_MachineClass_Static {
		return specs.MachineSetSpec_MachineClass_AllocationType_name[int32(c.allocationType)], nil
	}

	return c.value, nil
}

// Validate checks the machine set fields correctness.
func (machineset *MachineSet) Validate() error {
	if len(machineset.Machines) > 0 && machineset.MachineClass != nil {
		return fmt.Errorf("machine set can not have both machines and machine class defined")
	}

	return nil
}

// Translate the model.
func (machineset *MachineSet) translate(ctx TranslateContext, nameSuffix, roleLabel string) ([]resource.Resource, error) {
	id := omni.AdditionalWorkersResourceID(ctx.ClusterName, nameSuffix)

	machineSet := omni.NewMachineSet(resources.DefaultNamespace, id)
	machineSet.Metadata().Labels().Set(omni.LabelCluster, ctx.ClusterName)
	machineSet.Metadata().Labels().Set(roleLabel, "")

	machineSet.TypedSpec().Value.UpdateStrategy = specs.MachineSetSpec_Rolling

	resourceList := []resource.Resource{machineSet}

	if machineset.BootstrapSpec != nil {
		machineSet.TypedSpec().Value.BootstrapSpec = &specs.MachineSetSpec_BootstrapSpec{
			ClusterUuid: machineset.BootstrapSpec.ClusterUUID,
			Snapshot:    machineset.BootstrapSpec.Snapshot,
		}
	}

	if machineset.MachineClass != nil {
		machineSet.TypedSpec().Value.MachineClass = &specs.MachineSetSpec_MachineClass{
			Name:           machineset.MachineClass.Name,
			MachineCount:   machineset.MachineClass.Size.value,
			AllocationType: machineset.MachineClass.Size.allocationType,
		}
	} else {
		for _, machineID := range machineset.Machines {
			machineSetNode := omni.NewMachineSetNode(resources.DefaultNamespace, string(machineID), machineSet)

			_, locked := ctx.LockedMachines[machineID]
			if locked {
				machineSetNode.Metadata().Annotations().Set(omni.MachineLocked, "")
			}

			resourceList = append(resourceList, machineSetNode)
		}
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
