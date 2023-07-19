package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	nb "github.com/nautobot/go-nautobot/pkg/nautobot"
)

type nautobotProvider struct {
}

type nautobotProviderModel struct {
	Url   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

var _ provider.Provider = &nautobotProvider{}

type apiClient struct {
	Client     *nb.ClientWithResponses
	Server     string
	Token      *SecurityProviderNautobotToken
	BaseClient *nb.Client
}

func New() provider.Provider {
	return &nautobotProvider{}
}

func (p *nautobotProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nautobot"
}

func (p *nautobotProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional:   true,
				Validators: []validator.String{},
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *nautobotProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Nautobot provider.")
	var config nautobotProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Url.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Nautobot URL was not provided.",
			"The provider cannot create the Nautobot API client as there is an unknown configuration value for the Nautobot API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NAUTOBOT_URL environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Nautobot Token was not provided - IsUnknown.",
			"The provider cannot create the Nautobot API client as there is an unknown configuration value for the Nautobot token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NAUTOBOT_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	url := os.Getenv("NAUTOBOT_URL")
	token := os.Getenv("NAUTOBOT_TOKEN")

	if !config.Url.IsNull() {
		url = config.Url.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Nautobot URL was not provided.",
			"The provider cannot create the Nautobot API client as there is an unknown configuration value for the Nautobot API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NAUTOBOT_URL environment variable.",
		)
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Nautobot Token was not provided - is empty string.",
			"The provider cannot create the Nautobot API client as there is an unknown configuration value for the Nautobot token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NAUTOBOT_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: We want to check what errors we may get and how to handle them.
	// It looks like we'd just resp.Diagnostics.AddError and then return if errors same as above.
	client_token, _ := NewSecurityProviderNautobotToken(token)
	new_client, _ := nb.NewClientWithResponses(url, nb.WithRequestEditorFn((client_token.Intercept)))

	// ctx = tflog.SetField(ctx, "nautobot_url", url)
	// tflog.Debug(ctx, "Creating Nautobot client.")
	bc, _ := nb.NewClient(
		url,
		nb.WithRequestEditorFn(client_token.Intercept),
	)

	client := &apiClient{
		Client:     new_client,
		Server:     url,
		Token:      client_token,
		BaseClient: bc,
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *nautobotProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewManufacturersDataSource,
	}
}

func (p *nautobotProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
