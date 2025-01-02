package provider

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
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
}

func (p *awxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "awx"
	resp.Version = p.version
}

func (p *awxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "**Warning**: All v0.*.* releases are considered alpha and subject to breaking changes at any time.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				//MarkdownDescription: "Example provider attribute",
				Optional: true,
			},
			"token": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *awxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	token := os.Getenv("TOWER_OAUTH_TOKEN")
	endpoint := os.Getenv("TOWER_HOST")

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

	if token == "" {
		resp.Diagnostics.AddError(
			"Missing API Token Configuration",
			"While configuring the provider, the API token was not found in "+
				"the TOWER_OAUTH_TOKEN environment variable or provider "+
				"configuration block token attribute.",
		)
		// Not returning early allows the logic to collect all errors.
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

	// Example client configuration for data sources and resources
	httpclient := &http.Client{
		Timeout: 30 * time.Second,
	}

	client := new(AwxClient)

	client.client = httpclient
	client.endpoint = endpoint
	client.token = token

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *awxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewJobTemplateSurveyResource,
		NewJobTemplateResource,
		NewJobTemplateCredentialResource,
		NewLabelsResource,
		NewJobTemplateLabelsResource,
		NewJobTemplateInstanceGroupsResource,
		NewNotificationTemplatesResource,
		NewJobTemplateNotifTemplErrResource,
		NewWorkflowJobTemplatesResource,
	}
}

func (p *awxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewInstanceGroupsDataSource,
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
