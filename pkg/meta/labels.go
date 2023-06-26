// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package meta keeps Talos meta partition utils.
package meta

import "gopkg.in/yaml.v3"

// InitialLabels describes structure that is stored in the Talos metadata and keeps machine labels
// that are initially assigned to the machine when it connects to Omni.
type InitialLabels struct {
	Labels map[string]string `yaml:"machineInitialLabels"`
}

// Encode converts labels to the serialized value to be stored in the meta partition.
func (l InitialLabels) Encode() ([]byte, error) {
	return yaml.Marshal(l)
}

// ParseLabels reads label from the encoded metadata value.
func ParseLabels(data []byte) (*InitialLabels, error) {
	labels := &InitialLabels{}

	if err := yaml.Unmarshal(data, &labels); err != nil {
		return nil, err
	}

	return labels, nil
}
