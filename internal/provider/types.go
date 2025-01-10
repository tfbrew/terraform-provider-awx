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
