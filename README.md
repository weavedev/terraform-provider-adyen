# Terraform Provider Adyen

This repository is built on the template provided by: [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding).

This repository is a terraform provider for **Adyen**, containing:

- A resources and a data sources can be found in (`internal/provider/`),
- Examples can be found in (`examples/`) and generated documentation in (`docs/`),


## Currently supported resources
- [x] Webhook Merchant 
- [ ] Webhook Company
- [ ] More on the way...


## Usage 
#### Until the provider is public and in the terraform registry, we need to update our `.terraformrc` file with the following content:
```shell
provider_installation {

  dev_overrides {
      "hashicorp.com/edu/adyen" = "<GOBINPATH>"
  }

  direct {}
}
```

#### To fetch your `GOBINPATH` run:

```shell
go env GOBIN
```

#### If this command returns empty you can use the default `PATH` value instead: `Users/<yourusername>/go/bin`

#### Add the provider to your terraform project:
```hcl
terraform {
  required_providers {
    adyen = {
      source = "hashicorp.com/edu/adyen" //TODO: change to registry.terraform.io/weave/adyen
    }
  }
}

provider "adyen" {
  api_key          = "API_KEY"
  environment      = "test" // "live"
  merchant_account = "YOUR-ADYEN-MERCHANT-ACCOUNT"
  company_account  = "YOUR-ADYEN-COMPANY-ACCOUNT"
}
```
Development
===========
## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

## Using the provider

TODO

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation from root, run `go generate ./...`.

[//]: # (In order to run the full suite of Acceptance tests, run `make testacc`. )

[//]: # (*Note:* Acceptance tests create real resources, and often cost money to run.)
[//]: # ()
[//]: # (```shell)

[//]: # (make testacc)

[//]: # (```)
