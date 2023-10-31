package networkpolicy

import (
	"net"
	"strings"

	"github.com/kubescape/go-logger"
	"github.com/kubescape/go-logger/helpers"
	"github.com/kubescape/storage/pkg/apis/softwarecomposition"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	storageV1Beta1ApiVersion = "spdx.softwarecomposition.kubescape.io/v1beta1"
)

func GenerateNetworkPolicy(networkNeighbors softwarecomposition.NetworkNeighbors, knownServers []softwarecomposition.KnownServers) (softwarecomposition.GeneratedNetworkPolicy, error) {
	networkPolicy := softwarecomposition.NetworkPolicy{
		Kind:       "NetworkPolicy",
		APIVersion: "networking.k8s.io/v1",
		ObjectMeta: metav1.ObjectMeta{
			Name:      networkNeighbors.Name,
			Namespace: networkNeighbors.Namespace,
			Annotations: map[string]string{
				"generated-by": "kubescape",
			},
		},
	}

	if networkNeighbors.Spec.MatchLabels != nil {
		networkPolicy.Spec.PodSelector.MatchLabels = networkNeighbors.Spec.MatchLabels
	}

	if networkNeighbors.Spec.MatchExpressions != nil {
		networkPolicy.Spec.PodSelector.MatchExpressions = networkNeighbors.Spec.MatchExpressions
	}

	if len(networkNeighbors.Spec.Ingress) > 0 {
		networkPolicy.Spec.PolicyTypes = append(networkPolicy.Spec.PolicyTypes, "Ingress")
	}

	if len(networkNeighbors.Spec.Egress) > 0 {
		networkPolicy.Spec.PolicyTypes = append(networkPolicy.Spec.PolicyTypes, "Egress")
	}

	generatedNetworkPolicy := softwarecomposition.GeneratedNetworkPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GeneratedNetworkPolicy",
			APIVersion: storageV1Beta1ApiVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      networkNeighbors.Name,
			Namespace: networkNeighbors.Namespace,
			Labels:    networkNeighbors.Labels,
		},
	}

	for _, neighbor := range networkNeighbors.Spec.Ingress {

		ingressRules, policyRefs := generateIngressRule(neighbor, knownServers)

		generatedNetworkPolicy.PoliciesRef = append(generatedNetworkPolicy.PoliciesRef, policyRefs...)

		networkPolicy.Spec.Ingress = append(networkPolicy.Spec.Ingress, ingressRules)

	}

	for _, neighbor := range networkNeighbors.Spec.Egress {

		egressRules, policyRefs := generateEgressRule(neighbor, knownServers)

		generatedNetworkPolicy.PoliciesRef = append(generatedNetworkPolicy.PoliciesRef, policyRefs...)

		networkPolicy.Spec.Egress = append(networkPolicy.Spec.Egress, egressRules)

	}

	generatedNetworkPolicy.Spec = networkPolicy

	return generatedNetworkPolicy, nil
}

func generateEgressRule(neighbor softwarecomposition.NetworkNeighbor, knownServers []softwarecomposition.KnownServers) (softwarecomposition.NetworkPolicyEgressRule, []softwarecomposition.PolicyRef) {
	egressRule := softwarecomposition.NetworkPolicyEgressRule{}
	policyRefs := []softwarecomposition.PolicyRef{}

	if neighbor.PodSelector != nil {
		egressRule.To = append(egressRule.To, softwarecomposition.NetworkPolicyPeer{
			PodSelector: neighbor.PodSelector,
		})
	}

	if neighbor.NamespaceSelector != nil {
		// the ns label goes together with the pod label
		if len(egressRule.To) > 0 {
			egressRule.To[0].NamespaceSelector = neighbor.NamespaceSelector
		} else {
			// TOD0(DanielGrunberegerCA): is this a valid case?
			egressRule.To = append(egressRule.To, softwarecomposition.NetworkPolicyPeer{
				NamespaceSelector: neighbor.NamespaceSelector,
			})
		}
	}

	if neighbor.IPAddress != "" {
		isKnownServer := false
		// look if this IP is part of any known server
		for _, knownServer := range knownServers {
			_, subNet, err := net.ParseCIDR(knownServer.IPBlock)
			if err != nil {
				logger.L().Error("error parsing cidr", helpers.Error(err))
				continue
			}
			if subNet.Contains(net.ParseIP(neighbor.IPAddress)) {
				egressRule.To = append(egressRule.To, softwarecomposition.NetworkPolicyPeer{
					IPBlock: &softwarecomposition.IPBlock{
						CIDR: knownServer.IPBlock,
					},
				})
				isKnownServer = true

				policyRef := softwarecomposition.PolicyRef{
					Name:       knownServer.Name,
					OriginalIP: neighbor.IPAddress,
					IPBlock:    knownServer.IPBlock,
				}

				if knownServer.DNS != "" {
					policyRef.DNS = knownServer.DNS
				}

				policyRefs = append(policyRefs, policyRef)
				break
			}
		}

		if !isKnownServer {
			ipBlock := &softwarecomposition.IPBlock{CIDR: neighbor.IPAddress + "/32"}
			egressRule.To = append(egressRule.To, softwarecomposition.NetworkPolicyPeer{
				IPBlock: ipBlock,
			})

			if neighbor.DNS != "" {
				policyRefs = append(policyRefs, softwarecomposition.PolicyRef{
					Name:       neighbor.DNS,
					DNS:        neighbor.DNS,
					IPBlock:    ipBlock.CIDR,
					OriginalIP: neighbor.IPAddress,
				})
			}
		}
	}

	for _, networkPort := range neighbor.Ports {
		protocol := v1.Protocol(strings.ToUpper(string(networkPort.Protocol)))
		portInt32 := networkPort.Port

		egressRule.Ports = append(egressRule.Ports, softwarecomposition.NetworkPolicyPort{
			Protocol: &protocol,
			Port:     portInt32,
		})
	}

	return egressRule, policyRefs
}

