// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auth

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/resource/meta"
	"github.com/cosi-project/runtime/pkg/resource/protobuf"
	"github.com/cosi-project/runtime/pkg/resource/typed"
	"github.com/hashicorp/go-multierror"
	"golang.org/x/exp/slices"

	"github.com/siderolabs/omni-client/api/omni/specs"
	"github.com/siderolabs/omni-client/pkg/omni/resources"
	"github.com/siderolabs/omni-client/pkg/omni/resources/omni"
)

const (
	// AccessPolicyGroupPrefix is the special prefix used in the AccessPolicy rules to denote a group of users or clusters.
	AccessPolicyGroupPrefix = "group/"

	// AccessPolicyID is the ID of AccessPolicy resource.
	AccessPolicyID = "access-policy"

	// AccessPolicyType is the type of AccessPolicy resource.
	//
	// tsgen:AccessPolicyType
	AccessPolicyType = resource.Type("AccessPolicies.omni.sidero.dev")
)

// NewAccessPolicy creates new AccessPolicy resource.
func NewAccessPolicy() *AccessPolicy {
	return typed.NewResource[AccessPolicySpec, AccessPolicyExtension](
		resource.NewMetadata(resources.DefaultNamespace, AccessPolicyType, AccessPolicyID, resource.VersionUndefined),
		protobuf.NewResourceSpec(&specs.AccessPolicySpec{}),
	)
}

// AccessPolicy resource describes a user ACL.
type AccessPolicy = typed.Resource[AccessPolicySpec, AccessPolicyExtension]

// AccessPolicySpec wraps specs.AccessPolicySpec.
type AccessPolicySpec = protobuf.ResourceSpec[specs.AccessPolicySpec, *specs.AccessPolicySpec]

// AccessPolicyExtension providers auxiliary methods for AccessPolicy resource.
type AccessPolicyExtension struct{}

// ResourceDefinition implements [typed.Extension] interface.
func (AccessPolicyExtension) ResourceDefinition() meta.ResourceDefinitionSpec {
	return meta.ResourceDefinitionSpec{
		Type:             AccessPolicyType,
		Aliases:          []resource.Type{},
		DefaultNamespace: resources.DefaultNamespace,
		PrintColumns:     []meta.PrintColumn{},
	}
}

// CheckResult is the result of an access policy check.
type CheckResult struct {
	KubernetesImpersonateGroups []string
}

// ValidateAccessPolicy validates the given access policy by running all its tests.
func ValidateAccessPolicy(accessPolicy *AccessPolicy) error {
	var validationErrs error

	if accessPolicy.Metadata().ID() != AccessPolicyID {
		validationErrs = multierror.Append(validationErrs, fmt.Errorf(
			"access policy ID mismatch: expected %q, got %q",
			AccessPolicyID,
			accessPolicy.Metadata().ID(),
		))
	}

	if accessPolicy.Metadata().Namespace() != resources.DefaultNamespace {
		validationErrs = multierror.Append(validationErrs, fmt.Errorf(
			"access policy namespace mismatch: expected %q, got %q",
			resources.DefaultNamespace,
			accessPolicy.Metadata().Namespace(),
		))
	}

	for _, test := range accessPolicy.TypedSpec().Value.GetTests() {
		checkResult := checkAccessPolicy(accessPolicy, test.GetCluster(), test.GetUser())

		expectedImpersonateGroups := append([]string(nil), test.GetExpected().GetKubernetesImpersonateGroups()...)

		sort.Strings(expectedImpersonateGroups)
		sort.Strings(checkResult.KubernetesImpersonateGroups)

		if !slices.Equal(expectedImpersonateGroups, checkResult.KubernetesImpersonateGroups) {
			validationErrs = multierror.Append(validationErrs, fmt.Errorf(
				"access policy test %q failed: kubernetes impersonate groups mismatch: expected %v, got %v",
				test.GetName(),
				expectedImpersonateGroups,
				checkResult.KubernetesImpersonateGroups,
			))
		}
	}

	return validationErrs
}

// CheckAccessPolicy checks the given user against the given cluster, and returns the result of the check, containing
// which groups will be impersonated on the Kubernetes cluster access.
func CheckAccessPolicy(accessPolicy *AccessPolicy, cluster *omni.Cluster, userID string) CheckResult {
	return checkAccessPolicy(accessPolicy, cluster.Metadata().ID(), userID)
}

//nolint:gocognit
func checkAccessPolicy(accessPolicy *AccessPolicy, cluster, userID string) CheckResult {
	if len(accessPolicy.TypedSpec().Value.GetRules()) == 0 {
		return CheckResult{}
	}

	impersonateGroups := make([]string, 0, len(accessPolicy.TypedSpec().Value.GetRules()))

	for _, rule := range accessPolicy.TypedSpec().Value.GetRules() {
		userMatches := false

		for _, ruleUser := range rule.GetUsers() {
			if ruleUser == userID {
				userMatches = true

				break
			}

			if strings.HasPrefix(ruleUser, AccessPolicyGroupPrefix) {
				groupName := ruleUser[len(AccessPolicyGroupPrefix):]

				group, groupOk := accessPolicy.TypedSpec().Value.GetUserGroups()[groupName]
				if !groupOk {
					continue
				}

				for _, groupUser := range group.GetUsers() {
					if groupUser.GetName() == userID {
						userMatches = true

						break
					}
				}
			}
		}

		if !userMatches {
			continue
		}

		clusterMatches := false

		for _, ruleCluster := range rule.GetClusters() {
			if ruleCluster == cluster {
				clusterMatches = true

				break
			}

			if strings.HasPrefix(ruleCluster, AccessPolicyGroupPrefix) {
				groupName := ruleCluster[len(AccessPolicyGroupPrefix):]

				group, groupOk := accessPolicy.TypedSpec().Value.GetClusterGroups()[groupName]
				if !groupOk {
					continue
				}

				for _, groupCluster := range group.GetClusters() {
					if groupCluster.GetName() == cluster {
						clusterMatches = true

						break
					}
				}
			}
		}

		if !clusterMatches {
			continue
		}

		impersonateGroups = append(impersonateGroups, rule.GetKubernetes().GetImpersonate().GetGroups()...)
	}

	return CheckResult{
		KubernetesImpersonateGroups: impersonateGroups,
	}
}
