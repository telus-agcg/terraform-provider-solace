package provider

import (
	"context"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.ResourceType = aclProfileResourceType{}

type aclProfileClientConnectExceptionResourceType struct {
}

func (t aclProfileClientConnectExceptionResourceType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return NewResource[MsgVpnAclProfileClientConnectException](
		aclProfileClientConnectExceptionResource{provider: provider}), diags
}

func (t aclProfileClientConnectExceptionResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return MsgVpnAclProfileClientConnectExceptionSchema("msg_vpn_name", "acl_profile_name", "client_connect_exception_address"), nil
}

var _ solaceProviderResource[MsgVpnAclProfile] = aclProfileResource{}

type aclProfileClientConnectExceptionResource struct {
	provider
}

func (r aclProfileClientConnectExceptionResource) NewData() *MsgVpnAclProfileClientConnectException {
	return &MsgVpnAclProfileClientConnectException{}
}

func (r aclProfileClientConnectExceptionResource) Create(data *MsgVpnAclProfileClientConnectException, diag *diag.Diagnostics) (*http.Response, error) {
	_, httpResponse, err := r.Client.AclProfileApi.
		CreateMsgVpnAclProfileClientConnectException(r.Context, *data.MsgVpnName, *data.AclProfileName).
		Body(*data.ToApi()).
		Execute()
	return httpResponse, err
}

func (r aclProfileClientConnectExceptionResource) Read(data *MsgVpnAclProfileClientConnectException, diag *diag.Diagnostics) (*http.Response, error) {
	apiResponse, httpResponse, err := r.Client.AclProfileApi.
		GetMsgVpnAclProfileClientConnectException(
			r.Context, *data.MsgVpnName, *data.AclProfileName, *data.ClientConnectExceptionAddress).
		Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}

func (r aclProfileClientConnectExceptionResource) Update(_ *MsgVpnAclProfileClientConnectException, data *MsgVpnAclProfileClientConnectException, diag *diag.Diagnostics) (*http.Response, error) {
	diag.AddError("Update not supported", "Acl Profile Client Connect Exceptions cannot be updated")
	return nil, nil
}

func (r aclProfileClientConnectExceptionResource) Delete(data *MsgVpnAclProfileClientConnectException, diag *diag.Diagnostics) (*http.Response, error) {
	_, httpResponse, err := r.Client.AclProfileApi.
		DeleteMsgVpnAclProfileClientConnectException(
			r.Context, *data.MsgVpnName, *data.AclProfileName, *data.ClientConnectExceptionAddress).
		Execute()
	return httpResponse, err
}

var msgVpnAclProfileClientConnectExceptionImportRegexp *regexp.Regexp = regexp.MustCompile(
	"^([^\\s\\*\\?\\/]+)\\/([0-9a-zA-Z_\\-]+)\\/([^\\s]+)$")

func (r aclProfileClientConnectExceptionResource) Import(id string, data *MsgVpnAclProfileClientConnectException, diag *diag.Diagnostics) {
	match := msgVpnAclProfileClientConnectExceptionImportRegexp.FindStringSubmatch(id)
	if match != nil {
		data.MsgVpnName = &match[1]
		data.AclProfileName = &match[2]
		data.ClientConnectExceptionAddress = &match[3]
	} else {
		diag.AddError("Expected <vpn-name>/<acl-profile>/<cidr address>",
			id+" does not match "+msgVpnAclProfileClientConnectExceptionImportRegexp.String())
		return
	}
}
