package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ provider.ResourceType = aclProfileResourceType{}

type aclProfileResourceType struct {
}

func (t aclProfileResourceType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	solaceProvider, diags := convertProviderType(in)

	return NewResource[MsgVpnAclProfile](
		aclProfileResource{solaceProvider: solaceProvider}), diags
}

func (t aclProfileResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return MsgVpnAclProfileSchema("msg_vpn_name", "acl_profile_name"), nil
}

var _ solaceProviderResource[MsgVpnAclProfile] = aclProfileResource{}

type aclProfileResource struct {
	solaceProvider
}

func (r aclProfileResource) NewData() *MsgVpnAclProfile {
	return &MsgVpnAclProfile{}
}

func (r aclProfileResource) Create(data *MsgVpnAclProfile, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.AclProfileApi.
		CreateMsgVpnAclProfile(r.Context, *data.MsgVpnName).
		Body(*data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r aclProfileResource) Read(data *MsgVpnAclProfile, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.AclProfileApi.GetMsgVpnAclProfile(r.Context, *data.MsgVpnName, *data.AclProfileName)
	apiResponse, httpResponse, err := apiReq.Execute()
	if err == nil && apiResponse != nil && apiResponse.Data != nil {
		data.ToTF(apiResponse.Data)
	}
	return httpResponse, err
}

func (r aclProfileResource) Update(_ *MsgVpnAclProfile, data *MsgVpnAclProfile, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.AclProfileApi.
		UpdateMsgVpnAclProfile(r.Context, *data.MsgVpnName, *data.AclProfileName).
		Body(*data.ToApi())
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

func (r aclProfileResource) Delete(data *MsgVpnAclProfile, diag *diag.Diagnostics) (*http.Response, error) {
	apiReq := r.Client.AclProfileApi.DeleteMsgVpnAclProfile(r.Context, *data.MsgVpnName, *data.AclProfileName)
	_, httpResponse, err := apiReq.Execute()
	return httpResponse, err
}

var msgVpnAclProfileImportRegexp *regexp.Regexp = regexp.MustCompile(fmt.Sprintf(
	"^([^\\s%s]+)\\/([0-9a-zA-Z_\\-]+)$", regexp.QuoteMeta("*?/")))

func (r aclProfileResource) Import(id string, data *MsgVpnAclProfile, diag *diag.Diagnostics) {
	match := msgVpnAclProfileImportRegexp.FindStringSubmatch(id)
	if match != nil {
		data.MsgVpnName = &match[1]
		data.AclProfileName = &match[2]
	} else {
		diag.AddError("Expected <vpn-name>/<acl-profile>", id+" does not match "+msgVpnAclProfileImportRegexp.String())
		return
	}
}
