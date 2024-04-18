terraform {
  required_providers {
    adyen = {
      source = "hashicorp.com/edu/adyen"
    }
  }
}

provider "adyen" {
  api_key          = "API_KEY"
  environment      = "test" // "live"
  merchant_account = "WeaveAccountECOM"
  company_account  = "WeaveAccount"
}

resource "adyen_webhooks_merchant" "example_webhook" {
  webhooks_merchant = {
    type                               = "standard"
    url                                = "https://webhook.site/cb798fb3-7092-4cab-986b-f416fb04f92e"
    username                           = "YOUR_USER"
    password                           = "YOUR_PASSWORD_FROM_TERRAFORM"
    active                             = false
    communication_format               = "json"
    accepts_expired_certificate        = false
    accepts_self_signed_certificate    = true
    accepts_untrusted_root_certificate = true
    populate_soap_action_header        = false
  }
}
