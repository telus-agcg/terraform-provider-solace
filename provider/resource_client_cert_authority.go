package provider

import (
	"context"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ provider.ResourceType = clientCertAuthorityResourceType{}

type clientCertAuthorityResourceType struct {
}

func (t clientCertAuthorityResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	solaceProvider, diags := convertProviderType(in)

	return NewResource[ClientCertAuthority](
		clientCertAuthorityResource{solaceProvider: solaceProvider}), diags
}

func (t clientCertAuthorityResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return ClientCertAuthoritySchema("cert_authority_name"), nil
}

var _ solaceProviderResource[ClientCertAuthority] = clientCertAuthorityResource{}

type clientCertAuthorityResource struct {
	solaceProvider
}

func (r clientCertAuthorityResource) NewData() *ClientCertAuthority {
	return &ClientCertAuthority{}
}

func (r clientCertAuthorityResource) Create(data *ClientCertAuthority, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientCertAuthorityApi.CreateClientCertAuthority(r.Context).Body(*data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r clientCertAuthorityResource) Read(data *ClientCertAuthority, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientCertAuthorityApi.GetClientCertAuthority(r.Context, *data.CertAuthorityName)
	apiResponse, httpResponse, err := apiReq.Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}

func (r clientCertAuthorityResource) Update(_ *ClientCertAuthority, data *ClientCertAuthority, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientCertAuthorityApi.UpdateClientCertAuthority(r.Context, *data.CertAuthorityName).Body(*data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r clientCertAuthorityResource) Delete(data *ClientCertAuthority, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientCertAuthorityApi.DeleteClientCertAuthority(r.Context, *data.CertAuthorityName)
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

var clientCertAuthorityImportRegexp *regexp.Regexp = regexp.MustCompile("^([^\\*\\?\\/]+)$")

func (r clientCertAuthorityResource) Import(id string, data *ClientCertAuthority, diag *diag.Diagnostics) {
	match := clientCertAuthorityImportRegexp.FindStringSubmatch(id)
	if match != nil {
		data.CertAuthorityName = &match[1]
	} else {
		diag.AddError("Expected <vpn-name>", id+" does not match "+clientCertAuthorityImportRegexp.String())
		return
	}
}
