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
	Spec        []SurveySpecModel `tfsdk:"spec"` //TODO unsure
}

type SurveySpecModel struct {
	Max  types.Int32  `tfsdk:"max"`
	Min  types.Int32  `tfsdk:"min"`
	Type types.String `tfsdk:"type"`
	//Choices             types.ListType `tfsdk:"choices"` //TODO unsure
	Default             types.String `tfsdk:"default"`
	Required            types.Bool   `tfsdk:"required"`
	Variable            types.String `tfsdk:"variable"`
	QuestionName        types.String `tfsdk:"question_name"`
	QuestionDescription types.String `tfsdk:"question_description"`
}

type JobTemplateSurvey struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Spec        []SurveySpec `json:"spec"` // if this isn't any at end then unmarshal wont' work
}

type SurveySpec struct {
	Max  int    `json:"max"`
	Min  int    `json:"min"`
	Type string `json:"type"`
	//Choices  any    `json:"choices"`
	Default  any    `json:"default"`
	Required bool   `json:"required"`
	Variable string `json:"variable"`
	//`json:"new_question"`: true,
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
				// MarkdownDescription: "Example identifier",
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
			"spec": schema.ListNestedAttribute{ //TODO this whole block is not certain
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
						// "choices": schema.ListAttribute{ //TODO is this correct and match struct reciever?
						// 	ElementType:         types.StringType,
						// 	Optional:            true,
						// 	MarkdownDescription: "List of strings which define the choices users can make for multichoice or multiselect.",
						// 	Description:         "List of strings which define the choices users can make for multichoice or multiselect.",
						// },
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

	//var tempPointer *JobTemplateSurveyResource
	configureData := req.ProviderData.(*JobTemplateSurveyResource)
	//fmt.Println(tempPointer)
	// if !ok {
	// 	resp.Diagnostics.AddError(
	// 		"Unexpected Resource Configure Type",
	// 		fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	// 	)

	// 	return
	// }

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

	//	get Config for provider
	var configData awxProviderModel

	diags := req.Config.Get(ctx, &configData)
	resp.Diagnostics.Append(diags...)

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

	url := configData.Endpoint.ValueString() + fmt.Sprintf("/api/v2/job_templates/%d/survey_spec", id)

	// get body data for HTTP request
	var bodyData JobTemplateSurvey
	bodyData.Name = data.Name.ValueString()
	bodyData.Description = data.Description.ValueString()

	var specs []SurveySpec
	for _, spec := range data.Spec {
		specs = append(specs, SurveySpec{
			Type:                spec.Type.ValueString(),
			QuestionName:        spec.QuestionName.ValueString(),
			QuestionDescription: spec.QuestionDescription.ValueString(),
			Variable:            spec.Variable.ValueString(),
			Required:            spec.Required.ValueBool(),
			//TODO add the rest of the items like choices, default, etc...
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
	httpReq.Header.Add("Authorization", "Bearer"+configData.Token.ValueString())

	httpResp, err := r.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))

	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	// data.Id = types.StringValue("example-id")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
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

	//
	// do http request to get data
	//

	//	get Config for provider
	//var configData awxProviderModel
	//var temper schema.StringAttribute
	//temp := req.ProviderMeta.GetAttribute(ctx, path.Root("endpoint"), &configData)
	//temp := r.endpoint
	//temper := r.token
	//fmt.Println(temp)
	//fmt.Println(temper)
	//diags := req.ProviderMeta.Get(ctx, &configData)
	// resp.Diagnostics.Append(diags...)

	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// set url for create HTTP request
	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
	}

	url := r.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/survey_spec", id)

	// // get body data for HTTP request
	// var bodyData JobTemplateSurvey
	// bodyData.Name = data.Name.ValueString()
	// bodyData.Description = data.Description.ValueString()

	// var specs []SurveySpec
	// for _, spec := range data.Spec {
	// 	specs = append(specs, SurveySpec{
	// 		Type:                spec.Type.ValueString(),
	// 		QuestionName:        spec.QuestionName.ValueString(),
	// 		QuestionDescription: spec.QuestionDescription.ValueString(),
	// 		Variable:            spec.Variable.ValueString(),
	// 		Required:            spec.Required.ValueBool(),
	// 		//TODO add the rest of the items like choices, default, etc...
	// 	})
	// }

	//bodyData.Spec = specs

	// jsonData, err := json.Marshal(bodyData)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable marshal json",
	// 		fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
	// }

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
		//specModel.Choices
		//TODO fix choices on line above
		//specModel.Default = types.BoolValue(strconv.ParseBool(item.Default))
		specModel.Default = types.StringValue(item.Default.(string))
		specModel.Required = types.BoolValue(item.Required)
		specModel.QuestionName = types.StringValue(item.QuestionName)
		specModel.QuestionDescription = types.StringValue(item.QuestionDescription)
		specModel.Variable = types.StringValue(item.Variable)
		dataSpecs = append(dataSpecs, specModel)
	}

	data.Spec = dataSpecs

	//
	// then move data from request into data
	//
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
