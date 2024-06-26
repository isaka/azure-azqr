// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cr

import (
	"strings"

	"github.com/Azure/azqr/internal/scanners"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
)

// GetRules - Returns the rules for the ContainerRegistryScanner
func (a *ContainerRegistryScanner) GetRules() map[string]scanners.AzureRule {
	return map[string]scanners.AzureRule{
		"cr-001": {
			Id:             "cr-001",
			Category:       scanners.RulesCategoryMonitoringAndAlerting,
			Recommendation: "ContainerRegistry should have diagnostic settings enabled",
			Impact:         scanners.ImpactLow,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				service := target.(*armcontainerregistry.Registry)
				_, ok := scanContext.DiagnosticsSettings[strings.ToLower(*service.ID)]
				return !ok, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/container-registry/monitor-service",
		},
		"cr-002": {
			Id:             "cr-002",
			Category:       scanners.RulesCategoryHighAvailability,
			Recommendation: "ContainerRegistry should have availability zones enabled",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armcontainerregistry.Registry)
				zones := *i.Properties.ZoneRedundancy == armcontainerregistry.ZoneRedundancyEnabled
				return !zones, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/container-registry/zone-redundancy",
		},
		"cr-003": {
			Id:             "cr-003",
			Category:       scanners.RulesCategoryHighAvailability,
			Recommendation: "ContainerRegistry should have a SLA",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				return false, "99.95%"
			},
			Url: "https://www.azure.cn/en-us/support/sla/container-registry/",
		},
		"cr-004": {
			Id:             "cr-004",
			Category:       scanners.RulesCategorySecurity,
			Recommendation: "ContainerRegistry should have private endpoints enabled",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armcontainerregistry.Registry)
				pe := len(i.Properties.PrivateEndpointConnections) > 0
				return !pe, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/container-registry/container-registry-private-link",
		},
		"cr-005": {
			Id:             "cr-005",
			Category:       scanners.RulesCategoryHighAvailability,
			Recommendation: "ContainerRegistry SKU",
			Impact:         scanners.ImpactHigh,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armcontainerregistry.Registry)
				return false, string(*i.SKU.Name)
			},
			Url: "https://learn.microsoft.com/en-us/azure/container-registry/container-registry-skus",
		},
		"cr-006": {
			Id:             "cr-006",
			Category:       scanners.RulesCategoryGovernance,
			Recommendation: "ContainerRegistry Name should comply with naming conventions",
			Impact:         scanners.ImpactLow,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armcontainerregistry.Registry)
				caf := strings.HasPrefix(*c.Name, "cr")
				return !caf, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations",
		},
		"cr-007": {
			Id:             "cr-007",
			Category:       scanners.RulesCategorySecurity,
			Recommendation: "ContainerRegistry should have anonymous pull access disabled",
			Impact:         scanners.ImpactMedium,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armcontainerregistry.Registry)
				apull := c.Properties.AnonymousPullEnabled != nil && *c.Properties.AnonymousPullEnabled
				return apull, ""
			},
			Url: "https://learn.microsoft.com/azure/container-registry/anonymous-pull-access#configure-anonymous-pull-access",
		},
		"cr-008": {
			Id:             "cr-008",
			Category:       scanners.RulesCategorySecurity,
			Recommendation: "ContainerRegistry should have the Administrator account disabled",
			Impact:         scanners.ImpactMedium,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armcontainerregistry.Registry)
				admin := c.Properties.AdminUserEnabled != nil && *c.Properties.AdminUserEnabled
				return admin, ""
			},
			Url: "https://learn.microsoft.com/azure/container-registry/container-registry-authentication-managed-identity",
		},
		"cr-009": {
			Id:             "cr-009",
			Category:       scanners.RulesCategoryGovernance,
			Recommendation: "ContainerRegistry should have tags",
			Impact:         scanners.ImpactLow,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armcontainerregistry.Registry)
				return len(c.Tags) == 0, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/tag-resources?tabs=json",
		},
		"cr-010": {
			Id:             "cr-010",
			Category:       scanners.RulesCategoryGovernance,
			Recommendation: "ContainerRegistry should use retention policies",
			Impact:         scanners.ImpactMedium,
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armcontainerregistry.Registry)
				return c.Properties.Policies == nil ||
					c.Properties.Policies.RetentionPolicy == nil ||
					c.Properties.Policies.RetentionPolicy.Status == nil ||
					*c.Properties.Policies.RetentionPolicy.Status == armcontainerregistry.PolicyStatusDisabled, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/container-registry/container-registry-retention-policy",
		},
	}
}
