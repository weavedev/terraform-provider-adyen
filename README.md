# Terraform Provider Adyen

This repository is built on the template provided by: [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding).

This repository is a terraform provider for **Adyen**, containing:

- A resources and a data sources can be found in (`internal/provider/`),
- Examples can be found in (`examples/`) and generated documentation in (`docs/`),

## Currently supported resources
- [x] Webhook Merchant
- [ ] Webhook Company
- [ ] ???
- [ ] More on the way...


## Usage

### Adyen Test Customer Area
1. Go to [Adyen](https://docs.adyen.com/get-started-with-adyen/) and follow the instructions there to create a "test account" so you can get granted access to the "Test Customer Area".
2. Go to your [Test Customer Area](https://ca-test.adyen.com/) and login with your credentials.
3. Go to "**Developers**" -> "**API credentials**" -> "**Create new credential**".
4. Create API Credential as a "**Web service user**".
5. (Optional) Add a description to your API Credential.
6. Note/Copy your API key under "**Authentication**" and optionally edit scopes under "**Permissions**" --> "**Roles**".
7. Note your **"Merchant" & "Company" accounts** at the top left of your Dashboard.

### Add the provider to your terraform project:
```hcl
terraform {
  required_providers {
    adyen = {
      version = ">= 0.0.1"
      source = "weavedev/adyen"
    }
  }
}

provider "adyen" {
  api_key          = "API_KEY"                     // From Step 6
  environment      = "test"                        // Or "live"
  merchant_account = "YOUR-ADYEN-MERCHANT-ACCOUNT" // From Step 7
  company_account  = "YOUR-ADYEN-COMPANY-ACCOUNT"  // From Step 7
}

# Example resource
resource "adyen_webhooks_merchant" "example_webhook" {
  webhooks_merchant = {
    type                               = "standard"
    url                                = "https://webhook.site/etc-etc-etc"
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
```
Development
===========
## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.21

## Prepare Terraform for Local Provider Install

1. Find `GOBIN` path. Your path may vary depending on how your Go environment variables are configured:
    ```bash
    go env GOBIN
    ```

   If your GOBIN go environment variable is not set use the default path: `/Users/<Username>/go/bin`


2. Create a `.terraformrc` file in your home directory (`~`), then add the `dev_overrides` block below. Change the `<GOBIN_PATH>` to the value returned from the previous command: `go env GOBIN`
    ```ini
    provider_installation {
      dev_overrides {
        "registry.terraform.io/weavedev/adyen" = "<GOBIN_PATH>"
      }
      direct {}
    }
    ```

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install .
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.


## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install .`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation from root, run `go generate ./...`.

