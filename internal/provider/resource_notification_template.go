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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &NotificationTemplatesResource{}
var _ resource.ResourceWithImportState = &NotificationTemplatesResource{}

func NewNotificationTemplatesResource() resource.Resource {
	return &NotificationTemplatesResource{}
}

// NotificationTemplatesResource defines the resource implementation.
type NotificationTemplatesResource struct {
	client *AwxClient
}

// NotificationTemplatesResourceModel describes the resource data model.
type NotificationTemplatesResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	Organization              types.Int32  `tfsdk:"organization"`
	NotificationType          types.String `tfsdk:"notification_type"`
	NotificationConfiguration types.String `tfsdk:"notification_configuration"`
	Messages                  types.String `tfsdk:"messages"`
}

type NotificationTemplateAPI struct {
	Id                        int    `json:"id"`
	Name                      string `json:"name"`
	Description               string `json:"description,omitempty"`
	Organization              int    `json:"organization"`
	NotificationType          string `json:"notification_type"`
	NotificationConfiguration any    `json:"notification_configuration,omitempty"`
	Messages                  any    `json:"messages,omitempty"`
}

type SlackConfiguration struct {
	Channels  []string `json:"channels"`
	HexColors string   `json:"hex_color"`
	Token     string   `json:"token"`
}

type MessageValue struct {
	Body    string `json:"body"`
	Message string `json:"message"`
}

type Messages struct {
	Error            MessageValue            `json:"error"`
	Started          MessageValue            `json:"started"`
	Success          MessageValue            `json:"success"`
	WorkflowApproval map[string]MessageValue `json:"workflow_approval"`
}

func (r *NotificationTemplatesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification_template"
}

func (r *NotificationTemplatesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a notification template. These can be attached, by ID, to job templates, as an example usage.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"organization": schema.Int32Attribute{
				Required: true,
			},
			"notification_type": schema.StringAttribute{
				Required:    true,
				Description: "Only 'slack' is supported in this provider currently. Choose from: email, grafan, irc, mattermost, pagerduty, rocketchat, slack, twilio, webhook.",
			},
			"notification_configuration": schema.StringAttribute{
				Optional:    true,
				Description: "json. This value depends on the notification_type chosen. But, the value should be json. E.g. notification_configuration = jsonencode(blah blah blah). The AWX Tower API never returns a value for Token. So, this provider is coded to ignore changes to that field.",
			},
			"messages": schema.StringAttribute{
				Optional:    true,
				Description: "json",
			},
		},
	}
}

