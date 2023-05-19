// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package constants contains global backend constants.
package constants

const (
	// TalosRegistry is the default Talos repository URL.
	TalosRegistry = "ghcr.io/siderolabs/installer"

	// KubernetesRegistry is the default kubernetes repository URL.
	KubernetesRegistry = "ghcr.io/siderolabs/kubelet"
)

const (
	// PatchWeightInstallDisk is the weight of the install disk config patch.
	PatchWeightInstallDisk = 0
	// PatchBaseWeightCluster is the base weight for cluster patches.
	PatchBaseWeightCluster = 200
	// PatchBaseWeightMachineSet is the base weight for machine set patches.
	PatchBaseWeightMachineSet = 400
	// PatchBaseWeightClusterMachine is the base weight for cluster machine patches.
	PatchBaseWeightClusterMachine = 400
)

const (
	// DefaultAccessGroup specifies the default Kubernetes group asserted in the token claims.
	//
	// Later on we might want to customize access level based on authorization.
	DefaultAccessGroup = "system:masters"
)

// GRPCMaxMessageSize is the maximum message size for gRPC server.
const GRPCMaxMessageSize = 32 * 1024 * 1024
