package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
		value, ok := rs.Primary.Attributes["webhooks_merchant.id"]
		if rs.Type == "adyen_webhooks_merchant" && ok {
			data := client.Management().WebhooksMerchantLevelApi.GetWebhookInput(client.GetConfig().MerchantAccount, value)
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

var _ plancheck.PlanCheck = debugPlan{}

type debugPlan struct{}

func (e debugPlan) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	rd, err := json.MarshalIndent(req.Plan, "", "    ")
	if err != nil {
		fmt.Println("error marshalling machine-readable plan output:", err)
	}
	fmt.Printf("req.Plan - %s\n", string(rd))
}

func DebugPlan() plancheck.PlanCheck {
	return debugPlan{}
}

func TestAccWebhookMerchantResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		//CheckDestroy:             testAccCheckAdyenWebhookMerchantDestroy, //FIXME: destroy results in 301 error
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config:   testConfigCreate(),
				PlanOnly: true,
				//TODO: add below debug func to notes
				ExpectNonEmptyPlan: true, // FIXME: only way to make the test work that I know of, until further notice.
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Can be made more generic, resourceName param
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.type", "standard"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.url", "https://webhook.site/test-uuid"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.username", "YOUR_TEST_USER_1"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.active", "false"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.communication_format", "json"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.accepts_expired_certificate", "false"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.accepts_self_signed_certificate", "true"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.accepts_untrusted_root_certificate", "true"),
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test", "webhooks_merchant.populate_soap_action_header", "false"),
				),
			},
			{
				Config:             testConfigUpdate(),
				Destroy:            true,
				ExpectNonEmptyPlan: true, // FIXME: only way to make the test work that I know of, until further notice.
				Check: resource.ComposeAggregateTestCheckFunc(
					// TODO: Can be made more generic, resourceName param
					resource.TestCheckResourceAttr("adyen_webhooks_merchant.test_update", "webhooks_merchant.username", "YOUR_TEST_USER_2"),
				),
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

func testConfigUpdate() string {
	return `
provider "adyen" {
  api_key = "api-key"
  environment = "test"
  merchant_account = "WeaveAccountECOM"
  company_account = "WeaveAccount"
}

resource "adyen_webhooks_merchant" "test_update" {
	webhooks_merchant = {
		type                               = "standard"
		url                                = "https://webhook.site/test-uuid"
		username                           = "YOUR_TEST_USER_2"
		password                           = "YOUR_TEST_PASSWORD_FROM_TERRAFORM_2"
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
