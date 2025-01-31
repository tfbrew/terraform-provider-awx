package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urlParser "net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &JobTemplateDataSource{}

func NewJobTemplateDataSource() datasource.DataSource {
	return &JobTemplateDataSource{}
}

// JobTemplateDataSource defines the data source implementation.
type JobTemplateDataSource struct {
	client *AwxClient
}

func (d *JobTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job_template"
}

func (d *JobTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get job_template datasource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"job_type": schema.StringAttribute{
				Description: "Acceptable values are a choice of: `run`, `check`.",
				Computed:    true,
			},
			"inventory": schema.Int32Attribute{
				Computed:    true,
				Description: "ID number of the inventory to associate with the job template. Supply this or set `ask_inventory_on_launch = true`.",
			},
			"project": schema.Int32Attribute{
				Computed:    true,
				Description: "ID number of the project to associate with the job template",
			},
			"playbook": schema.StringAttribute{
				Computed:    true,
				Description: "Playbook name to be executed by this job",
			},
			"scm_branch": schema.StringAttribute{
				Computed: true,
			},
			"forks": schema.Int32Attribute{
				Computed: true,
			},
			"limit": schema.StringAttribute{
				Computed: true,
			},
			"verbosity": schema.Int32Attribute{
				Computed: true,
			},
			"extra_vars": schema.StringAttribute{
				Computed: true,
			},
			"job_tags": schema.StringAttribute{
				Computed: true,
			},
			"force_handlers": schema.BoolAttribute{
				Computed: true,
			},
			"skip_tags": schema.StringAttribute{
				Computed: true,
			},
			"start_at_tags": schema.StringAttribute{
				Computed: true,
			},
			"timeout": schema.Int32Attribute{
				Computed: true,
			},
			"use_fact_cache": schema.BoolAttribute{
				Computed: true,
			},
			"execution_environment": schema.Int32Attribute{
				Computed: true,
			},
			"host_config_key": schema.StringAttribute{
				Computed: true,
			},
			"ask_scm_branch_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_diff_mode_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_variables_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_limit_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_tags_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_skip_tags_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_job_type_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_verbosity_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_inventory_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_credential_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_execution_environment_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_labels_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_forks_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_job_slice_count_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_timeout_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"ask_instance_groups_on_launch": schema.BoolAttribute{
				Computed: true,
			},
			"survey_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"become_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"diff_mode": schema.BoolAttribute{
				Computed: true,
			},
			"allow_simultaneous": schema.BoolAttribute{
				Computed: true,
			},
			"custom_virtualenv": schema.StringAttribute{
				Computed: true,
			},
			"job_slice_count": schema.Int32Attribute{
				Computed: true,
			},
			"webhook_service": schema.StringAttribute{
				Computed: true,
			},
			"webhook_credential": schema.StringAttribute{
				Computed: true,
			},
			"prevent_instance_group_fallback": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (d *JobTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *JobTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data JobTemplateModel

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
		url = fmt.Sprintf("/api/v2/job_templates/%d/", id)
	}
	if !data.Name.IsNull() {
		// set url for read by name HTTP request
		name := urlParser.QueryEscape(data.Name.ValueString())
		url = fmt.Sprintf("/api/v2/job_templates/?name=%s", name)
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

	var responseData JobTemplateAPIModel

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
			Count   int                   `json:"count"`
			Results []JobTemplateAPIModel `json:"results"`
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
				"Incorrect number of execution_environments returned by name",
				fmt.Sprintf("Unable to read execution_environment as API returned %v execution_environments.", nameResult.Count))
			return
		}
	}

	idAsString := strconv.Itoa(responseData.Id)
	data.Id = types.StringValue(idAsString)

	data.Name = types.StringValue(responseData.Name)

	if responseData.Description != "" {
		data.Description = types.StringValue(responseData.Description)
	}
	if responseData.JobType != "" {
		data.JobType = types.StringValue(responseData.JobType)
	}
	if responseData.Inventory != 0 {
		data.Inventory = types.Int32Value(int32(responseData.Inventory))
	}
	if responseData.Project != 0 {
		data.Project = types.Int32Value(int32(responseData.Project))
	}
	if responseData.Playbook != "" {
		data.Playbook = types.StringValue(responseData.Playbook)
	}
	if responseData.ScmBranch != "" {
		data.ScmBranch = types.StringValue(responseData.ScmBranch)
	}
	if responseData.Forks != 0 {
		data.Forks = types.Int32Value(int32(responseData.Forks))
	}
	if responseData.Limit != "" {
		data.Limit = types.StringValue(responseData.Limit)
	}
	if responseData.Verbosity != 0 {
		data.Verbosity = types.Int32Value(int32(responseData.Verbosity))
	}
	if responseData.ExtraVars != "" {
		data.ExtraVars = types.StringValue(responseData.ExtraVars)
	}
	if responseData.JobTags != "" {
		data.JobTags = types.StringValue(responseData.JobTags)
	}
	if responseData.ForceHandlers {
		data.ForceHandlers = types.BoolValue(responseData.ForceHandlers)
	}
	if responseData.SkipTags != "" {
		data.SkipTags = types.StringValue(responseData.SkipTags)
	}
	if responseData.StartAtTask != "" {
		data.StartAtTask = types.StringValue(responseData.StartAtTask)
	}
	if responseData.Timeout != 0 {
		data.Timeout = types.Int32Value(int32(responseData.Timeout))
	}
	if responseData.UseFactCache {
		data.UseFactCache = types.BoolValue(responseData.UseFactCache)
	}
	if responseData.ExecutionEnvironment != 0 {
		data.ExecutionEnvironment = types.Int32Value(int32(responseData.ExecutionEnvironment))
	}
	if responseData.HostConfigKey != "" {
		data.HostConfigKey = types.StringValue(responseData.HostConfigKey)
	}
	if responseData.AskScmBranchOnLaunch {
		data.AskScmBranchOnLaunch = types.BoolValue(responseData.AskScmBranchOnLaunch)
	}
	if responseData.AskDiffModeOnLaunch {
		data.AskDiffModeOnLaunch = types.BoolValue(responseData.AskDiffModeOnLaunch)
	}
	if responseData.AskVariablesOnLaunch {
		data.AskVariablesOnLaunch = types.BoolValue(responseData.AskVariablesOnLaunch)
	}
	if responseData.AskLimitOnLaunch {
		data.AskLimitOnLaunch = types.BoolValue(responseData.AskLimitOnLaunch)
	}
	if responseData.AskTagsOnLaunch {
		data.AskTagsOnLaunch = types.BoolValue(responseData.AskTagsOnLaunch)
	}
	if responseData.AskSkipTagsOnLaunch {
		data.AskSkipTagsOnLaunch = types.BoolValue(responseData.AskSkipTagsOnLaunch)
	}
	if responseData.AskJobTypeOnLaunch {
		data.AskJobTypeOnLaunch = types.BoolValue(responseData.AskJobTypeOnLaunch)
	}
	if responseData.AskVerbosityOnLaunch {
		data.AskVerbosityOnLaunch = types.BoolValue(responseData.AskVerbosityOnLaunch)
	}
	if responseData.AskInventoryOnLaunch {
		data.AskInventoryOnLaunch = types.BoolValue(responseData.AskInventoryOnLaunch)
	}
	if responseData.AskCredentialOnLaunch {
		data.AskCredentialOnLaunch = types.BoolValue(responseData.AskCredentialOnLaunch)
	}
	if responseData.AskExecutionEnvironmenOnLaunch {
		data.AskExecutionEnvironmenOnLaunch = types.BoolValue(responseData.AskExecutionEnvironmenOnLaunch)
	}
	if responseData.AskLablesOnLaunch {
		data.AskLablesOnLaunch = types.BoolValue(responseData.AskLablesOnLaunch)
	}
	if responseData.AskForksOnLaunch {
		data.AskForksOnLaunch = types.BoolValue(responseData.AskForksOnLaunch)
	}
	if responseData.AskJobSliceCountOnLaunch {
		data.AskJobSliceCountOnLaunch = types.BoolValue(responseData.AskJobSliceCountOnLaunch)
	}
	if responseData.AskTimeoutOnLaunch {
		data.AskTimeoutOnLaunch = types.BoolValue(responseData.AskTimeoutOnLaunch)
	}
	if responseData.AskInstanceGroupsOnLaunch {
		data.AskInstanceGroupsOnLaunch = types.BoolValue(responseData.AskInstanceGroupsOnLaunch)
	}
	if responseData.SurveyEnabled {
		data.SurveyEnabled = types.BoolValue(responseData.SurveyEnabled)
	}
	if responseData.BecomeEnabled {
		data.BecomeEnabled = types.BoolValue(responseData.BecomeEnabled)
	}
	if responseData.DiffMode {
		data.DiffMode = types.BoolValue(responseData.DiffMode)
	}
	if responseData.AllowSimultaneous {
		data.AllowSimultaneous = types.BoolValue(responseData.AllowSimultaneous)
	}
	if responseData.CustomVirtualEnv == nil {
		data.CustomVirtualEnv = types.StringValue("")
	} else if responseData.CustomVirtualEnv != "" {
		if stringVirtualEnv, ok := responseData.CustomVirtualEnv.(string); ok {
			data.CustomVirtualEnv = types.StringValue(stringVirtualEnv)
		} else {
			resp.Diagnostics.AddError(
				"Invalid Type",
				"Expected responseData.CustomVirtualEnv to be a string",
			)
			return
		}
	}
	if responseData.JobSliceCount != 0 {
		data.JobSliceCount = types.Int32Value(int32(responseData.JobSliceCount))
	}
	if responseData.WebhookService != "" {
		data.WebhookService = types.StringValue(responseData.WebhookService)
	}
	if responseData.WebhookCredential == nil {
		data.WebhookCredential = types.StringValue("")
	} else if responseData.WebhookCredential != "" {
		if stringVirtualEnv, ok := responseData.WebhookCredential.(string); ok {
			data.WebhookCredential = types.StringValue(stringVirtualEnv)
		} else {
			resp.Diagnostics.AddError(
				"Invalid Type",
				"Expected responseData.WebhookCredential to be a string",
			)
			return
		}
	}
	if responseData.PreventInstanceGroupFallback {
		data.PreventInstanceGroupFallback = types.BoolValue(responseData.PreventInstanceGroupFallback)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
