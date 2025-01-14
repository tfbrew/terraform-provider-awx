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
	Organization int    `json:"organization"`
	Variables    string `json:"variables,omitempty"`
	Kind         string `json:"kind,omitempty"`
	HostFilter   string `json:"host_filter,omitempty"`
}

type ProjectModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Organization       types.Int32  `tfsdk:"organization"`
	ScmType            types.String `tfsdk:"scm_type"`
	Description        types.String `tfsdk:"description"`
	AllowOverride      types.String `tfsdk:"allow_override"`
	Credential         types.Int32  `tfsdk:"credential"`
	DefaultEnv         types.Int32  `tfsdk:"default_environment"`
	LocalPath          types.String `tfsdk:"local_path"`
	ScmBranch          types.String `tfsdk:"kind"`
	ScmClean           types.Bool   `tfsdk:"scm_clean"`
	ScmDelOnUpdate     types.Bool   `tfsdk:"scm_delete_on_update"`
	ScmRefSpec         types.String `tfsdk:"scm_refspec"`
	ScmTrackSubmodules types.Bool   `tfsdk:"scm_track_submodules"`
	ScmUpdCacheTimeout types.Int32  `tfsdk:"scm_update_cache_timeout"`
	ScmUpdOnLaunch     types.Bool   `tfsdk:"scm_update_on_launch"`
	ScmUrl             types.String `tfsdk:"scm_url"`
}

type ProjectAPIModel struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	Organization       int    `json:"organization"`
	ScmType            string `json:"scm_type"`
	Description        string `json:"description,omitempty"`
	AllowOverride      string `json:"allow_override,omitempty"`
	Credential         int    `json:"credential,omitempty"`
	DefaultEnv         int    `json:"default_environment,omitempty"`
	LocalPath          string `json:"local_path,omitempty"`
	ScmBranch          string `json:"scm_branch,omitempty"`
	ScmClean           bool   `json:"scm_clean,omitempty"`
	ScmDelOnUpdate     bool   `json:"scm_delete_on_update,omitempty"`
	ScmRefSpec         string `json:"scm_refspec,omitempty"`
	ScmTrackSubmodules bool   `json:"scm_track_submodules,omitempty"`
	ScmUpdCacheTimeout int    `json:"scm_update_cache_timeout"`
	ScmUpdOnLaunch     bool   `json:"scm_update_on_launch,omitempty"`
	ScmUrl             string `json:"scm_url,omitempty"`
}
