package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ProjectResource{}
var _ resource.ResourceWithImportState = &ProjectResource{}

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

// ProjectResource defines the resource implementation.
type ProjectResource struct {
	client *AwxClient
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Manage an AWX project.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Project ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Project name.",
				Required:    true,
			},
			"organization": schema.Int32Attribute{
				Description: "Organization ID for the project to live in.",
				Required:    true,
			},
			"scm_type": schema.StringAttribute{
				Description: "Type of SCM resource. Options: `manual`, `git`, `svn` `insights`, `archive`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"manual", "git", "svn", "insights", "archive"}...),
				},
			},
			"description": schema.StringAttribute{
				Description: "Project description.",
				Optional:    true,
			},
			"allow_override": schema.StringAttribute{
				Description: "Allow changing the Source Control branch or revision in a job template that uses this project.",
				Optional:    true,
			},
			"credential": schema.Int32Attribute{
				Description: "Source Control credential ID.",
				Optional:    true,
			},
			"default_environment": schema.Int32Attribute{
				Description: "The ID of the execution environment that will be used for jobs that use this project.",
				Optional:    true,
			},
			"local_path": schema.StringAttribute{
				Description: "The server playbook directory for manual projects.",
				Optional:    true,
			},
			"scm_branch": schema.StringAttribute{
				Description: "The branch name in source control.",
				Optional:    true,
			},
			"scm_clean": schema.BoolAttribute{
				Description: "Remove any local modifications prior to performing an update.",
				Optional:    true,
			},
			"scm_delete_on_update": schema.BoolAttribute{
				Description: "Delete the local repository in its entirety prior to performing an update. Depending on the size of the repository this may significantly increase the amount of time required to complete an update.",
				Optional:    true,
			},
			"scm_refspec": schema.StringAttribute{
				Description: "The refspec to use for the SCM resource.",
				Optional:    true,
			},
			"scm_track_submodules": schema.BoolAttribute{
				Description: "Track submodules latest commit on specified branch.",
				Optional:    true,
			},
			"scm_update_cache_timeout": schema.Int32Attribute{
				Description: "Cache Timeout to cache prior project syncs for a certain number of seconds. Only valid if scm_update_on_launch is to True, otherwise ignored.",
				Optional:    true,
			},
			"scm_update_on_launch": schema.BoolAttribute{
				Description: "Perform an update to the local repository before launching a job with this project.",
				Optional:    true,
			},
			"scm_url": schema.StringAttribute{
				Description: "Example URLs for Remote Archive Source Control include: `https://github.com/username/project/archive/v0.0.1.tar.gz` `https://github.com/username/project/archive/v0.0.2.zip`",
				Optional:    true,
			},
		},
	}
}

func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	configureData, ok := req.ProviderData.(*AwxClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = configureData
}

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProjectModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// set url for create HTTP request

	var bodyData ProjectAPIModel

	if !(data.Name.IsNull()) {
		bodyData.Name = data.Name.ValueString()
	}
	if !(data.Description.IsNull()) {
		bodyData.Description = data.Description.ValueString()
	}
	if !(data.Organization.IsNull()) {
		bodyData.Organization = int(data.Organization.ValueInt32())
	}
	if !(data.Variables.IsNull()) {
		bodyData.Variables = data.Variables.ValueString()
	}
	if !(data.Kind.IsNull()) {
		bodyData.Kind = data.Kind.ValueString()
	}
	if !(data.HostFilter.IsNull()) {
		bodyData.HostFilter = data.HostFilter.ValueString()
	}

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to marshal bodyData to json",
			fmt.Sprintf("bodyData: %+v.", bodyData))
		return
	}

	url := r.client.endpoint + "/api/v2/projects/"

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate create request",
			fmt.Sprintf("url: %v, data: %+v ", url, jsonData))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create project, got error: %s", err))
		return
	}
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Bad request status code",
			fmt.Sprintf("Expected 201, got %v.", httpResp.StatusCode))
		return
	}

	tmp := struct {
		Id int `json:"id"`
	}{}

	defer httpResp.Body.Close()
	httpRespBodyData, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get http response body to get newly created project ID",
			fmt.Sprintf("Error: %v", err))
		return
	}
	err = json.Unmarshal(httpRespBodyData, &tmp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unmarshal http response to get newly created project ID",
			fmt.Sprintf("Error: %v", err))
		return
	}

	idAsString := strconv.Itoa(tmp.Id)

	data.Id = types.StringValue(idAsString)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProjectModel

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
			fmt.Sprintf("Unable to convert id: %v.", data.Id))
		return
	}
	url := r.client.endpoint + fmt.Sprintf("/api/v2/projects/%d/", id)

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate read request",
			fmt.Sprintf("url: %v.", url))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read project, got error: %v", err))
		return
	}
	if httpResp.StatusCode != 200 && httpResp.StatusCode != 404 {
		defer httpResp.Body.Close()
		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable read http request response body.",
				err.Error())
			return
		}

		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v with message %s. ", httpResp.StatusCode, body))
		return
	}

	if httpResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	var responseData ProjectAPIModel

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read the http response data body",
			fmt.Sprintf("Body: %v.", body))
		return
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unmarshal json",
			fmt.Sprintf("bodyData: %+v.", body))
		return
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

	if !(data.Organization.IsNull() && responseData.Organization == 0) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization"), responseData.Organization)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !(data.Variables.IsNull() && responseData.Variables == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("variables"), responseData.Variables)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !(data.Kind.IsNull() && responseData.Kind == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("kind"), responseData.Kind)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !(data.HostFilter.IsNull() && responseData.HostFilter == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("host_filter"), responseData.HostFilter)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProjectModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// set url for create HTTP request
	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v.", data))
		return
	}

	var bodyData ProjectAPIModel

	if !(data.Name.IsNull()) {
		bodyData.Name = data.Name.ValueString()
	}
	if !(data.Description.IsNull()) {
		bodyData.Description = data.Description.ValueString()
	}
	if !(data.Organization.IsNull()) {
		bodyData.Organization = int(data.Organization.ValueInt32())
	}
	if !(data.Variables.IsNull()) {
		bodyData.Variables = data.Variables.ValueString()
	}
	if !(data.Kind.IsNull()) {
		bodyData.Kind = data.Kind.ValueString()
	}
	if !(data.HostFilter.IsNull()) {
		bodyData.HostFilter = data.HostFilter.ValueString()
	}

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to marshal bodyData to json",
			fmt.Sprintf("bodyData: %+v.", bodyData))
		return
	}

	url := r.client.endpoint + fmt.Sprintf("/api/v2/projects/%d/", id)

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, url, strings.NewReader(string(jsonData)))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate update request",
			fmt.Sprintf("url: %v, data: %+v ", url, jsonData))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update project, got error: %s", err))
		return
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Bad request status code",
			fmt.Sprintf("Expected 200, got %v.", httpResp.StatusCode))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProjectModel

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
			fmt.Sprintf("Unable to convert id: %v.", data.Id.ValueString()))
		return
	}
	url := r.client.endpoint + fmt.Sprintf("/api/v2/projects/%d/", id)

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate delete request",
			fmt.Sprintf("url: %v", url))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete project, got error: %s.", err))
		return
	}

	// 202 - accepted for deletion, 204 - success
	if httpResp.StatusCode != 202 && httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Bad request status code",
			fmt.Sprintf("Expected [202, 204], got %v.", httpResp.StatusCode))
		return
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
