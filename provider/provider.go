package provider

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"telusag/terraform-provider-solace/sempv2"

	rt "github.com/go-openapi/runtime/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ provider.Provider = &solaceProvider{}

type solaceProvider struct {
	version string

	Client  sempv2.APIClient
	Context context.Context
}

type providerData struct {
	Username *string `tfsdk:"username" env:"SEMP_USERNAME"`
	Password *string `tfsdk:"password" env:"SEMP_PASSWORD"`
	Scheme   *string `tfsdk:"scheme"   env:"SEMP_SCHEME"`
	Hostname *string `tfsdk:"hostname" env:"SEMP_HOSTNAME"`
	Insecure *bool   `tfsdk:"insecure" env:"SEMP_INSECURE"`
}

func (p *solaceProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "solace"
	resp.Version = p.version
}

func (p *solaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerData

	// Set default values where possible
	data.Scheme = ToPtr("https")
	data.Insecure = ToPtr(false)

	// Try to get Provider config from environment, if possible.
	resp.Diagnostics.Append(configureFromEnvironment(ctx, &data)...)

	// Now override that config with values from Terraform, if present.
	resp.Diagnostics.Append(configureFromTerraformConfig(ctx, req, &data)...)

	// Check if the proper values have now been provided
	if IsNilOrEmpty(data.Hostname) {
		resp.Diagnostics.AddError("Solace hostname is required", "Either set SEMP_HOSTNAME or provide 'hostname' in the provider configuration.")
	}
	if IsNilOrEmpty(data.Username) {
		resp.Diagnostics.AddError("Solace username is required", "Either set SEMP_USERNAME or provide 'username' in the provider configuration.")
	}
	if IsNilOrEmpty(data.Password) {
		resp.Diagnostics.AddError("Solace password is required", "Either set SEMP_PASSWORD or provide 'password' in the provider configuration.")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	p.Context = context.WithValue(context.Background(), sempv2.ContextBasicAuth,
		sempv2.BasicAuth{
			UserName: *data.Username,
			Password: *data.Password,
		},
	)

	config := sempv2.NewConfiguration()
	if !IsNilOrEmpty(data.Scheme) {
		config.Scheme = *data.Scheme
	}
	config.Host = *data.Hostname
	config.UserAgent = "Solace Terraform Provider"

	if data.Insecure != nil && *data.Insecure {
		httpClient, err := rt.TLSClient(
			rt.TLSClientOptions{InsecureSkipVerify: true})
		if err != nil {
			resp.Diagnostics.AddError("Unable to create HTTPS client", err.Error())
		}
		config.HTTPClient = httpClient
	}

	p.Client = *sempv2.NewAPIClient(config)

	resp.DataSourceData = p
	resp.ResourceData = p
}

// configureFromEnvironment read values from the environment variables configured in the struct tags and sets them in `data` if provided
func configureFromEnvironment(ctx context.Context, data interface{}) (diag diag.Diagnostics) {
	_value := reflect.ValueOf(data).Elem()
	_type := _value.Type()

	if _value.Kind() != reflect.Struct {
		diag.AddError("Error configuring from environment", "Attempt to configure into non-struct type: "+_value.Kind().String())
		return
	}

	for i := 0; i < _type.NumField(); i++ {
		typeField := _type.Field(i)
		valueField := _value.FieldByName(typeField.Name)
		envVar := typeField.Tag.Get("env")

		if envVar != "" {
			envVal := os.Getenv(envVar)
			if envVal == "" {
				// No environment variable on this field, ignore it
				continue
			}

			// Make sure this is a pointer field
			if valueField.Type().Kind() != reflect.Pointer {
				diag.AddError("Invalid type", fmt.Sprintf("Unsupported type '%v' for field '%v', only pointer types are supported",
					valueField.Type().Name(), typeField.Name))
				continue
			}

			// Get the kind of value the field is pointing to so we can
			// convert the env value to the correct type when assigning
			valueKind := valueField.Type().Elem().Kind()

			if valueKind == reflect.String {
				valueField.Set(reflect.ValueOf(&envVal))
			} else if valueKind == reflect.Bool {
				boolVal, err := strconv.ParseBool(envVal)
				if err == nil {
					valueField.Set(reflect.ValueOf(&boolVal))
				} else {
					diag.AddError(fmt.Sprintf("Invalid valid for %v", envVar), err.Error())
				}
			} else {
				diag.AddError("Unsupported type", fmt.Sprintf("%v has an invalid type: %v",
					typeField.Name, valueKind))
			}
		}
	}

	return
}

// configureFromTerraformConfig overrides values in 'data' if those same values were provided in the TF config
func configureFromTerraformConfig(ctx context.Context, req provider.ConfigureRequest, data *providerData) (diag diag.Diagnostics) {
	var dataFromTFConfig providerData
	req.Config.Get(ctx, &dataFromTFConfig)

	if !IsNilOrEmpty(dataFromTFConfig.Username) {
		data.Username = dataFromTFConfig.Username
	}
	if !IsNilOrEmpty(dataFromTFConfig.Password) {
		data.Password = dataFromTFConfig.Password
	}
	if !IsNilOrEmpty(dataFromTFConfig.Scheme) {
		data.Scheme = dataFromTFConfig.Scheme
	}
	if !IsNilOrEmpty(dataFromTFConfig.Hostname) {
		data.Hostname = dataFromTFConfig.Hostname
	}
	if dataFromTFConfig.Insecure != nil {
		data.Insecure = dataFromTFConfig.Insecure
	}

	return
}

func (p *solaceProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClientCertAuthorityResource,
		NewMsgVpnResource,
		NewMsgVpnAclProfileClientConnectExceptionResource,
		NewMsgVpnAclProfilePublishExceptionResource,
		NewMsgVpnAclProfileSubscribeExceptionResource,
		NewMsgVpnAclProfileResource,
		NewMsgVpnAuthenticationOauthProfileResource,
		NewMsgVpnClientProfileResource,
		NewMsgVpnClientUsernameResource,
		NewMsgVpnQueueSubscriptionResource,
		NewMsgVpnQueueResource,
	}
}

func (p *solaceProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewMsgVpnDataSource,
	}
}

func (p *solaceProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "Solace user with admin access (env: SEMP_USERNAME)",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password (env: SEMP_PASSWORD)",
			},
			"scheme": schema.StringAttribute{
				Optional:    true,
				Description: "URL scheme to use: http or https (env: SEMP_SCHEME)",
				Validators: []validator.String{
					stringvalidator.OneOf("http", "https"),
				},
			},
			"hostname": schema.StringAttribute{
				Optional:    true,
				Description: "Hostname for the Solace Event Broker (env: SEMP_HOSTNAME)",
			},
			"insecure": schema.BoolAttribute{
				Optional:    true,
				Description: "Ignore HTTPS certificate errors (env: SEMP_INSECURE)",
			},
		},
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &solaceProvider{
			version: version,
		}
	}
}
