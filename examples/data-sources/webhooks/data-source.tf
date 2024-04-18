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

data "adyen_webhooks_merchant" "example" {}

output "adyen_webhooks_merchant" {
  value = data.adyen_webhooks_merchant.example
}