func generateIngressRule(neighbor softwarecomposition.NetworkNeighbor, knownServers []softwarecomposition.KnownServers) (softwarecomposition.NetworkPolicyIngressRule, []softwarecomposition.PolicyRef) {
	ingressRule := softwarecomposition.NetworkPolicyIngressRule{}
	policyRefs := []softwarecomposition.PolicyRef{}

	if neighbor.PodSelector != nil {
		ingressRule.From = append(ingressRule.From, softwarecomposition.NetworkPolicyPeer{
			PodSelector: neighbor.PodSelector,
		})
	}
	if neighbor.NamespaceSelector != nil {
		// the ns label goes together with the pod label
		if len(ingressRule.From) > 0 {
			ingressRule.From[0].NamespaceSelector = neighbor.NamespaceSelector
		} else {
			// TOD0(DanielGrunberegerCA): is this a valid case?
			ingressRule.From = append(ingressRule.From, softwarecomposition.NetworkPolicyPeer{
				NamespaceSelector: neighbor.NamespaceSelector,
			})
		}
	}

	if neighbor.IPAddress != "" {
		isKnownServer := false
		// look if this IP is part of any known server
		for _, knownServer := range knownServers {
			_, subNet, err := net.ParseCIDR(knownServer.IPBlock)
			if err != nil {
				logger.L().Error("error parsing cidr", helpers.Error(err))
				continue
			}
			if subNet.Contains(net.ParseIP(neighbor.IPAddress)) {
				ingressRule.From = append(ingressRule.From, softwarecomposition.NetworkPolicyPeer{
					IPBlock: &softwarecomposition.IPBlock{
						CIDR: knownServer.IPBlock,
					},
				})
				isKnownServer = true

				policyRef := softwarecomposition.PolicyRef{
					Name:       knownServer.Name,
					OriginalIP: neighbor.IPAddress,
					IPBlock:    knownServer.IPBlock,
				}

				if knownServer.DNS != "" {
					policyRef.DNS = knownServer.DNS
				}

				policyRefs = append(policyRefs, policyRef)
				break
			}
		}

		if !isKnownServer {
			ipBlock := &softwarecomposition.IPBlock{CIDR: neighbor.IPAddress + "/32"}
			ingressRule.From = append(ingressRule.From, softwarecomposition.NetworkPolicyPeer{
				IPBlock: ipBlock,
			})

			if neighbor.DNS != "" {
				policyRefs = append(policyRefs, softwarecomposition.PolicyRef{
					Name:       neighbor.DNS,
					DNS:        neighbor.DNS,
					IPBlock:    ipBlock.CIDR,
					OriginalIP: neighbor.IPAddress,
				})
			}
		}
	}

	for _, networkPort := range neighbor.Ports {
		protocol := v1.Protocol(strings.ToUpper(string(networkPort.Protocol)))
		portInt32 := networkPort.Port

		ingressRule.Ports = append(ingressRule.Ports, softwarecomposition.NetworkPolicyPort{
			Protocol: &protocol,
			Port:     portInt32,
		})
	}

	return ingressRule, policyRefs
}
