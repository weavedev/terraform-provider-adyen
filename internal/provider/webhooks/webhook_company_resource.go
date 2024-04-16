package webhooks

import (
	"context"
	"fmt"
	"github.com/adyen/adyen-go-api-library/v9/src/adyen"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &webhooksCompanyResource{}
	_ resource.ResourceWithConfigure = &webhooksCompanyResource{}
)

// NewWebhooksCompanyResource is a helper function to simplify the provider implementation.
func NewWebhooksCompanyResource() resource.Resource {
	return &webhooksCompanyResource{}
}

// webhooksCompanyResource is the resource implementation.
type webhooksCompanyResource struct {
	client *adyen.APIClient
}

// Metadata returns the resource type name.
func (r *webhooksCompanyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhooks_company"
}

// Configure adds the provider configured client to the resource.
func (r *webhooksCompanyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*adyen.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *webhooksCompanyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *webhooksCompanyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

// Read refreshes the Terraform state with the latest data.
func (r *webhooksCompanyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *webhooksCompanyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *webhooksCompanyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
