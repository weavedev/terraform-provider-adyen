// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package webhooks

import (
	"context"
	"fmt"
	"github.com/adyen/adyen-go-api-library/v9/src/adyen"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &webhooksMerchantDataSource{}
	_ datasource.DataSourceWithConfigure = &webhooksMerchantDataSource{}
)

// webhooksMerchantDataSource defines the data source implementation.
type webhooksMerchantDataSource struct {
	client *adyen.APIClient
}

// webhooksMerchantDataSourceModel defines the data source implementation.
type webhooksMerchantDataSourceModel struct {
	Webhooks []webhooksModel `tfsdk:"webhooks_merchant"`
}

type webhooksModel struct {
	Links            webhookLinksModel  `tfsdk:"links"`
	ItemsTotal       types.Int64        `tfsdk:"items_total"`
	PagesTotal       types.Int64        `tfsdk:"pages_total"`
	AccountReference types.String       `tfsdk:"account_reference"`
	Data             []webhookDataModel `tfsdk:"data"`
}

type webhookLinksModel struct {
	First webhooksLinksHrefModel `tfsdk:"first"`
	Last  webhooksLinksHrefModel `tfsdk:"last"`
	Self  webhooksLinksHrefModel `tfsdk:"self"`
}

type webhooksLinksHrefModel struct {
	Href types.String `tfsdk:"href"`
}

type webhookDataModel struct {
	ID                              types.String            `tfsdk:"id"`
	Type                            types.String            `tfsdk:"type"`
	URL                             types.String            `tfsdk:"url"`
	Username                        types.String            `tfsdk:"username"`
	Description                     types.String            `tfsdk:"description"`
	HasPassword                     types.Bool              `tfsdk:"has_password"`
	Active                          types.Bool              `tfsdk:"active"`
	HasError                        types.Bool              `tfsdk:"has_error"`
	EncryptionProtocol              types.String            `tfsdk:"encryption_protocol"`
	CommunicationFormat             types.String            `tfsdk:"communication_format"`
	AcceptsExpiredCertificate       types.Bool              `tfsdk:"accepts_expired_certificate"`
	AcceptsSelfSignedCertificate    types.Bool              `tfsdk:"accepts_self_signed_certificate"`
	AcceptsUntrustedRootCertificate types.Bool              `tfsdk:"accepts_untrusted_root_certificate"`
	CertificateAlias                types.String            `tfsdk:"certificate_alias"`
	PopulateSoapActionHeader        types.Bool              `tfsdk:"populate_soap_action_header"`
	AdditionalSettings              webhookAdditionalModel  `tfsdk:"additional_settings"`
	Links                           []webhookDataLinksModel `tfsdk:"links"`
}

type webhookAdditionalModel struct {
	IncludeEventCodes []types.String `tfsdk:"include_event_codes"`
	ExcludeEventCodes []types.String `tfsdk:"exclude_event_codes"`
}

type webhookDataLinksModel struct {
	Self         webhooksLinksHrefModel `tfsdk:"self"`
	GenerateHmac webhooksLinksHrefModel `tfsdk:"generate_hmac"`
	Merchant     webhooksLinksHrefModel `tfsdk:"merchant"`
	TestWebhook  webhooksLinksHrefModel `tfsdk:"test_webhook"`
}

func NewWebhookMerchantDataSource() datasource.DataSource {
	return &webhooksMerchantDataSource{}
}

func (d *webhooksMerchantDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhooks_merchant"
}

