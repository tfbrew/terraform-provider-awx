// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &JobTemplateResource{}
var _ resource.ResourceWithImportState = &JobTemplateResource{}

func NewJobTemplateResource() resource.Resource {
	return &JobTemplateResource{}
}

// JobTemplateResource defines the resource implementation.
type JobTemplateResource struct {
	client *AwxClient
}

// JobTemplateResourceModel describes the resource data model.
type JobTemplateResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	Name                           types.String `tfsdk:"name"`
	Description                    types.String `tfsdk:"description"`
	JobType                        types.String `tfsdk:"job_type"`
	Inventory                      types.Int32  `tfsdk:"inventory"`
	Project                        types.Int32  `tfsdk:"project"`
	Playbook                       types.String `tfsdk:"playbook"`
	ScmBranch                      types.String `tfsdk:"scm_branch"`
	Forks                          types.Int32  `tfsdk:"forks"`
	Limit                          types.String `tfsdk:"limit"`
	Verbosity                      types.Int32  `tfsdk:"verbosity"`
	ExtraVars                      types.String `tfsdk:"extra_vars"`
	JobTags                        types.String `tfsdk:"job_tags"`
	ForceHandlers                  types.Bool   `tfsdk:"force_handlers"`
	SkipTags                       types.String `tfsdk:"skip_tags"`
	StartAtTask                    types.String `tfsdk:"start_at_tags"`
	Timeout                        types.Int32  `tfsdk:"timeout"`
	UseFactCache                   types.Bool   `tfsdk:"use_fact_cache"`
	Organization                   types.Int32  `tfsdk:"organization"`
	Status                         types.String `tfsdk:"status"`
	ExecutionEnvironment           types.Int32  `tfsdk:"execution_environment"`
	HostConfigKey                  types.String `tfsdk:"host_config_key"`
	AskScmBranchOnLaunch           types.Bool   `tfsdk:"ask_scm_branch_on_launch"`
	AskDiffModeOnLaunch            types.Bool   `tfsdk:"ask_diff_mode_on_launch"`
	AskVariablesOnLaunch           types.Bool   `tfsdk:"ask_variables_on_launch"`
	AskLimitOnLaunch               types.Bool   `tfsdk:"ask_limit_on_launch"`
	AskTagsOnLaunch                types.Bool   `tfsdk:"ask_tags_on_launch"`
	AskSkipTagsOnLaunch            types.Bool   `tfsdk:"ask_skip_tags_on_launch"`
	AskJobTypeOnLaunch             types.Bool   `tfsdk:"ask_job_type_on_launch"`
	AskVerbosityOnLaunch           types.Bool   `tfsdk:"ask_verbosity_on_launch"`
	AskInventoryOnLaunch           types.Bool   `tfsdk:"ask_inventory_on_launch"`
	AskCredentialOnLaunch          types.Bool   `tfsdk:"ask_credential_on_launch"`
	AskExecutionEnvironmenOnLaunch types.Bool   `tfsdk:"ask_execution_environment_on_launch"`
	AskLablesOnLaunch              types.Bool   `tfsdk:"ask_labels_on_launch"`
	AskForksOnLaunch               types.Bool   `tfsdk:"ask_forks_on_launch"`
	AskJobSliceCountOnLaunch       types.Bool   `tfsdk:"ask_job_slice_count_on_launch"`
	AskTimeoutOnLaunch             types.Bool   `tfsdk:"ask_timeout_on_launch"`
	AskInstanceGroupsOnLaunch      types.Bool   `tfsdk:"ask_instance_groups_on_launch"`
	SurveyEnabled                  types.Bool   `tfsdk:"survey_enabled"`
	BecomeEnabled                  types.Bool   `tfsdk:"become_enabled"`
	DiffMode                       types.Bool   `tfsdk:"diff_mode"`
	AllowSimultaneous              types.Bool   `tfsdk:"allow_simultaneous"`
	CustomVirtualEnv               types.String `tfsdk:"custom_virtualenv"`
	JobSliceCount                  types.Int32  `tfsdk:"job_slice_count"`
	WebhookService                 types.String `tfsdk:"webhook_service"`
	WebhookCredential              types.String `tfsdk:"webhook_credential"`
	PreventInstanceGroupFallback   types.Bool   `tfsdk:"prevent_instance_group_fallback"`
}

