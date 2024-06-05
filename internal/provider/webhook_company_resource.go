package provider

import (
	"context"
	"fmt"
	"github.com/adyen/adyen-go-api-library/v9/src/adyen"
	"github.com/adyen/adyen-go-api-library/v9/src/management"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &webhookCompanyResource{}
	_ resource.ResourceWithConfigure = &webhookCompanyResource{}
)

// webhookResource is the resource implementation.
type webhookCompanyResource struct {
	client         *adyen.APIClient
	companyAccount string
}

// NewWebhooksCompanyResource is a helper function to simplify the provider implementation.
func NewWebhooksCompanyResource(companyAccount string) resource.Resource {
	return &webhookCompanyResource{companyAccount: companyAccount}
}

// webhooksCompanyResourceModel maps the "webhooks_company" schema data for a resource.
type webhooksCompanyResourceModel struct {
	WebhooksCompany webhooksCompanyModel `tfsdk:"webhooks_company"`
}

type webhooksCompanyModel struct {
	ID                              types.String `tfsdk:"id"`
	Type                            types.String `tfsdk:"type"`
	URL                             types.String `tfsdk:"url"`
	Username                        types.String `tfsdk:"username"`
	Description                     types.String `tfsdk:"description"`
	HasPassword                     types.Bool   `tfsdk:"has_password"`
	Password                        types.String `tfsdk:"password"`
	Active                          types.Bool   `tfsdk:"active"`
	HasError                        types.Bool   `tfsdk:"has_error"`
	EncryptionProtocol              types.String `tfsdk:"encryption_protocol"`
	CommunicationFormat             types.String `tfsdk:"communication_format"`
	AcceptsExpiredCertificate       types.Bool   `tfsdk:"accepts_expired_certificate"`
	AcceptsSelfSignedCertificate    types.Bool   `tfsdk:"accepts_self_signed_certificate"`
	AcceptsUntrustedRootCertificate types.Bool   `tfsdk:"accepts_untrusted_root_certificate"`
	CertificateAlias                types.String `tfsdk:"certificate_alias"`
	PopulateSoapActionHeader        types.Bool   `tfsdk:"populate_soap_action_header"`
	Links                           types.Object `tfsdk:"links"`
	AdditionalSettings              types.Object `tfsdk:"additional_settings"`
	FilterMerchantAccountType       types.String `tfsdk:"filter_merchant_account_type"`
	FilterMerchantAccounts          types.List   `tfsdk:"filter_merchant_accounts"`
}

// Configure adds the provider configured client to the resource.
func (r *webhookCompanyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*adyen.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *adyen.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *webhookCompanyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhooks_company"
}