func (r *NotificationTemplatesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NotificationTemplatesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NotificationTemplatesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var bodyData NotificationTemplateAPI
	bodyData.Name = data.Name.ValueString()
	bodyData.Description = data.Description.ValueString()
	bodyData.Organization = int(data.Organization.ValueInt32())
	bodyData.NotificationType = data.NotificationType.ValueString()

	if !data.NotificationConfiguration.IsNull() {
		fieldToBytes := []byte(data.NotificationConfiguration.ValueString())

		slackConfig := new(SlackConfiguration)

		err := json.Unmarshal(fieldToBytes, &slackConfig)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to move Notification Config into json object",
				fmt.Sprintf("Error = %s ", err.Error()))
			return
		}

		bodyData.NotificationConfiguration = slackConfig
	}
	if !data.Messages.IsNull() {
		fieldToBytes := []byte(data.Messages.ValueString())

		messageData := new(Messages)

		err := json.Unmarshal(fieldToBytes, &messageData)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to move Messages into json object",
				fmt.Sprintf("Error = %s ", err.Error()))
			return
		}

		bodyData.Messages = messageData
	}
	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable marshal json",
			fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
		return
	}

	url := r.client.endpoint + "/api/v2/notification_templates/"

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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create notification_template, got error: %s", err))
		return
	}
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))
		return
	}

	tmp := struct {
		Id int `json:"id"`
	}{}

	defer httpResp.Body.Close()
	httpRepsBodyData, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get http response body",
			fmt.Sprintf("Error was %v", err))
		return
	}
	err = json.Unmarshal(httpRepsBodyData, &tmp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get unmarshall http response to grab ID",
			fmt.Sprintf("error was %v", err))
		return
	}

	idAsString := strconv.Itoa(tmp.Id)

	data.Id = types.StringValue(idAsString)

	// Once this object is created, the token value in the NotificationConfiguration field will never be returned as anythong
	// other than blank by the AWX Tower API. So, let's

	tflog.Trace(ctx, "created a resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotificationTemplatesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NotificationTemplatesResourceModel

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
		return
	}
	url := r.client.endpoint + fmt.Sprintf("/api/v2/notification_templates/%d/", id)

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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get data: %s", err))
	}
	if httpResp.StatusCode != 200 && httpResp.StatusCode != 404 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))
		return
	}

	if httpResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	var responseData NotificationTemplateAPI

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

	if !(data.NotificationType.IsNull() && responseData.NotificationType == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("notification_type"), responseData.NotificationType)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !(data.NotificationType.IsNull() && responseData.NotificationType == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("notification_type"), responseData.NotificationType)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Get the State's Token value provide that back into the State.
	// The api will always return a blank value for token on GET, I'm assuming as a safety measure.
	var stateNotifConfig types.String
	diags := req.State.GetAttribute(ctx, path.Root("notification_configuration"), &stateNotifConfig)
	if diags.HasError() {
		return
	}

	useResponseConfig := false
	var stateToken string

	if stateNotifConfig.IsNull() {
		useResponseConfig = true
	} else {
		stateSlackConfig := new(SlackConfiguration)

		err = json.Unmarshal([]byte(stateNotifConfig.ValueString()), &stateSlackConfig)
		if err != nil {
			resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
				"Unable to unmarshall plan's notification configuration into a go type for interogation."+err.Error(),
			)
			return
		}

		stateToken = stateSlackConfig.Token
	}

	// now build the state from the API response data for notification configuration (slack)
	slackConfig := new(SlackConfiguration)

	notificationConfig, ok := responseData.NotificationConfiguration.(map[string]any)
	if !ok {
		resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
			"Unexpected error in resource_notification_templates with responseData.NotificationConfiguration",
		)
		return
	}

	for k, v := range notificationConfig {
		if k == "hex_color" {
			if hexColor, ok := v.(string); ok {
				slackConfig.HexColors = hexColor
			} else {
				resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
					"Unexpected error in esource_notification_templates with. hex_color is not a string",
				)
				return
			}
		}
		if k == "token" {
			if useResponseConfig {
				if respToken, ok := v.(string); ok {
					slackConfig.Token = respToken
				} else {
					resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
						"Unexpected error in esource_notification_templates with. token is not a string",
					)
					return
				}
			} else {
				slackConfig.Token = stateToken
			}
		}
		if k == "channels" {
			if channels, ok := v.([]any); ok {
				channelList := make([]string, 0, len(channels))
				for _, ch := range channels {
					if channel, ok := ch.(string); ok {
						channelList = append(channelList, channel)
					} else {
						resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
							"Unexpected error in esource_notification_templates with. channel is not the right type.",
						)
						return
					}
				}
				slackConfig.Channels = channelList
			} else {
				resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
					"Unexpected error in esource_notification_templates with. channels is not the right type.",
				)
				return
			}
		}
	}

	config, err := json.Marshal(slackConfig)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to move Notification Config into json object",
			fmt.Sprintf("Error = %s ", err.Error()))
		return
	}

	if !(data.NotificationConfiguration.IsNull() && responseData.NotificationConfiguration == nil) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("notification_configuration"), string(config))...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	messages := new(Messages)

	if responseData.Messages != nil {

		foundOneMessageVal := false
		resp_msgs, ok := responseData.Messages.(map[string]any)
		if !ok {
			resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
				"Unexpected error in resource_notification_templates. resp_msgs not the right type",
			)
		}
		for k, v := range resp_msgs {
			if v == nil {
				continue
			} else {
				foundOneMessageVal = true
			}
			if k == "error" {

				if errorMsg, ok := v.(map[string]any); ok {
					if msg_body, ok := errorMsg["body"].(string); ok {
						if msg_message, ok := errorMsg["message"].(string); ok {
							messages.Error = MessageValue{Body: msg_body, Message: msg_message}
						} else {
							resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
								"Unexpected error in esource_notification_templates with. error msg_message is not the right type",
							)
						}
					} else {
						resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
							"Unexpected error in esource_notification_templates with. error msg_body is not the right type",
						)
					}
				} else {
					resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
						"Unexpected error in esource_notification_templates with. error msg is not the right type",
					)
					return
				}

			}
			if k == "started" {
				if startMsg, ok := v.(map[string]any); ok {
					if msg_body, ok := startMsg["body"].(string); ok {
						if msg_message, ok := startMsg["message"].(string); ok {
							messages.Started = MessageValue{Body: msg_body, Message: msg_message}
						} else {
							resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
								"Unexpected error in esource_notification_templates with. started msg_message is not the right type",
							)
						}
					} else {
						resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
							"Unexpected error in esource_notification_templates with. started msg_body is not the right type",
						)
					}
				} else {
					resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
						"Unexpected error in esource_notification_templates with. started msg is not the right type",
					)
					return
				}
			}
			if k == "success" {
				if successMsg, ok := v.(map[string]any); ok {
					if msg_body, ok := successMsg["body"].(string); ok {
						if msg_message, ok := successMsg["message"].(string); ok {
							messages.Success = MessageValue{Body: msg_body, Message: msg_message}
						} else {
							resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
								"Unexpected error in esource_notification_templates with. started msg_message is not the right type",
							)
						}
					} else {
						resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
							"Unexpected error in esource_notification_templates with. started msg_body is not the right type",
						)
					}
				} else {
					resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
						"Unexpected error in esource_notification_templates with. started msg is not the right type",
					)
					return
				}
			}
			if k == "workflow_approval" {
				if wkaParentValue, ok := v.(map[string]any); ok {
					wkaParent := wkaParentValue
					wka := make(map[string]MessageValue, len(wkaParent))
					for key, val := range wkaParent {

						if valMsg, ok := val.(map[string]any); ok {
							if msg_body, ok := valMsg["body"].(string); ok {
								if msg_message, ok := valMsg["message"].(string); ok {
									wka[key] = MessageValue{Body: msg_body, Message: msg_message}
								} else {
									resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
										"Unexpected error in esource_notification_templates with. val msg_message is not the right type",
									)
								}
							} else {
								resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
									"Unexpected error in esource_notification_templates with. val msg_body is not the right type",
								)
							}
						} else {
							resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
								"Unexpected error in esource_notification_templates with. val msg is not the right type",
							)
							return
						}

					}
					messages.WorkflowApproval = wka
				} else {
					resp.Diagnostics.AddError("Unexpected error in resource_notification_templates",
						"Unexpected error in resource_notification_templates. wkaParent is not the right type",
					)
				}
			}
		}

		msgJson, err := json.Marshal(messages)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to move Messages into json object",
				fmt.Sprintf("Error = %s ", err.Error()))
			return
		}

		if !(data.Messages.IsNull() && responseData.Messages == nil) && foundOneMessageVal {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("messages"), string(msgJson))...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}
}

