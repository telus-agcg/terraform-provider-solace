package provider

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func NewMsgVpnDataSource() datasource.DataSource {
	return &dataSource[MsgVpn]{spds: &msgVpnDataSource{}}
}

var _ solaceProviderDataSource[MsgVpn] = &msgVpnDataSource{}

type msgVpnDataSource struct {
	solaceProvider
}

func (r msgVpnDataSource) Name() string {
	return "msgvpn"
}

func (r msgVpnDataSource) Schema() schema.Schema {
	return MsgVpnDataSourceSchema("msg_vpn_name")
}

func (r *msgVpnDataSource) SetProvider(provider solaceProvider) {
	r.solaceProvider = provider
}

func (r msgVpnDataSource) NewData() *MsgVpn {
	return &MsgVpn{}
}

func (r msgVpnDataSource) Read(data *MsgVpn, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.MsgVpnApi.GetMsgVpn(r.Context, *data.MsgVpnName)
	apiResponse, httpResponse, err := apiReq.Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}
