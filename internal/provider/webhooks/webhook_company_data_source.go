// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package webhooks

import (
	"context"
	"fmt"
	"github.com/adyen/adyen-go-api-library/v9/src/adyen"
	"github.com/adyen/adyen-go-api-library/v9/src/management"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource              = &webhooksCompanyDataSource{}
	_ datasource.DataSourceWithConfigure = &webhooksCompanyDataSource{}
)

// webhookDataSource defines the data source implementation.
type webhooksCompanyDataSource struct {
	client         *adyen.APIClient
	companyAccount string `tfsdk:"company_account"`
}

// webhookDataSource defines the data source implementation.
type webhooksDataSourceModel struct {
	Webhooks []management.Webhook `tfsdk:"webhooks_company"`
}

func NewWebhookCompanyDataSource() datasource.DataSource {
	return &webhooksCompanyDataSource{}
}

func (d *webhooksCompanyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhooks_company"
}

func (d *webhooksCompanyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *webhooksCompanyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"webhooks_company": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"live": schema.StringAttribute{
							Computed:    true,
							Description: "",
						},
						"notification_items": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"notification_request_item": schema.SingleNestedAttribute{
										Computed: true,
										Attributes: map[string]schema.Attribute{
											"additional_data": schema.SingleNestedAttribute{
												Attributes: map[string]schema.Attribute{},
												Computed:   true,
											},
											"amount": schema.SingleNestedAttribute{
												Computed: true,
												Attributes: map[string]schema.Attribute{
													"currency": schema.StringAttribute{
														Computed: true,
													},
													"value": schema.Int64Attribute{
														Computed: true,
													},
												},
											},
											"event_code": schema.StringAttribute{
												Computed: true,
											},
											"event_date": schema.StringAttribute{
												Computed: true,
											},
											"merchant_account_code": schema.StringAttribute{
												Computed: true,
											},
											"merchant_reference": schema.StringAttribute{
												Computed: true,
											},
											"operations": schema.StringAttribute{
												Computed: true,
											},
											"original_reference": schema.StringAttribute{
												Computed: true,
											},
											"payment_method": schema.StringAttribute{
												Computed: true,
											},
											"psp_reference": schema.StringAttribute{
												Computed: true,
											},
											"reason": schema.StringAttribute{
												Computed: true,
											},
											"success": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
							},
						},
					},
				},
				Optional: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *webhooksCompanyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state webhooksDataSourceModel
	data := d.client.Management().WebhooksCompanyLevelApi.ListAllWebhooksInput(d.companyAccount)

	webhooksCompany, _, err := d.client.Management().WebhooksCompanyLevelApi.ListAllWebhooks(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Adyen Webhooks",
			err.Error(),
		)
		return
	}

	for _, webhooksCompany := range webhooksCompany.Data {
		state.Webhooks = append(state.Webhooks, webhooksCompany)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
