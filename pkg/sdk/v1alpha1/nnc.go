// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"context"
	"fmt"

	netopv1alpha1 "github.com/vmware-tanzu/net-operator-api/api/v1alpha1"

	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// NetworkNamespaceNetworkConfigurationOwner returns the corresponding NamespaceNetworkConfiguration
// that manages this given Network, and a boolean indicating if the label is present.
// If the boolean is false, this Network is not managed via a NamespaceNetworkConfiguration.
func NetworkNamespaceNetworkConfigurationOwner(network *netopv1alpha1.Network) (string, bool) {
	if network == nil {
		return "", false
	}

	labels := network.GetLabels()
	if labels == nil {
		return "", false
	}

	owner, ok := labels[netopv1alpha1.ManagedByNNCLabelKey]
	return owner, ok
}

// NetworksOwnedByNamespaceNetworkConfiguration retrieves a list of Network resources
// that are managed by the specified NamespaceNetworkConfiguration.
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
func NetworksOwnedByNamespaceNetworkConfiguration(ctx context.Context, c ctrlclient.Client, nncName, namespace string) ([]netopv1alpha1.Network, error) {
	var networkList netopv1alpha1.NetworkList

	listOpts := []ctrlclient.ListOption{
		ctrlclient.MatchingLabels{netopv1alpha1.ManagedByNNCLabelKey: nncName},
	}

	// Conditionally append the namespace filter.
	if namespace != "" {
		listOpts = append(listOpts, ctrlclient.InNamespace(namespace))
	}

	if err := c.List(ctx, &networkList, listOpts...); err != nil {
		return nil, fmt.Errorf("error listing Networks owned by NamespaceNetworkConfiguration '%s' in namespace '%s': %w", nncName, namespace, err)
	}

	return networkList.Items, nil
}
