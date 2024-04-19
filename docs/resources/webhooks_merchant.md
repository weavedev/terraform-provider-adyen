---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "adyen_webhooks_merchant Resource - adyen"
subcategory: ""
description: |-
  
---

# adyen_webhooks_merchant (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `webhooks_merchant` (Attributes) Manages a webhook on merchant level. (see [below for nested schema](#nestedatt--webhooks_merchant))

<a id="nestedatt--webhooks_merchant"></a>
### Nested Schema for `webhooks_merchant`

Required:

- `accepts_expired_certificate` (Boolean) Indicates if expired certificates are accepted.
- `accepts_self_signed_certificate` (Boolean) Indicates if self-signed certificates are accepted.
- `accepts_untrusted_root_certificate` (Boolean) Indicates if untrusted root certificates are accepted.
- `active` (Boolean) Indicates if the webhook is active.
- `communication_format` (String) The format of the communication (e.g., 'json').
- `password` (String, Sensitive) The password required for basic authentication.
- `type` (String) The type of the webhook.
- `url` (String) The URL the webhook will send requests to.
- `username` (String) The username required for basic authentication.

Optional:

- `certificate_alias` (String) The alias of the certificate.
- `encryption_protocol` (String) The encryption protocol used by the webhook.
- `populate_soap_action_header` (Boolean) Indicates if the SOAP action header should be populated.

Read-Only:

- `description` (String) A description of the webhook.
- `has_error` (Boolean) Indicates if there is an error with the webhook.
- `has_password` (Boolean) Indicates if the webhook is configured with a password.
- `id` (String) The unique identifier for the webhook.