// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/adyen/adyen-go-api-library/v9/src/adyen"
	"github.com/adyen/adyen-go-api-library/v9/src/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProviderWebhooksMerchant *schema.Provider

type AcceptanceSuite struct {
	suite.Suite
	client *adyen.APIClient
}

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"adyen": providerserver.NewProtocol6WithError(New("test")()),
}

func (s *AcceptanceSuite) SetupSuite() {
	conf := &common.Config{
		ApiKey:          os.Getenv("ADYEN_API_KEY"),
		Environment:     common.Environment(os.Getenv("ADYEN_API_ENVIRONMENT")),
		MerchantAccount: os.Getenv("ADYEN_API_MERCHANT_ACCOUNT"),
	}

	s.client = adyen.NewClient(conf)
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ADYEN_API_KEY"); v == "" {
		t.Fatal("ADYEN_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ADYEN_API_ENVIRONMENT"); v == "" {
		t.Fatal("ADYEN_API_ENVIRONMENT must be set for acceptance tests")
	}
	if v := os.Getenv("ADYEN_API_MERCHANT_ACCOUNT"); v == "" {
		t.Fatal("ADYEN_API_MERCHANT_ACCOUNT must be set for acceptance tests")
	}
}