// Schema defines the schema for the resource.
func (r *webhookCompanyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"webhooks_company": schema.SingleNestedAttribute{
				Description: "Subscribe to receive webhook notifications about events related to your company account.\n\n" +
					"You can add basic authentication to make sure the data is secure.\n\n" +
					"To make this request, your API credential must have the following roles:\n\nManagement API—Webhooks read and write",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "Unique identifier for this webhook.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(), // Required to use this when it is known that an unconfigured value will remain the same after a resource update.
						},
					},
					"type": schema.StringAttribute{
						Required: true,
						Description: "The type of webhook that is being created. Possible values are:\n\nstandard\naccount-settings-notification\n" +
							"banktransfer-notification\nboletobancario-notification\ndirectdebit-notification\nach-notification-of-change-notification\n" +
							"pending-notification\nideal-notification\nideal-pending-notification\nreport-notification\nrreq-notification\n" +
							"Find out more about standard notification webhooks and other types of notifications.",
					},
					"url": schema.StringAttribute{
						Required:    true,
						Description: "Public URL where webhooks will be sent, for example https://www.domain.com/webhook-endpoint.",
					},
					"username": schema.StringAttribute{
						Required:    true,
						Description: "Username to access the webhook URL.",
					},
					"password": schema.StringAttribute{
						Required:    true,
						Sensitive:   true,
						Description: "The password required for basic authentication.",
					},
					"has_password": schema.BoolAttribute{
						Computed:    true,
						Description: "Indicates if the webhook is password protected.",
					},
					"active": schema.BoolAttribute{
						Required:    true,
						Description: "Indicates if the webhook configuration is active. The field must be 'true' for Adyen to send webhooks about events related an account.",
					},
					"communication_format": schema.StringAttribute{
						Required:    true,
						Description: "Format or protocol for receiving webhooks. Possible values:\n\nsoap\nhttp\njson",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "Your description for this webhook configuration.",
					},
					"encryption_protocol": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Description: "SSL version to access the public webhook URL specified in the url field. " +
							"Possible values:\n\nTLSv1.3\nTLSv1.2\n & HTTP. HTTP is Only allowed on Test environment.\n" +
							"If not specified, the webhook will use sslVersion: TLSv1.2.",
					},
					"has_error": schema.BoolAttribute{
						Computed:    true,
						Description: "Indicates if the webhook configuration has errors that need troubleshooting. If the value is true, troubleshoot the configuration using the testing endpoint.",
					},
					"certificate_alias": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The alias of Adyen SSL certificate. When you receive a notification from Adyen, the alias from the HMAC signature will match this alias.",
					},
					"populate_soap_action_header": schema.BoolAttribute{
						Optional:    true,
						Description: "Indicates if the SOAP action header needs to be populated. Default value: false. Only applies if communicationFormat: soap.",
					},
					"accepts_expired_certificate": schema.BoolAttribute{
						Required:    true,
						Description: "Indicates if expired SSL certificates are accepted. Default value: false.",
					},
					"accepts_self_signed_certificate": schema.BoolAttribute{
						Required:    true,
						Description: "Indicates if self-signed SSL certificates are accepted. Default value: false.",
					},
					"accepts_untrusted_root_certificate": schema.BoolAttribute{
						Required:    true,
						Description: "Indicates if untrusted SSL certificates are accepted. Default value: false.",
					},
					"links": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"self": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"href": schema.StringAttribute{Computed: true},
								},
								Computed:    true,
								Description: "The API URL to the webhook itself.",
							},
							"generate_hmac": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"href": schema.StringAttribute{Computed: true},
								},
								Computed:    true,
								Description: "The API URL to generate an HMAC key for the webhook.",
							},
							"company": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"href": schema.StringAttribute{Computed: true},
								},
								Computed:    true,
								Description: "The API URL to the company account associated with the webhook.",
							},
							"test_webhook": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"href": schema.StringAttribute{Computed: true},
								},
								Computed:    true,
								Description: "The API URL to test the webhook.",
							},
						},
					},
					"additional_settings": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Additional shopper and transaction information to be included in your standard notifications.",
						Attributes: map[string]schema.Attribute{
							"properties": schema.MapAttribute{
								Computed:    true,
								ElementType: types.BoolType,
								Description: "Object containing boolean key-value pairs. " +
									"The key can be any standard webhook additional setting, and the value indicates if the setting is enabled. " +
									"For example, captureDelayHours: true means the standard notifications you get will contain the " +
									"number of hours remaining until the payment will be captured.",
							},
						},
					},
					"filter_merchant_account_type": schema.StringAttribute{
						Optional: true,
						Description: "Shows how merchant accounts are filtered when configuring the webhook.\n\n" +
							"Possible values:\n\nallAccounts : Includes all merchant accounts, and does not require specifying " +
							"filterMerchantAccounts.\nincludeAccounts : The webhook is configured for the merchant accounts listed in filterMerchantAccounts.\n" +
							"excludeAccounts : The webhook is not configured for the merchant accounts listed in filterMerchantAccounts.",
					},
					"filter_merchant_accounts": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Description: "A list of merchant account names that are included or excluded from receiving the webhook. " +
							"Inclusion or exclusion is based on the value defined for filterMerchantAccountType.\n\n" +
							"Required if filterMerchantAccountType is either:\n\nincludeAccounts\nexcludeAccounts\n" +
							"Not needed for filterMerchantAccountType: allAccounts.",
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *webhookCompanyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating adyen company webhook")

	// Retrieve values from the plan
	var plan webhooksCompanyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	createCompanyWebhookRequest := &management.CreateCompanyWebhookRequest{
		AcceptsExpiredCertificate:       plan.WebhooksCompany.AcceptsExpiredCertificate.ValueBoolPointer(),
		AcceptsSelfSignedCertificate:    plan.WebhooksCompany.AcceptsSelfSignedCertificate.ValueBoolPointer(),
		AcceptsUntrustedRootCertificate: plan.WebhooksCompany.AcceptsUntrustedRootCertificate.ValueBoolPointer(),
		Active:                          plan.WebhooksCompany.Active.ValueBool(),
		CommunicationFormat:             plan.WebhooksCompany.CommunicationFormat.ValueString(),
		Password:                        plan.WebhooksCompany.Password.ValueStringPointer(),
		PopulateSoapActionHeader:        plan.WebhooksCompany.PopulateSoapActionHeader.ValueBoolPointer(),
		Type:                            plan.WebhooksCompany.Type.ValueString(),
		Url:                             plan.WebhooksCompany.URL.ValueString(),
		Username:                        plan.WebhooksCompany.Username.ValueStringPointer(),
		FilterMerchantAccountType:       plan.WebhooksCompany.FilterMerchantAccountType.ValueString(),
	}

	var filterMerchantAccounts []attr.Value
	if len(createCompanyWebhookRequest.FilterMerchantAccounts) > 0 {
		filterMerchantAccounts = mapWebhooksCompanyFilterMerchantAccounts(createCompanyWebhookRequest.FilterMerchantAccounts)
	}

	// Create a new webhook
	webhookCompanyCreateRequest := r.client.
		Management().
		WebhooksCompanyLevelApi.
		SetUpWebhookInput(r.companyAccount).
		CreateCompanyWebhookRequest(*createCompanyWebhookRequest)
	webhookCompanyCreateResponse, _, err := r.client.
		Management().
		WebhooksCompanyLevelApi.
		SetUpWebhook(ctx, webhookCompanyCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating company webhook",
			"Could not create company webhook, unexpected error: "+err.Error(),
		)
		return
	}

	var properties map[string]attr.Value
	if webhookCompanyCreateResponse.AdditionalSettings.Properties != nil {
		properties = mapWebhooksAdditionalSettingsProperties(*webhookCompanyCreateResponse.AdditionalSettings.Properties)
	}

	// Map response body to schema and populate with attribute values
	plan.WebhooksCompany = webhooksCompanyModel{
		ID:                              types.StringPointerValue(webhookCompanyCreateResponse.Id),
		Description:                     types.StringPointerValue(webhookCompanyCreateResponse.Description),
		Type:                            types.StringValue(webhookCompanyCreateResponse.Type),
		URL:                             types.StringValue(webhookCompanyCreateResponse.Url),
		Username:                        types.StringPointerValue(webhookCompanyCreateResponse.Username),
		HasPassword:                     types.BoolPointerValue(webhookCompanyCreateResponse.HasPassword),
		Active:                          types.BoolValue(webhookCompanyCreateResponse.Active),
		HasError:                        types.BoolPointerValue(webhookCompanyCreateResponse.HasError),
		EncryptionProtocol:              types.StringPointerValue(webhookCompanyCreateResponse.EncryptionProtocol),
		CommunicationFormat:             types.StringValue(webhookCompanyCreateResponse.CommunicationFormat),
		AcceptsExpiredCertificate:       types.BoolPointerValue(webhookCompanyCreateResponse.AcceptsExpiredCertificate),
		AcceptsSelfSignedCertificate:    types.BoolPointerValue(webhookCompanyCreateResponse.AcceptsSelfSignedCertificate),
		AcceptsUntrustedRootCertificate: types.BoolPointerValue(webhookCompanyCreateResponse.AcceptsUntrustedRootCertificate),
		PopulateSoapActionHeader:        types.BoolPointerValue(webhookCompanyCreateResponse.PopulateSoapActionHeader),
		CertificateAlias:                types.StringPointerValue(webhookCompanyCreateResponse.CertificateAlias),
		Password:                        types.StringPointerValue(createCompanyWebhookRequest.Password), //FIXME: figure out how to hide this / or if not needed to hide
		Links: types.ObjectValueMust(linksAttributeMap, mapWebhooksLinks(
			webhookCompanyCreateResponse.Links.Self.Href,
			webhookCompanyCreateResponse.Links.GenerateHmac.Href,
			webhookCompanyCreateResponse.Links.Company.Href,
			webhookCompanyCreateResponse.Links.TestWebhook.Href),
		),
		AdditionalSettings: types.ObjectValueMust(additionalSettingsAttributeMapCompany, map[string]attr.Value{
			"properties": types.MapValueMust(types.BoolType, properties),
		}),
		FilterMerchantAccountType: types.StringValue(createCompanyWebhookRequest.FilterMerchantAccountType),
		FilterMerchantAccounts:    types.ListValueMust(types.StringType, filterMerchantAccounts),
	}

	// Set state with the fully populated webhookCompanyCreateResponse
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *webhookCompanyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state webhooksCompanyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data management.WebhooksCompanyLevelApiGetWebhookInput
	if r.companyAccount != "" && state.WebhooksCompany.ID.ValueString() != "" {
		data = r.client.Management().WebhooksCompanyLevelApi.GetWebhookInput(r.companyAccount, state.WebhooksCompany.ID.ValueString())
	}

	webhookCompanyGetRequest, _, _ := r.client.Management().WebhooksCompanyLevelApi.GetWebhook(ctx, data)
	_, ok := webhookCompanyGetRequest.GetIdOk()
	//TODO: check this state logic once over to check edgecases
	if !ok && state.WebhooksCompany.ID.ValueString() != "" {
		resp.State.RemoveResource(ctx)
		return
	}

	var properties map[string]attr.Value
	if webhookCompanyGetRequest.AdditionalSettings.Properties != nil {
		properties = mapWebhooksAdditionalSettingsProperties(*webhookCompanyGetRequest.AdditionalSettings.Properties)
	}

	var filterMerchantAccounts []attr.Value
	if len(webhookCompanyGetRequest.FilterMerchantAccounts) > 0 {
		filterMerchantAccounts = mapWebhooksCompanyFilterMerchantAccounts(webhookCompanyGetRequest.FilterMerchantAccounts)
	}

	state = webhooksCompanyResourceModel{
		webhooksCompanyModel{
			ID:                              types.StringPointerValue(webhookCompanyGetRequest.Id),
			Type:                            types.StringValue(webhookCompanyGetRequest.Type),
			URL:                             types.StringValue(webhookCompanyGetRequest.Url),
			Username:                        types.StringPointerValue(webhookCompanyGetRequest.Username),
			HasPassword:                     types.BoolPointerValue(webhookCompanyGetRequest.HasPassword),
			Active:                          types.BoolValue(webhookCompanyGetRequest.Active),
			HasError:                        types.BoolPointerValue(webhookCompanyGetRequest.HasError),
			EncryptionProtocol:              types.StringPointerValue(webhookCompanyGetRequest.EncryptionProtocol),
			CommunicationFormat:             types.StringValue(webhookCompanyGetRequest.CommunicationFormat),
			AcceptsExpiredCertificate:       types.BoolPointerValue(webhookCompanyGetRequest.AcceptsExpiredCertificate),
			AcceptsSelfSignedCertificate:    types.BoolPointerValue(webhookCompanyGetRequest.AcceptsSelfSignedCertificate),
			AcceptsUntrustedRootCertificate: types.BoolPointerValue(webhookCompanyGetRequest.AcceptsUntrustedRootCertificate),
			PopulateSoapActionHeader:        types.BoolPointerValue(webhookCompanyGetRequest.PopulateSoapActionHeader),
			CertificateAlias:                types.StringPointerValue(webhookCompanyGetRequest.CertificateAlias),
			Links: types.ObjectValueMust(linksAttributeMap, mapWebhooksLinks(
				webhookCompanyGetRequest.Links.Self.Href,
				webhookCompanyGetRequest.Links.GenerateHmac.Href,
				webhookCompanyGetRequest.Links.Company.Href,
				webhookCompanyGetRequest.Links.TestWebhook.Href),
			),
			AdditionalSettings: types.ObjectValueMust(additionalSettingsAttributeMapCompany, map[string]attr.Value{
				"properties": types.MapValueMust(types.BoolType, properties),
			}),
			FilterMerchantAccountType: types.StringPointerValue(webhookCompanyGetRequest.FilterMerchantAccountType),
			FilterMerchantAccounts:    types.ListValueMust(types.StringType, filterMerchantAccounts),
		},
	}

	tflog.Debug(ctx, "Reading company webhook...")

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *webhookCompanyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating adyen company webhook")

	// Retrieve values from the plan
	var plan webhooksCompanyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateCompanyWebhookRequest := &management.UpdateCompanyWebhookRequest{
		AcceptsExpiredCertificate:       plan.WebhooksCompany.AcceptsExpiredCertificate.ValueBoolPointer(),
		AcceptsSelfSignedCertificate:    plan.WebhooksCompany.AcceptsSelfSignedCertificate.ValueBoolPointer(),
		AcceptsUntrustedRootCertificate: plan.WebhooksCompany.AcceptsUntrustedRootCertificate.ValueBoolPointer(),
		Active:                          plan.WebhooksCompany.Active.ValueBoolPointer(),
		CommunicationFormat:             plan.WebhooksCompany.CommunicationFormat.ValueStringPointer(),
		Password:                        plan.WebhooksCompany.Password.ValueStringPointer(),
		PopulateSoapActionHeader:        plan.WebhooksCompany.PopulateSoapActionHeader.ValueBoolPointer(),
		Url:                             plan.WebhooksCompany.URL.ValueStringPointer(),
		Username:                        plan.WebhooksCompany.Username.ValueStringPointer(),
	}

	// Create a new webhook
	webhookCompanyUpdateRequest := r.client.
		Management().
		WebhooksCompanyLevelApi.
		UpdateWebhookInput(r.companyAccount, plan.WebhooksCompany.ID.ValueString()).
		UpdateCompanyWebhookRequest(*updateCompanyWebhookRequest)
	webhookCompanyUpdateResponse, _, err := r.client.
		Management().
		WebhooksCompanyLevelApi.
		UpdateWebhook(ctx, webhookCompanyUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating company webhook",
			"Could not create company webhook, unexpected error: "+err.Error(),
		)
		return
	}

	var properties map[string]attr.Value
	if webhookCompanyUpdateResponse.AdditionalSettings.Properties != nil {
		properties = mapWebhooksAdditionalSettingsProperties(*webhookCompanyUpdateResponse.AdditionalSettings.Properties)
	}

	// Map response body to schema and populate Computed attribute values
	plan.WebhooksCompany = webhooksCompanyModel{
		ID:                              types.StringPointerValue(webhookCompanyUpdateResponse.Id),
		Type:                            types.StringValue(webhookCompanyUpdateResponse.Type),
		URL:                             types.StringValue(webhookCompanyUpdateResponse.Url),
		Username:                        types.StringPointerValue(webhookCompanyUpdateResponse.Username),
		HasPassword:                     types.BoolPointerValue(webhookCompanyUpdateResponse.HasPassword),
		Active:                          types.BoolValue(webhookCompanyUpdateResponse.Active),
		HasError:                        types.BoolPointerValue(webhookCompanyUpdateResponse.HasError),
		EncryptionProtocol:              types.StringPointerValue(webhookCompanyUpdateResponse.EncryptionProtocol),
		CommunicationFormat:             types.StringValue(webhookCompanyUpdateResponse.CommunicationFormat),
		AcceptsExpiredCertificate:       types.BoolPointerValue(webhookCompanyUpdateResponse.AcceptsExpiredCertificate),
		AcceptsSelfSignedCertificate:    types.BoolPointerValue(webhookCompanyUpdateResponse.AcceptsSelfSignedCertificate),
		AcceptsUntrustedRootCertificate: types.BoolPointerValue(webhookCompanyUpdateResponse.AcceptsUntrustedRootCertificate),
		PopulateSoapActionHeader:        types.BoolPointerValue(webhookCompanyUpdateResponse.PopulateSoapActionHeader),
		CertificateAlias:                types.StringPointerValue(webhookCompanyUpdateResponse.CertificateAlias),
		Password:                        types.StringPointerValue(updateCompanyWebhookRequest.Password), //FIXME
		Links: types.ObjectValueMust(linksAttributeMap, mapWebhooksLinks(
			webhookCompanyUpdateResponse.Links.Self.Href,
			webhookCompanyUpdateResponse.Links.GenerateHmac.Href,
			webhookCompanyUpdateResponse.Links.Company.Href,
			webhookCompanyUpdateResponse.Links.TestWebhook.Href),
		),
		AdditionalSettings: types.ObjectValueMust(additionalSettingsAttributeMapCompany, map[string]attr.Value{
			"properties": types.MapValueMust(types.BoolType, properties),
		}),
	}

	// Set state with the fully populated webhookCreateRequest
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *webhookCompanyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state webhooksCompanyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	removeWebhookInput := r.client.Management().WebhooksCompanyLevelApi.RemoveWebhookInput(r.companyAccount, state.WebhooksCompany.ID.ValueString())
	_, err := r.client.Management().WebhooksCompanyLevelApi.RemoveWebhook(ctx, removeWebhookInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Webhooks Company",
			"Could not delete company webhook, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *webhookCompanyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
