package provider

import (
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func NewMsgVpnAclProfileClientConnectExceptionResource() resource.Resource {
	return &solaceResource[MsgVpnAclProfileClientConnectException]{spr: &aclProfileClientConnectExceptionResource{}}
}

var _ solaceProviderResource[MsgVpnAclProfileClientConnectException] = &aclProfileClientConnectExceptionResource{}

type aclProfileClientConnectExceptionResource struct {
	*solaceProvider
}

func (r aclProfileClientConnectExceptionResource) Name() string {
	return "aclprofile_client_connect_exception"
}

func (r aclProfileClientConnectExceptionResource) Schema() schema.Schema {
	return MsgVpnAclProfileClientConnectExceptionResourceSchema("msg_vpn_name", "acl_profile_name", "client_connect_exception_address")
}

func (r *aclProfileClientConnectExceptionResource) SetProvider(provider *solaceProvider) {
	r.solaceProvider = provider
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
