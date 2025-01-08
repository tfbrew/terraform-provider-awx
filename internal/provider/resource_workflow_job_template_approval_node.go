package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &WorkflowJobTemplateApprovalNode{}
var _ resource.ResourceWithImportState = &WorkflowJobTemplateApprovalNode{}

func NewWorkflowJobTemplateApprovalNodeResource() resource.Resource {
	return &WorkflowJobTemplateApprovalNode{}
}

// WorkflowJobTemplateApprovalNode defines the resource implementation.
type WorkflowJobTemplateApprovalNode struct {
	client *AwxClient
}

// type WorkflowJobTemplateApprovalNodeAPIMode struct {
// 	WorkflowJobTemplateID int `json:"workflow_job_template"`
// }

// WorkflowJobTemplateApprovalNodeModel describes the resource data model.
type WorkflowJobTemplateApprovalNodeModel struct {
	NodeId                types.String `tfsdk:"node_id"`
	ApprovalTemplateId    types.Int32  `tfsdk:"approval_template_id"`
	WorkflowJobTemplateId types.Int32  `tfsdk:"workflow_job_template_id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	Timeout               types.Int32  `tfsdk:"timeout"`
}

type WorkflowJobTmplNodeApprvCreateAPIModel struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Timeout     int    `json:"timeout"`
}

func (r *WorkflowJobTemplateApprovalNode) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_job_template_approval_node"
}

func (r *WorkflowJobTemplateApprovalNode) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "To add an approval node/step to an existing workflow job template a new node will be created and a small template will be added to that node. You'll need to create an awx_workflow_job_template_node_success resource (or always/failure) in order to make this approval node execute in the sequence you want your visualizer nodes to run.",

		Attributes: map[string]schema.Attribute{
			"node_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The unique ID for this node.",
			},
			"approval_template_id": schema.Int32Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
				Description: "The unique ID for this approval template.",
			},
			"workflow_job_template_id": schema.Int32Attribute{
				Required:    true,
				Description: "This the ID for the awx_workflow_job_template for which this approval node/template should be embedded.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of this template.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The name of this template.",
				Default:     stringdefault.StaticString(""),
				Computed:    true,
			},
			"timeout": schema.Int32Attribute{
				Description: "The number of seconds for timeout.",
				Required:    true,
			},
		},
	}
}

func (r *WorkflowJobTemplateApprovalNode) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WorkflowJobTemplateApprovalNode) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WorkflowJobTemplateApprovalNodeModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	/////////////////////////////////////////////////////
	// First create an empty node

	url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_templates/%d/workflow_nodes/", data.WorkflowJobTemplateId.ValueInt32())

	newWJTworkflowNode := struct {
		NodeType string `json:"node_type"`
	}{
		NodeType: "approval",
	}

	jsonData, err := json.Marshal(newWJTworkflowNode)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable marshal json",
			fmt.Sprintf("Unable to convert id: %+v. ", newWJTworkflowNode))
		return
	}

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

	data.NodeId = types.StringValue(idAsString)

	////////////////////////////////////////////////////
	// Now we create a new approval template on that node
	////////////////////////////////////////////////////
	var jsonBody WorkflowJobTmplNodeApprvCreateAPIModel

	jsonBody.Name = data.Name.ValueString()
	jsonBody.Description = data.Description.ValueString()
	jsonBody.Timeout = int(data.Timeout.ValueInt32())

	jsonData, err = json.Marshal(jsonBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable marshal json",
			fmt.Sprintf("Unable to convert id: %+v. ", jsonBody))
		return
	}

	url = r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_template_nodes/%d/create_approval_template/", tmp.Id)

	httpReq, err = http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	httpResp, err = r.client.client.Do(httpReq)
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

	tmp = struct {
		Id int `json:"id"`
	}{}

	defer httpResp.Body.Close()
	httpRespBodyData, err = io.ReadAll(httpResp.Body)
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

	data.ApprovalTemplateId = types.Int32Value(int32(tmp.Id))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkflowJobTemplateApprovalNode) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WorkflowJobTemplateApprovalNodeModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	/// read the node's workflow template ID first
	//set url for create HTTP request
	id, err := strconv.Atoi(data.NodeId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.NodeId.ValueString()))
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

	getNameFromResponse := struct {
		WorkflowJobTemplateId int `json:"workflow_job_template"`
		ApprovalTemplateId    int `json:"unified_job_template"`
	}{}

	defer httpResp.Body.Close()
	httpRespBodyData, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get http response body",
			fmt.Sprintf("Error was %v", err))
		return
	}
	err = json.Unmarshal(httpRespBodyData, &getNameFromResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get unmarshall http response to grab ID",
			fmt.Sprintf("error was %v", err))
		return
	}

	data.WorkflowJobTemplateId = types.Int32Value(int32(getNameFromResponse.WorkflowJobTemplateId))
	data.ApprovalTemplateId = types.Int32Value(int32(getNameFromResponse.ApprovalTemplateId))

	/// now read the node's template's data

	url = r.client.endpoint + fmt.Sprintf("/api/v2/workflow_approval_templates/%d/", getNameFromResponse.ApprovalTemplateId)

	// create HTTP request
	httpReq, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	httpResp, err = r.client.client.Do(httpReq)
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

	var readAPIResponse WorkflowJobTmplNodeApprvCreateAPIModel

	defer httpResp.Body.Close()
	httpRespBodyData, err = io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get http response body",
			fmt.Sprintf("Error was %v", err))
		return
	}
	err = json.Unmarshal(httpRespBodyData, &readAPIResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get unmarshall http response to grab ID",
			fmt.Sprintf("error was %v", err))
		return
	}

	data.Name = types.StringValue(readAPIResponse.Name)

	data.Description = types.StringValue(readAPIResponse.Description)
	data.Timeout = types.Int32Value(int32(readAPIResponse.Timeout))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *WorkflowJobTemplateApprovalNode) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WorkflowJobTemplateApprovalNodeModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Left Intentionally blank, as there is no API endpoint to delete a label.
func (r *WorkflowJobTemplateApprovalNode) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WorkflowJobTemplateApprovalNodeModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// set url for create HTTP request
	id, err := strconv.Atoi(data.NodeId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.NodeId.ValueString()))
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

func (r *WorkflowJobTemplateApprovalNode) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("node_id"), req, resp)
}
