// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &JobTemplateLabelsResource{}
var _ resource.ResourceWithImportState = &JobTemplateLabelsResource{}

func NewJobTemplateLabelsResource() resource.Resource {
	return &JobTemplateLabelsResource{}
}

// JobTemplateLabelsResource defines the resource implementation.
type JobTemplateLabelsResource struct {
	client *AwxClient
}

// JobTemplateLabelsResourceModel describes the resource data model.
type JobTemplateLabelsResourceModel struct {
	JobTemplateId types.String `tfsdk:"job_template_id"`
	LabelIDs      types.List   `tfsdk:"label_ids"`
}

type JTLabelsAPIRead struct {
	Count        int           `json:"count"`
	LabelResults []LabelResult `json:"results"`
}

type LabelResult struct {
	Id int `json:"id"`
}

type LabelDissasocBody struct {
	Id           int  `json:"id"`
	Disassociate bool `json:"disassociate"`
}

func (r *JobTemplateLabelsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jobtemplate_labels"
}

func (r *JobTemplateLabelsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The /api/v2/job_templates/{id}/labels/ returns all label objects associated to the template. But, when asked to associate a label or \
                              dissassociate a label, you must post a request once per label name.Therefore, I couldn't find a way to limit this resource to the 'one api call' \
                              principle. Instead, the terraform schema stores a list of associated labels. And, when creating or deleting or updated, it will make one api call PER \
                              list element. This allows the import function to work by only needing to pass in one job template ID to fill out the entire resource. If this was not done this way \
                              then when someone tries to to use the terraform plan -generate-config-out=./file.tf functionality it will create the resource block correctly. Otherwise, the \
                              -generate-config-out function would have to generate several resource blocks per tempalte id and it's not set up to do that, per my current awareness. As I'm writing this \
                              provider specifically so we can use the -generate-config-out option, I felt this was worth the price of breaking this principle. The downside seems to be that this means \
							  if one of the list element's api calls succeeds, but a subsequent list element's fails, the success of the first element's call is not magially un-done. \
							  So you'll perpas have to use refresh state functions in tf cli to resolve.`,
		Description: `The /api/v2/job_templates/{id}/labels/ returns all label objects associated to the template. But, when asked to associate a label or \
                              dissassociate a label, you must post a request once per label name.Therefore, I couldn't find a way to limit this resource to the 'one api call' \
                              principle. Instead, the terraform schema stores a list of associated labels. And, when creating or deleting or updated, it will make one api call PER \
                              list element. This allows the import function to work by only needing to pass in one job template ID to fill out the entire resource. If this was not done this way \
                              then when someone tries to to use the terraform plan -generate-config-out=./file.tf functionality it will create the resource block correctly. Otherwise, the \
                              -generate-config-out function would have to generate several resource blocks per tempalte id and it's not set up to do that, per my current awareness. As I'm writing this \
                              provider specifically so we can use the -generate-config-out option, I felt this was worth the price of breaking this principle. The downside seems to be that this means \
							  if one of the list element's api calls succeeds, but a subsequent list element's fails, the success of the first element's call is not magially un-done. \
							  So you'll perpas have to use refresh state functions in tf cli to resolve.`,
		Attributes: map[string]schema.Attribute{
			"job_template_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the containing Job Template.",
				MarkdownDescription: "The ID of the containing Job Template",
			},
			"label_ids": schema.ListAttribute{
				Required:            true,
				Description:         "An ordered list of label IDs associated to a particular Job Template. Create new labels first with awx_labels resource type.",
				MarkdownDescription: "An ordered list of label IDs associated to a particular Job Template. Create new labels first with awx_labels resource type.",
				ElementType:         types.Int32Type,
			},
		},
	}
}

func (r *JobTemplateLabelsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	configureData := req.ProviderData.(*AwxClient)

	r.client = configureData
}

func (r *JobTemplateLabelsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data JobTemplateLabelsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// set url for create HTTP request
	id, err := strconv.Atoi(data.JobTemplateId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.JobTemplateId.ValueString()))
	}

	var relatedIds []int

	diags := data.LabelIDs.ElementsAs(ctx, &relatedIds, false)
	if diags.HasError() {
		return
	}

	for _, val := range relatedIds {

		var bodyData Result
		bodyData.Id = val

		err := r.client.AssocJobTemplCredential(ctx, id, bodyData)
		if err != nil {
			resp.Diagnostics.AddError("Failed to associate credential.", err.Error())
			return
		}
	}

	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobTemplateLabelsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data JobTemplateLabelsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// set url for create HTTP request
	id, err := strconv.Atoi(data.JobTemplateId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Converting ID to Int failed", fmt.Sprintf("Converting the job template id %s to int failed.", data.JobTemplateId.ValueString()))
		return
	}

	url := r.client.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/labels/", id)

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
		if data.LabelIDs.IsNull() {
			tfRelatedIds = append(tfRelatedIds, v.Id)
		}
	}

	listValue, diags := types.ListValueFrom(ctx, types.Int32Type, tfRelatedIds)
	if diags.HasError() {
		return
	}

	data.LabelIDs = listValue

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobTemplateLabelsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data JobTemplateLabelsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.JobTemplateId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Converting ID to Int failed", fmt.Sprintf("Converting the job template id %s to int failed.", data.JobTemplateId.ValueString()))
		return
	}

	url := r.client.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/labels/", id)

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

	var responseData JTLabelsAPIRead

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

	ApiTfCredIds := make([]int, 0, responseData.Count)

	for _, v := range responseData.LabelResults {
		ApiTfCredIds = append(ApiTfCredIds, v.Id)
	}

	var PlanLabelIds []int
	diags := data.LabelIDs.ElementsAs(ctx, &PlanLabelIds, false)
	if diags.HasError() {
		return
	}

	// diassociate any credentials found currently via API call that
	//  are no longer in the plan
	for _, v := range ApiTfCredIds {
		if !slices.Contains(PlanLabelIds, v) {
			var bodyData LabelDissasocBody
			bodyData.Id = v

			err := r.client.DisassocJobTemplLabel(ctx, id, bodyData)
			if err != nil {
				resp.Diagnostics.AddError("Failed to disassociate label.", err.Error())
				return
			}
		}
	}
	// associate any credentials found in plan that weren't shown in API response
	for _, v := range PlanLabelIds {
		if !slices.Contains(ApiTfCredIds, v) {
			var bodyData LabelResult
			bodyData.Id = v

			err := r.client.AssocJobTemplLabel(ctx, id, bodyData)
			if err != nil {
				resp.Diagnostics.AddError("Failed to associate label.", err.Error())
				return
			}
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobTemplateLabelsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data JobTemplateLabelsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// set url for create HTTP request
	id, err := strconv.Atoi(data.JobTemplateId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.JobTemplateId.ValueString()))
	}

	var RelatedIds []int

	diags := data.LabelIDs.ElementsAs(ctx, &RelatedIds, false)
	if diags.HasError() {
		return
	}

	for _, val := range RelatedIds {

		var body DissasocBody

		body.Id = val
		body.Disassociate = true

		err := r.client.DisassocJobTemplCredential(ctx, id, body)
		if err != nil {
			resp.Diagnostics.AddError("Failed to disassociate label.", err.Error())
			return
		}
	}

}

func (r *JobTemplateLabelsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("job_template_id"), req, resp)
}
