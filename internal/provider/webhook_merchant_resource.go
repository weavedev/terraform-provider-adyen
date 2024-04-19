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

//TODO: check for more consistent naming of vars.

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
	webhookCreateRequest := r.client.
		Management().
		WebhooksMerchantLevelApi.
		SetUpWebhookInput(r.client.GetConfig().MerchantAccount).
		CreateMerchantWebhookRequest(*createMerchantWebhookRequest)
	webhookCreateResponse, _, err := r.client.
		Management().
		WebhooksMerchantLevelApi.
		SetUpWebhook(ctx, webhookCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating merchant webhook",
			"Could not create merchant webhook, unexpected error: "+err.Error(),
		)
		return
	}

	var includeEventCodes []attr.Value
	for _, code := range webhookCreateResponse.AdditionalSettings.IncludeEventCodes {
		includeEventCodes = append(includeEventCodes, types.StringValue(code))
	}

	var excludeEventCodes []attr.Value
	for _, code := range webhookCreateResponse.AdditionalSettings.ExcludeEventCodes {
		excludeEventCodes = append(excludeEventCodes, types.StringValue(code))
	}

	properties := make(map[string]attr.Value)
	if webhookCreateResponse.AdditionalSettings.Properties != nil {
		for k, v := range *webhookCreateResponse.AdditionalSettings.Properties {
			properties[k] = types.BoolValue(v)
		}
	}

	// Map response body to schema and populate with attribute values
	plan.WebhooksMerchant = webhooksMerchantModel{
		ID:                              types.StringPointerValue(webhookCreateResponse.Id),
		Description:                     types.StringPointerValue(webhookCreateResponse.Description),
		Type:                            types.StringValue(webhookCreateResponse.Type),
		URL:                             types.StringValue(webhookCreateResponse.Url),
		Username:                        types.StringPointerValue(webhookCreateResponse.Username),
		HasPassword:                     types.BoolPointerValue(webhookCreateResponse.HasPassword),
		Active:                          types.BoolValue(webhookCreateResponse.Active),
		HasError:                        types.BoolPointerValue(webhookCreateResponse.HasError),
		EncryptionProtocol:              types.StringPointerValue(webhookCreateResponse.EncryptionProtocol),
		CommunicationFormat:             types.StringValue(webhookCreateResponse.CommunicationFormat),
		AcceptsExpiredCertificate:       types.BoolPointerValue(webhookCreateResponse.AcceptsExpiredCertificate),
		AcceptsSelfSignedCertificate:    types.BoolPointerValue(webhookCreateResponse.AcceptsSelfSignedCertificate),
		AcceptsUntrustedRootCertificate: types.BoolPointerValue(webhookCreateResponse.AcceptsUntrustedRootCertificate),
		PopulateSoapActionHeader:        types.BoolPointerValue(webhookCreateResponse.PopulateSoapActionHeader),
		CertificateAlias:                types.StringPointerValue(webhookCreateResponse.CertificateAlias),
		Links: types.ObjectValueMust(
			map[string]attr.Type{
				"self": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"href": types.StringType,
					},
				},
				"generate_hmac": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"href": types.StringType,
					},
				},
				"merchant": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"href": types.StringType,
					},
				},
				"test_webhook": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"href": types.StringType,
					},
				},
			}, map[string]attr.Value{
				"self": types.ObjectValueMust(map[string]attr.Type{
					"href": types.StringType,
				}, map[string]attr.Value{
					"href": types.StringPointerValue(webhookCreateResponse.Links.Self.Href),
				}),
				"generate_hmac": types.ObjectValueMust(map[string]attr.Type{
					"href": types.StringType,
				}, map[string]attr.Value{
					"href": types.StringPointerValue(webhookCreateResponse.Links.Self.Href),
				}),
				"merchant": types.ObjectValueMust(map[string]attr.Type{
					"href": types.StringType,
				}, map[string]attr.Value{
					"href": types.StringPointerValue(webhookCreateResponse.Links.Self.Href),
				}),
				"test_webhook": types.ObjectValueMust(map[string]attr.Type{
					"href": types.StringType,
				}, map[string]attr.Value{
					"href": types.StringPointerValue(webhookCreateResponse.Links.Self.Href),
				}),
			}),
		Password: types.StringPointerValue(createMerchantWebhookRequest.Password), //FIXME: figure out how to hide this / or if not needed to hide
		AdditionalSettings: types.ObjectValueMust(map[string]attr.Type{
			"include_event_codes": types.ListType{
				ElemType: types.StringType,
			},
			"exclude_event_codes": types.ListType{
				ElemType: types.StringType,
			},
			"properties": types.MapType{
				ElemType: types.BoolType,
			},
		}, map[string]attr.Value{
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

	webhookMerchantData, _, err := r.client.Management().WebhooksMerchantLevelApi.GetWebhook(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Adyen Webhooks",
			err.Error(),
		)
		return
	}

	includeEventCodes := []attr.Value{}
	for _, code := range webhookMerchantData.AdditionalSettings.IncludeEventCodes {
		includeEventCodes = append(includeEventCodes, types.StringValue(code))
	}

	excludeEventCodes := []attr.Value{}
	for _, code := range webhookMerchantData.AdditionalSettings.ExcludeEventCodes {
		excludeEventCodes = append(excludeEventCodes, types.StringValue(code))
	}

	properties := make(map[string]attr.Value)
	if webhookMerchantData.AdditionalSettings.Properties != nil {
		for k, v := range *webhookMerchantData.AdditionalSettings.Properties {
			properties[k] = types.BoolValue(v)
		}
	}

	state = webhooksMerchantResourceModel{
		webhooksMerchantModel{
			ID:                              types.StringPointerValue(webhookMerchantData.Id),
			Type:                            types.StringValue(webhookMerchantData.Type),
			URL:                             types.StringValue(webhookMerchantData.Url),
			Username:                        types.StringPointerValue(webhookMerchantData.Username),
			HasPassword:                     types.BoolPointerValue(webhookMerchantData.HasPassword),
			Active:                          types.BoolValue(webhookMerchantData.Active),
			HasError:                        types.BoolPointerValue(webhookMerchantData.HasError),
			EncryptionProtocol:              types.StringPointerValue(webhookMerchantData.EncryptionProtocol),
			CommunicationFormat:             types.StringValue(webhookMerchantData.CommunicationFormat),
			AcceptsExpiredCertificate:       types.BoolPointerValue(webhookMerchantData.AcceptsExpiredCertificate),
			AcceptsSelfSignedCertificate:    types.BoolPointerValue(webhookMerchantData.AcceptsSelfSignedCertificate),
			AcceptsUntrustedRootCertificate: types.BoolPointerValue(webhookMerchantData.AcceptsUntrustedRootCertificate),
			PopulateSoapActionHeader:        types.BoolPointerValue(webhookMerchantData.PopulateSoapActionHeader),
			CertificateAlias:                types.StringPointerValue(webhookMerchantData.CertificateAlias),
			Links: types.ObjectValueMust(
				map[string]attr.Type{
					"self": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"href": types.StringType,
						},
					},
					"generate_hmac": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"href": types.StringType,
						},
					},
					"merchant": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"href": types.StringType,
						},
					},
					"test_webhook": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"href": types.StringType,
						},
					},
				}, map[string]attr.Value{
					"self": types.ObjectValueMust(map[string]attr.Type{
						"href": types.StringType,
					}, map[string]attr.Value{
						"href": types.StringPointerValue(webhookMerchantData.Links.Self.Href),
					}),
					"generate_hmac": types.ObjectValueMust(map[string]attr.Type{
						"href": types.StringType,
					}, map[string]attr.Value{
						"href": types.StringPointerValue(webhookMerchantData.Links.Self.Href),
					}),
					"merchant": types.ObjectValueMust(map[string]attr.Type{
						"href": types.StringType,
					}, map[string]attr.Value{
						"href": types.StringPointerValue(webhookMerchantData.Links.Self.Href),
					}),
					"test_webhook": types.ObjectValueMust(map[string]attr.Type{
						"href": types.StringType,
					}, map[string]attr.Value{
						"href": types.StringPointerValue(webhookMerchantData.Links.Self.Href),
					}),
				}),
			AdditionalSettings: types.ObjectValueMust(map[string]attr.Type{
				"include_event_codes": types.ListType{
					ElemType: types.StringType,
				},
				"exclude_event_codes": types.ListType{
					ElemType: types.StringType,
				},
				"properties": types.MapType{
					ElemType: types.BoolType,
				},
			}, map[string]attr.Value{
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
	webhookUpdateRequest := r.client.
		Management().
		WebhooksMerchantLevelApi.
		UpdateWebhookInput(r.client.GetConfig().MerchantAccount, plan.WebhooksMerchant.ID.ValueString()).
		UpdateMerchantWebhookRequest(*updateMerchantWebhookRequest)
	webhookUpdateResponse, _, err := r.client.
		Management().
		WebhooksMerchantLevelApi.
		UpdateWebhook(ctx, webhookUpdateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating merchant webhook",
			"Could not create merchant webhook, unexpected error: "+err.Error(),
		)
		return
	}

	includeEventCodes := []attr.Value{}
	for _, code := range webhookUpdateResponse.AdditionalSettings.IncludeEventCodes {
		includeEventCodes = append(includeEventCodes, types.StringValue(code))
	}

	excludeEventCodes := []attr.Value{}
	for _, code := range webhookUpdateResponse.AdditionalSettings.ExcludeEventCodes {
		excludeEventCodes = append(excludeEventCodes, types.StringValue(code))
	}

	properties := make(map[string]attr.Value)
	if webhookUpdateResponse.AdditionalSettings.Properties != nil {
		for k, v := range *webhookUpdateResponse.AdditionalSettings.Properties {
			properties[k] = types.BoolValue(v)
		}
	}

	// Map response body to schema and populate Computed attribute values
	plan.WebhooksMerchant = webhooksMerchantModel{
		ID:                              types.StringPointerValue(webhookUpdateResponse.Id),
		Type:                            types.StringValue(webhookUpdateResponse.Type),
		URL:                             types.StringValue(webhookUpdateResponse.Url),
		Username:                        types.StringPointerValue(webhookUpdateResponse.Username),
		HasPassword:                     types.BoolPointerValue(webhookUpdateResponse.HasPassword),
		Active:                          types.BoolValue(webhookUpdateResponse.Active),
		HasError:                        types.BoolPointerValue(webhookUpdateResponse.HasError),
		EncryptionProtocol:              types.StringPointerValue(webhookUpdateResponse.EncryptionProtocol),
		CommunicationFormat:             types.StringValue(webhookUpdateResponse.CommunicationFormat),
		AcceptsExpiredCertificate:       types.BoolPointerValue(webhookUpdateResponse.AcceptsExpiredCertificate),
		AcceptsSelfSignedCertificate:    types.BoolPointerValue(webhookUpdateResponse.AcceptsSelfSignedCertificate),
		AcceptsUntrustedRootCertificate: types.BoolPointerValue(webhookUpdateResponse.AcceptsUntrustedRootCertificate),
		PopulateSoapActionHeader:        types.BoolPointerValue(webhookUpdateResponse.PopulateSoapActionHeader),
		CertificateAlias:                types.StringPointerValue(webhookUpdateResponse.CertificateAlias),
		Password:                        types.StringPointerValue(updateMerchantWebhookRequest.Password),
		Links: types.ObjectValueMust(
			map[string]attr.Type{
				"self": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"href": types.StringType,
					},
				},
				"generate_hmac": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"href": types.StringType,
					},
				},
				"merchant": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"href": types.StringType,
					},
				},
				"test_webhook": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"href": types.StringType,
					},
				},
			}, map[string]attr.Value{
				"self": types.ObjectValueMust(map[string]attr.Type{
					"href": types.StringType,
				}, map[string]attr.Value{
					"href": types.StringPointerValue(webhookUpdateResponse.Links.Self.Href),
				}),
				"generate_hmac": types.ObjectValueMust(map[string]attr.Type{
					"href": types.StringType,
				}, map[string]attr.Value{
					"href": types.StringPointerValue(webhookUpdateResponse.Links.Self.Href),
				}),
				"merchant": types.ObjectValueMust(map[string]attr.Type{
					"href": types.StringType,
				}, map[string]attr.Value{
					"href": types.StringPointerValue(webhookUpdateResponse.Links.Self.Href),
				}),
				"test_webhook": types.ObjectValueMust(map[string]attr.Type{
					"href": types.StringType,
				}, map[string]attr.Value{
					"href": types.StringPointerValue(webhookUpdateResponse.Links.Self.Href),
				}),
			}),
		AdditionalSettings: types.ObjectValueMust(map[string]attr.Type{
			"include_event_codes": types.ListType{
				ElemType: types.StringType,
			},
			"exclude_event_codes": types.ListType{
				ElemType: types.StringType,
			},
			"properties": types.MapType{
				ElemType: types.BoolType,
			},
		}, map[string]attr.Value{
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

	data := r.client.Management().WebhooksMerchantLevelApi.RemoveWebhookInput(r.client.GetConfig().MerchantAccount, state.WebhooksMerchant.ID.ValueString())
	_, err := r.client.Management().WebhooksMerchantLevelApi.RemoveWebhook(ctx, data)
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
