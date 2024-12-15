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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &JobTemplateSurveyResource{}
var _ resource.ResourceWithImportState = &JobTemplateSurveyResource{}

func NewJobTemplateSurveyResource() resource.Resource {
	return &JobTemplateSurveyResource{}
}

// JobTemplateSurveyResource defines the resource implementation.
type JobTemplateSurveyResource struct {
	client   *http.Client
	endpoint string
	token    string
}

// JobTemplateSurveyResourceModel describes the resource data model.
type JobTemplateSurveyResourceModel struct {
	Id          types.String      `tfsdk:"id"`
	Name        types.String      `tfsdk:"name"`
	Description types.String      `tfsdk:"description"`
	Spec        []SurveySpecModel `tfsdk:"spec"`
}

type SurveySpecModel struct {
	Max                 types.Int32  `tfsdk:"max"`
	Min                 types.Int32  `tfsdk:"min"`
	Type                types.String `tfsdk:"type"`
	Choices             types.List   `tfsdk:"choices"`
	Default             types.String `tfsdk:"default"`
	Required            types.Bool   `tfsdk:"required"`
	Variable            types.String `tfsdk:"variable"`
	QuestionName        types.String `tfsdk:"question_name"`
	QuestionDescription types.String `tfsdk:"question_description"`
}

type JobTemplateSurvey struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Spec        []SurveySpec `json:"spec"`
}

type SurveySpec struct {
	Max                 int    `json:"max"`
	Min                 int    `json:"min"`
	Type                string `json:"type"`
	Choices             any    `json:"choices"`
	Default             any    `json:"default"`
	Required            bool   `json:"required"`
	Variable            string `json:"variable"`
	QuestionName        string `json:"question_name"`
	QuestionDescription string `json:"question_description"`
}

func (r *JobTemplateSurveyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jobtemplatesurvey"
}

func (r *JobTemplateSurveyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
			"spec": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"max": schema.Int32Attribute{
							Optional:            true,
							MarkdownDescription: "Maximum value, default 1024.",
							Description:         "Maximum value, default 1024.",
						},
						"min": schema.Int32Attribute{
							Optional:            true,
							MarkdownDescription: "Minimum value, default 0.",
							Description:         "Minimum value, default 1024.",
						},
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Must be one of the following: text, textarea, password, integer, float, multiplechoice, or multiselect.",
							Description:         "Must be one of the following: text, textarea, password, integer, float, multiplechoice, or multiselect.",
						},
						"question_name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Name of survey question.",
							Description:         "Name of survey question.",
						},
						"question_description": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Description of survey question.",
							Description:         "Description of survey question.",
						},
						"variable": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Variable name to store users answer to the survey question.",
							Description:         "Variable name to store users answer to the survey question.",
						},
						"required": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Set if the survey question is required, defaults to false.",
							Description:         "Set if the survey question is required, defaults to false.",
						},
						"default": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Default value for the survey question.",
							Description:         "Default value for the survey question.",
						},
						"choices": schema.ListAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							MarkdownDescription: "List of strings which define the choices users can make for multichoice or multiselect.",
							Description:         "List of strings which define the choices users can make for multichoice or multiselect.",
						},
					},
				},
			},
		},
	}
}

func (r *JobTemplateSurveyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	configureData := req.ProviderData.(*JobTemplateSurveyResource)

	r.client = configureData.client
	r.endpoint = configureData.endpoint
	r.token = configureData.token

}

func (r *JobTemplateSurveyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data JobTemplateSurveyResourceModel

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
	}

	url := r.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/survey_spec", id)

	// get body data for HTTP request
	var bodyData JobTemplateSurvey
	bodyData.Name = data.Name.ValueString()
	bodyData.Description = data.Description.ValueString()

	var specs []SurveySpec
	for _, spec := range data.Spec {

		// convert choices to slice of strings
		stringSlice := make([]string, 0, len(spec.Choices.Elements()))
		diag := spec.Choices.ElementsAs(ctx, &stringSlice, true)
		resp.Diagnostics.Append(diag...)

		if resp.Diagnostics.HasError() {
			return
		}

		// convert to interface{} type
		var finalList interface{}
		finalList = stringSlice

		specs = append(specs, SurveySpec{
			Type:                spec.Type.ValueString(),
			QuestionName:        spec.QuestionName.ValueString(),
			QuestionDescription: spec.QuestionDescription.ValueString(),
			Variable:            spec.Variable.ValueString(),
			Required:            spec.Required.ValueBool(),
			Max:                 int(spec.Max.ValueInt32()),
			Min:                 int(spec.Min.ValueInt32()),
			Choices:             finalList,
			Default:             spec.Default.ValueString(),
		})
	}

	bodyData.Spec = specs

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable marshal json",
			fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
	}

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.token)

	httpResp, err := r.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))

	}

	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobTemplateSurveyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data JobTemplateSurveyResourceModel

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

	url := r.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/survey_spec", id)

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.token)

	httpResp, err := r.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))

	}

	var responseData JobTemplateSurvey

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

	data.Name = types.StringValue(responseData.Name)
	data.Description = types.StringValue(responseData.Description)
	//data.Spec =

	var dataSpecs []SurveySpecModel
	for _, item := range responseData.Spec {
		specModel := SurveySpecModel{}
		specModel.Max = types.Int32Value(int32(item.Max))
		specModel.Min = types.Int32Value(int32(item.Min))
		specModel.Type = types.StringValue(item.Type)

		elements := make([]string, 0, len(item.Choices.([]any)))

		for _, v := range item.Choices.([]any) {
			elements = append(elements, v.(string))
		}

		listValue, diags := types.ListValueFrom(ctx, types.StringType, elements)
		if diags.HasError() {
			return
		}

		specModel.Choices = listValue

		specModel.Default = types.StringValue(item.Default.(string))
		specModel.Required = types.BoolValue(item.Required)
		specModel.QuestionName = types.StringValue(item.QuestionName)
		specModel.QuestionDescription = types.StringValue(item.QuestionDescription)
		specModel.Variable = types.StringValue(item.Variable)
		dataSpecs = append(dataSpecs, specModel)
	}

	data.Spec = dataSpecs

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobTemplateSurveyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data JobTemplateSurveyResourceModel

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

func (r *JobTemplateSurveyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data JobTemplateSurveyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *JobTemplateSurveyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
