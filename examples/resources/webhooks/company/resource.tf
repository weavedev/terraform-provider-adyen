terraform {
  required_providers {
    adyen = {
      version = ">= 0.0.1"
      source  = "weavedev/adyen"
    }
  }
}

provider "adyen" {
  api_key          = "YOUR_API_KEY"
  environment      = "test"
  merchant_account = "WeaveAccountECOM"
}

resource "adyen_webhooks_company" "example_webhook" {
  company_account  = "WeaveAccount"
  webhooks_company = {
    type                               = "standard"
    url                                = "https://webhook.site/cb798fb3-7092-4cab-986b-f416fb04f92e"
    username                           = "provider_tf"
    active                             = true
    communication_format               = "http"
    accepts_expired_certificate        = false
    accepts_self_signed_certificate    = true
    accepts_untrusted_root_certificate = true
    populate_soap_action_header        = false
    filter_merchant_account_type       = "includeAccounts"
    filter_merchant_accounts           = ["WeaveAccountECOM"]
  }
}
