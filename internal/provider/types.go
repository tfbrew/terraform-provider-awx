package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type JTChildAPIRead struct {
	Count   int           `json:"count"`
	Results []ChildResult `json:"results"`
}

type ChildResult struct {
	Id int `json:"id"`
}

type ChildAssocBody struct {
	Id        int  `json:"id"`
	Associate bool `json:"associate"`
}

type ChildDissasocBody struct {
	Id           int  `json:"id"`
	Disassociate bool `json:"disassociate"`
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

type OrganizationModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	CustomVirtualEnv types.String `tfsdk:"custom_virtualenv"`
	DefaultEnv       types.Int32  `tfsdk:"default_environment"`
	MaxHosts         types.Int32  `tfsdk:"max_hosts"`
}

type OrganizationAPIModel struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	CustomVirtualEnv string `json:"custom_virtualenv,omitempty"`
	DefaultEnv       int    `json:"default_environment,omitempty"`
	MaxHosts         int    `json:"max_hosts,omitempty"`
}

type InventoryModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Organization types.Int32  `tfsdk:"organization"`
	Variables    types.String `tfsdk:"variables"`
	Kind         types.String `tfsdk:"kind"`
	HostFilter   types.String `tfsdk:"host_filter"`
}

type InventoryAPIModel struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	Organization int    `json:"organization,omitempty"`
	Variables    string `json:"variables,omitempty"`
	Kind         string `json:"kind"`
	HostFilter   string `json:"host_filter"`
}