type JobTemplate struct {
	Id                             int    `json:"id,omitempty"`
	Name                           string `json:"name"`
	Description                    string `json:"description"`
	JobType                        string `json:"job_type"`
	Inventory                      int    `json:"inventory"`
	Project                        int    `json:"project"`
	Playbook                       string `json:"playbook"`
	ScmBranch                      string `json:"scm_branch"`
	Forks                          int    `json:"forks"`
	Limit                          string `json:"limit"`
	Verbosity                      int    `json:"verbosity"`
	ExtraVars                      string `json:"extra_vars"`
	JobTags                        string `json:"job_tags"`
	ForceHandlers                  bool   `json:"force_handlers"`
	SkipTags                       string `json:"skip_tags"`
	StartAtTask                    string `json:"start_at_tags"`
	Timeout                        int    `json:"timeout"`
	UseFactCache                   bool   `json:"use_fact_cache"`
	Organization                   int    `json:"organization"`
	Status                         string `json:"status"`
	ExecutionEnvironment           int    `json:"execution_environment"`
	HostConfigKey                  string `json:"host_config_key"`
	AskScmBranchOnLaunch           bool   `json:"ask_scm_branch_on_launch"`
	AskDiffModeOnLaunch            bool   `json:"ask_diff_mode_on_launch"`
	AskVariablesOnLaunch           bool   `json:"ask_variables_on_launch"`
	AskLimitOnLaunch               bool   `json:"ask_limit_on_launch"`
	AskTagsOnLaunch                bool   `json:"ask_tags_on_launch"`
	AskSkipTagsOnLaunch            bool   `json:"ask_skip_tags_on_launch"`
	AskJobTypeOnLaunch             bool   `json:"ask_job_type_on_launch"`
	AskVerbosityOnLaunch           bool   `json:"ask_verbosity_on_launch"`
	AskInventoryOnLaunch           bool   `json:"ask_inventory_on_launch"`
	AskCredentialOnLaunch          bool   `json:"ask_credential_on_launch"`
	AskExecutionEnvironmenOnLaunch bool   `json:"ask_execution_environment_on_launch"`
	AskLablesOnLaunch              bool   `json:"ask_labels_on_launch"`
	AskForksOnLaunch               bool   `json:"ask_forks_on_launch"`
	AskJobSliceCountOnLaunch       bool   `json:"ask_job_slice_count_on_launch"`
	AskTimeoutOnLaunch             bool   `json:"ask_timeout_on_launch"`
	AskInstanceGroupsOnLaunch      bool   `json:"ask_instance_groups_on_launch"`
	SurveyEnabled                  bool   `json:"survey_enabled"`
	BecomeEnabled                  bool   `json:"become_enabled"`
	DiffMode                       bool   `json:"diff_mode"`
	AllowSimultaneous              bool   `json:"allow_simultaneous"`
	CustomVirtualEnv               any    `json:"custom_virtualenv,omitempty"` //blank is returned by api as "custom_virtual": null (not "")
	JobSliceCount                  int    `json:"job_slice_count"`
	WebhookService                 string `json:"webhook_service"`
	WebhookCredential              any    `json:"webhook_credential,omitempty"` //blank is returned by api as "webhook_credentials": null (not "")
	PreventInstanceGroupFallback   bool   `json:"prevent_instance_group_fallback"`
}

func (r *JobTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job_template"
}

