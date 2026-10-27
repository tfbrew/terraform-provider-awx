package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	// AAP-AWX-DIFF-BEGIN: aap-only imports for the gateway and EDA IDs and the controller-only fields.
	// AAP-AWX-DIFF-END
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	// AAP-AWX-DIFF-BEGIN: aap-only imports for the gateway and EDA IDs and the controller-only fields.
	// AAP-AWX-DIFF-END
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	// AAP-AWX-DIFF-BEGIN: aap-only imports for the gateway and EDA IDs and the controller-only fields.
	// AAP-AWX-DIFF-END
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	// AAP-AWX-DIFF-BEGIN: aap-only imports for the gateway and EDA IDs and the controller-only fields.
	// AAP-AWX-DIFF-END
)

var _ resource.Resource = &OrganizationResource{}
var _ resource.ResourceWithImportState = &OrganizationResource{}

func NewOrganizationResource() resource.Resource {
	return &OrganizationResource{}
}

type OrganizationResource struct {
	client *providerClient
}

func (r *OrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (r *OrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manage an Automation Controller organization.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				// AAP-AWX-DIFF-BEGIN: aap documents that id is the controller ID.
				Computed: true,
				// AAP-AWX-DIFF-END
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// AAP-AWX-DIFF-BEGIN: aap also exposes the gateway and EDA organization IDs.
			// AAP-AWX-DIFF-END
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the organization.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Organization description.",
			},
			"max_hosts": schema.Int32Attribute{
				Optional:    true,
				Description: "Maximum number of hosts allowed to be managed by this organization.",
				Default:     int32default.StaticInt32(0),
				Computed:    true,
			},
			// AAP-AWX-DIFF-BEGIN: aap manages opa_query_path; awx manages default_environment (AAP 2.5+ dropped the organization default execution environment).
			"default_environment": schema.Int32Attribute{
				Optional:    true,
				Description: "The fallback execution environment that will be used for jobs inside of this organization if not explicitly assigned at the project, job template or workflow level.",
			},
			// AAP-AWX-DIFF-END
		},
	}
}

func (r *OrganizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	configureData, ok := req.ProviderData.(*providerClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = configureData
}

func (r *OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OrganizationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var bodyData OrganizationAPIModel

	if !(data.Name.IsNull()) {
		bodyData.Name = data.Name.ValueString()
	}
	if !(data.Description.IsNull()) {
		bodyData.Description = data.Description.ValueString()
	}
	if !(data.MaxHosts.IsNull()) {
		bodyData.MaxHosts = int(data.MaxHosts.ValueInt32())
	}
	// AAP-AWX-DIFF-BEGIN: awx also sends default_environment.
	if !(data.DefaultEnv.IsNull()) {
		bodyData.DefaultEnv = int(data.DefaultEnv.ValueInt32())
	}
	// AAP-AWX-DIFF-END

	url := "organizations/"
	// AAP-AWX-DIFF-BEGIN: aap creates organizations through the gateway API; awx uses /api/v2/.
	returnedData, _, err := r.client.CreateUpdateAPIRequest(ctx, http.MethodPost, url, bodyData, []int{201}, "")
	// AAP-AWX-DIFF-END
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making API http request",
			fmt.Sprintf("Error was: %s.", err.Error()))
		return
	}

	// AAP-AWX-DIFF-BEGIN: aap looks up the controller and EDA IDs of the new gateway organization (with retries); awx uses the returned ID.
	data.Id = types.StringValue(fmt.Sprintf("%v", returnedData["id"]))
	// AAP-AWX-DIFF-END

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	// AAP-AWX-DIFF-BEGIN: aap sets the controller-only fields (max_hosts, opa_query_path) through the controller API, since the gateway API ignores them.
	// AAP-AWX-DIFF-END
}

