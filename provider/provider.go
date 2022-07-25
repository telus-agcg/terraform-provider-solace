package provider

import (
	"context"
	"fmt"
	"log"
	"telusag/terraform-provider-solace/sempv2"
	"telusag/terraform-provider-solace/util"

	rt "github.com/go-openapi/runtime/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type provider struct {
	configured bool
	version    string

	Client  sempv2.APIClient
	Context context.Context
}

type providerData struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Scheme   types.String `tfsdk:"scheme"`
	Hostname types.String `tfsdk:"hostname"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	p.Context = context.WithValue(context.Background(), sempv2.ContextBasicAuth,
		sempv2.BasicAuth{
			UserName: data.Username.Value,
			Password: data.Password.Value,
		},
	)

	config := sempv2.NewConfiguration()
	if !data.Scheme.Null {
		config.Scheme = data.Scheme.Value
	}
	config.Host = data.Hostname.Value
	config.UserAgent = "Solace Terraform Provider"

	if !data.Insecure.Null && data.Insecure.Value {
		httpClient, err := rt.TLSClient(
			rt.TLSClientOptions{InsecureSkipVerify: true})
		if err != nil {
			log.Fatal("Unable to create HTTPS client")
		}
		config.HTTPClient = httpClient
	}

	p.Client = *sempv2.NewAPIClient(config)

	p.configured = true
}

func (p *provider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"solace_msgvpn":     msgVpnResourceType{},
		"solace_aclprofile": aclProfileResourceType{},
		"solace_aclprofile_client_connect_exception": aclProfileClientConnectExceptionResourceType{},
		"solace_aclprofile_publish_exception":        aclProfilePublishExceptionResourceType{},
		"solace_aclprofile_subscribe_exception":      aclProfileSubscribeExceptionResourceType{},
		"solace_clientprofile":                       clientProfileResourceType{},
		"solace_clientusername":                      clientUsernameResourceType{},
		"solace_queue":                               queueResourceType{},
		"solace_queue_subscription":                  queueSubscriptionResourceType{},
	}, nil
}

func (p *provider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"username": {
				Type:        types.StringType,
				Required:    true,
				Description: "Solace user with admin access",
			},
			"password": {
				Type:        types.StringType,
				Required:    true,
				Sensitive:   true,
				Description: "Password",
			},
			"scheme": {
				Type:        types.StringType,
				Optional:    true,
				Description: "URL scheme to use: http or https",
				Validators: []tfsdk.AttributeValidator{
					util.StringOneOfValidator("http", "https"),
				},
			},
			"hostname": {
				Type:        types.StringType,
				Required:    true,
				Description: "Hostname for the Solace Event Broker",
			},
			"insecure": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Ignore HTTPS certificate errors",
			},
		},
	}, nil
}

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in tfsdk.Provider) (provider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*provider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return provider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return provider{}, diags
	}

	return *p, diags
}
