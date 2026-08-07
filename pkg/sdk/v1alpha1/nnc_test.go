// Copyright (c) 2026 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package v1alpha1

import (
	"context"
	"testing"

	netopv1alpha1 "github.com/vmware-tanzu/net-operator-api/api/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrlfake "sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/stretchr/testify/require"
)

func testScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = netopv1alpha1.AddToScheme(scheme)
	return scheme
}

func TestNetworksOwnedByNamespaceNetworkConfiguration(t *testing.T) {
	netB := &netopv1alpha1.Network{ObjectMeta: metav1.ObjectMeta{
		Name:      "net-b",
		Namespace: "ns-1",
		Labels:    map[string]string{netopv1alpha1.ManagedByNNCLabelKey: "my-nnc"},
	}}
	netA := &netopv1alpha1.Network{ObjectMeta: metav1.ObjectMeta{
		Name:      "net-a",
		Namespace: "ns-1",
		Labels:    map[string]string{netopv1alpha1.ManagedByNNCLabelKey: "my-nnc"},
	}}
	netOtherNamespace := &netopv1alpha1.Network{ObjectMeta: metav1.ObjectMeta{
		Name:      "net-c",
		Namespace: "ns-2",
		Labels:    map[string]string{netopv1alpha1.ManagedByNNCLabelKey: "my-nnc"},
	}}
	netUnmanaged := &netopv1alpha1.Network{ObjectMeta: metav1.ObjectMeta{
		Name:      "net-unmanaged",
		Namespace: "ns-1",
		Labels:    map[string]string{netopv1alpha1.ManagedByNNCLabelKey: "other-nnc"},
	}}

	c := ctrlfake.NewClientBuilder().WithScheme(testScheme()).
		WithObjects(netB, netA, netOtherNamespace, netUnmanaged).Build()

	t.Run("filters by namespace and sorts by name", func(t *testing.T) {
		got, err := NetworksOwnedByNamespaceNetworkConfiguration(context.Background(), c, "my-nnc", "ns-1")
		require.NoError(t, err)
		require.Len(t, got, 2)
		require.Equal(t, "net-a", got[0].GetName())
		require.Equal(t, "net-b", got[1].GetName())
	})

	t.Run("empty namespace searches across all namespaces", func(t *testing.T) {
		got, err := NetworksOwnedByNamespaceNetworkConfiguration(context.Background(), c, "my-nnc", "")
		require.NoError(t, err)
		require.Len(t, got, 3)
	})

	t.Run("no matches returns empty, not nil", func(t *testing.T) {
		got, err := NetworksOwnedByNamespaceNetworkConfiguration(context.Background(), c, "unknown-nnc", "")
		require.NoError(t, err)
		require.Empty(t, got)
	})
}
