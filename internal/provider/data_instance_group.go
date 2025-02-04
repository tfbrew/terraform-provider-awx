package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urlParser "net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &InstanceGroupDataSource{}

func NewInstanceGroupDataSource() datasource.DataSource {
	return &InstanceGroupDataSource{}
}

// InstanceGroupDataSource defines the data source implementation.
type InstanceGroupDataSource struct {
	client *AwxClient
}

// InstanceGroupDataSourceModel describes the data source data model.
type InstanceGroupDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	MaxConcurrentJobs        types.Int32  `tfsdk:"max_concurrent_jobs"`
	MaxForks                 types.Int32  `tfsdk:"max_forks"`
	PolicyInstancePercentage types.Int32  `tfsdk:"policy_instance_percentage"`
	PolicyInstanceMinimum    types.Int32  `tfsdk:"policy_instance_minimum"`
}

type InstanceGroupDataSourceJson struct {
	Id                       int    `json:"id"`
	Name                     string `json:"name"`
	MaxConcurrentJobs        int    `json:"max_concurrent_jobs"`
	MaxForks                 int    `json:"max_forks"`
	PolicyInstancePercentage int    `json:"policy_instance_percentage"`
	PolicyInstanceMinimum    int    `json:"policy_instance_minimum"`
}

func (d *InstanceGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance_group"
}

func (d *InstanceGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get Instance Groups.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Instance Group ID.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Instance Group name.",
				Optional:    true,
			},
			"max_concurrent_jobs": schema.Int32Attribute{
				Description: "Maximum number of jobs to run concurrently on this group. Zero means no limit will be enforced.",
				Computed:    true,
			},
			"max_forks": schema.Int32Attribute{
				Description: "Maximum number of forks to allow across all jobs running concurrently on this group. Zero means no limit will be enforced.",
				Computed:    true,
			},
			"policy_instance_percentage": schema.Int32Attribute{
				Description: "Minimum percentage of all instances that will be automatically assigned to this group when new instances come online.",
				Computed:    true,
			},
			"policy_instance_minimum": schema.Int32Attribute{
				Description: "Minimum number of instances that will be automatically assigned to this group when new instances come online.",
				Computed:    true,
			},
		},
	}
}

func (d InstanceGroupDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

func (d *InstanceGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *InstanceGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InstanceGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var url string

	if !data.Id.IsNull() {
		// set url for read by id HTTP request
		id, err := strconv.Atoi(data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable convert id from string to int.",
				fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
			return
		}
		url = fmt.Sprintf("/api/v2/instance_groups/%d/", id)
	}
	if !data.Name.IsNull() {
		// set url for read by name HTTP request
		name := urlParser.QueryEscape(data.Name.ValueString())
		url = fmt.Sprintf("/api/v2/instance_groups/?name=%s", name)
	}

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

	var responseData InstanceGroupDataSourceJson

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
			Count   int                           `json:"count"`
			Results []InstanceGroupDataSourceJson `json:"results"`
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
				"Incorrect number of instance_groups returned by name",
				fmt.Sprintf("Unable to read instance_group as API returned %v instance_groups.", nameResult.Count))
			return
		}
	}

	idAsString := strconv.Itoa(responseData.Id)
	data.Id = types.StringValue(idAsString)

	data.Name = types.StringValue(responseData.Name)
	data.MaxConcurrentJobs = types.Int32Value(int32(responseData.MaxConcurrentJobs))
	data.MaxForks = types.Int32Value(int32(responseData.MaxForks))
	data.PolicyInstancePercentage = types.Int32Value(int32(responseData.PolicyInstancePercentage))
	data.PolicyInstanceMinimum = types.Int32Value(int32(responseData.PolicyInstanceMinimum))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
