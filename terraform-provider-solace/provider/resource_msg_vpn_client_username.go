package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.ResourceType = clientUsernameResourceType{}

type clientUsernameResourceType struct {
}

func (t clientUsernameResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return NewResource[MsgVpnClientUsername](
		clientUsernameResource{provider: provider}), diags
}

func (t clientUsernameResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return MsgVpnClientUsernameSchema("msg_vpn_name", "client_username"), nil
}

var _ solaceProviderResource[MsgVpnClientUsername] = clientUsernameResource{}

type clientUsernameResource struct {
	provider
}

func (r clientUsernameResource) NewData() *MsgVpnClientUsername {
	return &MsgVpnClientUsername{}
}

func (r clientUsernameResource) Create(data *MsgVpnClientUsername, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientUsernameApi.CreateMsgVpnClientUsername(r.Context, *data.MsgVpnName).Body(data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r clientUsernameResource) Read(data *MsgVpnClientUsername, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientUsernameApi.GetMsgVpnClientUsername(r.Context, *data.MsgVpnName, *data.ClientUsername)
	apiResponse, httpResponse, err := apiReq.Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}

func (r clientUsernameResource) Update(data *MsgVpnClientUsername, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientUsernameApi.UpdateMsgVpnClientUsername(r.Context, *data.MsgVpnName, *data.ClientUsername).Body(data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r clientUsernameResource) Delete(data *MsgVpnClientUsername, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.ClientUsernameApi.DeleteMsgVpnClientUsername(r.Context, *data.MsgVpnName, *data.ClientUsername)
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r clientUsernameResource) Import(*MsgVpnClientUsername, *diag.Diagnostics) {}
