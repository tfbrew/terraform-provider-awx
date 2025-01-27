package provider

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/providervalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure awxProvider satisfies various provider interfaces.
var _ provider.Provider = &awxProvider{}
var _ provider.ProviderWithFunctions = &awxProvider{}

// awxProvider defines the provider implementation.
type awxProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// awxProviderModel describes the provider data model.
type awxProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *awxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "awx"
	resp.Version = p.version
}

func (p *awxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "**Warning**: All v0 releases are considered alpha and subject to breaking changes at any time.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "URL for AWX (i.e. https://tower.example.com)",
				Required:    true,
			},
			"token": schema.StringAttribute{
				Description: "AWX access token (instead of username/password)",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "AWX username (instead of token)",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "AWX password (instead of token)",
				Optional:    true,
			},
		},
	}
}

func (p *awxProvider) ConfigValidators(ctx context.Context) []provider.ConfigValidator {
	return []provider.ConfigValidator{
		providervalidator.Conflicting(
			path.MatchRoot("token"),
			path.MatchRoot("username"),
		),
		providervalidator.Conflicting(
			path.MatchRoot("token"),
			path.MatchRoot("password"),
		),
		providervalidator.RequiredTogether(
			path.MatchRoot("username"),
			path.MatchRoot("password"),
		),
	}
}

func (p *awxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	token := os.Getenv("TOWER_OAUTH_TOKEN")
	endpoint := os.Getenv("TOWER_HOST")
	username := os.Getenv("TOWER_USERNAME")
	password := os.Getenv("TOWER_PASSWORD")
	auth := ""

	var data awxProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.Endpoint.ValueString() != "" {
		endpoint = data.Endpoint.ValueString()
	}

	if data.Token.ValueString() != "" {
		token = data.Token.ValueString()
	}

	if data.Username.ValueString() != "" {
		username = data.Username.ValueString()
	}

	if data.Password.ValueString() != "" {
		password = data.Password.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddError(
			"Missing API Endpoint Configuration",
			"While configuring the provider, the API endpoint hostname was not found in "+
				"the TOWER_HOST environment variable or provider "+
				"configuration block endpoint attribute.",
		)
		// Not returning early allows the logic to collect all errors.
	}

	if token != "" && (username != "" || password != "") {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			"The provider token (or TOWER_OAUTH_TOKEN) and either "+
				"username ( or TOWER_USERNAME) or password ( or TOWER_PASSWORD) were both set.",
		)
		return
	}

	if token != "" {
		auth = "Bearer" + " " + token
	} else if username != "" && password != "" {
		authString := username + ":" + password
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))
		auth = "Basic" + " " + encodedAuth
	}

	// Example client configuration for data sources and resources
	httpclient := &http.Client{
		Timeout: 30 * time.Second,
	}

	client := new(AwxClient)

	client.client = httpclient
	client.endpoint = endpoint
	client.auth = auth

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *awxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewInventorySourceResource,
		NewInventoryResource,
		NewJobTemplateCredentialResource,
		NewJobTemplateInstanceGroupsResource,
		NewJobTemplateLabelsResource,
		NewJobTemplateNotifTemplErrResource,
		NewJobTemplateResource,
		NewJobTemplateSurveyResource,
		NewLabelsResource,
		NewNotificationTemplatesResource,
		NewOrganizationResource,
		NewProjectResource,
		NewScheduleResource,
		NewWorkflowJobTemplatesResource,
		NewWorkflowJobTemplatesJobNodeResource,
		NewWorkflowJobTemplatesNodeLabelResource,
		NewWorkflowJobTemplatesNodeSuccessResource,
		NewWorkflowJobTemplatesNodeAlwaysResource,
		NewWorkflowJobTemplateApprovalNodeResource,
	}
}

func (p *awxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCredentialDataSource,
		NewExecutionEnvironmentDataSource,
		NewInventoryDataSource,
		NewInventorySourceDataSource,
		NewInstanceGroupDataSource,
		NewJobTemplateDataSource,
		NewOrganizationDataSource,
		NewProjectDataSource,
		NewScheduleDataSource,
	}
}

func (p *awxProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		//NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &awxProvider{
			version: version,
		}
	}
}