func (d *webhooksMerchantDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *webhooksMerchantDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"webhooks_merchant": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"links": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"first": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"href": schema.StringAttribute{Computed: true},
									},
								},
								"last": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"href": schema.StringAttribute{Computed: true},
									},
								},
								"self": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"href": schema.StringAttribute{Computed: true},
									},
								}},
							Description:         "Links to webhooks",
							MarkdownDescription: "Links to webhooks",
						},
						"items_total": schema.Int64Attribute{
							Computed: true,
						},
						"pages_total": schema.Int64Attribute{
							Computed: true,
						},
						"account_reference": schema.StringAttribute{
							Computed: true,
						},
						"data": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id":                                 schema.StringAttribute{Computed: true},
									"type":                               schema.StringAttribute{Computed: true},
									"url":                                schema.StringAttribute{Computed: true},
									"username":                           schema.StringAttribute{Computed: true},
									"description":                        schema.StringAttribute{Computed: true},
									"has_password":                       schema.BoolAttribute{Computed: true},
									"active":                             schema.BoolAttribute{Computed: true},
									"has_error":                          schema.BoolAttribute{Computed: true},
									"encryption_protocol":                schema.StringAttribute{Computed: true},
									"communication_format":               schema.StringAttribute{Computed: true},
									"accepts_expired_certificate":        schema.BoolAttribute{Computed: true},
									"accepts_self_signed_certificate":    schema.BoolAttribute{Computed: true},
									"accepts_untrusted_root_certificate": schema.BoolAttribute{Computed: true},
									"certificate_alias":                  schema.StringAttribute{Computed: true},
									"populate_soap_action_header":        schema.BoolAttribute{Computed: true},
									"additional_settings": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"include_event_codes": schema.ListAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
											"exclude_event_codes": schema.ListAttribute{
												Computed:    true,
												ElementType: types.StringType,
											},
										}},
									"links": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"self": schema.SingleNestedAttribute{
													Computed: true,
													Attributes: map[string]schema.Attribute{
														"href": schema.StringAttribute{Computed: true},
													},
												},
												"generate_hmac": schema.SingleNestedAttribute{
													Computed: true,
													Attributes: map[string]schema.Attribute{
														"href": schema.StringAttribute{Computed: true},
													},
												},
												"merchant": schema.SingleNestedAttribute{
													Computed: true,
													Attributes: map[string]schema.Attribute{
														"href": schema.StringAttribute{Computed: true},
													},
												},
												"test_webhook": schema.SingleNestedAttribute{
													Computed: true,
													Attributes: map[string]schema.Attribute{
														"href": schema.StringAttribute{Computed: true},
													},
												},
											}},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *webhooksMerchantDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state webhooksMerchantDataSourceModel
	data := d.client.Management().WebhooksMerchantLevelApi.ListAllWebhooksInput(d.client.GetConfig().MerchantAccount)

	webhooksMerchant, _, err := d.client.Management().WebhooksMerchantLevelApi.ListAllWebhooks(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Adyen Webhooks",
			err.Error(),
		)
		return
	}

	for _, webhookMerchantData := range webhooksMerchant.Data {
		links := webhookMerchantData.Links
		totalItems := webhooksMerchant.ItemsTotal
		pagesTotal := webhooksMerchant.PagesTotal

		webhookState := webhooksModel{
			Links: webhookLinksModel{
				Self: webhooksLinksHrefModel{
					Href: types.StringValue(*links.Self.Href),
				},
			},
			ItemsTotal:       types.Int64Value(int64(totalItems)),
			PagesTotal:       types.Int64Value(int64(pagesTotal)),
			AccountReference: types.StringValue(*webhooksMerchant.AccountReference),
		}

		for _, webhookData := range webhooksMerchant.Data {
			webhookState.Data = append(webhookState.Data, webhookDataModel{
				ID:                              types.StringValue(*webhookData.Id),
				Type:                            types.StringValue(webhookData.Type),
				URL:                             types.StringValue(webhookData.Url),
				Username:                        types.StringValue(*webhookData.Username),
				Description:                     types.StringValue(*webhookData.Description),
				HasPassword:                     types.BoolValue(*webhookData.HasPassword),
				Active:                          types.BoolValue(webhookData.Active),
				HasError:                        types.BoolValue(*webhookData.HasError),
				EncryptionProtocol:              types.StringValue(*webhookData.EncryptionProtocol),
				CommunicationFormat:             types.StringValue(webhookData.CommunicationFormat),
				AcceptsExpiredCertificate:       types.BoolValue(*webhookData.AcceptsExpiredCertificate),
				AcceptsSelfSignedCertificate:    types.BoolValue(*webhookData.AcceptsSelfSignedCertificate),
				AcceptsUntrustedRootCertificate: types.BoolValue(*webhookData.AcceptsUntrustedRootCertificate),
				CertificateAlias:                types.StringValue(*webhookData.CertificateAlias),
				PopulateSoapActionHeader:        types.BoolValue(*webhookData.PopulateSoapActionHeader),
			})

			for _, additionalSettings := range webhookState.Data {
				webhookState.Data = append(webhookState.Data, webhookDataModel{AdditionalSettings: webhookAdditionalModel{
					IncludeEventCodes: additionalSettings.AdditionalSettings.IncludeEventCodes,
					ExcludeEventCodes: additionalSettings.AdditionalSettings.ExcludeEventCodes,
				}})
			}

			webhookState.Data = append(webhookState.Data, webhookDataModel{
				Links: []webhookDataLinksModel{
					{
						Self: webhooksLinksHrefModel{
							Href: types.StringValue(*links.Self.Href),
						},
						GenerateHmac: webhooksLinksHrefModel{Href: types.StringValue(*links.GenerateHmac.Href)},
						Merchant:     webhooksLinksHrefModel{Href: types.StringValue(*links.Merchant.Href)},
						TestWebhook:  webhooksLinksHrefModel{Href: types.StringValue(*links.TestWebhook.Href)},
					},
				},
			})

		}
		state.Webhooks = append(state.Webhooks, webhookState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
