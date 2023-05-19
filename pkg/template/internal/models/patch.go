// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package models

import (
	"fmt"
	"os"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/hashicorp/go-multierror"
	"github.com/siderolabs/gen/pair"
	"gopkg.in/yaml.v3"

	"github.com/siderolabs/omni-client/pkg/omni/resources"
	"github.com/siderolabs/omni-client/pkg/omni/resources/omni"
)

// PatchList is a list of patches.
type PatchList []Patch

// Patch is a Talos machine configuration patch.
type Patch struct { //nolint:govet
	// Name of the patch.
	//
	// Optional for 'path' patches, mandatory for 'inline' patches.
	Name string `yaml:"name"`

	// File path to the file containing the patch.
	//
	// Mutually exclusive with `inline:`.
	File string `yaml:"file"`

	// Inline patch content.
	Inline map[string]any `yaml:"inline"`
}

// Validate the model.
func (l PatchList) Validate() error {
	var multiErr error

	for _, patch := range l {
		multiErr = joinErrors(multiErr, patch.Validate())
	}

	return multiErr
}

// Translate the list of patches into a list of resources.
func (l PatchList) Translate(prefix string, baseWeight int, labels ...pair.Pair[string, string]) ([]resource.Resource, error) {
	resources := make([]resource.Resource, 0, len(l))

	for i, patch := range l {
		r, err := patch.Translate(prefix, baseWeight+i, labels...)
		if err != nil {
			return nil, err
		}

		resources = append(resources, r)
	}

	return resources, nil
}

// Validate the model.
func (patch *Patch) Validate() error {
	var multiErr error

	name := patch.Name
	if name == "" {
		name = patch.File
	}

	if patch.File != "" && patch.Inline != nil {
		multiErr = multierror.Append(multiErr, fmt.Errorf("path and inline are mutually exclusive"))
	}

	if patch.File == "" && patch.Inline == nil {
		multiErr = multierror.Append(multiErr, fmt.Errorf("path or inline is required"))
	}

	switch {
	case patch.File != "":
		raw, err := os.ReadFile(patch.File)
		if err != nil {
			multiErr = multierror.Append(multiErr, fmt.Errorf("failed to access %q: %w", patch.File, err))
		} else {
			if err = omni.ValidateConfigPatch(string(raw)); err != nil {
				multiErr = multierror.Append(multiErr, fmt.Errorf("failed to validate patch %q: %w", patch.File, err))
			}
		}
	case patch.Inline != nil:
		if patch.Name == "" {
			multiErr = multierror.Append(multiErr, fmt.Errorf("name is required for inline patches"))
		}

		raw, err := yaml.Marshal(patch.Inline)
		if err != nil {
			multiErr = multierror.Append(multiErr, fmt.Errorf("failed to marshal inline patch %q: %w", name, err))
		} else {
			if err = omni.ValidateConfigPatch(string(raw)); err != nil {
				multiErr = multierror.Append(multiErr, fmt.Errorf("failed to validate inline patch %q: %w", name, err))
			}
		}
	}

	if multiErr != nil {
		return fmt.Errorf("patch %q is invalid: %w", name, multiErr)
	}

	return nil
}

// Translate the model into a resource.
func (patch *Patch) Translate(prefix string, weight int, labels ...pair.Pair[string, string]) (*omni.ConfigPatch, error) {
	name := patch.Name
	if name == "" {
		name = patch.File
	}

	id := fmt.Sprintf("%03d-%s-%s", weight, prefix, name)

	var (
		raw []byte
		err error
	)

	switch {
	case patch.File != "":
		raw, err = os.ReadFile(patch.File)
	case patch.Inline != nil:
		raw, err = yaml.Marshal(patch.Inline)
	default:
		panic("missing patch contents?")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read patch %q: %w", name, err)
	}

	patchResource := omni.NewConfigPatch(resources.DefaultNamespace, id, labels...)
	patchResource.Metadata().Annotations().Set("name", name)

	patchResource.TypedSpec().Value.Data = string(raw)

	return patchResource, nil
}