func (r *OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OrganizationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v.", data.Id))
		return
	}

	url := fmt.Sprintf("organizations/%d/", id)
	body, statusCode, err := r.client.GenericAPIRequest(ctx, http.MethodGet, url, nil, []int{200, 404}, "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making API http request",
			fmt.Sprintf("Error was: %s.", err.Error()))
		return
	}

	if statusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	var responseData OrganizationAPIModel

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unmarshal json",
			fmt.Sprintf("bodyData: %+v.", body))
		return
	}

	if !data.Name.IsNull() || responseData.Name != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), responseData.Name)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !data.Description.IsNull() || responseData.Description != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), responseData.Description)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("max_hosts"), responseData.MaxHosts)...)
	// AAP-AWX-DIFF-BEGIN: aap reads opa_query_path and looks up the gateway and EDA organization IDs; awx reads default_environment.

	if !data.DefaultEnv.IsNull() || responseData.DefaultEnv != 0 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("default_environment"), responseData.DefaultEnv)...)
	}
	// AAP-AWX-DIFF-END
}

func (r *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data OrganizationModel

	// AAP-AWX-DIFF-BEGIN: aap reads the prior state to find changed controller-only fields.
	// AAP-AWX-DIFF-END
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	// AAP-AWX-DIFF-BEGIN: aap reads the prior state to find changed controller-only fields.
	// AAP-AWX-DIFF-END
	if resp.Diagnostics.HasError() {
		return
	}

	// AAP-AWX-DIFF-BEGIN: aap updates and deletes organizations by their gateway ID; awx uses the /api/v2/ ID.
	id, err := strconv.Atoi(data.Id.ValueString())
	// AAP-AWX-DIFF-END
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v.", data.Id))
		return
	}

	var bodyData OrganizationAPIModel

	if !(data.Name.IsNull()) {
		bodyData.Name = data.Name.ValueString()
	}
	if !(data.Description.IsNull()) {
		bodyData.Description = data.Description.ValueString()
	}
	if !(data.MaxHosts.IsNull()) {
		bodyData.MaxHosts = int(data.MaxHosts.ValueInt32())
	}
	// AAP-AWX-DIFF-BEGIN: awx also sends default_environment.
	if !(data.DefaultEnv.IsNull()) {
		bodyData.DefaultEnv = int(data.DefaultEnv.ValueInt32())
	}
	// AAP-AWX-DIFF-END

	url := fmt.Sprintf("organizations/%d/", id)
	// AAP-AWX-DIFF-BEGIN: aap updates organizations through the gateway API; awx uses /api/v2/.
	_, _, err = r.client.CreateUpdateAPIRequest(ctx, http.MethodPut, url, bodyData, []int{200}, "")
	// AAP-AWX-DIFF-END
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making API update request",
			fmt.Sprintf("Error was: %s.", err.Error()))
		return
	}

	// AAP-AWX-DIFF-BEGIN: aap sets the controller-only fields (max_hosts, opa_query_path) through the controller API, since the gateway API ignores them.
	// AAP-AWX-DIFF-END
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	// AAP-AWX-DIFF-BEGIN: aap sets the controller-only fields (max_hosts, opa_query_path) through the controller API, since the gateway API ignores them.
	// AAP-AWX-DIFF-END
}

func (r *OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OrganizationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// AAP-AWX-DIFF-BEGIN: aap updates and deletes organizations by their gateway ID; awx uses the /api/v2/ ID.
	id, err := strconv.Atoi(data.Id.ValueString())
	// AAP-AWX-DIFF-END
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v.", data.Id.ValueString()))
		return
	}
	url := fmt.Sprintf("organizations/%d/", id)

	// AAP-AWX-DIFF-BEGIN: aap deletes organizations through the gateway API; awx uses /api/v2/.
	_, _, err = r.client.GenericAPIRequest(ctx, http.MethodDelete, url, nil, []int{202, 204}, "")
	// AAP-AWX-DIFF-END
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making API delete request",
			fmt.Sprintf("Error was: %s.", err.Error()))
		return
	}
}

func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
