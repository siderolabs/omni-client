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
	// DefaultAccessGroup specifies the default Kubernetes group asserted in the token claims if the user has modify access to the clusters.
	//
	// If not, the user will only have the groups specified in the ACLs (AccessPolicies) in the token claims (will be empty if there is no matching ACL).
	DefaultAccessGroup = "system:masters"
)

// GRPCMaxMessageSize is the maximum message size for gRPC server.
const GRPCMaxMessageSize = 32 * 1024 * 1024

// DisableValidation force disable resource validation on the Omni runtime for a particular resource (only for debug build).
const DisableValidation = "disable-validation"

const (
	// DiskConfigPatchPrefix is the prefix of machine install disk config patch.
	// tsgen:DiskConfigPatchPrefix
	DiskConfigPatchPrefix = "000"

	// EncryptionPatchPrefix is the prefix of the encryption config patch.
	// tsgen:EncryptionPatchPrefix
	EncryptionPatchPrefix = "950"
)

const (
	// InstallDiskConfigName human readable install disk config patch name annotation.
	// tsgen:InstallDiskConfigName
	InstallDiskConfigName = "install disk"

	// EncryptionConfigName human readable encryption config patch name annotation.
	// tsgen:EncryptionConfigName
	EncryptionConfigName = "disk encryption config"

	// InstallDiskConfigDescription description of disk config patch.
	// tsgen:InstallDiskConfigDescription
	InstallDiskConfigDescription = "Automatically generated config patch that defines the install disk"

	// EncryptionConfigDescription description of encryption config patch.
	// tsgen:EncryptionConfigDescription
	EncryptionConfigDescription = "Makes machine encrypt disks using Omni as a KMS server"
)
