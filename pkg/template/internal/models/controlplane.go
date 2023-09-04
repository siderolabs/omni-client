// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package models

import (
	"fmt"

	"github.com/cosi-project/runtime/pkg/resource"

	"github.com/siderolabs/omni-client/pkg/omni/resources/omni"
)

// KindControlPlane is ControlPlane model kind.
const KindControlPlane = "ControlPlane"

// ControlPlane describes Cluster controlplane nodes.
type ControlPlane struct {
	MachineSet `yaml:",inline"`
}

// Validate the model.
func (controlplane *ControlPlane) Validate() error {
	var multiErr error

	multiErr = joinErrors(multiErr, controlplane.Machines.Validate(), controlplane.Patches.Validate())

	if multiErr != nil {
		return fmt.Errorf("controlplane is invalid: %w", multiErr)
	}

	return nil
}

// Translate the model.
func (controlplane *ControlPlane) Translate(ctx TranslateContext) ([]resource.Resource, error) {
	return controlplane.translate(ctx, "control-planes", omni.LabelControlPlaneRole)
}

func init() {
	register[ControlPlane](KindControlPlane)
}