func (r *NotificationTemplatesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data NotificationTemplatesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
		return
	}

	var bodyData NotificationTemplateAPI

	bodyData.Name = data.Name.ValueString()

	if !(data.Description.IsNull() && data.Description.ValueString() != "") {
		bodyData.Description = data.Description.ValueString()
	}
	bodyData.Organization = int(data.Organization.ValueInt32())
	bodyData.NotificationType = data.NotificationType.ValueString()

	fieldToBytes := []byte(data.NotificationConfiguration.ValueString())

	slackConfig := new(SlackConfiguration)

	err = json.Unmarshal(fieldToBytes, &slackConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to move Notification Config into json object",
			fmt.Sprintf("Error = %s ", err.Error()))
		return
	}

	bodyData.NotificationConfiguration = slackConfig

	fieldToBytes = []byte(data.Messages.ValueString())

	messageData := new(Messages)

	err = json.Unmarshal(fieldToBytes, &messageData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to move Messages into json object",
			fmt.Sprintf("Error = %s ", err.Error()))
		return
	}

	bodyData.Messages = messageData

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable marshal json",
			fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
		return
	}

	url := r.client.endpoint + fmt.Sprintf("/api/v2/notification_templates/%d/", id)

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
			fmt.Sprintf("Expected 201, got %v with message %s. ", httpResp.StatusCode, body))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NotificationTemplatesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NotificationTemplatesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
		return
	}
	url := r.client.endpoint + fmt.Sprintf("/api/v2/notification_templates/%d/", id)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate delete request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete got error: %s", err))
		return
	}
	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 204, got %v. ", httpResp.StatusCode))
		return

	}

}

func (r *NotificationTemplatesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
