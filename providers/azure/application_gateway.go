// Copyright 2019 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azure

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/hashicorp/go-azure-helpers/authentication"
)

type ApplicationGatewayGenerator struct {
	AzureService
}

func (g ApplicationGatewayGenerator) createResources(ctx context.Context, iterator network.ApplicationGatewayListResultIterator) ([]terraformutils.Resource, error) {
	var resources []terraformutils.Resource
	for iterator.NotDone() {
		applicationGateways := iterator.Value()
		resources = append(resources, terraformutils.NewSimpleResource(
			*applicationGateways.ID,
			*applicationGateways.Name,
			"azurerm_application_gateway",
			g.ProviderName,
			[]string{}))
		if err := iterator.NextWithContext(ctx); err != nil {
			log.Println(err)
			return resources, err
		}
	}
	return resources, nil
}

func (g *ApplicationGatewayGenerator) InitResources() error {
	ctx := context.Background()
	subscriptionID := g.Args["config"].(authentication.Config).SubscriptionID
	resourceManagerEndpoint := g.Args["config"].(authentication.Config).CustomResourceManagerEndpoint
	applicationGatewaysClient := network.NewApplicationGatewaysClientWithBaseURI(resourceManagerEndpoint, subscriptionID)

	applicationGatewaysClient.Authorizer = g.Args["authorizer"].(autorest.Authorizer)

	var (
		output network.ApplicationGatewayListResultIterator
		err    error
	)

	if rg := g.Args["resource_group"].(string); rg != "" {
		output, err = applicationGatewaysClient.ListComplete(ctx, rg)
	} else {
		output, err = applicationGatewaysClient.ListAllComplete(ctx)
	}
	if err != nil {
		return err
	}
	g.Resources, err = g.createResources(ctx, output)
	return err
}
