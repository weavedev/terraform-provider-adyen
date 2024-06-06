variable "adyen_api_key" {
  description = "API key for the Adyen provider"
  type        = string
}

variable "environment" {
  description = "Environment of your Adyen Dashboard. Either 'live' or 'test' can be used."
  type        = string
}

variable "merchant_account" {
  description = "Merchant Account from your Adyen Dashboard"
  type        = string
}
variable "company_account" {
  description =  "Company Account from your Adyen Dashboard"
  type        = string
}
