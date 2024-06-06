package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"testing"
)

func testAccCheckAdyenWebhookCompanyDestroy(tfstate *terraform.State) error {
	suite := new(AcceptanceSuite)
	suite.SetupSuite()
	client := suite.client

	companyAccount := os.Getenv("ADYEN_API_COMPANY_ACCOUNT") //TODO: nil check
	if companyAccount == "" {
		return fmt.Errorf("received empty company account")
	}

	for _, rs := range tfstate.RootModule().Resources {
		value, ok := rs.Primary.Attributes["webhooks_company.id"]
		if rs.Type == "adyen_webhooks_company" && ok {
			data := client.Management().WebhooksCompanyLevelApi.GetWebhookInput(companyAccount, value)
			_, resp, err := client.Management().WebhooksCompanyLevelApi.GetWebhook(context.Background(), data)
			if resp.StatusCode == 422 { // 422 Unprocessable Entity error code from Adyen if resource does not exist.
				fmt.Printf("adyen_webhooks_company with id: '%s' does not exist and/or has been removed.\n", value)
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

func TestAccWebhookCompanyResource(t *testing.T) {
	resourceName := "adyen_webhooks_company.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAdyenWebhookCompanyDestroy,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config:             testProviderClientFromTmpl(t) + testConfigCreateCompanyWebhook(),
				ExpectNonEmptyPlan: true, // Creating a tf resource will propose changes, that's why this value is set to 'true'. Can be approached differently by using `PlanOnly: true`.
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "webhooks_company.type", "standard"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_company.url", "https://webhook.site/cb798fb3-7092-4cab-986b-f416fb04f92e"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_company.username", "provider_tf"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_company.active", "true"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_company.communication_format", "http"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_company.accepts_expired_certificate", "false"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_company.accepts_self_signed_certificate", "true"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_company.accepts_untrusted_root_certificate", "true"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_company.populate_soap_action_header", "false"),
					resource.TestCheckResourceAttr(resourceName, "webhooks_company.filter_merchant_account_type", "includeAccounts"),
					resource.TestCheckResourceAttr(
						resourceName, "webhooks_company.filter_merchant_accounts.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						resourceName, "webhooks_company.filter_merchant_accounts.*", "WeaveAccountECOM"),
				),
			},
		},
	})
}

func testConfigCreateCompanyWebhook() string {
	return `
	resource "adyen_webhooks_company" "test" {
		webhooks_company = {
			type                               = "standard"
			password 						   = "secretpassword"
			url                                = "https://webhook.site/cb798fb3-7092-4cab-986b-f416fb04f92e"
			username                           = "provider_tf"
			active                             = true
			communication_format               = "http"
			accepts_expired_certificate        = false
			accepts_self_signed_certificate    = true
			accepts_untrusted_root_certificate = true
			populate_soap_action_header        = false
			filter_merchant_account_type       = "includeAccounts"
  			filter_merchant_accounts  		   = ["WeaveAccountECOM"]
		}
	}
`
}
