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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &InventorySourceDataSource{}

func NewInventorySourceDataSource() datasource.DataSource {
	return &InventorySourceDataSource{}
}

// InventorySourceDataSource defines the data source implementation.
type InventorySourceDataSource struct {
	client *AwxClient
}

func (d *InventorySourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inventory_source"
}

func (d *InventorySourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get inventory_source datasource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Inventory Source ID.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Inventory Source name.",
				Computed:    true,
			},
			"inventory": schema.Int32Attribute{
				Description: "Inventory ID for the inventory source to be attached to.",
				Computed:    true,
			},
			"source": schema.StringAttribute{
				Description: "Type of SCM resource. Options: `scm`, `ec2`, `gce`, `azure_rm`, `vmware`, `satellite6`, `openstack`, `rhv`, `controller`, `insights`, `terraform`, `openshift_virtualization`.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "InventorySource description.",
				Computed:    true,
			},
			"credential": schema.Int32Attribute{
				Description: "Inventory source credential ID.",
				Computed:    true,
			},
			"execution_environment": schema.Int32Attribute{
				Description: "The ID of the execution environment this inventory source.",
				Computed:    true,
			},
			"source_path": schema.StringAttribute{
				Description: "(Inventory file) - The inventory file to be synced by this source.",
				Computed:    true,
			},
			"enabled_value": schema.StringAttribute{
				Description: "This field is ignored unless an Enabled Variable is set. If the enabled variable matches this value, the host will be enabled on import.",
				Computed:    true,
			},
			"enabled_var": schema.StringAttribute{
				Description: "Retrieve the enabled state from the given dict of host variables. The enabled variable may be specified using dot notation, e.g: 'foo.bar'",
				Computed:    true,
			},
			"host_filter": schema.StringAttribute{
				Description: "Regular expression where only matching host names will be imported. The filter is applied as a post-processing step after any inventory plugin filters are applied.",
				Computed:    true,
			},
			"overwrite": schema.BoolAttribute{
				Description: "If checked, any hosts and groups that were previously present on the external source but are now removed will be removed from the inventory. Hosts and groups that were not managed by the inventory source will be promoted to the next manually created group or if there is no manually created group to promote them into, they will be left in the `all` default group for the inventory. When not checked, local child hosts and groups not found on the external source will remain untouched by the inventory update process.",
				Computed:    true,
			},
			"overwrite_vars": schema.BoolAttribute{
				Description: "If checked, all variables for child groups and hosts will be removed and replaced by those found on the external source. When not checked, a merge will be performed, combining local variables with those found on the external source.",
				Computed:    true,
			},
			"update_on_launch": schema.BoolAttribute{
				Description: "Each time a job runs using this inventory, refresh the inventory from the selected source before executing job tasks.",
				Computed:    true,
			},
			"source_vars": schema.StringAttribute{
				Description: "Default value is `\"---\"`",
				Computed:    true,
			},
			"source_project": schema.Int32Attribute{
				Description: "The ID of the source project.",
				Computed:    true,
			},
			"scm_branch": schema.StringAttribute{
				Description: "Branch to use on inventory sync. Project default used if blank. Only allowed if project allow_override field is set to true.",
				Computed:    true,
			},
			"update_cache_timeout": schema.Int32Attribute{
				Description: "Time in seconds to consider an inventory sync to be current. During job runs and callbacks the task system will evaluate the timestamp of the latest sync. If it is older than Cache Timeout, it is not considered current, and a new inventory sync will be performed.",
				Computed:    true,
			},
			"verbosity": schema.Int32Attribute{
				Description: "Control the level of output Ansible will produce for inventory source update jobs. `0 - Warning`, `1 - Info`, `2 - Debug`",
				Computed:    true,
			},
		},
	}
}

func (d *InventorySourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InventorySourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InventorySourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var url string

	// set url for read by id HTTP request
	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int.",
			fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
		return
	}
	url = fmt.Sprintf("/api/v2/inventory_sources/%d/", id)

	httpResp, err := d.client.MakeHTTPRequestToAPI(ctx, http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"error making API http request",
			fmt.Sprintf("Error was: %s.", err.Error()))
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

	var responseData InventorySourceAPIModel

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unmarshall response body into object",
			fmt.Sprintf("Error =  %v.", err.Error()))
		return
	}

	idAsString := strconv.Itoa(responseData.Id)
	data.Id = types.StringValue(idAsString)

	data.Name = types.StringValue(responseData.Name)
	data.Inventory = types.Int32Value(int32(responseData.Inventory))
	data.Source = types.StringValue(responseData.Source)
	data.Verbosity = types.Int32Value(int32(responseData.Verbosity))

	if responseData.Credential != 0 {
		data.Credential = types.Int32Value(int32(responseData.Credential))
	}
	if responseData.Description != "" {
		data.Description = types.StringValue(responseData.Description)
	}
	if responseData.ExecutionEnvironment != 0 {
		data.ExecutionEnvironment = types.Int32Value(int32(responseData.ExecutionEnvironment))
	}
	if responseData.SourcePath != "" {
		data.SourcePath = types.StringValue(responseData.SourcePath)
	}
	if responseData.EnabledValue != "" {
		data.EnabledValue = types.StringValue(responseData.EnabledValue)
	}
	if responseData.EnabledVar != "" {
		data.EnabledVar = types.StringValue(responseData.EnabledVar)
	}
	if responseData.HostFilter != "" {
		data.HostFilter = types.StringValue(responseData.HostFilter)
	}
	if responseData.OverwriteVars {
		data.OverwriteVars = types.BoolValue(responseData.OverwriteVars)
	}
	if responseData.Overwrite {
		data.Overwrite = types.BoolValue(responseData.Overwrite)
	}
	if responseData.SourceVars != "" {
		data.SourceVars = types.StringValue(responseData.SourceVars)
	}
	if responseData.SourceProject != 0 {
		data.SourceProject = types.Int32Value(int32(responseData.SourceProject))
	}
	if responseData.ScmBranch != "" {
		data.ScmBranch = types.StringValue(responseData.ScmBranch)
	}
	if responseData.UpdateCacheTimeout != 0 {
		data.UpdateCacheTimeout = types.Int32Value(int32(responseData.UpdateCacheTimeout))
	}
	if responseData.UpdateOnLaunch {
		data.UpdateOnLaunch = types.BoolValue(responseData.UpdateOnLaunch)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
