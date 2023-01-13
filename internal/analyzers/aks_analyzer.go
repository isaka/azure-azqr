package analyzers

import (
	"context"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
)

// AKSAnalyzer - Analyzer for AKS Clusters
type AKSAnalyzer struct {
	diagnosticsSettings DiagnosticsSettings
	subscriptionID      string
	ctx                 context.Context
	cred                azcore.TokenCredential
	clustersClient      *armcontainerservice.ManagedClustersClient
	listClustersFunc    func(resourceGroupName string) ([]*armcontainerservice.ManagedCluster, error)
}

// Init - Initializes the AKSAnalyzer
func (a *AKSAnalyzer) Init(config ServiceAnalizerConfig) error {
	a.subscriptionID = config.SubscriptionID
	a.ctx = config.Ctx
	a.cred = config.Cred
	var err error 
	a.clustersClient, err = armcontainerservice.NewManagedClustersClient(config.SubscriptionID, config.Cred, nil)
	if err != nil {
		return err
	}
	a.diagnosticsSettings = DiagnosticsSettings{}
	err = a.diagnosticsSettings.Init(config.Ctx, config.Cred)
	if err != nil {
		return err
	}
	return nil
}

// Review - Analyzes all AKS Clusters in a Resource Group
func (a *AKSAnalyzer) Review(resourceGroupName string) ([]IAzureServiceResult, error) {
	log.Printf("Analyzing AKS Clusters in Resource Group %s", resourceGroupName)

	clusters, err := a.listClusters(resourceGroupName)
	if err != nil {
		return nil, err
	}
	results := []IAzureServiceResult{}
	for _, c := range clusters {
		hasDiagnostics, err := a.diagnosticsSettings.HasDiagnostics(*c.ID)
		if err != nil {
			return nil, err
		}

		zones := true
		for _, profile := range c.Properties.AgentPoolProfiles {
			if profile.AvailabilityZones == nil || (profile.AvailabilityZones != nil && len(profile.AvailabilityZones) <= 1) {
				zones = false
			}
		}

		sku := string(*c.SKU.Tier)
		sla := "None"
		if sku == "Paid" {
			sla = "99.9%"
			if zones {
				sla = "99.95%"
			}
		}

		privateEndpoints := false
		if c.Properties.APIServerAccessProfile != nil && *c.Properties.APIServerAccessProfile.EnablePrivateCluster {
			privateEndpoints = true
		}

		results = append(results, AzureServiceResult{
			SubscriptionID:     a.subscriptionID,
			ResourceGroup:      resourceGroupName,
			ServiceName:        *c.Name,
			SKU:                sku,
			SLA:                sla,
			Type:               *c.Type,
			Location:           *c.Location,
			CAFNaming:          strings.HasPrefix(*c.Name, "aks"),
			AvailabilityZones:  zones,
			PrivateEndpoints:   privateEndpoints,
			DiagnosticSettings: hasDiagnostics,
		})
	}
	return results, nil
}

func (a *AKSAnalyzer) listClusters(resourceGroupName string) ([]*armcontainerservice.ManagedCluster, error) {
	if a.listClustersFunc == nil {
		pager := a.clustersClient.NewListByResourceGroupPager(resourceGroupName, nil)

		clusters := make([]*armcontainerservice.ManagedCluster, 0)
		for pager.More() {
			resp, err := pager.NextPage(a.ctx)
			if err != nil {
				return nil, err
			}
			clusters = append(clusters, resp.Value...)
		}
		return clusters, nil
	}

	return a.listClustersFunc(resourceGroupName)
}
