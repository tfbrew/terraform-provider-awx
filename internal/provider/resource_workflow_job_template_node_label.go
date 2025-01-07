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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &WorkflowJobTemplatesNodeLabelResource{}
var _ resource.ResourceWithImportState = &WorkflowJobTemplatesNodeLabelResource{}

func NewWorkflowJobTemplatesNodeLabelResource() resource.Resource {
	return &WorkflowJobTemplatesNodeLabelResource{}
}

// WorkflowJobTemplatesNodeLabelResource defines the resource implementation.
type WorkflowJobTemplatesNodeLabelResource struct {
	client *AwxClient
}

// WorkflowJobTemplatesNodeLabelResourceModel describes the resource data model.
type WorkflowJobTemplatesNodeLabelResourceModel struct {
	NodeId   types.String `tfsdk:"node_id"`
	LabelIDs types.List   `tfsdk:"label_ids"`
}

func (r *WorkflowJobTemplatesNodeLabelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_job_template_node_label"
}

func (r *WorkflowJobTemplatesNodeLabelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Specify a node ID and then a list of the lable IDs that are associated to this node.",

		Attributes: map[string]schema.Attribute{
			"node_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the containing workflow job template node.",
			},
			"label_ids": schema.ListAttribute{
				Required:    true,
				Description: "An ordered list of label IDs associated to a particular Workflwo Job Template node. Create new labels first with awx_labels resource type.",
				ElementType: types.Int32Type,
			},
		},
	}
}

