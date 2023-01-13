package analyzers

import (
	"context"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/signalr/armsignalr"
)

// SignalRAnalyzer - Analyzer for SignalR
type SignalRAnalyzer struct {
	diagnosticsSettings DiagnosticsSettings
	subscriptionID      string
	ctx                 context.Context
	cred                azcore.TokenCredential
	signalrClient       *armsignalr.Client
	listSignalRFunc     func(resourceGroupName string) ([]*armsignalr.ResourceInfo, error)
}

// Init - Initializes the SignalRAnalyzer
func (c *SignalRAnalyzer) Init(config ServiceAnalizerConfig) error {
	c.subscriptionID = config.SubscriptionID
	c.ctx = config.Ctx
	c.cred = config.Cred
	var err error
	c.signalrClient, err = armsignalr.NewClient(config.SubscriptionID, config.Cred, nil)
	if err != nil {
		return err
	}
	c.diagnosticsSettings = DiagnosticsSettings{}
	err = c.diagnosticsSettings.Init(config.Ctx, config.Cred)
	if err != nil {
		return err
	}
	return nil
}

// Review - Analyzes all SignalR in a Resource Group
func (c *SignalRAnalyzer) Review(resourceGroupName string) ([]IAzureServiceResult, error) {
	log.Printf("Analyzing SignalR in Resource Group %s", resourceGroupName)

	signalr, err := c.listSignalR(resourceGroupName)
	if err != nil {
		return nil, err
	}
	results := []IAzureServiceResult{}
	for _, signalr := range signalr {
		hasDiagnostics, err := c.diagnosticsSettings.HasDiagnostics(*signalr.ID)
		if err != nil {
			return nil, err
		}

		sku := string(*signalr.SKU.Name)
		zones := false
		if strings.Contains(sku, "Premium") {
			zones = true
		}

		results = append(results, AzureServiceResult{
			SubscriptionID:     c.subscriptionID,
			ResourceGroup:      resourceGroupName,
			ServiceName:        *signalr.Name,
			SKU:                sku,
			SLA:                "99.9%",
			Type:               *signalr.Type,
			Location:           *signalr.Location,
			CAFNaming:          strings.HasPrefix(*signalr.Name, "sigr"),
			AvailabilityZones:  zones,
			PrivateEndpoints:   len(signalr.Properties.PrivateEndpointConnections) > 0,
			DiagnosticSettings: hasDiagnostics,
		})
	}
	return results, nil
}

func (c *SignalRAnalyzer) listSignalR(resourceGroupName string) ([]*armsignalr.ResourceInfo, error) {
	if c.listSignalRFunc == nil {
		pager := c.signalrClient.NewListByResourceGroupPager(resourceGroupName, nil)

		signalrs := make([]*armsignalr.ResourceInfo, 0)
		for pager.More() {
			resp, err := pager.NextPage(c.ctx)
			if err != nil {
				return nil, err
			}
			signalrs = append(signalrs, resp.Value...)
		}
		return signalrs, nil
	}

	return c.listSignalRFunc(resourceGroupName)
}
