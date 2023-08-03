// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package omni provides resources describing the Machines, Clusters, etc.
package omni

const (
	// SystemLabelPrefix is the prefix of all labels which are managed by the COSI controllers.
	// tsgen:SystemLabelPrefix.
	SystemLabelPrefix = "omni.sidero.dev/"
)

const (
	// Global Labels.

	// LabelControlPlaneRole indicates that the machine is a control plane.
	// tsgen:LabelControlPlaneRole
	LabelControlPlaneRole = SystemLabelPrefix + "role-controlplane"

	// LabelWorkerRole indicates that the machine is a worker.
	// tsgen:LabelWorkerRole
	LabelWorkerRole = SystemLabelPrefix + "role-worker"

	// LabelCluster defines the cluster relation label.
	// tsgen:LabelCluster
	LabelCluster = SystemLabelPrefix + "cluster"

	// LabelHostname defines machine hostname.
	// tsgen:LabelHostname
	LabelHostname = SystemLabelPrefix + "hostname"

	// LabelMachineSet defines the machine set relation label.
	// tsgen:LabelMachineSet
	LabelMachineSet = SystemLabelPrefix + "machine-set"

	// LabelClusterMachine defines the cluster machine relation label.
	// tsgen:LabelClusterMachine
	LabelClusterMachine = SystemLabelPrefix + "cluster-machine"

	// LabelMachine defines the machine relation label.
	// tsgen:LabelMachine
	LabelMachine = SystemLabelPrefix + "machine"

	// LabelSkipTeardown is the test label that configures machine set controller to skip teardown sequence for the cluster machine.
	LabelSkipTeardown = SystemLabelPrefix + "machine-set-skip-teardown"

	// LabelSystemPatch marks the patch as the system patch, so it shouldn't be editable by the user.
	// tsgen:LabelSystemPatch
	LabelSystemPatch = SystemLabelPrefix + "system-patch"

	// LabelExposedServiceAlias is the alias of the exposed service.
	// tsgen:LabelExposedServiceAlias
	LabelExposedServiceAlias = SystemLabelPrefix + "exposed-service-alias"
)

const (
	// MachineStatus labels.

	// MachineStatusLabelConnected is set if the machine is connected.
	// tsgen:MachineStatusLabelConnected
	MachineStatusLabelConnected = SystemLabelPrefix + "connected"

	// MachineStatusLabelDisconnected is set if the machine is disconnected.
	// tsgen:MachineStatusLabelDisconnected
	MachineStatusLabelDisconnected = SystemLabelPrefix + "disconnected"

	// MachineStatusLabelReportingEvents is set if the machine is reporting events.
	// tsgen:MachineStatusLabelReportingEvents
	MachineStatusLabelReportingEvents = SystemLabelPrefix + "reporting-events"

	// MachineStatusLabelAvailable is set if the machine is available to be added to a cluster.
	// tsgen:MachineStatusLabelAvailable
	MachineStatusLabelAvailable = SystemLabelPrefix + "available"

	// MachineStatusLabelArch describes the machine architecture.
	// tsgen:MachineStatusLabelArch
	MachineStatusLabelArch = SystemLabelPrefix + "arch"

	// MachineStatusLabelCPU describes the machine CPU.
	// tsgen:MachineStatusLabelCPU
	MachineStatusLabelCPU = SystemLabelPrefix + "cpu"

	// MachineStatusLabelCores describes the number of machine cores.
	// tsgen:MachineStatusLabelCores
	MachineStatusLabelCores = SystemLabelPrefix + "cores"

	// MachineStatusLabelMem describes the total memory available on the machine.
	// tsgen:MachineStatusLabelMem
	MachineStatusLabelMem = SystemLabelPrefix + "mem"

	// MachineStatusLabelStorage describes the total storage capacity of the machine.
	// tsgen:MachineStatusLabelStorage
	MachineStatusLabelStorage = SystemLabelPrefix + "storage"

	// MachineStatusLabelNet describes the machine network adapters speed.
	// tsgen:MachineStatusLabelNet
	MachineStatusLabelNet = SystemLabelPrefix + "net"

	// MachineStatusLabelPlatform describes the machine platform.
	// tsgen:MachineStatusLabelPlatform
	MachineStatusLabelPlatform = SystemLabelPrefix + "platform"

	// MachineStatusLabelRegion describes the machine region (for machines running in the clouds).
	// tsgen:MachineStatusLabelRegion
	MachineStatusLabelRegion = SystemLabelPrefix + "region"

	// MachineStatusLabelZone describes the machine zone (for machines running in the clouds).
	// tsgen:MachineStatusLabelZone
	MachineStatusLabelZone = SystemLabelPrefix + "zone"

	// MachineStatusLabelInstance describes the machine instance type (for machines running in the clouds).
	// tsgen:MachineStatusLabelInstance
	MachineStatusLabelInstance = SystemLabelPrefix + "instance"
)

const (
	// ClusterMachineStatus labels.

	// ClusterMachineStatusLabelNodeName is set to the node name.
	// tsgen:ClusterMachineStatusLabelNodeName
	ClusterMachineStatusLabelNodeName = SystemLabelPrefix + "node-name"
)

const (
	// Machine labels.

	// MachineAddressLabel is used for faster lookup of the machine by address.
	MachineAddressLabel = SystemLabelPrefix + "address"
)
