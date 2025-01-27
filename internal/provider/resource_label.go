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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &LabelsResource{}
var _ resource.ResourceWithImportState = &LabelsResource{}

func NewLabelsResource() resource.Resource {
	return &LabelsResource{}
}

// LabelsResource defines the resource implementation.
type LabelsResource struct {
	client *AwxClient
}

func (r *LabelsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_label"
}

func (r *LabelsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create a label resource. NOTE, the AWX Tower API does NOT contain a delete method for labels. Therefore, any terraform destroy actions will succeed in terraform and remove them from your state, but will NOT have actually removed the label from the AWX Tower instance.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Label ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Label name.",
				Required:    true,
			},
			"organization": schema.Int32Attribute{
				Description: "Organization ID for the label to live in.",
				Required:    true,
			},
		},
	}
}

func (r *LabelsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LabelsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LabelModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var bodyData LabelAPIModel

	bodyData.Name = data.Name.ValueString()
	bodyData.Organization = int(data.Organization.ValueInt32())

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable marshal json",
			fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
	}

	url := r.client.endpoint + "/api/v2/labels/"

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
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create label, got error: %s", err))
		return
	}
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
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
			"Unable to get http response body to get newly created label ID",
			fmt.Sprintf("Error: %v", err))
		return
	}
	err = json.Unmarshal(httpRespBodyData, &tmp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unmarshal http response to get newly created label ID",
			fmt.Sprintf("Error: %v", err))
		return
	}

	idAsString := strconv.Itoa(tmp.Id)

	data.Id = types.StringValue(idAsString)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LabelsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LabelModel

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
	url := r.client.endpoint + fmt.Sprintf("/api/v2/labels/%d/", id)

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
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read label, got error: %v", err))
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

	var responseData LabelAPIModel

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get all data out of the http response data body",
			fmt.Sprintf("Body got %v. ", body))
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unmarshall response body into object",
			fmt.Sprintf("Error =  %v. ", err.Error()))
	}

	if !(data.Name.IsNull() && responseData.Name == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), responseData.Name)...)
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

}

func (r *LabelsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data LabelModel

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

	var bodyData LabelAPIModel

	bodyData.Name = data.Name.ValueString()
	bodyData.Organization = int(data.Organization.ValueInt32())

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable marshal json",
			fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
	}

	url := r.client.endpoint + fmt.Sprintf("/api/v2/labels/%d/", id)

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, url, strings.NewReader(string(jsonData)))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", r.client.auth)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Left Intentionally blank, as there is no API endpoint to delete a label.
func (r *LabelsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LabelModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *LabelsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
