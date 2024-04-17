package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccCheckAdyenWebhookMerchantDestroy(tfstate *terraform.State) error {
	//TODO: can be simplified
	suite := new(AcceptanceSuite)
	suite.SetupSuite()
	client := suite.client

	for _, rs := range tfstate.RootModule().Resources {
		if rs.Type == "adyen_webhooks_merchant" {
			data := client.Management().WebhooksMerchantLevelApi.GetWebhookInput(client.GetConfig().MerchantAccount, rs.Primary.ID)
			_, resp, err := client.Management().WebhooksMerchantLevelApi.GetWebhook(context.Background(), data)
			if resp.StatusCode == 404 {
				fmt.Printf("adyen_webhooks_merchant with id %s has been removed\n", rs.Primary.ID)
				continue
			}
			if err != nil {
				return err
			}

			return fmt.Errorf("received response code with status %d", resp.StatusCode)
		}
	}
	return nil
}

func TestAccWebhookMerchantResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAdyenWebhookMerchantDestroy,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Can be made more generic
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.type", "standard"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.url", "https://webhook.site/test-uuid"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.username", "YOUR_TEST_USER"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.active", "true"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.communication_format", "json"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.accepts_expired_certificate", "false"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.accepts_self_signed_certificate", "true"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.accepts_untrusted_root_certificate", "true"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.populate_soap_action_header", "false"),
				),
			},
			{
				ResourceName:      "adyen_webhooks_merchant.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TODO: make test more generic by adding resourceName param + generic provider + generic tf resource
func testConfigCreate() string {
	return `
provider "adyen" {
  api_key = "api-key"
  environment = "test"
  merchant_account = "WeaveAccountECOM"
  company_account = "WeaveAccount"
}

resource "adyen_webhooks_merchant" "test" {
	webhooks_merchant = {
		type                               = "standard"
		url                                = "https://webhook.site/test-uuid"
		username                           = "YOUR_TEST_USER"
		password                           = "YOUR_TEST_PASSWORD_FROM_TERRAFORM"
		active                             = false
		communication_format               = "json"
		accepts_expired_certificate        = false
		accepts_self_signed_certificate    = true
		accepts_untrusted_root_certificate = true
		populate_soap_action_header        = false
	}
}
`
}
