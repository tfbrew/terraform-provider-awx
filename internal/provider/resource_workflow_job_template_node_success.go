package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &WorkflowJobTemplatesNodeSuccessResource{}
var _ resource.ResourceWithImportState = &WorkflowJobTemplatesNodeSuccessResource{}

func NewWorkflowJobTemplatesNodeSuccessResource() resource.Resource {
	return &WorkflowJobTemplatesNodeSuccessResource{}
}

// WorkflowJobTemplatesNodeSuccessResource defines the resource implementation.
type WorkflowJobTemplatesNodeSuccessResource struct {
	client *AwxClient
}

// WorkflowJobTemplatesNodeSuccessResourceModel describes the resource data model.
type WorkflowJobTemplatesNodeSuccessResourceModel struct {
	NodeId         types.String `tfsdk:"node_id"`
	SuccessNodeIds types.Set    `tfsdk:"success_node_ids"`
}

func (r *WorkflowJobTemplatesNodeSuccessResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_job_template_node_success"
}

func (r *WorkflowJobTemplatesNodeSuccessResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Specify a node ID and then a list of node IDs that should run when this one ends in success.",

		Attributes: map[string]schema.Attribute{
			"node_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the containing workflow job template node.",
			},
			"success_node_ids": schema.SetAttribute{
				Required:    true,
				Description: "An ordered list of Node IDs attached to this workflow template node that should run on successful completion of this node.",
				ElementType: types.Int32Type,
			},
		},
	}
}

func (r *WorkflowJobTemplatesNodeSuccessResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WorkflowJobTemplatesNodeSuccessResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WorkflowJobTemplatesNodeSuccessResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

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

	url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_template_nodes/%d/success_nodes/", id)

	var relatedIds []int

	diags := data.SuccessNodeIds.ElementsAs(ctx, &relatedIds, false)
	if diags.HasError() {
		return
	}

	for _, val := range relatedIds {

		var bodyData ChildAssocBody
		bodyData.Id = val
		bodyData.Associate = true

		err := r.client.AssocSuccessNode(ctx, bodyData, url)
		if err != nil {
			resp.Diagnostics.AddError("Failed to associate child.", err.Error())
			return
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *WorkflowJobTemplatesNodeSuccessResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WorkflowJobTemplatesNodeSuccessResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	//set url for create HTTP request
	id, err := strconv.Atoi(data.NodeId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.NodeId.ValueString()))
		return
	}
	url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_template_nodes/%d/success_nodes", id)

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

	var responseData JTCredentialAPIRead

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

	tfRelatedIds := make([]int, 0, responseData.Count)

	for _, v := range responseData.Results {
		tfRelatedIds = append(tfRelatedIds, v.Id)
	}

	listValue, diags := types.SetValueFrom(ctx, types.Int32Type, tfRelatedIds)
	if diags.HasError() {
		return
	}
	data.SuccessNodeIds = listValue

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkflowJobTemplatesNodeSuccessResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WorkflowJobTemplatesNodeSuccessResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.NodeId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Converting ID to Int failed", fmt.Sprintf("Converting the job template id %s to int failed.", data.NodeId.ValueString()))
		return
	}

	url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_template_nodes/%d/success_nodes/", id)

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
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))
		return
	}

	var responseData JTChildAPIRead

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

	ApiTfChildIds := make([]int, 0, responseData.Count)

	for _, v := range responseData.Results {
		ApiTfChildIds = append(ApiTfChildIds, v.Id)
	}

	var PlanChildIds []int
	diags := data.SuccessNodeIds.ElementsAs(ctx, &PlanChildIds, false)
	if diags.HasError() {
		return
	}

	// diassociate any chyildren found currently via API call that
	//  are no longer in the plan
	for _, v := range ApiTfChildIds {
		if !slices.Contains(PlanChildIds, v) {
			var bodyData ChildDissasocBody
			bodyData.Id = v

			err := r.client.DisassocJobTemplChild(ctx, bodyData, url)
			if err != nil {
				resp.Diagnostics.AddError("Failed to disassociate child.", err.Error())
				return
			}
		}
	}
	// associate any children found in plan that weren't shown in API response
	for _, v := range PlanChildIds {
		if !slices.Contains(ApiTfChildIds, v) {
			var bodyData ChildAssocBody
			bodyData.Id = v
			bodyData.Associate = true

			err := r.client.AssocSuccessNode(ctx, bodyData, url)
			if err != nil {
				resp.Diagnostics.AddError("Failed to associate child.", err.Error())
				return
			}
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Left Intentionally blank, as there is no API endpoint to delete a label.
func (r *WorkflowJobTemplatesNodeSuccessResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WorkflowJobTemplatesNodeSuccessResourceModel

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

	url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_template_nodes/%d/success_nodes/", id)

	var RelatedIds []int

	diags := data.SuccessNodeIds.ElementsAs(ctx, &RelatedIds, false)
	if diags.HasError() {
		return
	}

	for _, val := range RelatedIds {

		var body ChildDissasocBody

		body.Id = val
		body.Disassociate = true

		err := r.client.DisassocJobTemplChild(ctx, body, url)
		if err != nil {
			resp.Diagnostics.AddError("Failed to disassociate child.", err.Error())
			return
		}
	}
}

func (r *WorkflowJobTemplatesNodeSuccessResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("node_id"), req, resp)
}
