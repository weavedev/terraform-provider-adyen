terraform {
  required_providers {
    adyen = {
      source = "hashicorp.com/edu/adyen"
    }
  }
}

provider "adyen" {
}

data "adyen_webhook" "example" {}

