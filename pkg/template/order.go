// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package template

import (
	"fmt"
	"sort"

	"github.com/cosi-project/runtime/pkg/resource"

	"github.com/siderolabs/omni-client/pkg/omni/resources/omni"
)

// Canonical order of resources in the generated list.
var canonicalResourceOrder = map[resource.Type]int{
	omni.ClusterType:        1,
	omni.ConfigPatchType:    2,
	omni.MachineSetType:     3,
	omni.MachineSetNodeType: 4,
}

func sortResources[T any](s []T, mapper func(T) resource.Metadata) {
	sort.SliceStable(s, func(i, j int) bool {
		orderI := canonicalResourceOrder[mapper(s[i]).Type()]
		orderJ := canonicalResourceOrder[mapper(s[j]).Type()]

		if orderI == 0 {
			panic(fmt.Sprintf("unknown resource type %q", mapper(s[i]).Type()))
		}

		if orderJ == 0 {
			panic(fmt.Sprintf("unknown resource type %q", mapper(s[j]).Type()))
		}

		return orderI < orderJ
	})
}
