// Copyright (c) 2026 Broadcom. All Rights Reserved.
// Broadcom Confidential. The term "Broadcom" refers to Broadcom Inc.
// and/or its subsidiaries.

package v1alpha1

import (
	"context"
	"fmt"
	"sort"

	netopv1alpha1 "github.com/vmware-tanzu/net-operator-api/api/v1alpha1"

	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// NetworksOwnedByNamespaceNetworkConfiguration retrieves a list of Network resources
// that are managed by the specified NamespaceNetworkConfiguration, sorted by name.
//
// As NamespaceNetworkConfiguration's may be scoped to multiple Kubernetes namespaces, an optional
// namespace parameter is provided to help filter the query further. If an empty string ("") is
// provided, it will search for matching Networks across all namespaces at the cluster scope.
// Otherwise, it will only return Networks within the specified namespace.
//
// Example usage:
//
//	// Search across all namespaces
//	nets, err := NetworksOwnedByNamespaceNetworkConfiguration(ctx, c, "my-nnc", "")
//
//	// Search within a specific namespace
//	nets, err := NetworksOwnedByNamespaceNetworkConfiguration(ctx, c, "my-nnc", "default")
func NetworksOwnedByNamespaceNetworkConfiguration(ctx context.Context, c ctrlclient.Client, nncName, namespace string) ([]*netopv1alpha1.Network, error) {
	var networkList netopv1alpha1.NetworkList

	listOpts := []ctrlclient.ListOption{
		ctrlclient.MatchingLabels{netopv1alpha1.ManagedByNNCLabelKey: nncName},
	}

	// Conditionally append the namespace filter
	if namespace != "" {
		listOpts = append(listOpts, ctrlclient.InNamespace(namespace))
	}

	if err := c.List(ctx, &networkList, listOpts...); err != nil {
		return nil, fmt.Errorf("error listing Networks owned by NamespaceNetworkConfiguration '%s' in namespace '%s': %w", nncName, namespace, err)
	}

	networks := make([]*netopv1alpha1.Network, 0, len(networkList.Items))
	for i := range networkList.Items {
		networks = append(networks, &networkList.Items[i])
	}

	// To retain deterministic response, sort the networks before returning.
	sort.Slice(networks, func(i, j int) bool { return networks[i].GetName() < networks[j].GetName() })
	return networks, nil
}
