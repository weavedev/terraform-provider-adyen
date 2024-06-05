// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/adyen/adyen-go-api-library/v9/src/adyen"
	"github.com/adyen/adyen-go-api-library/v9/src/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
)

// Ensure adyenProvider satisfies various provider interfaces.
var _ provider.Provider = &adyenProvider{}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &adyenProvider{
			version: version,
		}
	}
}

// adyenProvider defines the provider implementation.
type adyenProvider struct {
	version string
}

// adyenProviderModel describes the provider data model.
type adyenProviderModel struct {
	ApiKey          types.String `tfsdk:"api_key"`
	Environment     types.String `tfsdk:"environment"`
	MerchantAccount types.String `tfsdk:"merchant_account"`
	CompanyAccount  types.String `tfsdk:"company_account"` //TODO: figure out how to use this globally
}

// Metadata returns the provider type name.
func (p *adyenProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "adyen"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *adyenProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The API Key for the Adyen API Client.",
			},
			"environment": schema.StringAttribute{
				MarkdownDescription: "The Development Environment for the Adyen API Client. Can be either 'live' or 'test'.",
				Required:            true,
			},
			"merchant_account": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The Merchant Account ID for the Adyen API Client.",
			},
			"company_account": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The Company Account ID for the Adyen API Client.",
			},
		},
	}
}

// Configure prepares an Adyen API client for data sources and resources.
func (p *adyenProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Adyen API client")

	var config adyenProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Adyen API Key",
			"The provider cannot create the Adyen API client as there is an unknown configuration value for the Adyen API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ADYEN_API_KEY environment variable.",
		)
	}

	if config.Environment.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("environment"),
			"Unknown Adyen API Environment",
			"The provider cannot create the Adyen API client as there is an unknown configuration value for the Adyen API Environment. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ADYEN_API_ENVIRONMENT environment variable.",
		)
	}

	if config.MerchantAccount.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("merchant_account"),
			"Unknown Adyen API Merchant Account",
			"The provider cannot create the Adyen API client as there is an unknown configuration value for the Adyen API Merchant Account. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ADYEN_API_MERCHANT_ACCOUNT environment variable.",
		)
	}

	if config.CompanyAccount.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("company_account"),
			"Unknown Adyen Company Account",
			"The provider cannot create the Adyen API client as there is an unknown configuration value for the Adyen API Company Account. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ADYEN_API_COMPANY_ACCOUNT environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("ADYEN_API_KEY")
	environment := os.Getenv("ADYEN_API_ENVIRONMENT")
	merchantAccount := os.Getenv("ADYEN_API_MERCHANT_ACCOUNT")
	companyAccount := os.Getenv("ADYEN_API_COMPANY_ACCOUNT")

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	if !config.Environment.IsNull() {
		environment = config.Environment.ValueString()
	}

	if !config.MerchantAccount.IsNull() {
		merchantAccount = config.MerchantAccount.ValueString()
	}

	if !config.CompanyAccount.IsNull() {
		companyAccount = config.CompanyAccount.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Missing Adyen API Key",
			"The provider cannot create the Adyen API client as there is a missing or empty value for the Adyen API Key. "+
				"Set the host value in the configuration or use the ADYEN_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if environment == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("environment"),
			"Missing Adyen API Environment",
			"The provider cannot create the Adyen API client as there is a missing or empty value for the Adyen API Environment. "+
				"Set the host value in the configuration or use the ADYEN_API_ENVIRONMENT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if merchantAccount == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("merchant_account"),
			"Missing Adyen API Merchant Account",
			"The provider cannot create the Adyen API client as there is a missing or empty value for the Adyen API Merchant Account. "+
				"Set the host value in the configuration or use the ADYEN_API_MERCHANT_ACCOUNT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if companyAccount == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("company_account"),
			"Missing Adyen API Company Account",
			"The provider cannot create the Adyen API client as there is a missing or empty value for the Adyen API Company Account. "+
				"Set the host value in the configuration or use the ADYEN_API_COMPANY_ACCOUNT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "adyen_apikey", apiKey)
	ctx = tflog.SetField(ctx, "adyen_environment", environment)
	ctx = tflog.SetField(ctx, "adyen_merchant_account", merchantAccount)
	ctx = tflog.SetField(ctx, "adyen_company_account", companyAccount)

	// Add a filter to mask the apikey since it is sensitive information about the environment.
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "adyen_apikey", "adyen_merchant_account", "adyen_company_account")

	client := adyen.NewClient(&common.Config{
		ApiKey:          apiKey,
		Environment:     common.Environment(environment),
		MerchantAccount: merchantAccount,
	})

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Adyen API client", map[string]any{"success": true})
}

// Resources defines the resources implemented in the provider.
func (p *adyenProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWebhooksMerchantResource,
		NewWebhooksCompanyResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *adyenProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
