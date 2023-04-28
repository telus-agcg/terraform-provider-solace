package provider

import (
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewMsgVpnResource() resource.Resource {
	return &solaceResource[MsgVpn]{spr: &msgVpnResource{}}
}

var _ solaceProviderResource[MsgVpn] = &msgVpnResource{}

type msgVpnResource struct {
	*solaceProvider
}

func (r msgVpnResource) Name() string {
	return "msgvpn"
}

func (r msgVpnResource) Schema() schema.Schema {
	return MsgVpnResourceSchema("msg_vpn_name")
}

func (r *msgVpnResource) SetProvider(provider *solaceProvider) {
	r.solaceProvider = provider
}

func (r msgVpnResource) NewData() *MsgVpn {
	return &MsgVpn{}
}

func (r msgVpnResource) Create(data *MsgVpn, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.MsgVpnApi.CreateMsgVpn(r.Context).Body(*data.ToApi())
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

func (r msgVpnResource) Update(_ *MsgVpn, data *MsgVpn, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.MsgVpnApi.UpdateMsgVpn(r.Context, *data.MsgVpnName).Body(*data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r msgVpnResource) Delete(data *MsgVpn, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.MsgVpnApi.DeleteMsgVpn(r.Context, *data.MsgVpnName)
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

var msgVpnImportRegexp *regexp.Regexp = regexp.MustCompile(`^([^\*\?\/]+)$`)

func (r msgVpnResource) Import(id string, data *MsgVpn, diag *diag.Diagnostics) {
	match := msgVpnImportRegexp.FindStringSubmatch(id)
	if match != nil {
		data.MsgVpnName = &match[1]
	} else {
		diag.AddError("Expected <vpn-name>", id+" does not match "+msgVpnImportRegexp.String())
		return
	}
}
