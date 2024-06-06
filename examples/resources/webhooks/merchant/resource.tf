terraform {
  required_providers {
    adyen = {
      version = ">= 0.0.1"
      source  = "weavedev/adyen"
    }
  }
}

provider "adyen" {
  api_key          = var.adyen_api_key
  environment      = var.environment
  merchant_account = var.merchant_account
  company_account  = var.company_account //TODO: remove? keep?
}

resource "adyen_webhooks_merchant" "example_webhook" {
  webhooks_merchant = {
    type                               = "standard"
    url                                = "https://webhook.site/cb798fb3-7092-4cab-986b-f416fb04f92e"
    username                           = "YOUR_USER"
    password                           = "YOUR_PASSWORD"
    active                             = false
    communication_format               = "json"
    accepts_expired_certificate        = false
    accepts_self_signed_certificate    = true
    accepts_untrusted_root_certificate = true
    populate_soap_action_header        = false
  }
}
