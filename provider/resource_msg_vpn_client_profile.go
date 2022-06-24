package provider

import (
	"context"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.ResourceType = clientProfileResourceType{}

type clientProfileResourceType struct {
}

func (t clientProfileResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return NewResource[MsgVpnClientProfile](
		clientProfileResource{provider: provider}), diags
}

func (t clientProfileResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return MsgVpnClientProfileSchema("msg_vpn_name", "client_profile_name"), nil
}

var _ solaceProviderResource[MsgVpnClientProfile] = clientProfileResource{}

type clientProfileResource struct {
	provider
}

func (r clientProfileResource) NewData() *MsgVpnClientProfile {
	return &MsgVpnClientProfile{}
}

func (r clientProfileResource) Create(data *MsgVpnClientProfile, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientProfileApi.CreateMsgVpnClientProfile(r.Context, *data.MsgVpnName).Body(*data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r clientProfileResource) Read(data *MsgVpnClientProfile, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientProfileApi.GetMsgVpnClientProfile(r.Context, *data.MsgVpnName, *data.ClientProfileName)
	apiResponse, httpResponse, err := apiReq.Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}

func (r clientProfileResource) Update(_ *MsgVpnClientProfile, data *MsgVpnClientProfile, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientProfileApi.UpdateMsgVpnClientProfile(r.Context, *data.MsgVpnName, *data.ClientProfileName).Body(*data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r clientProfileResource) Delete(data *MsgVpnClientProfile, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientProfileApi.DeleteMsgVpnClientProfile(r.Context, *data.MsgVpnName, *data.ClientProfileName)
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

var msgVpnClientProfileImportRegexp *regexp.Regexp = regexp.MustCompile(
	"^([^\\s\\*\\?\\/]+)\\/([0-9a-zA-Z_\\-]+)$")

func (r clientProfileResource) Import(id string, data *MsgVpnClientProfile, diag *diag.Diagnostics) {
	match := msgVpnClientProfileImportRegexp.FindStringSubmatch(id)
	if match != nil {
		data.MsgVpnName = &match[1]
		data.ClientProfileName = &match[2]
	} else {
		diag.AddError("Expected <vpn-name>/<client-profile>", id+" does not match "+msgVpnClientProfileImportRegexp.String())
		return
	}
}
