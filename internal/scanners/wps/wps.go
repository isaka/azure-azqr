// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package wps

import (
	"github.com/Azure/azqr/internal/azqr"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/webpubsub/armwebpubsub"
)

// WebPubSubScanner - Scanner for WebPubSub
type WebPubSubScanner struct {
	config *azqr.ScannerConfig
	client *armwebpubsub.Client
}

// Init - Initializes the WebPubSubScanner
func (c *WebPubSubScanner) Init(config *azqr.ScannerConfig) error {
	c.config = config
	var err error
	c.client, err = armwebpubsub.NewClient(config.SubscriptionID, config.Cred, config.ClientOptions)
	return err
}

// Scan - Scans all WebPubSub in a Resource Group
func (c *WebPubSubScanner) Scan(scanContext *azqr.ScanContext) ([]azqr.AzqrServiceResult, error) {
	azqr.LogSubscriptionScan(c.config.SubscriptionID, c.ResourceTypes()[0])

	WebPubSub, err := c.listWebPubSub()
	if err != nil {
		return nil, err
	}
	engine := azqr.RecommendationEngine{}
	rules := c.GetRecommendations()
	results := []azqr.AzqrServiceResult{}

	for _, w := range WebPubSub {
		rr := engine.EvaluateRecommendations(rules, w, scanContext)

		results = append(results, azqr.AzqrServiceResult{
			SubscriptionID:   c.config.SubscriptionID,
			SubscriptionName: c.config.SubscriptionName,
			ResourceGroup:    azqr.GetResourceGroupFromResourceID(*w.ID),
			ServiceName:      *w.Name,
			Type:             *w.Type,
			Location:         *w.Location,
			Recommendations:  rr,
		})
	}
	return results, nil
}

func (c *WebPubSubScanner) listWebPubSub() ([]*armwebpubsub.ResourceInfo, error) {
	pager := c.client.NewListBySubscriptionPager(nil)

	WebPubSubs := make([]*armwebpubsub.ResourceInfo, 0)
	for pager.More() {
		resp, err := pager.NextPage(c.config.Ctx)
		if err != nil {
			return nil, err
		}
		WebPubSubs = append(WebPubSubs, resp.Value...)
	}
	return WebPubSubs, nil
}

func (c *WebPubSubScanner) ResourceTypes() []string {
	return []string{"Microsoft.SignalRService/webPubSub"}
}
