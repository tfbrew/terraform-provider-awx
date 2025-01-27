package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &JobTemplateSurveyResource{}
var _ resource.ResourceWithImportState = &JobTemplateSurveyResource{}

func NewJobTemplateSurveyResource() resource.Resource {
	return &JobTemplateSurveyResource{}
}

// JobTemplateSurveyResource defines the resource implementation.
type JobTemplateSurveyResource struct {
	client *AwxClient
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
	Choices             any    `json:"choices,omitempty"`
	Default             any    `json:"default"`
	Required            bool   `json:"required"`
	Variable            string `json:"variable"`
	QuestionName        string `json:"question_name"`
	QuestionDescription string `json:"question_description"`
}

func (r *JobTemplateSurveyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job_template_survey_spec"
}

func (r *JobTemplateSurveyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Associate survey specs to an existing Job Template.",

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
							Optional:    true,
							Description: "Maximum value, default `1024`.",
						},
						"min": schema.Int32Attribute{
							Optional:    true,
							Description: "Minimum value, default `1024`.",
						},
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Must be one of the following: `text`, `textarea`, `password`, `integer`, `float`, `multiplechoice`, or `multiselect`.",
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"text", "textarea", "password", "integer", "float", "multiplechoice", "multiselect"}...),
							},
						},
						"question_name": schema.StringAttribute{
							Required:    true,
							Description: "Name of survey question.",
						},
						"question_description": schema.StringAttribute{
							Required:    true,
							Description: "Description of survey question.",
						},
						"variable": schema.StringAttribute{
							Required:    true,
							Description: "Variable name to store users answer to the survey question.",
						},
						"required": schema.BoolAttribute{
							Optional:    true,
							Description: "Set if the survey question is required, defaults to `false`.",
						},
						"default": schema.StringAttribute{
							Default:     stringdefault.StaticString(""),
							Optional:    true,
							Computed:    true,
							Description: "Default value for the survey question. Supply a value of \"\" when you want no default value, even for type values that are non-text-based.",
						},
						"choices": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "List of strings which define the choices users can make for multichoice or multiselect.",
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

	url := r.client.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/survey_spec", id)

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
		var finalList interface{} = stringSlice

		specs = append(specs, SurveySpec{
			Type:                spec.Type.ValueString(),
			QuestionName:        spec.QuestionName.ValueString(),
			QuestionDescription: spec.QuestionDescription.ValueString(),
			Variable:            spec.Variable.ValueString(),
			Required:            spec.Required.ValueBool(),
			Max:                 int(spec.Max.ValueInt32()),
			Min:                 int(spec.Min.ValueInt32()),
			Choices:             finalList,
			Default:             spec.Default,
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
	httpReq.Header.Add("Authorization", r.client.auth)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))

	}

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

	url := r.client.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/survey_spec", id)

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", r.client.auth)

	httpResp, err := r.client.client.Do(httpReq)
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

	var dataSpecs []SurveySpecModel
	for _, item := range responseData.Spec {
		specModel := SurveySpecModel{}
		specModel.Max = types.Int32Value(int32(item.Max))
		specModel.Min = types.Int32Value(int32(item.Min))
		specModel.Type = types.StringValue(item.Type)

		choiceType := reflect.TypeOf(item.Choices)

		if choiceType != nil {

			itemChoiceKind := reflect.TypeOf(item.Choices).Kind()

			if itemChoiceKind == reflect.Slice {

				choices, ok := item.Choices.([]any)
				if !ok {
					resp.Diagnostics.AddError("Unexpected error in resource_jobtemplate_survey",
						"Unexpected error in resource_jobtemplate_survey",
					)
				}

				elements := make([]string, 0, len(choices))

				for _, v := range choices {
					if strValue, ok := v.(string); ok {
						elements = append(elements, strValue)
					} else {
						resp.Diagnostics.AddError("Unexpected error in resource_jobtemplate_survey",
							"Unexpected error in resource_jobtemplate_survey",
						)
						return
					}
				}

				listValue, diags := types.ListValueFrom(ctx, types.StringType, elements)
				if diags.HasError() {
					return
				}

				specModel.Choices = listValue
			} else {
				specModel.Choices = types.ListNull(types.StringType)
			}
		} else {
			specModel.Choices = types.ListNull(types.StringType)
		}

		itemType := reflect.TypeOf(item.Default)
		if itemType != nil {

			itemDefaultKind := reflect.TypeOf(item.Default).Kind()
			switch itemDefaultKind {
			case reflect.Float64:

				if defaultValue, ok := item.Default.(float64); ok {
					specModel.Default = types.StringValue(fmt.Sprint(defaultValue))
				} else {
					resp.Diagnostics.AddError("Unexpected error in resource_jobtemplate_survey",
						"Unexpected error in resource_jobtemplate_survey",
					)
				}

			default:

				if defaultValue, ok := item.Default.(string); ok {
					specModel.Default = types.StringValue(defaultValue)
				} else {
					resp.Diagnostics.AddError("Unexpected error in resource_jobtemplate_survey",
						"Unexpected error in resource_jobtemplate_survey",
					)
				}
			}
		}
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

// Left intentinally "blank" (as initialized by clone of template scaffold) as these resources is replace by schema plan modifiers.
func (r *JobTemplateSurveyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	url := r.client.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/survey_spec", id)

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
		var finalList interface{} = stringSlice

		specBuilt := SurveySpec{
			Type:                spec.Type.ValueString(),
			QuestionName:        spec.QuestionName.ValueString(),
			QuestionDescription: spec.QuestionDescription.ValueString(),
			Variable:            spec.Variable.ValueString(),
			Required:            spec.Required.ValueBool(),
			Max:                 int(spec.Max.ValueInt32()),
			Min:                 int(spec.Min.ValueInt32()),
			Choices:             finalList,
		}

		stringTypes := []string{"text", "textarea", "multiplechoice", "multipleselect", "password"}
		numberTypes := []string{"integer", "float"}

		switch {
		case slices.Contains(stringTypes, specBuilt.Type):
			specBuilt.Default = spec.Default.ValueString()
		case slices.Contains(numberTypes, specBuilt.Type) && spec.Default.ValueString() != "":
			defaultNumber, err := strconv.Atoi(spec.Default.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("uanble to convert to integer", err.Error())
				return
			}
			specBuilt.Default = defaultNumber
		default:
			specBuilt.Default = ""
		}

		specs = append(specs, specBuilt)
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
	httpReq.Header.Add("Authorization", r.client.auth)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))

	}

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

	// set url for create HTTP request
	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
	}

	url := r.client.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/survey_spec", id)

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate delete request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", r.client.auth)

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

func (r *JobTemplateSurveyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