func (r *JobTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		//TODO fix description on schema and markdown descr
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"job_type": schema.StringAttribute{
				Required:            true,
				Description:         "Acceptable values are a choice of: run, or check.",
				MarkdownDescription: "Acceptable values are a choice of: run, or check.",
			},
			"inventory": schema.Int32Attribute{
				Optional: true,
			},
			"project": schema.Int32Attribute{
				Optional: true,
			},
			"playbook": schema.StringAttribute{
				Optional: true,
			},
			"scm_branch": schema.StringAttribute{
				Optional: true,
			},
			"forks": schema.Int32Attribute{
				Optional: true,
			},
			"limit": schema.StringAttribute{
				Optional: true,
			},
			"verbosity": schema.Int32Attribute{
				Optional: true,
			},
			"extra_vars": schema.StringAttribute{
				Optional: true,
			},
			"job_tags": schema.StringAttribute{
				Optional: true,
			},
			"force_handlers": schema.BoolAttribute{
				Optional: true,
			},
			"skip_tags": schema.StringAttribute{
				Optional: true,
			},
			"start_at_tags": schema.StringAttribute{
				Optional: true,
			},
			"timeout": schema.Int32Attribute{
				Optional: true,
			},
			"use_fact_cache": schema.BoolAttribute{
				Optional: true,
			},
			"organization": schema.Int32Attribute{
				Optional: true,
			},
			"status": schema.StringAttribute{
				Optional: true,
			},
			"execution_environment": schema.Int32Attribute{
				Optional: true,
			},
			"host_config_key": schema.StringAttribute{
				Optional: true,
			},
			"ask_scm_branch_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_diff_mode_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_variables_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_limit_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_tags_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_skip_tags_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_job_type_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_verbosity_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_inventory_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_credential_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_execution_environment_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_labels_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_forks_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_job_slice_count_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_timeout_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"ask_instance_groups_on_launch": schema.BoolAttribute{
				Optional: true,
			},
			"survey_enabled": schema.BoolAttribute{
				Optional: true,
			},
			"become_enabled": schema.BoolAttribute{
				Optional: true,
			},
			"diff_mode": schema.BoolAttribute{
				Optional: true,
			},
			"allow_simultaneous": schema.BoolAttribute{
				Optional: true,
			},
			"custom_virtualenv": schema.StringAttribute{
				Optional: true,
			},
			"job_slice_count": schema.Int32Attribute{
				Optional: true,
			},
			"webhook_service": schema.StringAttribute{
				Optional: true,
			},
			"webhook_credential": schema.StringAttribute{
				Optional: true,
			},
			"prevent_instance_group_fallback": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (r *JobTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	configureData := req.ProviderData.(*AwxClient)

	r.client = configureData
}

