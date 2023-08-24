// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package models

import (
	"fmt"
	"slices"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/hashicorp/go-multierror"
	"github.com/siderolabs/gen/xslices"

	"github.com/siderolabs/omni-client/pkg/omni/resources/omni"
)

// List implements a list of models, essentially it's a template.
type List []Model

// Validate the set of models as a complete template.
//
// Each model should be valid, but also the set of models should be complete.
func (l List) Validate() error {
	var multiErr error

	for _, model := range l {
		multiErr = joinErrors(multiErr, model.Validate())
	}

	// complete template should contain 1 cluster, 1 controlplane, 1 worker
	// and machines mentioned in the controlplane or worker models
	var (
		clusterCount         int
		controlplaneCount    int
		workersCount         int
		controlPlaneMachines MachineIDList
		workerMachines       MachineIDList
	)

	lockedMachines := make(map[MachineID]struct{})

	for _, model := range l {
		switch m := model.(type) {
		case *Cluster:
			clusterCount++
		case *ControlPlane:
			controlplaneCount++

			controlPlaneMachines = append(controlPlaneMachines, m.Machines...)
		case *Workers:
			workersCount++

			workerMachines = append(workerMachines, m.Machines...)
		case *Machine:
			if m.Locked {
				lockedMachines[m.Name] = struct{}{}
			}
		}
	}

	if clusterCount != 1 {
		multiErr = multierror.Append(multiErr, fmt.Errorf("template should contain 1 cluster, got %d", clusterCount))
	}

	if controlplaneCount != 1 {
		multiErr = multierror.Append(multiErr, fmt.Errorf("template should contain 1 controlplane, got %d", controlplaneCount))
	}

	if workersCount != 1 {
		multiErr = multierror.Append(multiErr, fmt.Errorf("template should contain 1 workers, got %d", workersCount))
	}

	cpMachinesSet := xslices.ToSet(controlPlaneMachines)
	intersection := xslices.Filter(workerMachines, func(id MachineID) bool {
		_, ok := cpMachinesSet[id]

		return ok
	})

	if len(intersection) != 0 {
		multiErr = multierror.Append(multiErr, fmt.Errorf("machines %v are used in both controlplane and workers", intersection))
	}

	allMachines := append(slices.Clone(controlPlaneMachines), workerMachines...)
	allMachinesSet := xslices.ToSet(allMachines)

	for _, model := range l {
		if machine, ok := model.(*Machine); ok {
			if _, ok := allMachinesSet[machine.Name]; !ok {
				multiErr = multierror.Append(multiErr, fmt.Errorf("machine %q is not used in controlplane or workers", machine.Name))
			}
		}
	}

	for _, cpMachine := range controlPlaneMachines {
		if _, ok := lockedMachines[cpMachine]; ok {
			multiErr = multierror.Append(multiErr, fmt.Errorf("machine %q is locked and used in controlplane", cpMachine))
		}
	}

	return multiErr
}

// Translate a set of models (template) to a set of Omni resources.
//
// Translate assumes that the template is valid.
func (l List) Translate() ([]resource.Resource, error) {
	context := TranslateContext{
		LockedMachines: make(map[MachineID]struct{}),
	}

	for _, model := range l {
		switch m := model.(type) {
		case *Cluster:
			context.ClusterName = m.Name
		case *Machine:
			if m.Locked {
				context.LockedMachines[m.Name] = struct{}{}
			}
		}
	}

	var (
		multiErr      error
		resourcesList []resource.Resource
	)

	for _, model := range l {
		resources, err := model.Translate(context)
		if err != nil {
			multiErr = multierror.Append(multiErr, err)

			continue
		}

		resourcesList = append(resourcesList, resources...)
	}

	// perform additional validation:
	// - all resources except for cluster itself should have a cluster label
	// - all resources should have a unique ID
	resourceIDs := map[resource.Type]map[resource.ID]struct{}{}

	for _, r := range resourcesList {
		if _, ok := resourceIDs[r.Metadata().Type()]; !ok {
			resourceIDs[r.Metadata().Type()] = map[resource.ID]struct{}{}
		}

		if _, ok := resourceIDs[r.Metadata().Type()][r.Metadata().ID()]; ok {
			multiErr = multierror.Append(multiErr, fmt.Errorf("resource ID duplicate %s", resource.String(r)))
		}

		resourceIDs[r.Metadata().Type()][r.Metadata().ID()] = struct{}{}

		if r.Metadata().Type() == omni.ClusterType {
			continue
		}

		if l, _ := r.Metadata().Labels().Get(omni.LabelCluster); l != context.ClusterName {
			multiErr = multierror.Append(multiErr, fmt.Errorf("resource %q is missing cluster label", r.Metadata().ID()))
		}
	}

	return resourcesList, multiErr
}

// ClusterName retruns the name of the cluster in the template.
func (l List) ClusterName() (string, error) {
	for _, model := range l {
		if cluster, ok := model.(*Cluster); ok {
			return cluster.Name, nil
		}
	}

	return "", fmt.Errorf("cluster model not found")
}
