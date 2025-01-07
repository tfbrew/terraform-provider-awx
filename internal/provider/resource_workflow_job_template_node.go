package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &WorkflowJobTemplatesNodeResource{}
var _ resource.ResourceWithImportState = &WorkflowJobTemplatesNodeResource{}

func NewWorkflowJobTemplatesNodeResource() resource.Resource {
	return &WorkflowJobTemplatesNodeResource{}
}

// WorkflowJobTemplatesNodeResource defines the resource implementation.
type WorkflowJobTemplatesNodeResource struct {
	client *AwxClient
}

// WorkflowJobTemplatesNodeResourceModel describes the resource data model.
type WorkflowJobTemplatesNodeResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	WorkflowJobId          types.Int32  `tfsdk:"workflow_job_template_id"`
	UnifiedJobTemplateId   types.Int32  `tfsdk:"unified_job_template"`
	Inventory              types.Int32  `tfsdk:"inventory"`
	ExtraData              types.String `tfsdk:"extra_data"`
	ScmBranch              types.String `tfsdk:"scm_branch"`
	JobType                types.String `tfsdk:"job_type"`
	JobTags                types.String `tfsdk:"job_tags"`
	SkipTags               types.String `tfsdk:"skip_tags"`
	Limit                  types.String `tfsdk:"limit"`
	DiffMode               types.Bool   `tfsdk:"diff_mode"`
	Verbosity              types.Int32  `tfsdk:"verbosity"`
	AllParentsMustConverge types.Bool   `tfsdk:"all_parents_must_converge"`
	Identifier             types.String `tfsdk:"identifier"`
}

type WorkflowJobTemplateNodeAPIModel struct {
	Id                     int    `json:"id"`
	WorkflowJobId          int    `json:"workflow_job_template"`
	UnifiedJobTemplateId   int    `json:"unified_job_template"`
	Inventory              int    `json:"inventory"`
	ExtraData              any    `json:"extra_data,omitempty"`
	ScmBranch              string `json:"scm_branch,omitempty"`
	JobType                string `json:"job_type,omitempty"`
	JobTags                string `json:"job_tags,omitempty"`
	SkipTags               string `json:"skip_tags,omitempty"`
	Limit                  string `json:"limit,omitempty"`
	DiffMode               bool   `json:"diff_mode,omitempty"`
	Verbosity              int    `json:"verbosity,omitempty"`
	AllParentsMustConverge bool   `json:"all_parents_must_converge"`
	Identifier             string `json:"identifier"`
}

func (r *WorkflowJobTemplatesNodeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_job_template_node"
}

func (r *WorkflowJobTemplatesNodeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A node has its own ID, but must be associated to a workflow job template ID. And, the node must specify which actual job template this node will run, which is set using the attribute unified_job_template.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The unique ID for this node",
			},
			"workflow_job_template_id": schema.Int32Attribute{
				Required:    true,
				Description: "This is the ID of the workflow job template to which this node is embedded.",
			},
			"unified_job_template": schema.Int32Attribute{
				Required:    true,
				Description: "This is the ID of the actual job template that should be executed as this node.",
			},
			"inventory": schema.Int32Attribute{
				Optional:    true,
				Description: "This attribute is set to optional. However, creating new nodes may not work without providing this value. This provider was set up marking this optional so that you can import existing nodes from your AWX tower environment that were created without specficying inventory. Something that doesn't appear allowed on more current versions of AWX.",
			},
			"extra_data": schema.StringAttribute{
				Optional:    true,
				Description: "JSON Key/value pairs, wrap in jsonencode().",
			},
			"scm_branch": schema.StringAttribute{
				Optional: true,
			},
			"job_type": schema.StringAttribute{
				Optional: true,
			},
			"job_tags": schema.StringAttribute{
				Optional: true,
			},
			"skip_tags": schema.StringAttribute{
				Optional: true,
			},
			"limit": schema.StringAttribute{
				Optional: true,
			},
			"diff_mode": schema.BoolAttribute{
				Optional: true,
			},
			"verbosity": schema.Int32Attribute{
				Optional: true,
			},
			"all_parents_must_converge": schema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Defaults to false.",
			},
			"identifier": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The unique identifier for this node, set automatically by API when creating a new one.",
			},
		},
	}
}

func (r *WorkflowJobTemplatesNodeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WorkflowJobTemplatesNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WorkflowJobTemplatesNodeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var bodyData WorkflowJobTemplateNodeAPIModel

	if !data.WorkflowJobId.IsNull() {
		bodyData.WorkflowJobId = int(data.WorkflowJobId.ValueInt32())
	}
	if !data.UnifiedJobTemplateId.IsNull() {
		bodyData.UnifiedJobTemplateId = int(data.UnifiedJobTemplateId.ValueInt32())
	}
	if !data.Inventory.IsNull() {
		bodyData.Inventory = int(data.Inventory.ValueInt32())
	}
	if !data.ExtraData.IsNull() {
		extraDataMap := new(map[string]any)
		err := json.Unmarshal([]byte(data.ExtraData.ValueString()), &extraDataMap)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable unmarshal map to json",
				fmt.Sprintf("Unable to convert id: %+v. ", data.ExtraData))
			return
		}

		bodyData.ExtraData = extraDataMap
	}
	if !data.ScmBranch.IsNull() {
		bodyData.ScmBranch = data.ScmBranch.ValueString()
	}
	if !data.JobType.IsNull() {
		bodyData.JobType = data.JobType.ValueString()
	}
	if !data.JobTags.IsNull() {
		bodyData.JobTags = data.JobTags.ValueString()
	}
	if !data.SkipTags.IsNull() {
		bodyData.SkipTags = data.SkipTags.ValueString()
	}
	if !data.Limit.IsNull() {
		bodyData.Limit = data.Limit.ValueString()
	}
	if !data.DiffMode.IsNull() {
		bodyData.DiffMode = data.DiffMode.ValueBool()
	}
	if !data.Verbosity.IsNull() {
		bodyData.Verbosity = int(data.Verbosity.ValueInt32())
	}
	if !data.AllParentsMustConverge.IsNull() {
		bodyData.AllParentsMustConverge = data.AllParentsMustConverge.ValueBool()
	}
	if !data.Identifier.IsNull() {
		bodyData.Identifier = data.Identifier.ValueString()
	}

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable marshal json",
			fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
		return
	}

	url := r.client.endpoint + "/api/v2/workflow_job_template_nodes/"

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}
	if httpResp.StatusCode != 201 {
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
			fmt.Sprintf("Expected 201, got %v with message %s. ", httpResp.StatusCode, body))
		return
	}

	tmp := struct {
		Id int `json:"id"`
	}{}

	defer httpResp.Body.Close()
	httpRespBodyData, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get http response body",
			fmt.Sprintf("Error was %v", err))
		return
	}
	err = json.Unmarshal(httpRespBodyData, &tmp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get unmarshall http response to grab ID",
			fmt.Sprintf("error was %v", err))
		return
	}

	idAsString := strconv.Itoa(tmp.Id)

	data.Id = types.StringValue(idAsString)

	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkflowJobTemplatesNodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WorkflowJobTemplatesNodeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	//set url for create HTTP request
	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
		return
	}
	url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_template_nodes/%d/", id)

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
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

	var responseData WorkflowJobTemplateNodeAPIModel

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Uanble to get all data out of the http response data body",
			fmt.Sprintf("Body got %v. ", body))
		return
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Uanble unmarshall response body into object",
			fmt.Sprintf("Error =  %v. ", err.Error()))
		return
	}

	if !(data.WorkflowJobId.IsNull() && responseData.WorkflowJobId == 0) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("workflow_job_template_id"), responseData.WorkflowJobId)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !(data.UnifiedJobTemplateId.IsNull() && responseData.UnifiedJobTemplateId == 0) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("unified_job_template"), responseData.UnifiedJobTemplateId)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !(data.Inventory.IsNull() && responseData.Inventory == 0) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("inventory"), responseData.Inventory)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	rawExtraData := responseData.ExtraData

	rawType := reflect.TypeOf(rawExtraData)

	if rawType.Kind() == reflect.Map {

		rawExtraData, ok := rawExtraData.(map[string]any)

		if !ok {
			resp.Diagnostics.AddError("unable to cast", "Unable to cast Extra Data as map[string]any")
			return
		}

		if len(rawExtraData) == 0 {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("extra_data"), "{}")...)
			if resp.Diagnostics.HasError() {
				return
			}
		} else {
			tempMap := make(map[string]any, len(rawExtraData))
			for k, v := range rawExtraData {
				tempMap[k] = v
			}
			tempJson, err := json.Marshal(tempMap)
			if err != nil {
				resp.Diagnostics.AddError("marshall issue", "Unable to marshall extra data into json for storage.")
				return
			}
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("extra_data"), string(tempJson))...)
			if resp.Diagnostics.HasError() {
				return
			}

		}

	}

	if rawType.Kind() == reflect.String {

		if extraDataProper, ok := responseData.ExtraData.(string); ok {
			if !(data.ExtraData.IsNull() && (extraDataProper == "" || extraDataProper == "{}")) {
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("extra_data"), responseData.ExtraData)...)
				if resp.Diagnostics.HasError() {
					return
				}
			}
		} else {
			resp.Diagnostics.AddError("fail to cast any as string",
				"ExtraData couldnt be interpreted as a string.",
			)
			return
		}
	}

	if !(data.ScmBranch.IsNull() && (responseData.ScmBranch == "")) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scm_branch"), responseData.ScmBranch)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !(data.JobType.IsNull() && (responseData.JobType == "")) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("job_type"), responseData.JobType)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !(data.JobTags.IsNull() && (responseData.JobTags == "")) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("job_tags"), responseData.JobTags)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !(data.SkipTags.IsNull() && (responseData.SkipTags == "")) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("skip_tags"), responseData.SkipTags)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !(data.Limit.IsNull() && (responseData.Limit == "")) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("limit"), responseData.Limit)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !(data.DiffMode.IsNull() && (responseData.DiffMode)) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("diff_mode"), responseData.DiffMode)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !(data.Verbosity.IsNull() && (responseData.Verbosity == 0)) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("verbosity"), responseData.Verbosity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	// if !(data.AllParentsMustConverge.IsNull() && (responseData.AllParentsMustConverge)) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("all_parents_must_converge"), responseData.AllParentsMustConverge)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// }
	if !(data.Identifier.IsNull() && (responseData.Identifier == "")) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier"), responseData.Identifier)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

}

func (r *WorkflowJobTemplatesNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WorkflowJobTemplatesNodeResourceModel

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
			fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
		return
	}

	var bodyData WorkflowJobTemplateNodeAPIModel
	bodyData.WorkflowJobId = int(data.WorkflowJobId.ValueInt32())
	bodyData.UnifiedJobTemplateId = int(data.UnifiedJobTemplateId.ValueInt32())
	bodyData.Inventory = int(data.Inventory.ValueInt32())

	// Generate a go type that fits into the any var so that when ALL
	//  bodyData fields are set with go types, we call Marshall to generate entire JSON
	extraDataMap := new(map[string]any)
	err = json.Unmarshal([]byte(data.ExtraData.ValueString()), &extraDataMap)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable unmarshal map to json",
			fmt.Sprintf("Unable to convert id: %+v. ", data.ExtraData))
		return
	}
	bodyData.ExtraData = extraDataMap

	bodyData.ScmBranch = data.ScmBranch.ValueString()
	bodyData.JobType = data.JobType.ValueString()
	bodyData.JobTags = data.JobTags.ValueString()
	bodyData.SkipTags = data.SkipTags.ValueString()
	bodyData.Limit = data.Limit.ValueString()
	bodyData.DiffMode = data.DiffMode.ValueBool()
	bodyData.Verbosity = int(data.Verbosity.ValueInt32())
	bodyData.AllParentsMustConverge = data.AllParentsMustConverge.ValueBool()
	bodyData.Identifier = data.Identifier.ValueString()

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable marshal json",
			fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
		return
	}

	url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_template_nodes/%d/", id)

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, url, strings.NewReader(string(jsonData)))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}
	if httpResp.StatusCode != 200 {
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

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Left Intentionally blank, as there is no API endpoint to delete a label.
func (r *WorkflowJobTemplatesNodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WorkflowJobTemplatesNodeResourceModel

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
	url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_template_nodes/%d/", id)

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
	if httpResp.StatusCode != 204 {
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
			fmt.Sprintf("Expected 204, got %v with message %s. ", httpResp.StatusCode, body))
		return

	}
}

func (r *WorkflowJobTemplatesNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
