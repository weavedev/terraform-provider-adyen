package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func testAccCheckAdyenWebhookMerchantDestroy(tfstate *terraform.State) error {
	suite := new(AcceptanceSuite)
	suite.SetupSuite()
	client := suite.client

	for _, rs := range tfstate.RootModule().Resources {
		value, ok := rs.Primary.Attributes["webhooks_merchant.id"]
		if rs.Type == "adyen_webhooks_merchant" && ok {
			data := client.Management().WebhooksMerchantLevelApi.GetWebhookInput(client.GetConfig().MerchantAccount, value)
			_, resp, err := client.Management().WebhooksMerchantLevelApi.GetWebhook(context.Background(), data)
			if resp.StatusCode == 422 { // 422 Unprocessable Entity error code from Adyen if resource does not exist.
				fmt.Printf("adyen_webhooks_merchant with id: '%s' does not exist and/or has been removed.\n", value)
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
	resourceName := "adyen_webhooks_merchant.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAdyenWebhookMerchantDestroy,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config:             testProviderClientFromTmpl(t) + testConfigCreateMerchantWebhook(),
				ExpectNonEmptyPlan: true, // Creating a tf resource will propose changes, that's why this value is set to 'true'. Can be approached differently by using `PlanOnly: true`.
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "webhooks_merchant.type", "standard"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_merchant.url", "https://webhook.site/test-uuid"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_merchant.username", "YOUR_TEST_USER_1"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_merchant.active", "false"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_merchant.communication_format", "json"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_merchant.accepts_expired_certificate", "false"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_merchant.accepts_self_signed_certificate", "true"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_merchant.accepts_untrusted_root_certificate", "true"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_merchant.populate_soap_action_header", "false"),
				),
			},
		},
	})
}

func testConfigCreateMerchantWebhook() string {
	return `
	resource "adyen_webhooks_merchant" "test" {
		webhooks_merchant = {
			type                               = "standard"
			url                                = "https://webhook.site/test-uuid"
			username                           = "YOUR_TEST_USER_1"
			password                           = "YOUR_TEST_PASSWORD_FROM_TERRAFORM_1"
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