func (r *JobTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data JobTemplateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// // set url for create HTTP request
	// id, err := strconv.Atoi(data.Id.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable convert id from string to int",
	// 		fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
	// }

	// url := r.client.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/survey_spec", id)

	// // get body data for HTTP request
	// var bodyData JobTemplate
	// bodyData.Name = data.Name.ValueString()
	// bodyData.Description = data.Description.ValueString()

	// jsonData, err := json.Marshal(bodyData)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable marshal json",
	// 		fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
	// }

	// // create HTTP request
	// httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to generate request",
	// 		fmt.Sprintf("Unable to gen url: %v. ", url))
	// }

	// httpReq.Header.Add("Content-Type", "application/json")
	// httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	// httpResp, err := r.client.client.Do(httpReq)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	// }
	// if httpResp.StatusCode != 200 {
	// 	resp.Diagnostics.AddError(
	// 		"Bad request status code.",
	// 		fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))

	// }

	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data JobTemplateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// set url for create HTTP request
	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
	}
	url := r.client.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/", id)

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))

	}

	var responseData JobTemplate

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Uanble to get all data out of the http response data body",
			fmt.Sprintf("Body got %v. ", body))
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Uanble unmarshall response body into object",
			fmt.Sprintf("Error =  %v. ", err.Error()))
	}

	//responseID := fmt.Sprint(responseData.Id)
	//data.Id = types.StringValue(responseID)
	data.Name = types.StringValue(responseData.Name)
	data.Description = types.StringValue(responseData.Description)
	data.JobType = types.StringValue(responseData.JobType)
	data.Inventory = types.Int32Value(int32(responseData.Inventory))
	data.Project = types.Int32Value(int32(responseData.Project))
	data.Playbook = types.StringValue(responseData.Playbook)
	data.ScmBranch = types.StringValue(responseData.ScmBranch)
	data.Forks = types.Int32Value(int32(responseData.Forks))
	data.Limit = types.StringValue(responseData.Limit)
	data.Verbosity = types.Int32Value(int32(responseData.Verbosity))
	data.ExtraVars = types.StringValue(responseData.ExtraVars)
	data.JobTags = types.StringValue(responseData.JobTags)
	data.ForceHandlers = types.BoolValue(responseData.ForceHandlers)
	//data.SkipTags = types.StringValue(responseData.SkipTags)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("skip_tags"), responseData.SkipTags)...)
	if resp.Diagnostics.HasError() {
		return
	}
	//resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	//if resp.Diagnostics.HasError() {
	//	return

	//data.StartAtTask = types.StringValue(responseData.StartAtTask)
	//data.StartAtTask.IsNull() & responseData.StartAtTask == ""
	data.Timeout = types.Int32Value(int32(responseData.Timeout))
	data.UseFactCache = types.BoolValue(responseData.UseFactCache)
	data.Organization = types.Int32Value(int32(responseData.Organization))
	data.Status = types.StringValue(responseData.Status)
	data.ExecutionEnvironment = types.Int32Value(int32(responseData.ExecutionEnvironment))
	data.HostConfigKey = types.StringValue(responseData.HostConfigKey)
	data.AskScmBranchOnLaunch = types.BoolValue(responseData.AskScmBranchOnLaunch)
	data.AskDiffModeOnLaunch = types.BoolValue(responseData.AskDiffModeOnLaunch)
	data.AskVariablesOnLaunch = types.BoolValue(responseData.AskVariablesOnLaunch)
	data.AskLimitOnLaunch = types.BoolValue(responseData.AskLimitOnLaunch)
	data.AskTagsOnLaunch = types.BoolValue(responseData.AskTagsOnLaunch)
	data.AskSkipTagsOnLaunch = types.BoolValue(responseData.AskSkipTagsOnLaunch)
	data.AskJobTypeOnLaunch = types.BoolValue(responseData.AskJobTypeOnLaunch)
	data.AskVerbosityOnLaunch = types.BoolValue(responseData.AskVerbosityOnLaunch)
	data.AskInventoryOnLaunch = types.BoolValue(responseData.AskInventoryOnLaunch)
	data.AskCredentialOnLaunch = types.BoolValue(responseData.AskCredentialOnLaunch)
	data.AskExecutionEnvironmenOnLaunch = types.BoolValue(responseData.AskExecutionEnvironmenOnLaunch)
	data.AskLablesOnLaunch = types.BoolValue(responseData.AskLablesOnLaunch)
	data.AskForksOnLaunch = types.BoolValue(responseData.AskForksOnLaunch)
	data.AskJobSliceCountOnLaunch = types.BoolValue(responseData.AskJobSliceCountOnLaunch)
	data.AskTimeoutOnLaunch = types.BoolValue(responseData.AskTimeoutOnLaunch)
	data.AskInstanceGroupsOnLaunch = types.BoolValue(responseData.AskInstanceGroupsOnLaunch)
	data.SurveyEnabled = types.BoolValue(responseData.SurveyEnabled)
	data.BecomeEnabled = types.BoolValue(responseData.BecomeEnabled)
	data.DiffMode = types.BoolValue(responseData.DiffMode)
	data.AllowSimultaneous = types.BoolValue(responseData.AllowSimultaneous)
	if responseData.CustomVirtualEnv == nil {
		data.CustomVirtualEnv = types.StringValue("")
	} else {
		data.CustomVirtualEnv = types.StringValue(responseData.CustomVirtualEnv.(string))
	}
	data.JobSliceCount = types.Int32Value(int32(responseData.JobSliceCount))
	data.WebhookService = types.StringValue(responseData.WebhookService)
	if responseData.WebhookCredential == nil {
		data.WebhookCredential = types.StringValue("")
	} else {
		data.WebhookCredential = types.StringValue(responseData.WebhookCredential.(string))
	}
	data.PreventInstanceGroupFallback = types.BoolValue(responseData.PreventInstanceGroupFallback)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Left intentinally "blank" (as initialized by clone of template scaffold) as these resources is replace by schema plan modifiers
func (r *JobTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data JobTemplateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data JobTemplateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// set url for create HTTP request
	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
	}
	url := r.client.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/", id)

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate delete request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete got error: %s", err))
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))

	}
}

func (r *JobTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