func (r *WorkflowJobTemplatesNodeLabelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WorkflowJobTemplatesNodeLabelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WorkflowJobTemplatesNodeLabelResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// var bodyData WorkflowJobTemplateNodeAPIModel

	// if !data.WorkflowJobId.IsNull() {
	// 	bodyData.WorkflowJobId = int(data.WorkflowJobId.ValueInt32())
	// }
	// if !data.UnifiedJobTemplateId.IsNull() {
	// 	bodyData.UnifiedJobTemplateId = int(data.UnifiedJobTemplateId.ValueInt32())
	// }
	// if !data.Inventory.IsNull() {
	// 	bodyData.Inventory = int(data.Inventory.ValueInt32())
	// }
	// if !data.ExtraData.IsNull() {
	// 	extraDataMap := new(map[string]any)
	// 	err := json.Unmarshal([]byte(data.ExtraData.ValueString()), &extraDataMap)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Unable unmarshal map to json",
	// 			fmt.Sprintf("Unable to convert id: %+v. ", data.ExtraData))
	// 		return
	// 	}

	// 	bodyData.ExtraData = extraDataMap
	// }
	// if !data.ScmBranch.IsNull() {
	// 	bodyData.ScmBranch = data.ScmBranch.ValueString()
	// }
	// if !data.JobType.IsNull() {
	// 	bodyData.JobType = data.JobType.ValueString()
	// }
	// if !data.JobTags.IsNull() {
	// 	bodyData.JobTags = data.JobTags.ValueString()
	// }
	// if !data.SkipTags.IsNull() {
	// 	bodyData.SkipTags = data.SkipTags.ValueString()
	// }
	// if !data.Limit.IsNull() {
	// 	bodyData.Limit = data.Limit.ValueString()
	// }
	// if !data.DiffMode.IsNull() {
	// 	bodyData.DiffMode = data.DiffMode.ValueBool()
	// }
	// if !data.Verbosity.IsNull() {
	// 	bodyData.Verbosity = int(data.Verbosity.ValueInt32())
	// }
	// if !data.AllParentsMustConverge.IsNull() {
	// 	bodyData.AllParentsMustConverge = data.AllParentsMustConverge.ValueBool()
	// }
	// if !data.Identifier.IsNull() {
	// 	bodyData.Identifier = data.Identifier.ValueString()
	// }

	// jsonData, err := json.Marshal(bodyData)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable marshal json",
	// 		fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
	// 	return
	// }

	// url := r.client.endpoint + "/api/v2/workflow_job_template_nodes/"

	// // create HTTP request
	// httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to generate request",
	// 		fmt.Sprintf("Unable to gen url: %v. ", url))
	// 	return
	// }

	// httpReq.Header.Add("Content-Type", "application/json")
	// httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	// httpResp, err := r.client.client.Do(httpReq)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	// 	return
	// }
	// if httpResp.StatusCode != 201 {
	// 	defer httpResp.Body.Close()
	// 	body, err := io.ReadAll(httpResp.Body)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Unable read http request response body.",
	// 			err.Error())
	// 		return
	// 	}

	// 	resp.Diagnostics.AddError(
	// 		"Bad request status code.",
	// 		fmt.Sprintf("Expected 201, got %v with message %s. ", httpResp.StatusCode, body))
	// 	return
	// }

	// tmp := struct {
	// 	Id int `json:"id"`
	// }{}

	// defer httpResp.Body.Close()
	// httpRespBodyData, err := io.ReadAll(httpResp.Body)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to get http response body",
	// 		fmt.Sprintf("Error was %v", err))
	// 	return
	// }
	// err = json.Unmarshal(httpRespBodyData, &tmp)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to get unmarshall http response to grab ID",
	// 		fmt.Sprintf("error was %v", err))
	// 	return
	// }

	// idAsString := strconv.Itoa(tmp.Id)

	// data.Id = types.StringValue(idAsString)

	// tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkflowJobTemplatesNodeLabelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WorkflowJobTemplatesNodeLabelResourceModel

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
	url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_template_nodes/%d/labels", id)

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

	listValue, diags := types.ListValueFrom(ctx, types.Int32Type, tfRelatedIds)
	if diags.HasError() {
		return
	}

	data.LabelIDs = listValue
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkflowJobTemplatesNodeLabelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WorkflowJobTemplatesNodeLabelResourceModel

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
	// 	return
	// }

	// var bodyData WorkflowJobTemplateNodeAPIModel
	// bodyData.WorkflowJobId = int(data.WorkflowJobId.ValueInt32())
	// bodyData.UnifiedJobTemplateId = int(data.UnifiedJobTemplateId.ValueInt32())
	// bodyData.Inventory = int(data.Inventory.ValueInt32())

	// // Generate a go type that fits into the any var so that when ALL
	// //  bodyData fields are set with go types, we call Marshall to generate entire JSON
	// extraDataMap := new(map[string]any)
	// err = json.Unmarshal([]byte(data.ExtraData.ValueString()), &extraDataMap)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable unmarshal map to json",
	// 		fmt.Sprintf("Unable to convert id: %+v. ", data.ExtraData))
	// 	return
	// }
	// bodyData.ExtraData = extraDataMap

	// bodyData.ScmBranch = data.ScmBranch.ValueString()
	// bodyData.JobType = data.JobType.ValueString()
	// bodyData.JobTags = data.JobTags.ValueString()
	// bodyData.SkipTags = data.SkipTags.ValueString()
	// bodyData.Limit = data.Limit.ValueString()
	// bodyData.DiffMode = data.DiffMode.ValueBool()
	// bodyData.Verbosity = int(data.Verbosity.ValueInt32())
	// bodyData.AllParentsMustConverge = data.AllParentsMustConverge.ValueBool()
	// bodyData.Identifier = data.Identifier.ValueString()

	// jsonData, err := json.Marshal(bodyData)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable marshal json",
	// 		fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
	// 	return
	// }

	// url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_template_nodes/%d/", id)

	// // create HTTP request
	// httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, url, strings.NewReader(string(jsonData)))
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to generate request",
	// 		fmt.Sprintf("Unable to gen url: %v. ", url))
	// 	return
	// }

	// httpReq.Header.Add("Content-Type", "application/json")
	// httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	// httpResp, err := r.client.client.Do(httpReq)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	// 	return
	// }
	// if httpResp.StatusCode != 200 {
	// 	defer httpResp.Body.Close()
	// 	body, err := io.ReadAll(httpResp.Body)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Unable read http request response body.",
	// 			err.Error())
	// 		return
	// 	}

	// 	resp.Diagnostics.AddError(
	// 		"Bad request status code.",
	// 		fmt.Sprintf("Expected 200, got %v with message %s. ", httpResp.StatusCode, body))
	// 	return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Left Intentionally blank, as there is no API endpoint to delete a label.
func (r *WorkflowJobTemplatesNodeLabelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// var data WorkflowJobTemplatesNodeLabelResourceModel

	// // Read Terraform prior state data into the model
	// resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	// // set url for create HTTP request
	// id, err := strconv.Atoi(data.NodeId.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable convert id from string to int",
	// 		fmt.Sprintf("Unable to convert id: %v. ", data.NodeId.ValueString()))
	// }
	// url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_template_nodes/%d/", id)

	// // create HTTP request
	// httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to generate delete request",
	// 		fmt.Sprintf("Unable to gen url: %v. ", url))
	// }

	// httpReq.Header.Add("Content-Type", "application/json")
	// httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	// httpResp, err := r.client.client.Do(httpReq)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete got error: %s", err))
	// }
	// if httpResp.StatusCode != 204 {
	// 	defer httpResp.Body.Close()
	// 	body, err := io.ReadAll(httpResp.Body)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Unable read http request response body.",
	// 			err.Error())
	// 		return
	// 	}

	// 	resp.Diagnostics.AddError(
	// 		"Bad request status code.",
	// 		fmt.Sprintf("Expected 204, got %v with message %s. ", httpResp.StatusCode, body))
	// 	return

	// }
}

func (r *WorkflowJobTemplatesNodeLabelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("node_id"), req, resp)
}
