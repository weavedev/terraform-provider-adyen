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
	_ resource.Resource              = &webhookMerchantResource{}
	_ resource.ResourceWithConfigure = &webhookMerchantResource{}
)

// webhookResource is the resource implementation.
type webhookMerchantResource struct {
	client *adyen.APIClient
}

// NewWebhooksMerchantResource is a helper function to simplify the provider implementation.
func NewWebhooksMerchantResource() resource.Resource {
	return &webhookMerchantResource{}
}

// webhooksMerchantResourceModel maps the "webhooks_merchant" schema data for a resource.
type webhooksMerchantResourceModel struct {
	WebhooksMerchant webhooksMerchantModel `tfsdk:"webhooks_merchant"`
}

type webhooksMerchantModel struct {
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
}

// Configure adds the provider configured client to the resource.
func (r *webhookMerchantResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *webhookMerchantResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhooks_merchant"
}

// Schema defines the schema for the resource.
func (r *webhookMerchantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"webhooks_merchant": schema.SingleNestedAttribute{
				Description: "Subscribe to receive webhook notifications about events related to your merchant account.\n\n" +
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
							"merchant": schema.SingleNestedAttribute{
								Attributes: map[string]schema.Attribute{
									"href": schema.StringAttribute{Computed: true},
								},
								Computed:    true,
								Description: "The API URL to the merchant account associated with the webhook.",
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
							"include_event_codes": schema.ListAttribute{
								Computed:    true,
								ElementType: types.StringType,
								Description: "Object containing list of event codes for which the notification will be sent.",
							},
							"exclude_event_codes": schema.ListAttribute{
								Computed:    true,
								ElementType: types.StringType,
								Description: "Object containing list of event codes for which the notification will NOT be sent.",
							},
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *webhookMerchantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating adyen merchant webhook")

	// Retrieve values from the plan
	var plan webhooksMerchantResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	createMerchantWebhookRequest := &management.CreateMerchantWebhookRequest{
		AcceptsExpiredCertificate:       plan.WebhooksMerchant.AcceptsExpiredCertificate.ValueBoolPointer(),
		AcceptsSelfSignedCertificate:    plan.WebhooksMerchant.AcceptsSelfSignedCertificate.ValueBoolPointer(),
		AcceptsUntrustedRootCertificate: plan.WebhooksMerchant.AcceptsUntrustedRootCertificate.ValueBoolPointer(),
		Active:                          plan.WebhooksMerchant.Active.ValueBool(),
		CommunicationFormat:             plan.WebhooksMerchant.CommunicationFormat.ValueString(),
		Password:                        plan.WebhooksMerchant.Password.ValueStringPointer(),
		PopulateSoapActionHeader:        plan.WebhooksMerchant.PopulateSoapActionHeader.ValueBoolPointer(),
		Type:                            plan.WebhooksMerchant.Type.ValueString(),
		Url:                             plan.WebhooksMerchant.URL.ValueString(),
		Username:                        plan.WebhooksMerchant.Username.ValueStringPointer(),
	}

	// Create a new webhook
	webhookMerchantCreateRequest := r.client.
		Management().
		WebhooksMerchantLevelApi.
		SetUpWebhookInput(r.client.GetConfig().MerchantAccount).
		CreateMerchantWebhookRequest(*createMerchantWebhookRequest)
	webhookMerchantCreateResponse, _, err := r.client.
		Management().
		WebhooksMerchantLevelApi.
		SetUpWebhook(ctx, webhookMerchantCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating merchant webhook",
			"Could not create merchant webhook, unexpected error: "+err.Error(),
		)
		return
	}

	includeEventCodes := mapWebhooksAdditionalSettingsEventCodes(webhookMerchantCreateResponse.AdditionalSettings.IncludeEventCodes)
	excludeEventCodes := mapWebhooksAdditionalSettingsEventCodes(webhookMerchantCreateResponse.AdditionalSettings.ExcludeEventCodes)
	var properties map[string]attr.Value
	if webhookMerchantCreateResponse.AdditionalSettings.Properties != nil {
		properties = mapWebhooksAdditionalSettingsProperties(*webhookMerchantCreateResponse.AdditionalSettings.Properties)
	}

	// Map response body to schema and populate with attribute values
	plan.WebhooksMerchant = webhooksMerchantModel{
		ID:                              types.StringPointerValue(webhookMerchantCreateResponse.Id),
		Description:                     types.StringPointerValue(webhookMerchantCreateResponse.Description),
		Type:                            types.StringValue(webhookMerchantCreateResponse.Type),
		URL:                             types.StringValue(webhookMerchantCreateResponse.Url),
		Username:                        types.StringPointerValue(webhookMerchantCreateResponse.Username),
		HasPassword:                     types.BoolPointerValue(webhookMerchantCreateResponse.HasPassword),
		Active:                          types.BoolValue(webhookMerchantCreateResponse.Active),
		HasError:                        types.BoolPointerValue(webhookMerchantCreateResponse.HasError),
		EncryptionProtocol:              types.StringPointerValue(webhookMerchantCreateResponse.EncryptionProtocol),
		CommunicationFormat:             types.StringValue(webhookMerchantCreateResponse.CommunicationFormat),
		AcceptsExpiredCertificate:       types.BoolPointerValue(webhookMerchantCreateResponse.AcceptsExpiredCertificate),
		AcceptsSelfSignedCertificate:    types.BoolPointerValue(webhookMerchantCreateResponse.AcceptsSelfSignedCertificate),
		AcceptsUntrustedRootCertificate: types.BoolPointerValue(webhookMerchantCreateResponse.AcceptsUntrustedRootCertificate),
		PopulateSoapActionHeader:        types.BoolPointerValue(webhookMerchantCreateResponse.PopulateSoapActionHeader),
		CertificateAlias:                types.StringPointerValue(webhookMerchantCreateResponse.CertificateAlias),
		Password:                        types.StringPointerValue(createMerchantWebhookRequest.Password), //FIXME: figure out how to hide this / or if not needed to hide
		Links: types.ObjectValueMust(linksAttributeMap, mapWebhooksLinks(
			webhookMerchantCreateResponse.Links.Self.Href,
			webhookMerchantCreateResponse.Links.GenerateHmac.Href,
			webhookMerchantCreateResponse.Links.Merchant.Href,
			webhookMerchantCreateResponse.Links.TestWebhook.Href),
		),
		AdditionalSettings: types.ObjectValueMust(additionalSettingsAttributeMapMerchant, map[string]attr.Value{
			"include_event_codes": types.ListValueMust(types.StringType, includeEventCodes),
			"exclude_event_codes": types.ListValueMust(types.StringType, excludeEventCodes),
			"properties":          types.MapValueMust(types.BoolType, properties),
		}),
	}

	// Set state with the fully populated webhookMerchantCreateResponse
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *webhookMerchantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state webhooksMerchantResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data management.WebhooksMerchantLevelApiGetWebhookInput
	if r.client.GetConfig().MerchantAccount != "" && state.WebhooksMerchant.ID.ValueString() != "" {
		data = r.client.Management().WebhooksMerchantLevelApi.GetWebhookInput(r.client.GetConfig().MerchantAccount, state.WebhooksMerchant.ID.ValueString())
	}

	webhookMerchantGetRequest, _, _ := r.client.Management().WebhooksMerchantLevelApi.GetWebhook(ctx, data)
	_, ok := webhookMerchantGetRequest.GetIdOk()
	//TODO: check this state logic once over to check edgecases
	if !ok && state.WebhooksMerchant.ID.ValueString() != "" {
		resp.State.RemoveResource(ctx)
		return
	}

	includeEventCodes := mapWebhooksAdditionalSettingsEventCodes(webhookMerchantGetRequest.AdditionalSettings.IncludeEventCodes)
	excludeEventCodes := mapWebhooksAdditionalSettingsEventCodes(webhookMerchantGetRequest.AdditionalSettings.ExcludeEventCodes)
	var properties map[string]attr.Value
	if webhookMerchantGetRequest.AdditionalSettings.Properties != nil {
		properties = mapWebhooksAdditionalSettingsProperties(*webhookMerchantGetRequest.AdditionalSettings.Properties)
	}

	state = webhooksMerchantResourceModel{
		webhooksMerchantModel{
			ID:                              types.StringPointerValue(webhookMerchantGetRequest.Id),
			Type:                            types.StringValue(webhookMerchantGetRequest.Type),
			URL:                             types.StringValue(webhookMerchantGetRequest.Url),
			Username:                        types.StringPointerValue(webhookMerchantGetRequest.Username),
			HasPassword:                     types.BoolPointerValue(webhookMerchantGetRequest.HasPassword),
			Active:                          types.BoolValue(webhookMerchantGetRequest.Active),
			HasError:                        types.BoolPointerValue(webhookMerchantGetRequest.HasError),
			EncryptionProtocol:              types.StringPointerValue(webhookMerchantGetRequest.EncryptionProtocol),
			CommunicationFormat:             types.StringValue(webhookMerchantGetRequest.CommunicationFormat),
			AcceptsExpiredCertificate:       types.BoolPointerValue(webhookMerchantGetRequest.AcceptsExpiredCertificate),
			AcceptsSelfSignedCertificate:    types.BoolPointerValue(webhookMerchantGetRequest.AcceptsSelfSignedCertificate),
			AcceptsUntrustedRootCertificate: types.BoolPointerValue(webhookMerchantGetRequest.AcceptsUntrustedRootCertificate),
			PopulateSoapActionHeader:        types.BoolPointerValue(webhookMerchantGetRequest.PopulateSoapActionHeader),
			CertificateAlias:                types.StringPointerValue(webhookMerchantGetRequest.CertificateAlias),
			Links: types.ObjectValueMust(linksAttributeMap, mapWebhooksLinks(
				webhookMerchantGetRequest.Links.Self.Href,
				webhookMerchantGetRequest.Links.GenerateHmac.Href,
				webhookMerchantGetRequest.Links.Merchant.Href,
				webhookMerchantGetRequest.Links.TestWebhook.Href),
			),
			AdditionalSettings: types.ObjectValueMust(additionalSettingsAttributeMapMerchant, map[string]attr.Value{
				"include_event_codes": types.ListValueMust(types.StringType, includeEventCodes),
				"exclude_event_codes": types.ListValueMust(types.StringType, excludeEventCodes),
				"properties":          types.MapValueMust(types.BoolType, properties),
			}),
		},
	}

	tflog.Debug(ctx, "Reading merchant webhook...")

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *webhookMerchantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Updating adyen merchant webhook")

	// Retrieve values from the plan
	var plan webhooksMerchantResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updateMerchantWebhookRequest := &management.UpdateMerchantWebhookRequest{
		AcceptsExpiredCertificate:       plan.WebhooksMerchant.AcceptsExpiredCertificate.ValueBoolPointer(),
		AcceptsSelfSignedCertificate:    plan.WebhooksMerchant.AcceptsSelfSignedCertificate.ValueBoolPointer(),
		AcceptsUntrustedRootCertificate: plan.WebhooksMerchant.AcceptsUntrustedRootCertificate.ValueBoolPointer(),
		Active:                          plan.WebhooksMerchant.Active.ValueBoolPointer(),
		CommunicationFormat:             plan.WebhooksMerchant.CommunicationFormat.ValueStringPointer(),
		Password:                        plan.WebhooksMerchant.Password.ValueStringPointer(),
		PopulateSoapActionHeader:        plan.WebhooksMerchant.PopulateSoapActionHeader.ValueBoolPointer(),
		Url:                             plan.WebhooksMerchant.URL.ValueStringPointer(),
		Username:                        plan.WebhooksMerchant.Username.ValueStringPointer(),
	}

	// Create a new webhook
	webhookMerchantUpdateRequest := r.client.
		Management().
		WebhooksMerchantLevelApi.
		UpdateWebhookInput(r.client.GetConfig().MerchantAccount, plan.WebhooksMerchant.ID.ValueString()).
		UpdateMerchantWebhookRequest(*updateMerchantWebhookRequest)
	webhookMerchantUpdateResponse, _, err := r.client.
		Management().
		WebhooksMerchantLevelApi.
		UpdateWebhook(ctx, webhookMerchantUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating merchant webhook",
			"Could not create merchant webhook, unexpected error: "+err.Error(),
		)
		return
	}

	includeEventCodes := mapWebhooksAdditionalSettingsEventCodes(webhookMerchantUpdateResponse.AdditionalSettings.IncludeEventCodes)
	excludeEventCodes := mapWebhooksAdditionalSettingsEventCodes(webhookMerchantUpdateResponse.AdditionalSettings.ExcludeEventCodes)
	var properties map[string]attr.Value
	if webhookMerchantUpdateResponse.AdditionalSettings.Properties != nil {
		properties = mapWebhooksAdditionalSettingsProperties(*webhookMerchantUpdateResponse.AdditionalSettings.Properties)
	}

	// Map response body to schema and populate Computed attribute values
	plan.WebhooksMerchant = webhooksMerchantModel{
		ID:                              types.StringPointerValue(webhookMerchantUpdateResponse.Id),
		Type:                            types.StringValue(webhookMerchantUpdateResponse.Type),
		URL:                             types.StringValue(webhookMerchantUpdateResponse.Url),
		Username:                        types.StringPointerValue(webhookMerchantUpdateResponse.Username),
		HasPassword:                     types.BoolPointerValue(webhookMerchantUpdateResponse.HasPassword),
		Active:                          types.BoolValue(webhookMerchantUpdateResponse.Active),
		HasError:                        types.BoolPointerValue(webhookMerchantUpdateResponse.HasError),
		EncryptionProtocol:              types.StringPointerValue(webhookMerchantUpdateResponse.EncryptionProtocol),
		CommunicationFormat:             types.StringValue(webhookMerchantUpdateResponse.CommunicationFormat),
		AcceptsExpiredCertificate:       types.BoolPointerValue(webhookMerchantUpdateResponse.AcceptsExpiredCertificate),
		AcceptsSelfSignedCertificate:    types.BoolPointerValue(webhookMerchantUpdateResponse.AcceptsSelfSignedCertificate),
		AcceptsUntrustedRootCertificate: types.BoolPointerValue(webhookMerchantUpdateResponse.AcceptsUntrustedRootCertificate),
		PopulateSoapActionHeader:        types.BoolPointerValue(webhookMerchantUpdateResponse.PopulateSoapActionHeader),
		CertificateAlias:                types.StringPointerValue(webhookMerchantUpdateResponse.CertificateAlias),
		Password:                        types.StringPointerValue(updateMerchantWebhookRequest.Password), //FIXME
		Links: types.ObjectValueMust(linksAttributeMap, mapWebhooksLinks(
			webhookMerchantUpdateResponse.Links.Self.Href,
			webhookMerchantUpdateResponse.Links.GenerateHmac.Href,
			webhookMerchantUpdateResponse.Links.Merchant.Href,
			webhookMerchantUpdateResponse.Links.TestWebhook.Href),
		),
		AdditionalSettings: types.ObjectValueMust(additionalSettingsAttributeMapMerchant, map[string]attr.Value{
			"include_event_codes": types.ListValueMust(types.StringType, includeEventCodes),
			"exclude_event_codes": types.ListValueMust(types.StringType, excludeEventCodes),
			"properties":          types.MapValueMust(types.BoolType, properties),
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
func (r *webhookMerchantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state webhooksMerchantResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	removeWebhookInput := r.client.Management().WebhooksMerchantLevelApi.RemoveWebhookInput(r.client.GetConfig().MerchantAccount, state.WebhooksMerchant.ID.ValueString())
	_, err := r.client.Management().WebhooksMerchantLevelApi.RemoveWebhook(ctx, removeWebhookInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Webhooks Merchant",
			"Could not delete merchant webhook, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *webhookMerchantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
