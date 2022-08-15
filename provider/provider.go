package provider

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"telusag/terraform-provider-solace/sempv2"
	"telusag/terraform-provider-solace/util"

	rt "github.com/go-openapi/runtime/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type solaceProvider struct {
	configured bool
	version    string

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

	p.configured = true
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

func (p *solaceProvider) GetResources(ctx context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return map[string]provider.ResourceType{
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

func (p *solaceProvider) GetDataSources(ctx context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	return map[string]provider.DataSourceType{
		"solace_msgvpn": msgVpnDataSourceType{},
	}, nil
}

func (p *solaceProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"username": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Solace user with admin access (env: SEMP_USERNAME)",
			},
			"password": {
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
				Description: "Password (env: SEMP_PASSWORD)",
			},
			"scheme": {
				Type:        types.StringType,
				Optional:    true,
				Description: "URL scheme to use: http or https (env: SEMP_SCHEME)",
				Validators: []tfsdk.AttributeValidator{
					util.StringOneOfValidator("http", "https"),
				},
			},
			"hostname": {
				Type:        types.StringType,
				Optional:    true,
				Description: "Hostname for the Solace Event Broker (env: SEMP_HOSTNAME)",
			},
			"insecure": {
				Type:        types.BoolType,
				Optional:    true,
				Description: "Ignore HTTPS certificate errors (env: SEMP_INSECURE)",
			},
		},
	}, nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &solaceProvider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in provider.Provider) (solaceProvider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*solaceProvider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return solaceProvider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return solaceProvider{}, diags
	}

	return *p, diags
}
