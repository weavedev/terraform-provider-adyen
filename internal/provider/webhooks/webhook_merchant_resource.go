package webhooks

import (
	"context"
	"fmt"
	"github.com/adyen/adyen-go-api-library/v9/src/adyen"
	"github.com/adyen/adyen-go-api-library/v9/src/management"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	//AdditionalSettings              webhooksMerchantAdditionalSettingsModel `tfsdk:"additional_settings"`
	//Links                           webhooksMerchantLinksModel              `tfsdk:"links"`
}

type webhooksMerchantAdditionalSettingsModel struct {
	IncludeEventCodes []types.String `tfsdk:"include_event_codes"`
	ExcludeEventCodes []types.String `tfsdk:"exclude_event_codes"`
	Properties        types.Map      `tfsdk:"properties"`
}

type webhooksMerchantLinksModel struct {
	Self         webhooksLinksHrefModel `tfsdk:"self"`
	GenerateHmac webhooksLinksHrefModel `tfsdk:"generate_hmac"`
	Merchant     webhooksLinksHrefModel `tfsdk:"merchant"`
	TestWebhook  webhooksLinksHrefModel `tfsdk:"test_webhook"`
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
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "The unique identifier for the webhook.",
					},
					"type": schema.StringAttribute{
						Required:    true,
						Description: "The type of the webhook.",
					},
					"url": schema.StringAttribute{
						Required:    true,
						Description: "The URL the webhook will send requests to.",
					},
					"username": schema.StringAttribute{
						Required:    true,
						Description: "The username required for basic authentication.",
					},
					"password": schema.StringAttribute{
						Sensitive:   true,
						Required:    true,
						Description: "The password required for basic authentication.",
					},
					"has_password": schema.BoolAttribute{
						Computed:    true,
						Description: "Indicates if the webhook is configured with a password.",
					},
					"active": schema.BoolAttribute{
						Required:    true,
						Description: "Indicates if the webhook is active.",
					},
					"communication_format": schema.StringAttribute{
						Required:    true,
						Description: "The format of the communication (e.g., 'json').",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "A description of the webhook.",
					},
					"encryption_protocol": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The encryption protocol used by the webhook.",
					},
					"has_error": schema.BoolAttribute{
						Computed:    true,
						Description: "Indicates if there is an error with the webhook.",
					},
					"certificate_alias": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The alias of the certificate.",
					},
					"populate_soap_action_header": schema.BoolAttribute{
						Optional:    true,
						Description: "Indicates if the SOAP action header should be populated.",
					},
					"accepts_expired_certificate": schema.BoolAttribute{
						Required:    true,
						Description: "Indicates if expired certificates are accepted.",
					},
					"accepts_self_signed_certificate": schema.BoolAttribute{
						Required:    true,
						Description: "Indicates if self-signed certificates are accepted.",
					},
					"accepts_untrusted_root_certificate": schema.BoolAttribute{
						Required:    true,
						Description: "Indicates if untrusted root certificates are accepted.",
					},
					//"additional_settings": schema.SingleNestedAttribute{
					//	Computed: true,
					//	Attributes: map[string]schema.Attribute{
					//		"properties": schema.MapAttribute{
					//			Computed:    true,
					//			ElementType: types.BoolType,
					//		},
					//		"include_event_codes": schema.ListAttribute{
					//			Computed:    true,
					//			ElementType: types.StringType,
					//		},
					//		"exclude_event_codes": schema.ListAttribute{
					//			Computed:    true,
					//			ElementType: types.StringType,
					//		},
					//	},
					//},
					//"links": schema.SingleNestedAttribute{
					//	Computed: true,
					//	Attributes: map[string]schema.Attribute{
					//		"self": schema.SingleNestedAttribute{
					//			Computed: true,
					//			Attributes: map[string]schema.Attribute{
					//				"href": schema.StringAttribute{Computed: true},
					//			},
					//		},
					//		"generate_hmac": schema.SingleNestedAttribute{
					//			Computed: true,
					//			Attributes: map[string]schema.Attribute{
					//				"href": schema.StringAttribute{Computed: true},
					//			},
					//		},
					//		"merchant": schema.SingleNestedAttribute{
					//			Computed: true,
					//			Attributes: map[string]schema.Attribute{
					//				"href": schema.StringAttribute{Computed: true},
					//			},
					//		},
					//		"test_webhook": schema.SingleNestedAttribute{
					//			Computed: true,
					//			Attributes: map[string]schema.Attribute{
					//				"href": schema.StringAttribute{Computed: true},
					//			},
					//		},
					//	},
					//},
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

	// Map response body to schema and populate Computed attribute values
	plan.WebhooksMerchant = webhooksMerchantModel{
		ID:                              types.StringPointerValue(webhookCreateResponse.Id),
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
		Password:                        types.StringPointerValue(createMerchantWebhookRequest.Password),
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
	data := r.client.Management().WebhooksMerchantLevelApi.ListAllWebhooksInput(r.client.GetConfig().MerchantAccount)

	listWebhooksMerchant, _, err := r.client.Management().WebhooksMerchantLevelApi.ListAllWebhooks(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Adyen Webhooks",
			err.Error(),
		)
		return
	}

	for _, webhookMerchantState := range listWebhooksMerchant.Data {
		state = webhooksMerchantResourceModel{
			webhooksMerchantModel{
				ID:                              types.StringValue(*webhookMerchantState.Id),
				Type:                            types.StringValue(webhookMerchantState.Type),
				URL:                             types.StringValue(webhookMerchantState.Url),
				Username:                        types.StringValue(*webhookMerchantState.Username),
				HasPassword:                     types.BoolValue(*webhookMerchantState.HasPassword),
				Active:                          types.BoolValue(webhookMerchantState.Active),
				HasError:                        types.BoolValue(*webhookMerchantState.HasError),
				EncryptionProtocol:              types.StringValue(*webhookMerchantState.EncryptionProtocol),
				CommunicationFormat:             types.StringValue(webhookMerchantState.CommunicationFormat),
				AcceptsExpiredCertificate:       types.BoolValue(*webhookMerchantState.AcceptsExpiredCertificate),
				AcceptsSelfSignedCertificate:    types.BoolValue(*webhookMerchantState.AcceptsSelfSignedCertificate),
				AcceptsUntrustedRootCertificate: types.BoolValue(*webhookMerchantState.AcceptsUntrustedRootCertificate),
				PopulateSoapActionHeader:        types.BoolValue(*webhookMerchantState.PopulateSoapActionHeader),
				CertificateAlias:                types.StringValue(*webhookMerchantState.CertificateAlias),
			},
		}
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *webhookMerchantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *webhookMerchantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
