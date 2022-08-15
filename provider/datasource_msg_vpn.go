package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ provider.DataSourceType = msgVpnDataSourceType{}

type msgVpnDataSourceType struct {
}

func (t msgVpnDataSourceType) NewDataSource(ctx context.Context, in provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	solaceProvider, diags := convertProviderType(in)

	return NewDataSource[MsgVpn](
		msgVpnDataSource{solaceProvider: solaceProvider}), diags
}

func (t msgVpnDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return MsgVpnSchema("msg_vpn_name"), nil
}

var _ solaceProviderDataSource[MsgVpn] = msgVpnDataSource{}

type msgVpnDataSource struct {
	solaceProvider
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
