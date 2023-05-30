// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auth_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/cosi-project/runtime/pkg/resource/protobuf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/siderolabs/omni-client/pkg/omni/resources"
	"github.com/siderolabs/omni-client/pkg/omni/resources/auth"
	"github.com/siderolabs/omni-client/pkg/omni/resources/omni"
)

//go:embed testdata/acl-valid.yaml
var aclValidRaw []byte

//go:embed testdata/acl-invalid-metadata.yaml
var aclInvalidMetadataRaw []byte

func TestCheckAccessPolicy(t *testing.T) {
	accessPolicy := getAccessPolicy(t, aclValidRaw)

	assert.ElementsMatch(t, []string{"k8s-group-1", "k8s-group-2"}, auth.CheckAccessPolicy(accessPolicy,
		omni.NewCluster(resources.DefaultNamespace, "cluster-group-1-cluster-1"), "user-group-1-user-1").KubernetesImpersonateGroups)
	assert.ElementsMatch(t, []string{"k8s-group-1", "k8s-group-2"}, auth.CheckAccessPolicy(accessPolicy,
		omni.NewCluster(resources.DefaultNamespace, "cluster-group-1-cluster-1"), "user-group-1-user-2").KubernetesImpersonateGroups)
	assert.ElementsMatch(t, []string{"k8s-group-1", "k8s-group-2"}, auth.CheckAccessPolicy(accessPolicy,
		omni.NewCluster(resources.DefaultNamespace, "cluster-group-1-cluster-1"), "standalone-user-1").KubernetesImpersonateGroups)
	assert.ElementsMatch(t, []string{"k8s-group-1", "k8s-group-2"}, auth.CheckAccessPolicy(accessPolicy,
		omni.NewCluster(resources.DefaultNamespace, "standalone-cluster-1"), "standalone-user-1").KubernetesImpersonateGroups)

	assert.Empty(t, auth.CheckAccessPolicy(accessPolicy, omni.NewCluster(resources.DefaultNamespace, "cluster-group-1-cluster-1"), "user-group-2-user-1").KubernetesImpersonateGroups)

	assert.ElementsMatch(t, []string{"k8s-group-3", "k8s-group-4"}, auth.CheckAccessPolicy(accessPolicy,
		omni.NewCluster(resources.DefaultNamespace, "cluster-group-2-cluster-1"), "user-group-2-user-1").KubernetesImpersonateGroups)
	assert.ElementsMatch(t, []string{"k8s-group-3", "k8s-group-4"}, auth.CheckAccessPolicy(accessPolicy,
		omni.NewCluster(resources.DefaultNamespace, "cluster-group-2-cluster-1"), "user-group-2-user-2").KubernetesImpersonateGroups)
	assert.ElementsMatch(t, []string{"k8s-group-3", "k8s-group-4"}, auth.CheckAccessPolicy(accessPolicy,
		omni.NewCluster(resources.DefaultNamespace, "cluster-group-2-cluster-1"), "user-group-2-user-3").KubernetesImpersonateGroups)
	assert.ElementsMatch(t, []string{"k8s-group-3", "k8s-group-4"}, auth.CheckAccessPolicy(accessPolicy,
		omni.NewCluster(resources.DefaultNamespace, "cluster-group-2-cluster-1"), "standalone-user-2").KubernetesImpersonateGroups)
	assert.ElementsMatch(t, []string{"k8s-group-3", "k8s-group-4"}, auth.CheckAccessPolicy(accessPolicy,
		omni.NewCluster(resources.DefaultNamespace, "standalone-cluster-2"), "standalone-user-2").KubernetesImpersonateGroups)

	assert.Empty(t, auth.CheckAccessPolicy(accessPolicy, omni.NewCluster(resources.DefaultNamespace, "cluster-group-2-cluster-1"), "user-group-1-user-1").KubernetesImpersonateGroups)
}

func TestValidateAccessPolicy(t *testing.T) {
	accessPolicy := getAccessPolicy(t, aclValidRaw)

	err := auth.ValidateAccessPolicy(accessPolicy)
	assert.NoError(t, err)

	accessPolicy.TypedSpec().Value.Tests[0].Expected.KubernetesImpersonateGroups = []string{"k8s-group-1", "k8s-group-2", "k8s-group-3"}
	accessPolicy.TypedSpec().Value.Tests[1].Expected.KubernetesImpersonateGroups = []string{"k8s-group-3", "k8s-group-4", "k8s-group-5"}

	err = auth.ValidateAccessPolicy(accessPolicy)
	assert.ErrorContains(t, err, "2 errors occurred")
	assert.ErrorContains(t, err, `access policy test "test-1" failed`)
	assert.ErrorContains(t, err, `access policy test "test-2" failed`)
}

func TestValidateAccessPolicyInvalidMetadata(t *testing.T) {
	accessPolicy := getAccessPolicy(t, aclInvalidMetadataRaw)

	err := auth.ValidateAccessPolicy(accessPolicy)
	assert.ErrorContains(t, err, "2 errors occurred")
	assert.ErrorContains(t, err, `access policy ID mismatch`)
	assert.ErrorContains(t, err, `access policy namespace mismatch`)
}

func getAccessPolicy(t *testing.T, raw []byte) *auth.AccessPolicy {
	dec := yaml.NewDecoder(bytes.NewReader(raw))

	var res protobuf.YAMLResource

	err := dec.Decode(&res)
	require.NoError(t, err)

	policy, ok := res.Resource().(*auth.AccessPolicy)
	require.True(t, ok, "resource is not an AccessPolicy")

	return policy
}
