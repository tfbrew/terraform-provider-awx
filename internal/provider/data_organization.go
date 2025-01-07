package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &OrganizationDataSource{}

func NewOrganizationDataSource() datasource.DataSource {
	return &OrganizationDataSource{}
}

// OrganizationDataSource defines the data source implementation.
type OrganizationDataSource struct {
	client *AwxClient
}

// OrganizationDataSourceModel describes the data source data model.
type OrganizationDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	CustomVirtualEnv types.String `tfsdk:"custom_virtualenv"`
	DefaultEnv       types.Int32  `tfsdk:"default_environment"`
	MaxHosts         types.Int32  `tfsdk:"max_hosts"`
}

type OrganizationDataSourceJson struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	CustomVirtualEnv string `json:"custom_virtualenv"`
	DefaultEnv       int    `json:"default_environment"`
	MaxHosts         int    `json:"max_hosts"`
}

func (d *OrganizationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (d *OrganizationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get organization datasource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Organization ID",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Organization Name",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Organization description.",
				Computed:    true,
			},
			"custom_virtualenv": schema.StringAttribute{
				Description: "Local absolute file path containing a custom Python virtualenv to use.",
				Computed:    true,
			},
			"default_environment": schema.Int32Attribute{
				Description: "The fallback execution environment that will be used for jobs inside of this organization if not explicitly assigned at the project, job template or workflow level.",
				Computed:    true,
			},
			"max_hosts": schema.Int32Attribute{
				Description: "Maximum number of hosts allowed to be managed by this organization.",
				Computed:    true,
			},
		},
	}
}

func (d *OrganizationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	configureData, ok := req.ProviderData.(*AwxClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = configureData
}

func (d *OrganizationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var url string

	if !data.Id.IsNull() && !data.Name.IsNull() {
		resp.Diagnostics.AddError(
			"Too many inputs for datasource.",
			fmt.Sprintf("Both id: %v and name: %v were provided for organization datasource. Please provide only one.", data.Id.ValueString(), data.Name.ValueString()))
		return
	} else if !data.Id.IsNull() {
		// set url for read by id HTTP request
		id, err := strconv.Atoi(data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable convert id from string to int.",
				fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
			return
		}
		url = d.client.endpoint + fmt.Sprintf("/api/v2/organizations/%d/", id)
	} else if !data.Name.IsNull() {
		// set url for read by name HTTP request
		name := data.Name.ValueString()
		url = d.client.endpoint + fmt.Sprintf("/api/v2/organizations/?name=%s", name)
	} else {
		resp.Diagnostics.AddError(
			"ID or Name of organization required.",
			"Neither id or name of organization provided.")
		return
	}

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v.", url))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+d.client.token)

	httpResp, err := d.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read organization, got error: %s", err))
		return
	}
	if httpResp.StatusCode != 200 && httpResp.StatusCode != 404 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))
		return
	}

	if httpResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read the http response data body",
			fmt.Sprintf("Body: %v.", body))
		return
	}

	var responseData OrganizationDataSourceJson

	if !data.Id.IsNull() && data.Name.IsNull() {
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to unmarshall response body into object",
				fmt.Sprintf("Error =  %v.", err.Error()))
			return
		}
	}
	// If looking up by name, check that there is only one response and extract it.
	if data.Id.IsNull() && !data.Name.IsNull() {
		nameResult := struct {
			Count   int                          `json:"count"`
			Results []OrganizationDataSourceJson `json:"results"`
		}{}
		err = json.Unmarshal(body, &nameResult)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to unmarshall response body into object",
				fmt.Sprintf("Error:  %v.", err.Error()))
			return
		}
		if nameResult.Count == 1 {
			responseData = nameResult.Results[0]
		} else {
			resp.Diagnostics.AddError(
				"Incorrect number of organizations returned by name",
				fmt.Sprintf("Unable to read organization as API returned %v organizations.", nameResult.Count))
			return
		}
	}

	if !(data.Id.IsNull() && responseData.Id == 0) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), fmt.Sprintf("%v", responseData.Id))...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !(data.Name.IsNull() && responseData.Name == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), responseData.Name)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !(data.Description.IsNull() && responseData.Description == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), responseData.Description)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !(data.CustomVirtualEnv.IsNull() && responseData.CustomVirtualEnv == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("custom_virtualenv"), responseData.CustomVirtualEnv)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !(data.DefaultEnv.IsNull() && responseData.DefaultEnv == 0) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("default_environment"), responseData.DefaultEnv)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("max_hosts"), responseData.MaxHosts)...)

}
