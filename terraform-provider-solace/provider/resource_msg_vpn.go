package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.ResourceType = msgVpnResourceType{}

type msgVpnResourceType struct {
}

func (t msgVpnResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return NewResource[MsgVpn](
		msgVpnResource{provider: provider}), diags
}

func (t msgVpnResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return MsgVpnSchema("msg_vpn_name"), nil
}

var _ solaceProviderResource[MsgVpn] = msgVpnResource{}

type msgVpnResource struct {
	provider
}

func (r msgVpnResource) NewData() *MsgVpn {
	return &MsgVpn{}
}

func (r msgVpnResource) Create(data *MsgVpn, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.MsgVpnApi.CreateMsgVpn(r.Context).Body(data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r msgVpnResource) Read(data *MsgVpn, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.MsgVpnApi.GetMsgVpn(r.Context, *data.MsgVpnName)
	apiResponse, httpResponse, err := apiReq.Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}

func (r msgVpnResource) Update(data *MsgVpn, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.MsgVpnApi.UpdateMsgVpn(r.Context, *data.MsgVpnName).Body(data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r msgVpnResource) Delete(data *MsgVpn, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.MsgVpnApi.DeleteMsgVpn(r.Context, *data.MsgVpnName)
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r msgVpnResource) Import(*MsgVpn, *diag.Diagnostics) {}
