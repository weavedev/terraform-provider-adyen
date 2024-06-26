# Terraform Provider Adyen

This repository is built on the template provided by: [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). The template repository built on the [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk) can be found at [terraform-provider-scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding).

This repository is a terraform provider for **Adyen**, containing:

- A resources and a data sources can be found in (`internal/provider/`),
- Examples can be found in (`examples/`) and generated documentation in (`docs/`),
- To generate or update documentation from root, run `go generate ./...`.

## Adyen Resources Implementation Roadmap
### Management API
- Webhooks:
  - [x] Webhook Merchant
  - [x] Webhook Company
#### 
- Users
   - [ ] Users Merchant
   - [ ] Users Company
####
- API Credentials
   - [ ] My API Credentials
   - [ ] Company API Credentials
   - [ ] Merchant API Credentials
####
   - Terminal
     - [ ] Actions
     - [ ] Settings
     - [ ] Orders
####
- Account:
   - [ ] Account Merchant
   - [ ] Account Company
   - [ ] Account Store
####
   - [ ] Payment Methods
   - [ ] Payout Settings
   - [ ] Allowed Origins


## Provider Setup and Usage

### Clone the repository

First, clone the `terraform-provider-adyen` repository:

1. Open your terminal or command prompt.
2. Run the following command to clone the repository:

    ```sh
    git clone https://github.com/weavedev/terraform-provider-adyen.git
    ```

3. Navigate to the cloned directory:

    ```sh
    cd terraform-provider-adyen
    ```
   
### Setup Credentials

#### Adyen Test Customer Area
1. Go to [Adyen](https://docs.adyen.com/get-started-with-adyen/) and follow the instructions there to create a "test account" so you can get granted access to the "Test Customer Area".
2. Go to your [Test Customer Area](https://ca-test.adyen.com/) and login with your credentials.
3. Go to "**Developers**" -> "**API credentials**" -> "**Create new credential**".
4. Create API Credential as a "**Web service user**".
5. (Optional) Add a description to your API Credential.
6. Note/Copy your API key under "**Authentication**" and optionally edit scopes under "**Permissions**" --> "**Roles**".
7. Note your **"Merchant" & "Company" accounts** at the top left of your Dashboard.

#### Add the provider to your terraform project:
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
  api_key     = "<api_key>"                       // From Step 6
  environment = "<environment>"                   // Or "live"
  merchant_account  = "<merchant_account>"        // From Step 7
  company_account = "<company_account>"           // From Step 7
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

## Developing the Provider

### Prepare Terraform for Local Provider Install

1. Find `GOBIN` path. Your path may vary depending on how your Go environment variables are configured:
    ```bash
    go env GOBIN
    ```

   If your `GOBIN` go environment variable is not set use the default path: `/Users/<Username>/go/bin`

2. Create a `.terraformrc` file in your home directory (`~`), **if necessary**, then add the `dev_overrides` block below. Change the `<GOBIN_PATH>` to the value returned from the previous command: `go env GOBIN`. **Note: This override ensures that while developing you are using your locally compiled provider**.
    ```ini
    provider_installation {
      dev_overrides {
        "registry.terraform.io/weavedev/adyen" = "<GOBIN_PATH>"
      }
      direct {}
    }
    ```
3. You should be getting this "Warning" from Terraform after adding your `dev_overrides` to the `.terraformrc` file: 
   ```
      Warning: Provider development overrides are in effect
   
      The following provider development overrides are set in the CLI configuration:
        - weavedev/adyen in /Users/tolgaakyazi/go/bin
      
      The behavior may therefore not match any released version of the provider and applying changes may cause the state to become
      incompatible with published releases.
   ```
4. Now you are ready to locally develop and test new data-sources, resources or functions for this provider!  

## Contributing to the Provider 

After setting up your local provider install, you can start contributing to the provider.
1. You should, at this point, have a folder called `terraform-provider-adyen`, `cd` into that folder and create a new branch
    ```bash
    cd terraform-provider-adyen
    git checkout -b feat/<your-feature-name>
    ```

2. Make the necessary changes to the codebase. Be sure to follow the project's coding standards and best practices.
3. Add your changes to a new commit
   ```bash
    git add .
   ```
4. Commit the staged changes with a descriptive message
    ```bash
    git commit -m 'Clear description of your changes/additions'
    ```
5. Push your changes to the repository 
    ```bash
    git push origin feat/<your-feature-name>
    ```
6. Finally, navigate to the repository on GitHub and open a pull request from your new branch. Provide a clear description of your changes and any relevant context.
    

## PR Reviewing Process
- Your pull request will be reviewed by one of the maintainers.
- You might be asked to make additional changes.
- Once the changes are approved, the build & tests pass, your pull request will be merged.

## Additional Tips
- Keep your changes focused and concise.
- Ensure your code follows the project's style guidelines.
- Generate documentation as needed (`go generate ./...`).
- Write acceptance tests for any new functionality (see: `webhook_merchant_resource_test.go`). Without tests, your PR will not be approved. 

# Creating a Provider Release

To create a release for your provider, follow these steps. The GitHub Action will trigger and create a release whenever a new valid version tag is pushed to the repository. Ensure your Terraform provider versions follow the Semantic Versioning standard (vMAJOR.MINOR.PATCH).

1. **Stage your changes**  
   Add your changes to a new commit:
   ```bash
   git add .
   ```
2. **Commit your changes**
   Commit the staged changes with a descriptive message:
   ```bash
    git commit -m 'Add docs, goreleaser, and GH actions'
   ```
3. **Create a new tag**
   Create a new tag following the Semantic Versioning standard (e.g. `v0.0.1`):
      ```bash
       git tag v0.0.1
   ```
4. **Push the tag to GitHub**
   Push the new tag to GitHub to trigger the release action:
    ```bash
   git push origin v0.0.1
    ```
   Once the tag is pushed, GitHub Actions will automatically trigger the workflow to create a new release for your provider.

#### Thank you for contributing! We appreciate your help in improving our provider.



